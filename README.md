# MERSİN ÜNİVERSİTESİ — MEZUN TAKİP VE BİLGİ SİSTEMİ (ALUMNI TRACKING SYSTEM)

Bu doküman, Mersin Üniversitesi Mühendislik Fakültesi staj programı kapsamında kurumun ihtiyaçlarına yönelik olarak tasarlanan ve geliştirilen **Mezun Takip ve Anket Sistemi'nin** teknik altyapısını, mimari kararlarını ve kurulum süreçlerini içermektedir.

## 1. PROJENİN KAPSAMI VE AMACI (NE YAPTIK?)

Mersin Üniversitesi mezunlarının kariyer yollarını izlemek, istihdam sürelerini analiz etmek ve üniversitede verilen eğitimin sektörel karşılığını istatistiksel olarak ölçmek amacıyla **Tam Yığın (Full-Stack) bir web uygulaması** geliştirilmiştir. 

Proje, manuel olarak yürütülen veya üçüncü parti (Google Forms vb.) yazılımlarla sağlanan anket süreçlerini, üniversitenin kendi sunucularında barındırabileceği, veri güvenliğini sağlayan ve özel istatistikler üretebilen kapalı bir ekosisteme taşımaktadır.

---

## 2. MİMARİ TASARIM KARARLARI (NEDEN VE NASIL YAPTIK?)

Uygulama; ölçeklenebilirlik, veri tutarlılığı ve bağımsız geliştirme/dağıtım (loose coupling) ilkeleri göz önünde bulundurularak "İstemci-Sunucu İzolasyonu" (Client-Server Separation) mimarisine göre tasarlanmıştır.

### 2.1. Sunucu Katmanı (Backend API & Worker)
*   **Teknoloji:** Go (Golang) & Gin Web Framework
*   **Neden Seçildi?** Go dilinin eşzamanlılık (concurrency) yetenekleri, düşük bellek tüketimi ve yüksek I/O performansı nedeniyle tercih edilmiştir. Özellikle arka planda çalışacak olan "Mail Worker" sisteminin ana HTTP iş parçacıklarını (thread) bloklamaması için Go Goroutines kullanılmıştır.
*   **Nasıl Kurgulandı?** RESTful API mimarisi standartlarına (JSON payload, HTTP statü kodları) uygun olarak yapılandırılmıştır. Veritabanı işlemleri için GORM kullanılarak SQL Injection zafiyetleri önlenmiş ve nesne-ilişkisel eşleme (ORM) sağlanmıştır.

### 2.2. İstemci Katmanı (Frontend)
*   **Teknoloji:** React (Vite altyapısı), TypeScript, Mantine UI
*   **Neden Seçildi?** 50 soruluk karmaşık bir anketin durum yönetimini (state management) sayfa yenilenmeden yapabilmek için React (Single Page Application) kullanılmıştır. Tip güvenliğini (type safety) sağlamak ve çalışma zamanı (runtime) hatalarını minimuma indirmek amacıyla TypeScript tercih edilmiştir.
*   **Nasıl Kurgulandı?** Kullanıcı deneyimini (UX) artırmak adına anket, dinamik adımlar (Stepper) halinde bölünmüştür. Her adım geçişinde veriler asenkron (AJAX/Fetch) olarak sunucuya iletilerek, kullanıcının anketi yarıda bırakması durumunda "veri kaybı" yaşanmasının önüne geçilmiştir.

### 2.3. Sistem Dağıtımı ve İzolasyon (DevOps)
*   **Teknoloji:** Docker, Docker Compose, Nginx
*   **Neden Seçildi?** Projenin farklı geliştirici bilgisayarlarında ("Benim bilgisayarımda çalışıyordu" problemini önlemek) ve sunucularda birebir aynı ortamda çalışmasını garanti altına almak için Docker kullanılmıştır.
*   **Nasıl Kurgulandı?** Veritabanı, Backend ve Frontend için ayrı konteynerler (containers) oluşturulmuştur. Frontend istekleri Nginx üzerinden "Reverse Proxy" (Ters Vekil Sunucu) ile yönetilerek API'ye yönlendirilmiş, böylece olası CORS problemlerinin önüne geçilmiş ve yük dağıtımı için altyapı hazırlanmıştır.

---

## 3. TEMEL SİSTEM MODÜLLERİ

1.  **Dinamik Anket ve Veri Kurtarma (Stepper/Wizard):** Mezun anketi doldururken, bağlantı kopsa dahi veriler adım bazlı (step-by-step) kaydedilir.
2.  **Terk (Bounce Rate) Analizi:** Yöneticiler, anketi yarıda bırakan mezunların istatistiksel olarak en çok hangi soruda/adımda sistemi terk ettiğini analiz edebilir.
3.  **Asenkron E-Posta Kuyruğu (Mail Worker):** Mezunlara gönderilecek davet veya hatırlatma e-postaları ana akışı yavaşlatmamak adına bir kuyruğa (outbox table) yazılır. Arka planda çalışan ayrı bir Goroutine, 15 saniyede bir bu kuyruğu tarayarak SMTP üzerinden iletimi asenkron olarak gerçekleştirir.
4.  **JWT Tabanlı Güvenlik:** Yönetim paneline erişim, durumsuz (stateless) ve zaman aşımı mekanizmasına sahip JSON Web Token (JWT) ile güvence altına alınmıştır.

---

## 4. KURULUM VE TEST ORTAMI (LOCAL DEVELOPMENT)

Projeyi yerel ortamda çalıştırmak için aşağıdaki adımların sırasıyla uygulanması gerekmektedir.

### 4.1. Deponun Klonlanması
```bash
git clone [https://github.com/zeynepbetulyilmaz/mezun-anket-sistemi.git](https://github.com/zeynepbetulyilmaz/mezun-anket-sistemi.git)
cd mezun-anket-sistemi