# Mersin Üniversitesi Mezun Takip Anket Sistemi — Teknik Mimari

**Yığın:** Go (backend) + React (frontend) — Prisma/BullMQ/Node.js geçmediğim için not: mesajınızın sonunda yapıştırdığınız örnek promptlar Node.js/Prisma/BullMQ üzerinden yazılmış görünüyor. Ancak en başta net biçimde belirttiğiniz gereksinim **Go + React + DB-üzeri kuyruk (harici kuyruk teknolojisi yok)** olduğundan, bu doküman o gereksinimi esas alır. Go tarafında Prisma'nın karşılığı olarak **GORM** (veya `ent`) kullanılacak.

Tasarım referansı olarak `tanitim.mersin.edu.tr` incelendi: kurumsal ana renk **`#2c3e7a`** (koyu lacivert) olarak öne çıkıyor; UI kütüphanesi teması bu tonu temel alacak.

---

## 1. Genel Mimari

```
┌────────────┐      HTTPS/JSON      ┌───────────────────┐        ┌──────────────┐
│  React SPA │ ───────────────────► │   Go API (Gin/Fiber)│ ─────► │  PostgreSQL  │
│ (Vite+TS)  │ ◄─────────────────── │  + GORM             │ ◄───── │              │
└────────────┘                      └─────────┬───────────┘        └──────┬───────┘
                                               │                            │
                                               │ goroutine (ticker, 30s)    │
                                               ▼                            │
                                     ┌────────────────────┐                 │
                                     │  Mail Worker        │◄────────────────┘
                                     │  (email_outbox      │  SELECT ... FOR UPDATE SKIP LOCKED
                                     │  tablosunu okur/     │
                                     │  SMTP ile gönderir)  │
                                     └────────────────────┘
```

- Tek binary: API sunucusu + arka plan mail worker'ı aynı Go process içinde ayrı bir goroutine olarak çalışır. Ayrı bir Redis/RabbitMQ/BullMQ **yok**; kuyruk PostgreSQL tablosu üzerinde tutulur.
- OBS'den gelen veri, ayrı bir **import endpoint'i / batch job** ile `graduates` tablosuna yazılır (aşağıda §2).

---

## 2. Veri Modeli (GORM)

### 2.1 OBS'den gelen mezun verisi (`graduates`)

OBS'den **soyad gelmeyecek**, kimlik alanı **hash string** olacak. Diğer alanlar (ad, bölüm, mezuniyet yılı vb.) öngörülen formatta gelecek.

```go
// internal/domain/graduate.go
type Graduate struct {
    ID              uint      `gorm:"primaryKey"`
    OBSHashID       string    `gorm:"uniqueIndex;size:128;not null"` // OBS'den gelen hash — girişte token/lookup anahtarı
    FirstName       string    `gorm:"size:100;not null"`
    // Soyad OBS'den gelmiyor: sistemde tutulmuyor, gerekiyorsa mezun ilk girişte
    // kendi rızasıyla profilini tamamlarken girer (opsiyonel, ayrı tablo: graduate_profile_extra)
    FacultyName     string    `gorm:"size:150"`
    DepartmentName  string    `gorm:"size:150"`
    GraduationYear  int
    StudentNoHash   string    `gorm:"size:128;index"` // gerekiyorsa öğrenci no da hash'li
    EmailEnc        []byte    `gorm:"type:bytea"`     // AES-GCM ile şifreli (bkz. §3)
    EmailNonce      []byte    `gorm:"type:bytea"`
    PhoneEnc        []byte    `gorm:"type:bytea"`
    PhoneNonce      []byte    `gorm:"type:bytea"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

> `soyad` alanı DB şemasında yok — KVKK'daki "veri minimizasyonu" ilkesine uygun: OBS zaten göndermiyor, sistem de talep etmiyor, sadece anket ve karşılama metinlerinde ada + bölüme + yıla göre kişiselleştirme yapılıyor.

### 2.2 Anket şeması (5 kategori, 50 soru — statik + esnek)

Sorular kod içinde sabit değil, DB'de tutulur ki ileride değiştirilebilsin:

