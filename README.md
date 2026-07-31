# Mersin Üniversitesi — Mezun Takip Anket Sistemi

Go (Gin + GORM + PostgreSQL) backend ve React (Vite + TS + Mantine) frontend'den oluşan,
50 soruluk mezun anketini 5 adımlı bir Stepper/Wizard olarak sunan tam yığın uygulama.

```
mezun-anket-sistemi/
├── backend/     Go API, DB-üzeri asenkron mail kuyruğu, şifreleme, seed
├── frontend/    React SPA - kişiselleştirilmiş karşılama, stepper anket, admin dashboard
└── mezun-anket-mimari.md   Ayrıntılı mimari/tasarım dokümanı
```

## Hızlı başlangıç

```bash
# 1) Veritabanı
docker run --name mezun-anket-db -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=mezun_anket -p 5432:5432 -d postgres:16

# 2) Backend
cd backend
cp .env.example .env
# .env içine: openssl rand -hex 32  çıktısını ENCRYPTION_KEY olarak yapıştırın
go mod tidy
go run ./cmd/api
# ayrı bir terminalde ilk admin kullanıcısını oluşturun:
go run ./cmd/createadmin -username=admin -password=guclu-bir-sifre -role=admin

# 3) Frontend
cd ../frontend
npm install
npm run dev
```

Frontend: http://localhost:5173 · Backend: http://localhost:8080

## Öne çıkan mimari kararlar

- **Veri minimizasyonu**: OBS'den soyad gelmiyor, sistem de tutmuyor; kimlik alanı hash string.
- **Şifreleme**: Telefon/e-posta DB'de AES-256-GCM ile şifreli (`bytea`); decrypt sadece
  yetkili servis çağrılarında (mail gönderimi, admin görüntüleme) yapılır.
- **Harici kuyruk yok**: Asenkron mail gönderimi tamamen PostgreSQL üzerinde
  (`email_outboxes` tablosu + `SELECT ... FOR UPDATE SKIP LOCKED`) — Redis/RabbitMQ/BullMQ
  gerekmez.
- **Standart hata formatı**: Tüm API yanıtları `{success, data}` veya
  `{success:false, error:{code,message,details}}` zarfında; 400/401/403/404/409/500
  kodları frontend'de tek bir axios interceptor'da yönetilir.
- **Bounce oranını düşürme**: 50 soru tek sayfada değil, 5 kategoriye bölünmüş Stepper +
  üstte ilerleme çubuğu; her adım autosave edilir, kaldığı adımdan devam edilebilir.

## Derleme/test notu

Bu geliştirme ortamında internet erişimi olmadığından (`go mod tidy`, `npm install`
paket indirme gerektirir) kod bu ortamda derlenip test edilememiştir. Lütfen kendi
makinenizde `go build ./...` ve `npm run build` ile doğrulayın; küçük tip/import
düzeltmeleri gerekebilir.

Ayrıntılı mimari, veri modeli, API uç nokta listesi ve tasarım gerekçeleri için
`mezun-anket-mimari.md` dosyasına bakınız.
