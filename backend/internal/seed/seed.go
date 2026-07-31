// Package seed: uygulama ilk ayağa kalktığında anket kategorilerini ve
// soruları DB'ye yükler. TRUNCATE ile eski verileri temizler.
package seed

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"mezun-anket-backend/internal/domain"
)

type questionSeed struct {
	Code          string
	Text          string
	AnswerType    string
	Options       string // JSON dizi, choice tipleri için
	TargetFaculty string // Boş ise herkese, dolu ise sadece o fakülteye görünür
}

func Run(db *gorm.DB) error {
	db.Exec("TRUNCATE TABLE survey_questions, survey_categories RESTART IDENTITY CASCADE")
	//db.Exec("TRUNCATE TABLE survey_responses, survey_answers, survey_questions, survey_categories RESTART IDENTITY CASCADE")
	categories := []struct {
		Order     int
		Title     string
		Questions []questionSeed
	}{
		{1, "A. Demografik ve Genel Bilgiler", []questionSeed{
			{"A1", "1. Mezun olurken kariyer planlamanızdaki öncelikli hedefiniz neydi?", "single_choice", `["Özel Sektör (Kurumsal)", "Kamu Kurumu", "Kendi İşimi Kurmak / Girişimcilik", "Akademik Kariyer"]`, ""},
			{"A2", "2. Mezuniyet not ortalamanız (GANO) hangi aralıktaydı?", "single_choice", `["3.50 - 4.00", "3.00 - 3.49", "2.50 - 2.99", "2.00 - 2.49"]`, ""},
			{"A3", "3. Üniversite eğitiminiz süresince herhangi bir öğrenci değişim programından (Erasmus, Farabi, Mevlana vb.) faydalandınız mı?", "single_choice", `["Evet", "Hayır"]`, ""},
			{"A4", "4. Öğrencilik döneminizde KYK, üniversite veya özel kurumlardan burs desteği aldınız mı?", "single_choice", `["Evet", "Hayır"]`, ""},
			{"A5", "5. Üniversitemizin düzenleyeceği mezun etkinliklerinden (kariyer günleri, mezun buluşmaları vb.) haberdar olmak istiyor musunuz?", "single_choice", `["Evet, e-posta almak istiyorum", "Hayır, istemiyorum"]`, ""},
		}},
		{2, "B. Eğitim Deneyimi Değerlendirmesi", []questionSeed{
			{"B6", "6. Aldığınız eğitimin teorik bilgi açısından yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["Çok İyi","İyi","Orta","Kötü","Çok Kötü"]`, ""},
			{"B7", "7. Aldığınız eğitimin uygulama/pratik beceri kazandırma açısından yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["Çok İyi","İyi","Orta","Kötü","Çok Kötü"]`, ""},
			{"B8", "8. Bölümünüzdeki öğretim elemanlarının akademik yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["Çok İyi","İyi","Orta","Kötü","Çok Kötü"]`, ""},
			{"B9", "9. Ders programının güncel sektör ihtiyaçlarına uygunluğunu nasıl değerlendirirsiniz?", "single_choice", `["Çok İyi","İyi","Orta","Kötü","Çok Kötü"]`, ""},
			{"B10", "10. Staj/uygulama olanaklarının yeterliliği hakkında ne düşünüyorsunuz?", "single_choice", `["Çok İyi","İyi","Orta","Kötü","Çok Kötü"]`, ""},
			{"B11", "11. Kütüphane, laboratuvar ve teknik altyapı olanaklarını nasıl değerlendirirsiniz?", "single_choice", `["Çok İyi","İyi","Orta","Kötü","Çok Kötü"]`, ""},
			{"B12", "12. Yabancı dil eğitiminin yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["Çok İyi","İyi","Orta","Kötü","Çok Kötü"]`, ""},
			{"B13", "13. Sosyal, kültürel ve sportif faaliyetlerin çeşitliliğini nasıl değerlendirirsiniz?", "single_choice", `["Çok İyi","İyi","Orta","Kötü","Çok Kötü"]`, ""},
			{"B14", "14. Akademik danışmanlık hizmetlerinden ne ölçüde memnun kaldınız?", "single_choice", `["Çok İyi","İyi","Orta","Kötü","Çok Kötü"]`, ""},
			{"B15", "15. Üniversitedeki eğitim kalitesini genel olarak nasıl değerlendirirsiniz?", "single_choice", `["Çok İyi","İyi","Orta","Kötü","Çok Kötü"]`, ""},
		}},
		{3, "C. İstihdam / İş Deneyimi", []questionSeed{
			{"C16", "16. Mezuniyet sonrası ilk işinizi bulma süreniz ne kadar sürdü?", "single_choice", `["Henüz bulamadım","İlk 3 ay","3-6 ay","6-12 ay","1 yıldan fazla"]`, ""},
			{"C17", "17. Şu anda bir işte çalışıyor musunuz?", "single_choice", `["Evet","Hayır"]`, ""},
			{"C18", "18. Çalıştığınız sektör, mezun olduğunuz alanla ne ölçüde örtüşüyor?", "single_choice", `["Tamamen","Kısmen","Hiç"]`, ""},
			{"C19", "19. İlk işinizi bulma sürecinde en çok hangi kaynak/yöntem etkili oldu?", "single_choice", `["İş İlanları (Kariyer.net, vb.)", "Staj Yaptığım Kurum", "Üniversite Kariyer Merkezi / Hocalar", "Tanıdık / Referans", "Kendi Girişimim"]`, ""},
			{"C20", "20. Şu anki pozisyonunuz/unvanınız hangi seviyeye daha yakın?", "single_choice", `["Yeni Başlayan (Junior / Uzman Yrd.)", "Orta Seviye (Mid-level / Uzman)", "Kıdemli (Senior / Yönetici)", "Serbest Çalışan (Freelancer)", "İş Yeri Sahibi"]`, ""},
			{"C21", "21. Aylık ortalama gelir aralığınız nedir?", "single_choice", `["Asgari Ücret", "Asgari Ücret - 1.5 Katı", "1.5 - 2 Katı", "2 - 3 Katı", "3 Katı ve Üzeri"]`, ""},
			{"C22", "22. İşyerinizde kariyer gelişimi ve terfi olanaklarını nasıl değerlendirirsiniz?", "single_choice", `["Çok İyi","İyi","Orta","Kötü","Çok Kötü"]`, ""},
			{"C23", "23. Mezun olduğunuz üniversite diplomasının iş bulma sürecinde katkısı oldu mu?", "single_choice", `["Evet","Hayır"]`, ""},
			{"C24", "24. Çalışma hayatında en çok zorlandığınız konu(lar) nelerdir?", "single_choice", `["Teorik bilginin pratiğe uymaması", "İş-yaşam dengesini kuramamak", "Düşük ücret / Uzun mesai", "İletişim ve takım çalışması", "Zorluk yaşamadım"]`, ""},
			{"C25", "25. Kendi işinizi kurdunuz mu / girişimcilik faaliyetinde bulundunuz mu?", "single_choice", `["Evet","Hayır"]`, ""},
		}},
		{4, "D. Mezuniyet Sonrası Kariyer Gelişimi", []questionSeed{
			{"D26", "26. Lisansüstü eğitime (yüksek lisans/doktora) devam ettiniz mi / ediyor musunuz?", "single_choice", `["Evet","Hayır"]`, ""},
			{"D27", "27. Mesleki sertifika, kurs veya ek eğitim aldınız mı?", "single_choice", `["Evet","Hayır"]`, ""},
			{"D28", "28. Mevcut işinizden memnuniyet düzeyiniz nedir?", "single_choice", `["Çok Memnunum","Memnunum","Kararsızım","Memnun Değilim","Hiç Memnun Değilim"]`, ""},
			{"D29", "29. Aldığınız eğitim, mesleki hedeflerinize ulaşmanızda ne ölçüde katkı sağladı?", "single_choice", `["Çok Katkısı Oldu","Biraz Katkısı Oldu","Hiç Katkısı Olmadı"]`, ""},
			{"D30", "30. Kariyeriniz boyunca kaç farklı işyerinde çalıştınız?", "single_choice", `["1 (Sadece ilk/mevcut işim)", "2 - 3", "4 - 5", "5'ten fazla", "Henüz çalışmadım"]`, ""},
			{"D31", "31. Yurt dışında çalışma/eğitim deneyiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, ""},
			{"D32", "32. Mezuniyet sonrası mesleki ağınızı (network) genişletmede üniversite bağlantılarınız yardımcı oldu mu?", "single_choice", `["Evet","Hayır"]`, ""},
			{"D33", "33. Şu anki işinizde, aldığınız eğitimle doğrudan ilgili olmayan en çok hangi beceriye ihtiyaç duydunuz?", "single_choice", `["Yabancı Dil", "Bilişim ve Yazılım Araçları", "İletişim, İkna ve Sunum", "Liderlik ve Ekip Yönetimi", "Finans ve Yönetim", "İhtiyaç Duymadım"]`, ""},
			{"D34", "34. Önümüzdeki 5 yıl içindeki kariyer hedefiniz nedir?", "single_choice", `["Mevcut sektörümde yönetici pozisyonuna yükselmek", "Kendi işimi kurmak / Girişimci olmak", "Akademik kariyer (Yüksek Lisans/Doktora) yapmak", "Yurt dışına gitmek / Yurt dışında çalışmak", "Farklı bir sektöre veya mesleğe geçiş yapmak", "Kamu kurumuna atanmak"]`, ""},
			{"D35", "35. Üniversitenin mezunlarla iletişimini (mezun takip sistemi, etkinlikler vb.) yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, ""},
		}},
		{5, "E. Üniversiteye Genel Değerlendirme ve Öneriler", []questionSeed{
			{"E36", "36. Üniversitenizi/bölümünüzü başkalarına tavsiye eder misiniz?", "single_choice", `["Kesinlikle Evet","Evet","Kararsızım","Hayır","Kesinlikle Hayır"]`, ""},
			{"E37", "37. Üniversitenizden mezun olduğunuz için gurur duyuyor musunuz?", "single_choice", `["Evet","Hayır"]`, ""},
			{"E38", "38. Eğitim sürecinde eksik olduğunu düşündüğünüz konular nelerdir?", "single_choice", `["Pratik ve uygulama eksikliği", "Güncel teknoloji/yazılım eğitimleri", "Yabancı dil eğitimi", "Sektörel staj olanakları", "Akademik/Teorik derinlik", "Eksik bir konu yoktu"]`, ""},
			{"E39", "39. Müfredatta hangi konuların/derslerin eklenmesini veya güçlendirilmesini önerirsiniz?", "single_choice", `["Sektörel uygulamalar ve saha çalışmaları", "Yapay zeka ve dijitalleşme", "Girişimcilik ve proje yönetimi", "Kariyer planlama ve iletişim", "Müfredat yeterliydi"]`, ""},
			{"E40", "40. Üniversitenize ve bölümünüze iletmek istediğiniz diğer görüş ve önerileriniz nelerdir?", "single_choice", `["Mezunlarla iletişim güçlendirilmeli", "Kariyer merkezi daha aktif olmalı", "Fiziki şartlar (kampüs, lab vb.) iyileştirilmeli", "Akademik kadro desteklenmeli", "Herhangi bir ek önerim yok"]`, ""},
		}},
		{6, "F. Alanınıza Özgü Sorular", []questionSeed{
			// 1. Denizcilik Fakültesi
			{"F41", "41. Mezuniyet sonrası gemi adamlığı cüzdanı/yeterlilik belgesi (STCW) aldınız mı?", "single_choice", `["Evet","Hayır"]`, "Denizcilik Fakültesi"},
			{"F41", "42. Şu anda denizde mi yoksa karada mı (kıyı ofisi/liman/lojistik) çalışıyorsunuz?", "single_choice", `["Denizde (Gemide)", "Karada (Kıyı ofisi, liman vb.)", "Farklı bir sektörde", "Çalışmıyorum"]`, "Denizcilik Fakültesi"},
			{"F41", "43. Gemide görev yapıyorsanız hangi unvanla (kaptan, çarkçı, zabit vb.) çalışıyorsunuz?", "single_choice", `["Kaptan / 1. Zabit", "Çarkçıbaşı / Mühendis", "Vardiya Zabiti", "Stajyer / Diğer", "Gemide çalışmıyorum"]`, "Denizcilik Fakültesi"},
			{"F41", "44. Uluslararası gemi seferlerine katıldınız mı?", "single_choice", `["Evet","Hayır"]`, "Denizcilik Fakültesi"},
			{"F41", "45. Aldığınız denizcilik eğitiminin uluslararası STCW standartlarına uygunluğunu nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Denizcilik Fakültesi"},
			{"F41", "46. Liman işletmeciliği, lojistik veya denizcilik sigortacılığı gibi kara kariyerlerine yöneldiniz mi?", "single_choice", `["Evet","Hayır"]`, "Denizcilik Fakültesi"},
			{"F41", "47. Simülatör ve gemi uygulama eğitimlerinin yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Denizcilik Fakültesi"},
			{"F41", "48. Denizcilik sektöründeki iş bulma sürecinde üniversite-sektör işbirliğinin (staj, gemi stajı vb.) katkısını nasıl değerlendirirsiniz?", "single_choice", `["Çok Fazla","Biraz","Hiç"]`, "Denizcilik Fakültesi"},
			{"F41", "49. Yabancı dil (özellikle denizcilik İngilizcesi) yeterliliğinizin mesleki hayatınıza katkısı oldu mu?", "single_choice", `["Evet","Hayır"]`, "Denizcilik Fakültesi"},
			{"F41", "50. Denizcilik sektöründeki kariyer olanakları hakkında ne düşünüyorsunuz?", "single_choice", `["Çok avantajlı ve çeşitli", "Maaşlar iyi ama yıpratıcı", "Karada iş bulmak zor", "Sektör giderek daralıyor", "Fikrim yok"]`, "Denizcilik Fakültesi"},

			// 2. Diş Hekimliği Fakültesi
			{"F41", "41. Mezuniyet sonrası özel klinik mi, kamu kurumu mu yoksa kendi muayenehanenizi mi tercih ettiniz?", "single_choice", `["Kamu Kurumu (Devlet Hastanesi vb.)", "Özel Klinik/Hastane", "Kendi Muayenehanem", "Akademi / Üniversite", "Çalışmıyorum"]`, "Diş Hekimliği Fakültesi"},
			{"F41", "42. Diş hekimliğinde uzmanlık eğitimine (DUS) başladınız mı / tamamladınız mı?", "single_choice", `["Evet","Hayır"]`, "Diş Hekimliği Fakültesi"},
			{"F41", "43. Klinik becerilerinizi geliştirmede fakültedeki uygulamalı eğitimin yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Diş Hekimliği Fakültesi"},
			{"F41", "44. Hangi branşta (ortodonti, cerrahi, protez vb.) çalışıyorsunuz veya uzmanlaşmayı düşünüyorsunuz?", "single_choice", `["Ağız, Diş ve Çene Cerrahisi", "Ortodonti", "Protetik Diş Tedavisi", "Endodonti / Restoratif", "Uzmanlaşmayı düşünmüyorum / Pratisyen"]`, "Diş Hekimliği Fakültesi"},
			{"F41", "45. Mesleki uygulamalarınızda en çok karşılaştığınız zorluklar nelerdir?", "single_choice", `["Yoğun hasta çalışma saatleri", "Malzeme ve ekipman maliyetleri", "Hasta iletişimi ve beklentileri", "Fiziksel yıpranma (ergonomi)", "Zorluk yaşamıyorum"]`, "Diş Hekimliği Fakültesi"},
			{"F41", "46. Diş hekimliği fakültesindeki klinik/hasta uygulama sayısını yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Diş Hekimliği Fakültesi"},
			{"F41", "47. Mezuniyet sonrası bir işletme (klinik) açma sürecinde girişimcilik bilgisine ihtiyaç duydunuz mu?", "single_choice", `["Evet","Hayır"]`, "Diş Hekimliği Fakültesi"},
			{"F41", "48. Diş hekimliğinde teknolojik gelişmeleri (dijital diş hekimliği, CAD/CAM vb.) takip etme konusunda kendinizi yeterli görüyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Diş Hekimliği Fakültesi"},
			{"F41", "49. Diş hekimliği alanında yurt dışında çalışma/eğitim deneyiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "Diş Hekimliği Fakültesi"},
			{"F41", "50. Fakültenizin mezunlar arası mesleki dayanışma ve bilgi paylaşımını yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Diş Hekimliği Fakültesi"},

			// 3. Eczacılık Fakültesi
			{"F41", "41. Şu anda serbest eczane, hastane eczanesi, ilaç sanayi veya akademi alanlarından hangisinde çalışıyorsunuz?", "single_choice", `["Serbest Eczane", "Hastane Eczanesi", "İlaç Sanayi / Endüstri", "Akademi / Araştırma", "Farklı Sektör / Çalışmıyorum"]`, "Eczacılık Fakültesi"},
			{"F41", "42. Kendi eczanenizi açtıysanız, bu süreçte karşılaştığınız zorluklar nelerdir?", "single_choice", `["Sermaye ve finansman eksikliği", "Mevzuat ve bürokrasi", "Uygun lokasyon bulamamak", "Rekabet", "Kendi eczanemi açmadım"]`, "Eczacılık Fakültesi"},
			{"F41", "43. İlaç sektöründe (Ar-Ge, üretim, pazarlama) çalışıyorsanız, eğitiminizin bu alana katkısını nasıl değerlendirirsiniz?", "single_choice", `["Çok yeterliydi", "Kısmen yeterliydi", "Pratik/Staj eksikti", "Tamamen yetersizdi", "İlaç sektöründe çalışmıyorum"]`, "Eczacılık Fakültesi"},
			{"F41", "44. Klinik eczacılık uygulamalarına yönelik eğitiminizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Eczacılık Fakültesi"},
			{"F41", "45. Mesleki uzmanlık alanı (klinik eczacılık, toksikoloji vb.) için lisansüstü eğitim aldınız mı / almayı düşünüyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Eczacılık Fakültesi"},
			{"F41", "46. Reçete değerlendirme ve ilaç etkileşimi konusundaki eğitiminizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Eczacılık Fakültesi"},
			{"F41", "47. Farmasötik teknoloji ve laboratuvar uygulamalarına yönelik altyapıyı yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Eczacılık Fakültesi"},
			{"F41", "48. Mesleki mevzuat (eczacılık kanunu, ruhsatlandırma vb.) konusundaki bilginizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Eczacılık Fakültesi"},
			{"F41", "49. Eczacılık mesleğinde dijitalleşme ve e-reçete sistemlerine uyum sürecinizi nasıl değerlendirirsiniz?", "single_choice", `["Çok hızlı ve kolay oldu", "Biraz zorlandım ama alıştım", "Sistemler çok sorunlu", "Bu alanda çalışmıyorum"]`, "Eczacılık Fakültesi"},
			{"F41", "50. Eczacılık sektöründeki istihdam olanaklarını nasıl değerlendiriyorsunuz?", "single_choice", `["İş bulmak çok kolay", "Eskisi kadar kolay değil (Kontenjan artışı vb.)", "Sadece belli alanlarda açık var", "İstihdam çok zor"]`, "Eczacılık Fakültesi"},

			// 4. Eğitim Fakültesi
			{"F41", "41. Mezuniyet sonrası MEB'e atandınız mı / öğretmenlik yapıyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Eğitim Fakültesi"},
			{"F41", "42. Özel sektörde (özel okul, dershane, kurs merkezi) mi yoksa kamuda mı görev yapıyorsunuz?", "single_choice", `["Kamu (MEB)", "Özel Okul / Kolej", "Kurs Merkezi / Dershane", "Özel Ders / Freelance", "Çalışmıyorum / Farklı sektör"]`, "Eğitim Fakültesi"},
			{"F41", "43. KPSS/öğretmenlik atama sürecinde aldığınız eğitimin yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Eğitim Fakültesi"},
			{"F41", "44. Öğretmenlik uygulaması (staj) derslerinin mesleğe hazırlık açısından yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Eğitim Fakültesi"},
			{"F41", "45. Sınıf yönetimi becerileriniz konusunda kendinizi ne kadar yeterli hissediyorsunuz?", "single_choice", `["Yeterli","Kısmen","Yetersiz"]`, "Eğitim Fakültesi"},
			{"F41", "46. Alan bilginizin (branşınıza özgü) güncel müfredata uygunluğunu nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Eğitim Fakültesi"},
			{"F41", "47. Eğitim teknolojilerini (akıllı tahta, dijital içerik vb.) derslerinizde kullanma konusunda kendinizi yeterli görüyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Eğitim Fakültesi"},
			{"F41", "48. Özel eğitim/kaynaştırma öğrencilerine yönelik bilginizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Eğitim Fakültesi"},
			{"F41", "49. Lisansüstü eğitime (yüksek lisans/doktora - eğitim bilimleri) yöneldiniz mi?", "single_choice", `["Evet","Hayır"]`, "Eğitim Fakültesi"},
			{"F41", "50. Öğretmenlik mesleğinde kariyer memnuniyetinizi nasıl değerlendirirsiniz?", "single_choice", `["Çok Yüksek", "Yüksek", "Orta", "Düşük", "Çok Düşük"]`, "Eğitim Fakültesi"},

			// 5. Fen Fakültesi
			{"F41", "41. Mezuniyet sonrası akademik kariyere (araştırma görevliliği, doktora vb.) yöneldiniz mi?", "single_choice", `["Evet","Hayır"]`, "Fen Fakültesi"},
			{"F41", "42. Özel sektörde Ar-Ge, laboratuvar veya analiz uzmanı olarak mı çalışıyorsunuz?", "single_choice", `["Evet","Hayır"]`, "Fen Fakültesi"},
			{"F41", "43. Öğretmenlik formasyonu aldıysanız, öğretmenlik mesleğine geçiş yaptınız mı?", "single_choice", `["Evet","Hayır","Almadım"]`, "Fen Fakültesi"},
			{"F41", "44. Laboratuvar ve uygulamalı ders altyapısının yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Fen Fakültesi"},
			{"F41", "45. Aldığınız temel bilim eğitiminin (matematik, fizik, kimya, biyoloji vb.) iş hayatınıza katkısını nasıl değerlendirirsiniz?", "single_choice", `["Çok Fazla","Kısmen","Hiç"]`, "Fen Fakültesi"},
			{"F41", "46. Bilimsel araştırma yöntemleri konusundaki eğitiminizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Fen Fakültesi"},
			{"F41", "47. Veri analizi/istatistik yazılımları (R, Python, SPSS vb.) konusunda kendinizi yeterli görüyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Fen Fakültesi"},
			{"F41", "48. Kamu kurumlarında (TÜBİTAK, kalite kontrol laboratuvarları vb.) istihdam deneyiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "Fen Fakültesi"},
			{"F41", "49. Alanınızla ilgili bilimsel yayın/proje üretme deneyiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "Fen Fakültesi"},
			{"F41", "50. Fen Fakültesi mezunu olarak iş bulma sürecinde en çok zorlandığınız konu neydi?", "single_choice", `["Kadro/İstihdam yetersizliği", "Özel sektörde unvan karşılığının olmaması", "Tecrübe istenmesi", "Mülakat süreçleri", "Zorlanmadım"]`, "Fen Fakültesi"},

			// 6. Güzel Sanatlar Fakültesi
			{"F41", "41. Mezuniyet sonrası sanatsal üretim/serbest sanatçılık ile mi yoksa kurumsal bir işte mi çalışıyorsunuz?", "single_choice", `["Serbest Sanatçı / Kendi Atölyem", "Kurumsal Firma (Tasarımcı, Sanat Yönetmeni vb.)", "Akademi / Eğitim Kurumu", "Sanat Dışı Bir Sektör", "Çalışmıyorum"]`, "Güzel Sanatlar Fakültesi"},
			{"F41", "42. Kendi sergi/proje/atölye çalışmalarınızı sürdürüyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Güzel Sanatlar Fakültesi"},
			{"F41", "43. Eğitim aldığınız disiplinle (resim, heykel, seramik, grafik tasarım vb.) ilgili sektörde istihdam edildiniz mi?", "single_choice", `["Evet","Hayır"]`, "Güzel Sanatlar Fakültesi"},
			{"F41", "44. Atölye ve uygulamalı derslerin mesleki hazırlığa katkısını nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Güzel Sanatlar Fakültesi"},
			{"F41", "45. Dijital tasarım araçları (Photoshop, Illustrator vb.) konusundaki eğitiminizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Güzel Sanatlar Fakültesi"},
			{"F41", "46. Sanat eserlerinizi/ürünlerinizi pazarlama ve sergileme konusunda kendinizi yeterli hissediyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Güzel Sanatlar Fakültesi"},
			{"F41", "47. Freelance/proje bazlı çalışma deneyiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "Güzel Sanatlar Fakültesi"},
			{"F41", "48. Öğretmenlik formasyonu aldıysanız, görsel sanatlar öğretmenliğine yöneldiniz mi?", "single_choice", `["Evet","Hayır","Almadım"]`, "Güzel Sanatlar Fakültesi"},
			{"F41", "49. Sanatsal girişimcilik (kendi atölyenizi/stüdyonuzu açma) konusunda destek/bilgi ihtiyacı duydunuz mu?", "single_choice", `["Evet","Hayır"]`, "Güzel Sanatlar Fakültesi"},
			{"F41", "50. Güzel Sanatlar mezunu olarak istihdam sürecinde en çok karşılaştığınız zorluk neydi?", "single_choice", `["Sanata verilen değerin düşüklüğü", "Düzenli gelir / Maaş güvencesizliği", "Sektörde tanıdık/çevre (network) gerekliliği", "Malzeme/Atölye maliyetleri", "Zorluk yaşamadım"]`, "Güzel Sanatlar Fakültesi"},

			// 7. Hemşirelik Fakültesi
			{"F41", "41. Mezuniyet sonrası kamu hastanesi, özel hastane veya başka bir kurumda mı çalışıyorsunuz?", "single_choice", `["Kamu Hastanesi", "Özel Hastane / Klinik", "Aile Sağlığı Merkezi (ASM)", "İş Yeri Hemşireliği / Diğer", "Çalışmıyorum"]`, "Hemşirelik Fakültesi"},
			{"F41", "42. Hangi klinik alanda (yoğun bakım, ameliyathane, pediatri vb.) görev yapıyorsunuz?", "single_choice", `["Servis Hemşireliği (Dahiliye, Cerrahi vb.)", "Yoğun Bakım", "Ameliyathane", "Pediatri / Kadın Doğum", "Acil / Diğer"]`, "Hemşirelik Fakültesi"},
			{"F41", "43. Klinik uygulama derslerinin mesleğe hazırlık açısından yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Hemşirelik Fakültesi"},
			{"F41", "44. Hemşirelikte uzmanlık alanına yönelik lisansüstü eğitim aldınız mı / almayı düşünüyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Hemşirelik Fakültesi"},
			{"F41", "45. Acil durum ve kriz yönetimi becerileriniz konusunda kendinizi ne kadar yeterli hissediyorsunuz?", "single_choice", `["Yeterli","Kısmen","Yetersiz"]`, "Hemşirelik Fakültesi"},
			{"F41", "46. Hasta iletişimi ve empati becerileri konusundaki eğitiminizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Hemşirelik Fakültesi"},
			{"F41", "47. Vardiyalı çalışma düzeninin mesleki/özel hayatınıza etkisini nasıl değerlendirirsiniz?", "single_choice", `["Hiç olumsuz etkisi yok", "Biraz zorluyor ama yönetebiliyorum", "Özel hayatımı ve sağlığımı çok olumsuz etkiliyor", "Sürekli gündüz çalışıyorum (Vardiyalı değilim)"]`, "Hemşirelik Fakültesi"},
			{"F41", "48. Yurt dışında hemşirelik yapma (yeterlilik/denklik vb.) girişiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "Hemşirelik Fakültesi"},
			{"F41", "49. Mesleki tükenmişlik (burnout) yaşadınız mı, bu konuda destek aldınız mı?", "single_choice", `["Evet","Hayır"]`, "Hemşirelik Fakültesi"},
			{"F41", "50. Hemşirelik mesleğinde kariyer memnuniyetinizi nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Hemşirelik Fakültesi"},

			// 8. İktisadi ve İdari Bilimler Fakültesi
			{"F41", "41. Mezuniyet sonrası kamu (memur/uzman) mı yoksa özel sektörde mi çalışıyorsunuz?", "single_choice", `["Kamu (Memur, Uzman, Müfettiş vb.)", "Özel Sektör (Kurumsal Firma)", "Özel Sektör (KOBİ / Yerel İşletme)", "Kendi İşim", "Çalışmıyorum"]`, "İktisadi ve İdari Bilimler Fakültesi"},
			{"F41", "42. Çalıştığınız departman/pozisyon (finans, muhasebe, insan kaynakları, pazarlama vb.) nedir?", "single_choice", `["Muhasebe / Finans / Denetim", "İnsan Kaynakları / İdari İşler", "Pazarlama / Satış / Dış Ticaret", "Yönetim / Planlama / Operasyon", "Diğer / Kendi İşim"]`, "İktisadi ve İdari Bilimler Fakültesi"},
			{"F41", "43. KPSS/kamu personeli sınavlarına yönelik eğitiminizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "İktisadi ve İdari Bilimler Fakültesi"},
			{"F41", "44. Muhasebe, finans veya ekonomi derslerinin sektördeki uygulamalara uygunluğunu nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "İktisadi ve İdari Bilimler Fakültesi"},
			{"F41", "45. Excel, ERP veya finansal analiz yazılımları konusunda kendinizi yeterli görüyor musunuz?", "single_choice", `["Evet","Hayır"]`, "İktisadi ve İdari Bilimler Fakültesi"},
			{"F41", "46. Staj/işletmelerle işbirliği (kariyer günleri, örnek olay çalışmaları vb.) deneyiminizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "İktisadi ve İdari Bilimler Fakültesi"},
			{"F41", "47. Kendi işinizi kurma/girişimcilik konusunda eğitiminizin katkısını nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "İktisadi ve İdari Bilimler Fakültesi"},
			{"F41", "48. Uluslararası ticaret/ihracat alanında çalışma deneyiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "İktisadi ve İdari Bilimler Fakültesi"},
			{"F41", "49. Lisansüstü eğitime (MBA, işletme/iktisat yüksek lisansı vb.) yöneldiniz mi?", "single_choice", `["Evet","Hayır"]`, "İktisadi ve İdari Bilimler Fakültesi"},
			{"F41", "50. İİBF mezunu olarak istihdam sürecinde en çok hangi becerilere ihtiyaç duyduğunuzu fark ettiniz?", "single_choice", `["Yabancı Dil (İngilizce)", "Veri Analizi ve Excel / ERP Programları", "Mülakat ve İletişim Becerileri", "İş Deneyimi (Staj vb.)", "Zorlanmadım"]`, "İktisadi ve İdari Bilimler Fakültesi"},

			// 9. İlahiyat Fakültesi
			{"F41", "41. Mezuniyet sonrası Diyanet İşleri Başkanlığı'na mı yoksa MEB'e (Din Kültürü öğretmenliği) mi atandınız?", "single_choice", `["Diyanet İşleri Başkanlığı", "Milli Eğitim Bakanlığı (MEB)", "Akademi / Üniversite", "Özel Sektör", "Atanamadım / Çalışmıyorum"]`, "İlahiyat Fakültesi"},
			{"F41", "42. Din görevlisi/vaiz/imam-hatip olarak mı yoksa öğretmen olarak mı görev yapıyorsunuz?", "single_choice", `["İmam-Hatip / Müezzin-Kayyım", "Kur'an Kursu Öğreticisi / Vaiz", "Din Kültürü ve Ahlak Bilgisi Öğretmeni", "İdari Görev / Uzman", "Çalışmıyorum / Farklı Sektör"]`, "İlahiyat Fakültesi"},
			{"F41", "43. Arapça ve dini ilimler eğitiminizin mesleki hayatınıza katkısını nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "İlahiyat Fakültesi"},
			{"F41", "44. Hitabet ve topluluk önünde konuşma becerileriniz konusunda kendinizi ne kadar yeterli hissediyorsunuz?", "single_choice", `["Yeterli","Kısmen","Yetersiz"]`, "İlahiyat Fakültesi"},
			{"F41", "45. Din eğitimi ve öğretimi derslerinin (formasyon) yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "İlahiyat Fakültesi"},
			{"F41", "46. Dini danışmanlık/rehberlik hizmetleri konusunda kendinizi yeterli görüyor musunuz?", "single_choice", `["Evet","Hayır"]`, "İlahiyat Fakültesi"},
			{"F41", "47. Toplumsal ve kültürel çeşitlilik konusundaki eğitiminizi mesleğinizde yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "İlahiyat Fakültesi"},
			{"F41", "48. Akademik kariyere (yüksek lisans/doktora - ilahiyat alanında) yöneldiniz mi?", "single_choice", `["Evet","Hayır"]`, "İlahiyat Fakültesi"},
			{"F41", "49. Yurt dışında (din görevlisi olarak) görev alma deneyiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "İlahiyat Fakültesi"},
			{"F41", "50. İlahiyat mezunu olarak istihdam sürecinde en çok karşılaştığınız zorluk neydi?", "single_choice", `["KPSS/DHBT puanlarının yüksekliği", "Mülakat süreçleri", "Kadro sayısının yetersizliği", "Atanılan bölgenin zorlukları", "Zorlanmadım"]`, "İlahiyat Fakültesi"},

			// 10. İletişim Fakültesi
			{"F41", "41. Mezuniyet sonrası medya (TV, gazete, dijital medya), halkla ilişkiler/reklam veya kurumsal iletişim alanlarından hangisinde çalışıyorsunuz?", "single_choice", `["Geleneksel Medya (TV, Radyo, Gazete)", "Dijital Medya / Sosyal Medya", "Halkla İlişkiler / Kurumsal İletişim", "Reklam / Prodüksiyon Ajansı", "İletişim Dışı Sektör / Çalışmıyorum"]`, "İletişim Fakültesi"},
			{"F41", "42. Gazetecilik, halkla ilişkiler veya radyo-tv-sinema gibi hangi bölümden mezun oldunuz ve alanınızla ilgili işe yerleştiniz mi?", "single_choice", `["RTS mezunuyum, alanımda çalışıyorum", "Gazetecilik mezunuyum, alanımda çalışıyorum", "HİT mezunuyum, alanımda çalışıyorum", "Farklı bir sektörde çalışıyorum", "Çalışmıyorum"]`, "İletişim Fakültesi"},
			{"F41", "43. Dijital/sosyal medya yönetimi becerileriniz konusunda kendinizi ne kadar yeterli hissediyorsunuz?", "single_choice", `["Yeterli","Kısmen","Yetersiz"]`, "İletişim Fakültesi"},
			{"F41", "44. Kamera, kurgu, prodüksiyon gibi teknik/uygulamalı derslerin yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "İletişim Fakültesi"},
			{"F41", "45. İçerik üretimi (yazılı, görsel, video) konusunda kendinizi yeterli görüyor musunuz?", "single_choice", `["Evet","Hayır"]`, "İletişim Fakültesi"},
			{"F41", "46. Freelance/proje bazlı iletişim çalışmaları (içerik üreticiliği, influencer vb.) yaptınız mı?", "single_choice", `["Evet","Hayır"]`, "İletişim Fakültesi"},
			{"F41", "47. Kriz iletişimi ve kurumsal itibar yönetimi konusundaki eğitiminizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "İletişim Fakültesi"},
			{"F41", "48. Medya sektöründeki hızlı teknolojik değişime uyum sağlamada kendinizi ne kadar yeterli hissediyorsunuz?", "single_choice", `["Yeterli","Kısmen","Yetersiz"]`, "İletişim Fakültesi"},
			{"F41", "49. Stajlarınızın medya/iletişim sektöründe iş bulmanıza katkısı oldu mu?", "single_choice", `["Evet","Hayır"]`, "İletişim Fakültesi"},
			{"F41", "50. İletişim Fakültesi mezunu olarak istihdam sürecinde en çok hangi becerilere ihtiyaç duyduğunuzu fark ettiniz?", "single_choice", `["Dijital Pazarlama / SEO", "Kurgu / Grafik Tasarım Programları", "İçerik Üretimi / Metin Yazarlığı", "Yabancı Dil ve Çevre (Network)", "Zorlanmadım"]`, "İletişim Fakültesi"},

			// 11. İnsan ve Toplum Bilimleri Fakültesi
			{"F41", "41. Mezun olduğunuz bölüm (psikoloji, sosyoloji, tarih, felsefe vb.) alanıyla ilgili bir işte mi çalışıyorsunuz?", "single_choice", `["Evet","Hayır"]`, "İnsan ve Toplum Bilimleri Fakültesi"},
			{"F41", "42. Kamu kurumu, STK veya özel sektörden hangisinde istihdam edildiniz?", "single_choice", `["Kamu Kurumu", "Özel Sektör", "Sivil Toplum Kuruluşu (STK)", "Akademi / Araştırma Merkezi", "Çalışmıyorum"]`, "İnsan ve Toplum Bilimleri Fakültesi"},
			{"F41", "43. Alanınızla ilgili lisansüstü eğitime (yüksek lisans/doktora) yöneldiniz mi?", "single_choice", `["Evet","Hayır"]`, "İnsan ve Toplum Bilimleri Fakültesi"},
			{"F41", "44. Araştırma yöntemleri ve bilimsel yazım becerileriniz konusunda kendinizi ne kadar yeterli hissediyorsunuz?", "single_choice", `["Yeterli","Kısmen","Yetersiz"]`, "İnsan ve Toplum Bilimleri Fakültesi"},
			{"F41", "45. Sosyal bilimler eğitiminizin analitik düşünme becerinize katkısını nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "İnsan ve Toplum Bilimleri Fakültesi"},
			{"F41", "46. Psikolojik danışmanlık/sosyal hizmet gibi uygulamalı alanlarda staj deneyiminiz yeterli oldu mu?", "single_choice", `["Evet","Hayır"]`, "İnsan ve Toplum Bilimleri Fakültesi"},
			{"F41", "47. Öğretmenlik formasyonu aldıysanız, ilgili branşta öğretmenliğe yöneldiniz mi?", "single_choice", `["Evet","Hayır","Almadım"]`, "İnsan ve Toplum Bilimleri Fakültesi"},
			{"F41", "48. Alanınızla doğrudan ilgili olmayan sektörlerde (insan kaynakları, eğitim danışmanlığı vb.) istihdam edildiyseniz, bu geçişte zorlandınız mı?", "single_choice", `["Evet","Hayır","İlgili Sektördeyim"]`, "İnsan ve Toplum Bilimleri Fakültesi"},
			{"F41", "49. Toplumsal konularla ilgili proje/STK çalışmalarına katıldınız mı?", "single_choice", `["Evet","Hayır"]`, "İnsan ve Toplum Bilimleri Fakültesi"},
			{"F41", "50. İnsan ve Toplum Bilimleri mezunu olarak istihdam sürecinde en çok karşılaştığınız zorluk neydi?", "single_choice", `["Bölüme özgü net bir meslek tanımı olmaması", "Özel sektörde talep azlığı", "Kamu alımlarının (KPSS) yetersizliği", "Düşük başlangıç maaşları", "Zorlanmadım"]`, "İnsan ve Toplum Bilimleri Fakültesi"},

			// 12. Mimarlık Fakültesi
			{"F41", "41. Mezuniyet sonrası serbest mimarlık bürosu, kamu kurumu veya inşaat şirketinde mi çalışıyorsunuz?", "single_choice", `["Özel Mimarlık/Tasarım Bürosu", "İnşaat / Şantiye / Proje Şirketi", "Kamu Kurumu (Bakanlık, Belediye vb.)", "Kendi Mimarlık Ofisim", "Çalışmıyorum / Farklı Sektör"]`, "Mimarlık Fakültesi"},
			{"F41", "42. Mimar olarak SMM (Serbest Mimarlık Müşavirlik) yetkisi/oda kaydınızı aldınız mı?", "single_choice", `["Evet","Hayır"]`, "Mimarlık Fakültesi"},
			{"F41", "43. Proje çizim yazılımları (AutoCAD, Revit, SketchUp vb.) konusundaki eğitiminizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Mimarlık Fakültesi"},
			{"F41", "44. Stüdyo/atölye derslerinin mesleki uygulamaya hazırlık açısından yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Mimarlık Fakültesi"},
			{"F41", "45. Şantiye/uygulama deneyiminiz mesleğinize ne ölçüde katkı sağladı?", "single_choice", `["Çok Fazla","Biraz","Hiç"]`, "Mimarlık Fakültesi"},
			{"F41", "46. Kentsel tasarım/restorasyon gibi alanlarda uzmanlaşmayı düşünüyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Mimarlık Fakültesi"},
			{"F41", "47. Sürdürülebilir/yeşil bina tasarımı konusundaki bilgi düzeyinizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Mimarlık Fakültesi"},
			{"F41", "48. Kendi büronuzu açma (girişimcilik) sürecinde zorlandığınız konular oldu mu?", "single_choice", `["Müşteri / Çevre bulamamak", "Sermaye ve vergi yükleri", "Piyasadaki yoğun rekabet", "Kendi büromu açmadım / Zorlanmadım"]`, "Mimarlık Fakültesi"},
			{"F41", "49. Yurt dışında mimarlık projelerinde çalışma deneyiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "Mimarlık Fakültesi"},
			{"F41", "50. Mimarlık mesleğindeki istihdam olanaklarını ve rekabeti nasıl değerlendiriyorsunuz?", "single_choice", `["İstihdam çok zor, rekabet aşırı yüksek", "Çevreye bağlı olarak iş bulunabiliyor", "Nitelikli personel için olanaklar iyi", "Şartlar ve maaşlar çok düşük"]`, "Mimarlık Fakültesi"},

			// 13. Mühendislik Fakültesi
			{"F41", "41. Mezun olduğunuz mühendislik bölümüyle (inşaat, elektrik-elektronik, bilgisayar, makine vb.) doğrudan ilgili bir işte mi çalışıyorsunuz?", "single_choice", `["Evet","Hayır"]`, "Mühendislik Fakültesi"},
			{"F41", "42. Kamu kurumu, özel sektör veya kendi işiniz olmak üzere hangi alanda istihdam edildiniz?", "single_choice", `["Özel Sektör (Ulusal Firma)", "Özel Sektör (Uluslararası/Çok Uluslu Firma)", "Kamu Kurumu", "Kendi Şirketim", "Çalışmıyorum"]`, "Mühendislik Fakültesi"},
			{"F41", "43. Mühendislik yazılımları/programlama becerileriniz konusunda kendinizi ne kadar yeterli hissediyorsunuz?", "single_choice", `["Yeterli","Kısmen","Yetersiz"]`, "Mühendislik Fakültesi"},
			{"F41", "44. Laboratuvar ve tasarım projesi derslerinin mesleki uygulamaya hazırlık açısından yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Mühendislik Fakültesi"},
			{"F41", "45. Staj deneyiminizin iş bulma sürecinize katkısını nasıl değerlendirirsiniz?", "single_choice", `["Çok Fazla","Kısmen","Hiç"]`, "Mühendislik Fakültesi"},
			{"F41", "46. Meslek odası (İnşaat Mühendisleri Odası, Makine Mühendisleri Odası vb.) üyeliği/yetki belgesi aldınız mı?", "single_choice", `["Evet","Hayır"]`, "Mühendislik Fakültesi"},
			{"F41", "47. Sektördeki teknolojik gelişmelere (yapay zeka, otomasyon, yeni yazılımlar vb.) uyum sağlama konusunda kendinizi yeterli görüyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Mühendislik Fakültesi"},
			{"F41", "48. Proje yönetimi becerileriniz konusundaki eğitiminizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Mühendislik Fakültesi"},
			{"F41", "49. Yurt dışında mühendislik projelerinde çalışma deneyiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "Mühendislik Fakültesi"},
			{"F41", "50. Mühendislik mesleğindeki istihdam olanaklarını ve ücret memnuniyetinizi nasıl değerlendiriyorsunuz?", "single_choice", `["İstihdam kolay, ücretler tatmin edici", "İstihdam var ancak başlangıç ücretleri düşük", "Kariyer ilerledikçe şartlar iyileşiyor", "Sektörde ciddi işsizlik var"]`, "Mühendislik Fakültesi"},

			// 14. Müzik ve Sahne Sanatları Fakültesi
			{"F41", "41. Mezuniyet sonrası profesyonel müzisyen/sanatçı olarak mı yoksa öğretmen/eğitmen olarak mı çalışıyorsunuz?", "single_choice", `["Sadece Sanatçı / İcracı (Sahne, Orkestra)", "Sadece Eğitmen / Öğretmen", "Hem Sahne Hem Eğitmenlik", "Sanat Dışı Bir Sektör", "Çalışmıyorum"]`, "Müzik ve Sahne Sanatları Fakültesi"},
			{"F41", "42. Bir orkestra, tiyatro, opera/bale topluluğu veya kurumda kadrolu/proje bazlı görev alıyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Müzik ve Sahne Sanatları Fakültesi"},
			{"F41", "43. Enstrüman/ses eğitiminin mesleki performansınıza katkısını nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Müzik ve Sahne Sanatları Fakültesi"},
			{"F41", "44. Sahne performansı ve topluluk önünde icra becerileriniz konusunda kendinizi ne kadar yeterli hissediyorsunuz?", "single_choice", `["Yeterli","Kısmen","Yetersiz"]`, "Müzik ve Sahne Sanatları Fakültesi"},
			{"F41", "45. Müzik/sahne sanatları öğretmenliği formasyonu aldıysanız, bu alanda istihdam edildiniz mi?", "single_choice", `["Evet","Hayır","Almadım"]`, "Müzik ve Sahne Sanatları Fakültesi"},
			{"F41", "46. Kendi proje/konser/gösterinizi organize etme (sanatsal girişimcilik) deneyiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "Müzik ve Sahne Sanatları Fakültesi"},
			{"F41", "47. Serbest/freelance sanatçılık sürecinde gelir istikrarı konusunda zorluk yaşadınız mı?", "single_choice", `["Evet","Hayır"]`, "Müzik ve Sahne Sanatları Fakültesi"},
			{"F41", "48. Ses/enstrüman kayıt teknolojileri konusundaki bilgi düzeyinizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Müzik ve Sahne Sanatları Fakültesi"},
			{"F41", "49. Yurt dışında sahne alma/eğitim deneyiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "Müzik ve Sahne Sanatları Fakültesi"},
			{"F41", "50. Müzik ve Sahne Sanatları mezunu olarak istihdam sürecinde en çok karşılaştığınız zorluk neydi?", "single_choice", `["Düzenli kadro / Sigorta eksikliği", "Özel ders ve sahne bulma zorluğu", "Sektördeki haksız rekabet", "Enstrüman ve ekipman maliyetleri", "Zorlanmadım"]`, "Müzik ve Sahne Sanatları Fakültesi"},

			// 15. Sağlık Bilimleri Fakültesi
			{"F41", "41. Mezun olduğunuz bölümle (fizyoterapi, beslenme ve diyetetik, sosyal hizmet, odyoloji vb.) ilgili bir işte mi çalışıyorsunuz?", "single_choice", `["Evet","Hayır"]`, "Sağlık Bilimleri Fakültesi"},
			{"F41", "42. Kamu hastanesi, özel klinik/merkez veya kendi işletmenizde mi görev yapıyorsunuz?", "single_choice", `["Kamu Hastanesi / Kurumu", "Özel Hastane", "Özel Klinik / Rehabilitasyon Merkezi", "Kendi İşletmem / Kliniğim", "Çalışmıyorum"]`, "Sağlık Bilimleri Fakültesi"},
			{"F41", "43. Klinik uygulama/staj derslerinin mesleğe hazırlık açısından yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Sağlık Bilimleri Fakültesi"},
			{"F41", "44. Hasta/danışan ile birebir çalışma becerileriniz konusunda kendinizi ne kadar yeterli hissediyorsunuz?", "single_choice", `["Yeterli","Kısmen","Yetersiz"]`, "Sağlık Bilimleri Fakültesi"},
			{"F41", "45. Alanınızla ilgili lisansüstü eğitime (yüksek lisans/doktora) yöneldiniz mi?", "single_choice", `["Evet","Hayır"]`, "Sağlık Bilimleri Fakültesi"},
			{"F41", "46. Multidisipliner sağlık ekibi içinde çalışma deneyiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "Sağlık Bilimleri Fakültesi"},
			{"F41", "47. Kendi kliniğinizi/merkezinizi açma sürecinde (girişimcilik) zorlandığınız konular oldu mu?", "single_choice", `["Sermaye yetersizliği", "Resmi prosedürler ve mevzuat", "Müşteri/Hasta portföyü oluşturmak", "Kendi merkezimi açmadım / Zorlanmadım"]`, "Sağlık Bilimleri Fakültesi"},
			{"F41", "48. Mesleki sertifika/ek eğitim (örn. manuel terapi, klinik beslenme sertifikası vb.) aldınız mı?", "single_choice", `["Evet","Hayır"]`, "Sağlık Bilimleri Fakültesi"},
			{"F41", "49. Yurt dışında mesleğinizi icra etme (denklik, yeterlilik vb.) girişiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "Sağlık Bilimleri Fakültesi"},
			{"F41", "50. Sağlık Bilimleri Fakültesi mezunu olarak istihdam sürecinde en çok karşılaştığınız zorluk neydi?", "single_choice", `["Kamu atamalarının (kadroların) azlığı", "Özel sektörde düşük maaş teklifleri", "İş tecrübesi istenmesi", "Mezun sayısının fazlalığı", "Zorlanmadım"]`, "Sağlık Bilimleri Fakültesi"},

			// 16. Spor Bilimleri Fakültesi
			{"F41", "41. Mezuniyet sonrası beden eğitimi öğretmenliği, antrenörlük veya spor yöneticiliği alanlarından hangisinde çalışıyorsunuz?", "single_choice", `["Beden Eğitimi Öğretmenliği", "Antrenörlük / Fitness Eğitmenliği", "Spor Yöneticiliği / Organizasyon", "Akademi / Üniversite", "Çalışmıyorum / Farklı Sektör"]`, "Spor Bilimleri Fakültesi"},
			{"F41", "42. Bir spor kulübü, fitness merkezi, okul veya kamu kurumunda mı görev yapıyorsunuz?", "single_choice", `["Okul / MEB", "Özel Spor Kulübü / Fitness Merkezi", "Kamu Spor Kurumları (GSB vb.)", "Kendi Spor Merkezim", "Diğer"]`, "Spor Bilimleri Fakültesi"},
			{"F41", "43. Antrenörlük/branş eğitiminizin mesleki uygulamaya hazırlık açısından yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Spor Bilimleri Fakültesi"},
			{"F41", "44. Egzersiz fizyolojisi ve performans analizi konusundaki bilgi düzeyinizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Spor Bilimleri Fakültesi"},
			{"F41", "45. Spor yönetimi/organizasyonu becerileriniz konusunda kendinizi ne kadar yeterli hissediyorsunuz?", "single_choice", `["Yeterli","Kısmen","Yetersiz"]`, "Spor Bilimleri Fakültesi"},
			{"F41", "46. Kendi spor merkezinizi/kulübünüzü açma (girişimcilik) sürecinde zorlandığınız konular oldu mu?", "single_choice", `["Maliyet ve kira yükseklikleri", "Belge ve ruhsat işlemleri", "Üye/Müşteri bulmak", "Kendi merkezimi açmadım"]`, "Spor Bilimleri Fakültesi"},
			{"F41", "47. Rekreasyon/engelli sporları gibi özel alanlarda çalışma deneyiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "Spor Bilimleri Fakültesi"},
			{"F41", "48. Ulusal/uluslararası spor organizasyonlarında görev aldınız mı?", "single_choice", `["Evet","Hayır"]`, "Spor Bilimleri Fakültesi"},
			{"F41", "49. Mesleki sertifika (branş antrenörlüğü, hakemlik vb.) aldınız mı?", "single_choice", `["Evet","Hayır"]`, "Spor Bilimleri Fakültesi"},
			{"F41", "50. Spor Bilimleri mezunu olarak istihdam sürecinde en çok karşılaştığınız zorluk neydi?", "single_choice", `["Özel sektörde güvencesiz çalışma", "MEB/Kamu atamalarının zorluğu", "Sektör dışı alaylı antrenörlerin rekabeti", "Düşük ücret politikası", "Zorlanmadım"]`, "Spor Bilimleri Fakültesi"},

			// 17. Su Ürünleri Fakültesi
			{"F41", "41. Mezuniyet sonrası su ürünleri yetiştiriciliği (akuakültür), avcılık-işleme teknolojisi veya kamu kurumunda mı çalışıyorsunuz?", "single_choice", `["Yetiştiricilik (Çiftlik/Kuluçkahane)", "İşleme ve Kalite Kontrol", "Avcılık ve Donanım", "Kamu Kurumu (Bakanlık, Enstitü)", "Çalışmıyorum / Farklı Sektör"]`, "Su Ürünleri Fakültesi"},
			{"F41", "42. Kendi işletmenizi (balık çiftliği, işleme tesisi vb.) kurdunuz mu?", "single_choice", `["Evet","Hayır"]`, "Su Ürünleri Fakültesi"},
			{"F41", "43. Su ürünleri yetiştiriciliği ve hastalıkları konusundaki uygulamalı eğitiminizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Su Ürünleri Fakültesi"},
			{"F41", "44. Deniz/tatlı su ekosistemleri ve çevre yönetimi konusundaki bilgi düzeyinizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Su Ürünleri Fakültesi"},
			{"F41", "45. Laboratuvar ve saha uygulamalarının mesleki hazırlığa katkısını nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Su Ürünleri Fakültesi"},
			{"F41", "46. Su ürünleri işleme ve kalite kontrol standartları konusundaki eğitiminizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Su Ürünleri Fakültesi"},
			{"F41", "47. Balıkçılık teknolojisi ve avcılık yönetimi konusunda kendinizi yeterli görüyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Su Ürünleri Fakültesi"},
			{"F41", "48. Sektördeki sürdürülebilirlik ve çevresel düzenlemelere uyum konusunda bilgi düzeyinizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Su Ürünleri Fakültesi"},
			{"F41", "49. Yurt dışında su ürünleri sektöründe çalışma deneyiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "Su Ürünleri Fakültesi"},
			{"F41", "50. Su Ürünleri Fakültesi mezunu olarak istihdam sürecinde en çok karşılaştığınız zorluk neydi?", "single_choice", `["İş yerlerinin (çiftliklerin) şehre uzaklığı", "Sektörün belli bölgelerde (Ege/Akdeniz) sınırlı olması", "Kamu atamalarının yetersizliği", "Çalışma koşullarının zorluğu", "Zorlanmadım"]`, "Su Ürünleri Fakültesi"},

			// 18. Tıp Fakültesi
			{"F41", "41. Mezuniyet sonrası pratisyen hekim olarak mı çalışıyorsunuz yoksa uzmanlık eğitimine (TUS) devam mı ediyorsunuz?", "single_choice", `["Pratisyen Hekim (Devlet Hizmet Yükümlülüğü)", "Asistan Hekim (TUS ile uzmanlık eğitimi)", "Uzman Hekim", "Akademi / Temel Bilimler", "Farklı Durum (Yurtdışı vb.)"]`, "Tıp Fakültesi"},
			{"F41", "42. Hangi branşta uzmanlık yapıyorsunuz veya yapmayı planlıyorsunuz?", "single_choice", `["Dahili Branşlar", "Cerrahi Branşlar", "Temel Tıp Bilimleri", "Pratisyen kalmayı tercih ediyorum"]`, "Tıp Fakültesi"},
			{"F41", "43. Kamu hastanesi, üniversite hastanesi veya özel sektörde mi görev yapıyorsunuz?", "single_choice", `["Kamu Hastanesi (EAH vb.)", "Üniversite Hastanesi", "Özel Hastane / Klinik", "Aile Sağlığı Merkezi (ASM)", "Diğer"]`, "Tıp Fakültesi"},
			{"F41", "44. Klinik staj ve pratik uygulama derslerinin mesleğe hazırlık açısından yeterliliğini nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Tıp Fakültesi"},
			{"F41", "45. Acil durum yönetimi ve klinik karar verme becerileriniz konusunda kendinizi ne kadar yeterli hissediyorsunuz?", "single_choice", `["Yeterli","Kısmen","Yetersiz"]`, "Tıp Fakültesi"},
			{"F41", "46. Nöbet/yoğun çalışma temposunun mesleki ve özel hayatınıza etkisini nasıl değerlendirirsiniz?", "single_choice", `["Katlanılabilir seviyede, işimin bir parçası", "Özel hayatımı ciddi şekilde kısıtlıyor", "Tükenmişliğe (Burnout) yol açıyor", "Nöbetli çalışmıyorum"]`, "Tıp Fakültesi"},
			{"F41", "47. Akademik kariyere (araştırma görevliliği, akademisyenlik) yöneldiniz mi?", "single_choice", `["Evet","Hayır"]`, "Tıp Fakültesi"},
			{"F41", "48. Yurt dışında çalışma/uzmanlık (denklik süreçleri dahil) girişiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "Tıp Fakültesi"},
			{"F41", "49. Mesleki tükenmişlik (burnout) yaşadınız mı, bu konuda destek aldınız mı?", "single_choice", `["Evet","Hayır"]`, "Tıp Fakültesi"},
			{"F41", "50. Tıp Fakültesi mezunu olarak kariyer memnuniyetinizi nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Tıp Fakültesi"},

			// 19. Turizm Fakültesi
			{"F41", "41. Mezuniyet sonrası otelcilik, seyahat acenteciliği, yiyecek-içecek işletmeciliği veya rehberlik alanlarından hangisinde çalışıyorsunuz?", "single_choice", `["Otel / Konaklama İşletmeciliği", "Seyahat Acentesi / Tur Operatörü", "Turist Rehberliği", "Yiyecek ve İçecek (Gastronomi)", "Çalışmıyorum / Farklı Sektör"]`, "Turizm Fakültesi"},
			{"F41", "42. Ulusal veya uluslararası bir turizm işletmesinde mi görev yapıyorsunuz?", "single_choice", `["Ulusal","Uluslararası"]`, "Turizm Fakültesi"},
			{"F41", "43. Turist rehberliği yapıyorsanız, profesyonel rehberlik belgesi/lisansınız var mı?", "single_choice", `["Evet","Hayır","Rehber Değilim"]`, "Turizm Fakültesi"},
			{"F41", "44. Staj (zorunlu turizm stajı) deneyiminizin mesleğe hazırlık açısından katkısını nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Turizm Fakültesi"},
			{"F41", "45. Yabancı dil yeterliliğinizin turizm sektöründeki kariyerinize katkısını nasıl değerlendirirsiniz?", "single_choice", `["İyi","Orta","Kötü"]`, "Turizm Fakültesi"},
			{"F41", "46. Konaklama/işletme yönetimi becerileriniz konusunda kendinizi ne kadar yeterli hissediyorsunuz?", "single_choice", `["Yeterli","Kısmen","Yetersiz"]`, "Turizm Fakültesi"},
			{"F41", "47. Kendi turizm işletmenizi (acente, butik otel vb.) kurma girişiminiz oldu mu?", "single_choice", `["Evet","Hayır"]`, "Turizm Fakültesi"},
			{"F41", "48. Turizm sektöründeki mevsimsellik/iş güvencesi konusunda ne gibi zorluklar yaşadınız?", "single_choice", `["Sadece yaz sezonunda iş bulabiliyorum", "Sürekli kadroya geçmekte zorlandım", "Kış aylarında gelir düşüklüğü", "Zorluk yaşamadım / Tam zamanlıyım"]`, "Turizm Fakültesi"},
			{"F41", "49. Dijital pazarlama ve online rezervasyon sistemleri konusundaki bilgi düzeyinizi yeterli buluyor musunuz?", "single_choice", `["Evet","Hayır"]`, "Turizm Fakültesi"},
			{"F41", "50. Turizm Fakültesi mezunu olarak istihdam sürecinde en çok karşılaştığınız zorluk neydi?", "single_choice", `["Uzun mesai saatleri ve düşük ücret", "Sektördeki krizlerden hızlı etkilenme", "Mezun olmayanların sektörde çoğunlukta olması", "Yabancı dil eksikliği", "Zorlanmadım"]`, "Turizm Fakültesi"},
		}},
	}

	for _, cat := range categories {
		category := domain.SurveyCategory{Order: cat.Order, Title: cat.Title}
		if err := db.Create(&category).Error; err != nil {
			return err
		}
		for i, q := range cat.Questions {
			uniqueCode := q.Code
			if q.TargetFaculty != "" {
				uniqueCode = fmt.Sprintf("F%d", 41+i)
			}

			question := domain.SurveyQuestion{
				CategoryID:    category.ID,
				Order:         i + 1,
				Code:          uniqueCode,
				Text:          q.Text,
				AnswerType:    q.AnswerType,
				OptionsJSON:   q.Options,
				Required:      true,
				TargetFaculty: nil,
			}
			
			if q.TargetFaculty != "" {
				tf := q.TargetFaculty
				question.TargetFaculty = &tf
			}

			if err := db.Create(&question).Error; err != nil {
				return err
			}
		}
	}

	log.Println("[seed] Tüm kategoriler ve fakültelere özel filtreli sorular başarıyla yüklendi")
	return nil
}