```go
type SurveyCategory struct {
    ID    uint   `gorm:"primaryKey"`
    Order int    `gorm:"not null"`      // 1..5 (Stepper adım sırası)
    Title string `gorm:"size:150"`      // "Üniversite Deneyimi" vb.
}

type SurveyQuestion struct {
    ID          uint   `gorm:"primaryKey"`
    CategoryID  uint   `gorm:"index;not null"`
    Order       int    `gorm:"not null"`   // kategori içi sıra
    Code        string `gorm:"size:10"`    // "Q01".."Q50"
    Text        string `gorm:"type:text"`
    AnswerType  string `gorm:"size:20"`    // "scale_1_10" | "single_choice" | "multi_choice" | "text" | "duration_months"
    OptionsJSON string `gorm:"type:jsonb"` // choice sorularının seçenekleri
    Required    bool   `gorm:"default:true"`
}

type SurveyResponse struct {
    ID          uint      `gorm:"primaryKey"`
    GraduateID  uint      `gorm:"uniqueIndex;not null"` // kişi başı tek anket
    Status      string    `gorm:"size:20;default:'in_progress'"` // in_progress|completed
    CurrentStep int       `gorm:"default:1"`  // kaldığı adım — yarım bırakırsa devam edebilsin
    StartedAt   time.Time
    CompletedAt *time.Time
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type SurveyAnswer struct {
    ID          uint   `gorm:"primaryKey"`
    ResponseID  uint   `gorm:"index:idx_resp_q,unique"`
    QuestionID  uint   `gorm:"index:idx_resp_q,unique"`
    ValueText   string `gorm:"type:text"`   // her tip cevap burada normalize edilir (JSON/string)
    UpdatedAt   time.Time
}
```

Bu normalize yapı sayesinde 50 soruyu 50 ayrı kolon açmak yerine tek `survey_answers` tablosunda tutuyoruz; hem şema değişikliklerine dayanıklı hem de admin tarafında `GROUP BY question_id` ile kolayca istatistik çıkar (pasta/bar grafik).

### 2.3 Şifreli iletişim bilgisi & KVKK

