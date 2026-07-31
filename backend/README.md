# Mezun Anket Sistemi — Backend (Go)

Gin + GORM + PostgreSQL. Harici kuyruk teknolojisi (Redis/RabbitMQ) kullanılmaz;
asenkron mail gönderimi `email_outboxes` tablosu üzerinden, uygulama içindeki
bir goroutine (`internal/mail/worker.go`) tarafından yapılır.

## Kurulum

```bash
cd backend
cp .env.example .env
# .env içindeki ENCRYPTION_KEY için:
openssl rand -hex 32
# çıkan değeri .env dosyasındaki ENCRYPTION_KEY'e yapıştırın

go mod tidy   # bağımlılıkları indirir (internet bağlantısı gerekir)
```

PostgreSQL'i çalıştırın (örnek, Docker ile):

```bash
docker run --name mezun-anket-db -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=mezun_anket -p 5432:5432 -d postgres:16
```

## Çalıştırma

```bash
# .env dosyasını yükleyip çalıştırmak için `direnv` veya `godotenv` kullanabilir,
# ya da değişkenleri doğrudan export edebilirsiniz.
export $(grep -v '^#' .env | xargs)
go run ./cmd/api
```

İlk açılışta:
1. `AutoMigrate` tüm tabloları oluşturur.
2. `internal/seed` 5 kategori + 50 soruyu bir kere yükler.
3. Mail worker goroutine'i arka planda başlar (her 15 sn'de bir `email_outboxes` tablosunu tarar).

## İlk yönetici kullanıcısını oluşturma

```bash
go run ./cmd/createadmin -username=admin -password=guclu-bir-sifre -role=admin
```

## OBS'den mezun verisi içe aktarma

```
POST /api/v1/admin/graduates/import
Authorization: Bearer <admin-jwt>

{
  "sendInvites": true,
  "graduates": [
    {
      "obsHashId": "a1b2c3...",
      "firstName": "Ahmet",
      "facultyName": "Mühendislik Fakültesi",
      "departmentName": "Bilgisayar Mühendisliği",
      "graduationYear": 2021,
      "studentNoHash": "d4e5f6...",
      "email": "ahmet@example.com",
      "phone": "5xx xxx xx xx"
    }
  ]
}
```

- `soyad` alanı yok — OBS zaten göndermiyor, sistem de talep etmiyor.
- `email`/`phone` düz metin olarak gönderilir, backend'e ulaştığı anda AES-256-GCM ile şifrelenip öyle saklanır; hiçbir zaman düz metin DB'ye yazılmaz.
- `sendInvites: true` ise her mezun için tek kullanımlık giriş linki üretilir ve davet e-postası `email_outboxes` kuyruğuna eklenir.

## Uç noktalar

Tüm uç noktaların listesi ve standart hata formatı için ana mimari dokümanına (`mezun-anket-mimari.md`) bakınız.

## Test / doğrulama notu

Bu ortamda internet erişimi (dolayısıyla `go mod tidy`/`go build`) çalıştırılamadığından
kod bu ortamda derlenip test edilememiştir. Kodu kendi makinenizde/CI'de
`go build ./...` ile derleyip küçük düzeltmeler gerekip gerekmediğini kontrol edin.
