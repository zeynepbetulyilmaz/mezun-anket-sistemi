# Mezun Anket Sistemi — Frontend (React + Vite + TS + Mantine)

## Kurulum

```bash
cd frontend
npm install
npm run dev
```

`vite.config.ts` içindeki proxy ayarı sayesinde `/api/*` istekleri geliştirme
sırasında otomatik olarak `http://localhost:8080` adresindeki Go API'sine yönlenir.

## Sayfa akışı

1. `/giris?token=...` — mezuna e-postayla gönderilen tek kullanımlık link
2. `/hosgeldin` — kişiselleştirilmiş karşılama ("Hoşgeldin {Ad}, {Yıl} {Bölüm} mezunumuz!")
3. `/anket` — 5 adımlı Stepper + üstte ilerleme çubuğu, her adımda autosave
4. `/tesekkurler` — tamamlama ekranı
5. `/admin/giris` → `/admin` (Panel) / `/admin/mezun-ekle` (Mezun Ekle) — kullanıcı adı/şifre
   ile giriş sonrası ortak üst menüyü (`AdminLayout.tsx`) paylaşan yönetim sayfaları:
   - **Panel**: sektör/çalışma modeli dağılımı, tamamlanma oranı, adım bazlı terk analizi (Recharts)
   - **Mezun Ekle**: OBS verisini CSV/Excel dosyasıyla toplu ya da manuel form ile tek tek
     `POST /api/v1/admin/graduates/import` uç noktasına gönderir (`ImportGraduates.tsx`)

## Tasarım

Kurumsal renk (`#2c3e7a`, tanitim.mersin.edu.tr temel alınarak) `src/theme.ts` içinde
Mantine paletine uyarlanmıştır. Kurum logosu eklemek için `public/meu-logo.png` dosyasını
yerleştirip `TokenLogin.tsx`/`Welcome.tsx` içindeki yer tutucuya bağlayabilirsiniz.

## Not

Bu ortamda internet erişimi olmadığından `npm install`/`npm run build` bu ortamda
çalıştırılıp doğrulanamamıştır. Kendi makinenizde `npm install` sonrası
`npm run build` ile tip hatası olup olmadığını kontrol edin.