- `EmailEnc/PhoneEnc` DB'de **hiçbir zaman düz metin tutulmaz**.
- Uygulama içinde (mail gönderimi, admin'in "iletişim bilgisi görüntüleme" ekranı gibi *yetkili* akışlarda) decrypt edilip kullanılır, response cache'lenmez.
- Anahtar (`ENCRYPTION_KEY`) ortam değişkeninden/secret manager'dan okunur, koda gömülmez.

### 2.4 Admin & yetkilendirme

```go
type AdminUser struct {
    ID           uint   `gorm:"primaryKey"`
    Username     string `gorm:"uniqueIndex;size:100"`
    PasswordHash string `gorm:"size:255"` // bcrypt
    Role         string `gorm:"size:20"`  // "admin" | "viewer"
}
```

### 2.5 DB-üzerinde asenkron mail kuyruğu

```go
type EmailOutbox struct {
    ID          uint       `gorm:"primaryKey"`
    ToEmailEnc  []byte     `gorm:"type:bytea"`   // hedef adres de şifreli saklanır
    ToEmailNonce []byte    `gorm:"type:bytea"`
    Subject     string     `gorm:"size:255"`
    Body        string     `gorm:"type:text"`
    Status      string     `gorm:"size:20;default:'pending'"` // pending|processing|sent|failed
    Attempts    int        `gorm:"default:0"`
    LastError   string     `gorm:"type:text"`
    LockedAt    *time.Time
    SendAfter   time.Time  `gorm:"index"` // retry/backoff için
    CreatedAt   time.Time
    SentAt      *time.Time
}
```

Worker döngüsü (basit, harici bağımlılık yok):

```go
func (w *MailWorker) Run(ctx context.Context) {
    ticker := time.NewTicker(15 * time.Second)
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            w.processBatch(ctx)
        }
    }
}

func (w *MailWorker) processBatch(ctx context.Context) {
    var jobs []EmailOutbox
    err := w.db.Transaction(func(tx *gorm.DB) error {
        // PostgreSQL: SKIP LOCKED sayesinde birden fazla instance çalışsa bile
        // aynı satırı iki worker birden almaz -> ekstra kuyruk teknolojisine gerek kalmaz
        if err := tx.Raw(`
            SELECT * FROM email_outboxes
            WHERE status = 'pending' AND send_after <= now()
            ORDER BY id
            LIMIT 20
            FOR UPDATE SKIP LOCKED
        `).Scan(&jobs).Error; err != nil {
            return err
        }
        ids := idsOf(jobs)
        return tx.Model(&EmailOutbox{}).Where("id IN ?", ids).
            Updates(map[string]any{"status": "processing", "locked_at": time.Now()}).Error
    })
    if err != nil { /* log */ return }

    for _, job := range jobs {
        email, _ := crypto.Decrypt(job.ToEmailEnc, job.ToEmailNonce)
        if err := w.smtp.Send(email, job.Subject, job.Body); err != nil {
            w.db.Model(&job).Updates(map[string]any{
                "status": "pending", "attempts": job.Attempts + 1,
                "last_error": err.Error(),
                "send_after": time.Now().Add(backoff(job.Attempts)),
            })
            continue
        }
        w.db.Model(&job).Updates(map[string]any{"status": "sent", "sent_at": time.Now()})
    }
}
```

`SELECT ... FOR UPDATE SKIP LOCKED` PostgreSQL'in yerleşik özelliği; harici kuyruk sunucusu kurmadan, birden fazla API instance'ı çalışsa bile güvenli asenkron işleme sağlar.

---

## 3. Şifreleme Yaklaşımı (AES-256-GCM)

```go
// internal/crypto/crypto.go
func Encrypt(plain string, key []byte) (cipherText, nonce []byte, err error) {
    block, err := aes.NewCipher(key) // key: 32 byte
    if err != nil { return nil, nil, err }
    gcm, err := cipher.NewGCM(block)
    if err != nil { return nil, nil, err }
    nonce = make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil { return nil, nil, err }
    cipherText = gcm.Seal(nil, nonce, []byte(plain), nil)
    return cipherText, nonce, nil
}

func Decrypt(cipherText, nonce []byte, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil { return "", err }
    gcm, err := cipher.NewGCM(block)
    if err != nil { return "", err }
    plain, err := gcm.Open(nil, nonce, cipherText, nil)
    if err != nil { return "", err }
    return string(plain), nil
}
```

- Bu fonksiyonlar `service` katmanında (GORM `Hooks`'ta değil — hook'larda anahtar/servis erişimi karmaşıklaşır) kullanılır: `GraduateService.Create()`, `.GetDecrypted()` gibi.
- Decrypt edilmiş değer **hiçbir zaman** loglara veya API response'larına ham biçimde yazılmaz; admin panelinde iletişim bilgisi görüntüleme ayrı, loglanan bir yetkili işlemdir (audit log önerilir: kim, ne zaman, kimin bilgisini gördü).

---

## 4. Kimlik Doğrulama & Kişiselleştirilmiş Giriş

Mezun, OBS'den türetilen **tekil token** (ör. `OBSHashID` üzerinden üretilmiş imzalı bir link/kod) ile giriş yapar — şifre yönetimi yok:

1. Üniversite, mezuna "anketi doldurmak için" bir bağlantı gönderir: `https://anket.mersin.edu.tr/giris?token=...`
2. Backend token'ı doğrular (JWT, kısa ömürlü + tek kullanımlık DB kaydıyla eşleştirilmiş), `Graduate` kaydını bulur.
3. Frontend `GET /api/v1/me` ile ad, bölüm, mezuniyet yılı bilgisini çeker ve **"Hoşgeldin {FirstName}, {GraduationYear} yılı {DepartmentName} mezunumuz!"** karşılama ekranını gösterir.
4. `SurveyResponse.CurrentStep` sayesinde yarıda bırakan mezun tekrar girdiğinde kaldığı adımdan devam eder (bounce oranını azaltan kritik detay).

Admin girişi ayrı: kullanıcı adı/şifre + bcrypt + JWT, rol bazlı yetki (`admin`/`viewer`).

---

## 5. API Uç Noktaları & Standart Hata Formatı

### 5.1 Standart response zarfı

```json
// Başarılı
{ "success": true, "data": { ... } }

// Hatalı
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Bazı alanlar eksik veya hatalı.",
    "details": [{ "field": "Q07", "message": "Bu soru zorunludur." }]
  }
}
```

| HTTP | `error.code`         | Kullanım alanı                                   |
|------|----------------------|---------------------------------------------------|
| 400  | `VALIDATION_ERROR`   | Form/step doğrulama hataları                       |
| 401  | `UNAUTHORIZED`       | Token yok/geçersiz/süresi dolmuş                   |
| 403  | `FORBIDDEN`          | Yetkisiz rol (viewer'ın admin işlemi denemesi vb.) |
| 404  | `NOT_FOUND`          | Mezun/anket kaydı bulunamadı                       |
| 409  | `CONFLICT`           | Anket zaten tamamlanmış, tekrar submit denemesi    |
| 500  | `INTERNAL_ERROR`     | Beklenmeyen sunucu hatası                          |

Go tarafında merkezi bir middleware ile tüm handler'lar bu zarfa sarılır; her handler sadece `AppError{Code, HTTPStatus, Message, Details}` döner, middleware JSON'a çevirir. React tarafında tek bir axios interceptor bu `error.code`'a göre kullanıcıya toast/redirect gösterir (401 → login'e at, 403 → "yetkiniz yok" ekranı, 400 → alan bazlı form hatası).

### 5.2 Uç nokta listesi

**Mezun tarafı**
| Method | Path | Açıklama |
|---|---|---|
| POST | `/api/v1/auth/token-login` | Token doğrulama, session/JWT üretir |
| GET  | `/api/v1/me` | Kişiselleştirme için ad/bölüm/yıl |
| GET  | `/api/v1/survey/structure` | 5 kategori + 50 soru (choice seçenekleriyle) |
| GET  | `/api/v1/survey/response` | Var olan cevapları + `currentStep` döner (kaldığı yerden devam) |
| PUT  | `/api/v1/survey/response/step/{stepNo}` | O adımdaki cevapları kaydeder (autosave / "İleri" butonu) |
| POST | `/api/v1/survey/response/complete` | Son adım sonrası anketi `completed` yapar, teşekkür maili kuyruğa eklenir |

**Admin tarafı**
| Method | Path | Açıklama |
|---|---|---|
| POST | `/api/v1/admin/login` | Admin girişi |
| GET  | `/api/v1/admin/stats/overview` | Toplam katılım, tamamlanma oranı, adım bazlı terk oranı |
| GET  | `/api/v1/admin/stats/sector-distribution` | Soru 21 bazlı sektör dağılımı (pasta grafik verisi) |
| GET  | `/api/v1/admin/stats/salary-satisfaction` | Soru 24 bazlı gelir memnuniyeti (bar grafik) |
| GET  | `/api/v1/admin/stats/question/{code}` | Herhangi bir sorunun dağılımı (genel amaçlı) |
| GET  | `/api/v1/admin/graduates/import-status` | OBS import geçmişi |
| POST | `/api/v1/admin/graduates/import` | OBS'den toplu veri alma (CSV/JSON) |

---

## 6. Frontend Mimarisi (React)

- **Stack:** Vite + React + TypeScript, `react-hook-form` + `zod` (adım bazlı validasyon), `react-router` (adımlar arası URL senkronu, ör. `/anket/2`).
- **UI kütüphanesi:** **Mantine** veya **Chakra UI** — ikisi de tema (`theme.colors.primary`) üzerinden kurumsal renge (`#2c3e7a`) kolayca uyarlanır, hazır `Stepper`/`Progress` bileşenleri var (özellikle Mantine'in `Stepper` bileşeni bu iş için birebir).
- **Sayfa akışı:**
  1. `/giris` — token doğrulama (otomatik, link tıklanınca)
  2. `/hosgeldin` — kişiselleştirilmiş karşılama ("Hoşgeldin Ahmet, 2021 Bilgisayar Mühendisliği mezunumuz!")
  3. `/anket` — **5 adımlı Stepper**, üstte sabit **Progress Bar** (`tamamlanan soru / 50`)
     - Adım 1: Üniversite Deneyimi (S1–S10)
     - Adım 2: İstihdam Süreci (S11–S20)
     - Adım 3: Çalışma Şartları (S21–S30)
     - Adım 4: Eğitim Uyumu (S31–S40)
     - Adım 5: Kariyer Planları (S41–S50)
     - Her adımda "İleri" tıklanınca o adımın soruları `zod` şemasıyla doğrulanır, sonra `PUT /survey/response/step/{n}` çağrılır (adım bazlı autosave → tarayıcı kapansa da veri kaybolmaz).
  4. `/tesekkurler` — tamamlama ekranı.
- **Mobile-first**: Mantine/Chakra grid'i tek sütun mobilde, çok sütun masaüstünde; soru başına tek kart, büyük dokunma alanları (özellikle 1-10 ölçek soruları için buton grubu, dropdown değil).

### Örnek Stepper iskeleti (Mantine)

```tsx
function SurveyWizard() {
  const [active, setActive] = useState(response.currentStep - 1);
  const answeredCount = useAnsweredCount(); // toplam cevaplanan soru sayısı

  return (
    <>
      <Progress value={(answeredCount / 50) * 100} radius="xl" />
      <Stepper active={active} onStepClick={setActive} allowNextStepsSelect={false}>
        {categories.map((cat) => (
          <Stepper.Step key={cat.id} label={cat.title}>
            <StepForm category={cat} onNext={() => setActive((s) => s + 1)} />
          </Stepper.Step>
        ))}
      </Stepper>
    </>
  );
}
```

---

## 7. Admin Dashboard (Recharts)

- `PieChart` — Soru 21 (sektör) dağılımı
- `BarChart` — Soru 27 (haftalık mesai) ortalaması, bölüm bazlı kırılım
- `RadialBarChart`/`Gauge` — anket tamamlanma oranı ve adım bazlı terk (bounce) noktaları — hangi adımda en çok bırakıldığını görmek `SurveyResponse.CurrentStep` dağılımından çıkarılır.
- Tüm agregasyonlar backend'de `GROUP BY` ile hesaplanır, ham cevap/iletişim verisi asla frontend'e dönmez (KVKK — sadece istatistiksel/anonim veri).

---

## 8. KVKK Uyumluluğu Notları

- Soyad hiçbir katmanda tutulmaz (OBS zaten göndermiyor).
- Telefon/e-posta DB'de şifreli; decrypt sadece yetkili, loglanan işlemlerde.
- Anket sonuçları hiçbir dış servise gönderilmez; sadece kendi PostgreSQL'inizde durur.
- Aydınlatma metni + açık rıza onayı, mezun ilk girişte (karşılama ekranında) tek seferlik gösterilip onaylatılmalı ve onay zaman damgası `GraduateConsent` tablosunda saklanmalı (öneri, şemaya eklenebilir).
- Audit log: admin'in kişisel veriye (decrypt edilmiş iletişim bilgisi) her erişimi kayıt altına alınmalı.

---

## 9. Klasör Yapısı (öneri)

```
backend/
  cmd/api/main.go
  internal/
    domain/          (Graduate, SurveyResponse, EmailOutbox ...)
    repository/      (GORM sorguları)
    service/         (iş mantığı, şifreleme çağrıları)
    handler/         (HTTP handler'lar, Gin/Fiber)
    middleware/       (auth, error-envelope, logging)
    mail/            (SMTP client + worker)
    crypto/
  migrations/
frontend/
  src/
    pages/ (Giris, Hosgeldin, Anket, Tesekkurler, Admin/*)
    components/ (Stepper adımları, sorular tip bazlı render eden <QuestionField />)
    api/ (axios instance + interceptor)
    theme.ts (#2c3e7a bazlı Mantine teması)
```

---

Bu doküman; veri modelini, güvenlik/şifreleme yaklaşımını, DB-üzeri asenkron mail kuyruğunu, hata yönetimi standardını ve stepper tabanlı frontend akışını uçtan uca kapsıyor. İsterseniz bir sonraki adımda bu şemadan gerçek çalışan bir `schema.sql`/GORM migration seti ve örnek Gin handler'ları olan başlangıç kod iskeletini de oluşturabilirim.
