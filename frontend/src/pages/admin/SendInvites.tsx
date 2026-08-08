import { useState } from 'react';
import { Card, Title, Text, Button, Checkbox, Select, Stack, Notification, Group, ThemeIcon, Container } from '@mantine/core';
import { IconMail, IconCheck, IconX } from '@tabler/icons-react';
import { api } from '../../api/client';

export default function SendInvites() {
    const [sendToAll, setSendToAll] = useState(false);
    const [facultyId, setFacultyId] = useState<string | null>(null);
    const [departmentName, setDepartmentName] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);
    const [status, setStatus] = useState<{ type: 'success' | 'error', message: string } | null>(null);

    const mockFaculties = [
        { value: '1', label: 'Denizcilik Fakültesi' },
        { value: '2', label: 'Diş Hekimliği Fakültesi' },
        { value: '3', label: 'Eczacılık Fakültesi' },
        { value: '4', label: 'Eğitim Fakültesi' },
        { value: '5', label: 'Fen Fakültesi' },
        { value: '6', label: 'Güzel Sanatlar Fakültesi' },
        { value: '7', label: 'Hemşirelik Fakültesi' },
        { value: '8', label: 'İktisadi ve İdari Bilimler Fakültesi' },
        { value: '9', label: 'İlahiyat Fakültesi' },
        { value: '10', label: 'İletişim Fakültesi' },
        { value: '11', label: 'İnsan ve Toplum Bilimleri Fakültesi' },
        { value: '12', label: 'Mimarlık Fakültesi' },
        { value: '13', label: 'Mühendislik Fakültesi' },
        { value: '14', label: 'Müzik ve Sahne Sanatları Fakültesi' },
        { value: '15', label: 'Sağlık Bilimleri Fakültesi' },
        { value: '16', label: 'Spor Bilimleri Fakültesi' },
        { value: '17', label: 'Su Ürünleri Fakültesi' },
        { value: '18', label: 'Tıp Fakültesi' },
        { value: '19', label: 'Turizm Fakültesi' }
    ];

    const mockDepartments = [
        // Denizcilik (1)
        { value: 'Denizcilik İşletmeleri Yönetimi', label: 'Denizcilik İşletmeleri Yönetimi', facultyId: '1' },
        { value: 'Deniz Ulaştırma İşletme Mühendisliği', label: 'Deniz Ulaştırma İşletme Mühendisliği', facultyId: '1' },

        // Diş Hekimliği (2)
        { value: 'Diş Hekimliği', label: 'Diş Hekimliği', facultyId: '2' },

        // Eczacılık (3)
        { value: 'Eczacılık', label: 'Eczacılık', facultyId: '3' },

        // Eğitim (4)
        { value: 'Rehberlik ve Psikolojik Danışmanlık', label: 'Rehberlik ve Psikolojik Danışmanlık', facultyId: '4' },
        { value: 'Sınıf Öğretmenliği', label: 'Sınıf Öğretmenliği', facultyId: '4' },
        { value: 'İngilizce Öğretmenliği', label: 'İngilizce Öğretmenliği', facultyId: '4' },
        { value: 'Okul Öncesi Öğretmenliği', label: 'Okul Öncesi Öğretmenliği', facultyId: '4' },
        { value: 'İlköğretim Matematik Öğretmenliği', label: 'İlköğretim Matematik Öğretmenliği', facultyId: '4' },
        { value: 'Türkçe Öğretmenliği', label: 'Türkçe Öğretmenliği', facultyId: '4' },
        { value: 'Özel Eğitim Öğretmenliği', label: 'Özel Eğitim Öğretmenliği', facultyId: '4' },

        // Fen (5)
        { value: 'Biyoloji', label: 'Biyoloji', facultyId: '5' },
        { value: 'Fizik', label: 'Fizik', facultyId: '5' },
        { value: 'Kimya', label: 'Kimya', facultyId: '5' },
        { value: 'Matematik', label: 'Matematik', facultyId: '5' },

        // Güzel Sanatlar (6)
        { value: 'Resim', label: 'Resim', facultyId: '6' },
        { value: 'Heykel', label: 'Heykel', facultyId: '6' },
        { value: 'Grafik', label: 'Grafik', facultyId: '6' },
        { value: 'Seramik', label: 'Seramik', facultyId: '6' },
        { value: 'Tekstil ve Moda Tasarımı', label: 'Tekstil ve Moda Tasarımı', facultyId: '6' },

        // Hemşirelik (7)
        { value: 'Hemşirelik', label: 'Hemşirelik', facultyId: '7' },

        // İİBF (8)
        { value: 'İşletme', label: 'İşletme', facultyId: '8' },
        { value: 'İktisat', label: 'İktisat', facultyId: '8' },
        { value: 'Kamu Yönetimi', label: 'Kamu Yönetimi', facultyId: '8' },
        { value: 'Uluslararası İlişkiler', label: 'Uluslararası İlişkiler', facultyId: '8' },
        { value: 'Maliye', label: 'Maliye', facultyId: '8' },

        // İlahiyat (9)
        { value: 'İlahiyat', label: 'İlahiyat', facultyId: '9' },
        { value: 'İslami İlimler', label: 'İslami İlimler', facultyId: '9' },

        // İletişim (10)
        { value: 'Gazetecilik', label: 'Gazetecilik', facultyId: '10' },
        { value: 'Radyo, Televizyon ve Sinema', label: 'Radyo, Televizyon ve Sinema', facultyId: '10' },
        { value: 'Halkla İlişkiler ve Reklamcılık', label: 'Halkla İlişkiler ve Reklamcılık', facultyId: '10' },

        // İnsan ve Toplum Bilimleri (11)
        { value: 'Psikoloji', label: 'Psikoloji', facultyId: '11' },
        { value: 'Sosyoloji', label: 'Sosyoloji', facultyId: '11' },
        { value: 'Tarih', label: 'Tarih', facultyId: '11' },
        { value: 'Felsefe', label: 'Felsefe', facultyId: '11' },
        { value: 'Türk Dili ve Edebiyatı', label: 'Türk Dili ve Edebiyatı', facultyId: '11' },
        { value: 'İngiliz Dil Bilimi', label: 'İngiliz Dil Bilimi', facultyId: '11' },
        { value: 'Çeviribilim', label: 'Çeviribilim', facultyId: '11' },

        // Mimarlık (12)
        { value: 'Mimarlık', label: 'Mimarlık', facultyId: '12' },
        { value: 'Şehir ve Bölge Planlama', label: 'Şehir ve Bölge Planlama', facultyId: '12' },
        { value: 'İç Mimarlık', label: 'İç Mimarlık', facultyId: '12' },

        // Mühendislik (13)
        { value: 'Bilgisayar Mühendisliği', label: 'Bilgisayar Mühendisliği', facultyId: '13' },
        { value: 'Elektrik-Elektronik Mühendisliği', label: 'Elektrik-Elektronik Mühendisliği', facultyId: '13' },
        { value: 'Makine Mühendisliği', label: 'Makine Mühendisliği', facultyId: '13' },
        { value: 'Çevre Mühendisliği', label: 'Çevre Mühendisliği', facultyId: '13' },
        { value: 'Gıda Mühendisliği', label: 'Gıda Mühendisliği', facultyId: '13' },
        { value: 'İnşaat Mühendisliği', label: 'İnşaat Mühendisliği', facultyId: '13' },
        { value: 'Harita Mühendisliği', label: 'Harita Mühendisliği', facultyId: '13' },
        { value: 'Metalurji ve Malzeme Mühendisliği', label: 'Metalurji ve Malzeme Mühendisliği', facultyId: '13' },

        // Müzik ve Sahne Sanatları (14)
        { value: 'Müzik', label: 'Müzik', facultyId: '14' },
        { value: 'Sahne Sanatları', label: 'Sahne Sanatları', facultyId: '14' },

        // Sağlık Bilimleri (15)
        { value: 'Ebelik', label: 'Ebelik', facultyId: '15' },
        { value: 'Fizyoterapi ve Rehabilitasyon', label: 'Fizyoterapi ve Rehabilitasyon', facultyId: '15' },
        { value: 'Beslenme ve Diyetetik', label: 'Beslenme ve Diyetetik', facultyId: '15' },
        { value: 'Sağlık Yönetimi', label: 'Sağlık Yönetimi', facultyId: '15' },

        // Spor Bilimleri (16)
        { value: 'Beden Eğitimi ve Spor Öğretmenliği', label: 'Beden Eğitimi ve Spor Öğretmenliği', facultyId: '16' },
        { value: 'Antrenörlük Eğitimi', label: 'Antrenörlük Eğitimi', facultyId: '16' },
        { value: 'Spor Yöneticiliği', label: 'Spor Yöneticiliği', facultyId: '16' },
        { value: 'Rekreasyon', label: 'Rekreasyon', facultyId: '16' },

        // Su Ürünleri (17)
        { value: 'Su Ürünleri Mühendisliği', label: 'Su Ürünleri Mühendisliği', facultyId: '17' },

        // Tıp (18)
        { value: 'Tıp', label: 'Tıp', facultyId: '18' },

        // Turizm (19)
        { value: 'Turizm İşletmeciliği', label: 'Turizm İşletmeciliği', facultyId: '19' },
        { value: 'Turizm Rehberliği', label: 'Turizm Rehberliği', facultyId: '19' },
        { value: 'Gastronomi ve Mutfak Sanatları', label: 'Gastronomi ve Mutfak Sanatları', facultyId: '19' }
    ];

    const filteredDepartments = mockDepartments
        .filter(dept => dept.facultyId === facultyId)
        .map(dept => ({ value: dept.value, label: dept.label }));

    const handleFacultyChange = (value: string | null) => {
        setFacultyId(value);
        setDepartmentName(null);
    };

    const handleSendMails = async () => {
        if (!sendToAll && !departmentName) {
            setStatus({ type: 'error', message: 'Lütfen bir bölüm seçiniz veya "Tüm Mezunlara Gönder" seçeneğini işaretleyiniz.' });
            return;
        }

        setLoading(true);
        setStatus(null);

        try {
            // Backend'e artık sayısal ID değil, doğrudan bölümün adını (string) gönderiyoruz
            const response = await api.post('/admin/mail/send-invites', {
                send_to_all: sendToAll,
                department_name: departmentName || "",
            });

            setStatus({ type: 'success', message: response.data?.message || 'Davetiyeler başarıyla mail kuyruğuna eklendi!' });

            setFacultyId(null);
            setDepartmentName(null);
            setSendToAll(false);
        } catch (error: any) {
            setStatus({
                type: 'error',
                message: error.response?.data?.error || 'İşlem sırasında bir hata oluştu. Lütfen tekrar deneyin.',
            });
        } finally {
            setLoading(false);
        }
    };

    return (
        <Container size="md" mt="xl">
            <Card shadow="sm" padding="xl" radius="md" withBorder>
                <Stack gap="lg">
                    <Group style={{ borderBottom: '1px solid #eee' }} pb="md">
                        <ThemeIcon size="lg" radius="md" color="orange.6" variant="light">
                            <IconMail size={24} />
                        </ThemeIcon>
                        <Title order={3} c="meuBlue.9">Mezun Davetiye Yönetimi</Title>
                    </Group>

                    <Text c="dimmed" size="sm">
                        Bu panel üzerinden anket davetiyelerini filtreleyerek mezunlara iletebilirsiniz.
                        E-postalar arka planda asenkron olarak (15 saniyede bir) gönderilecektir.
                    </Text>

                    <Checkbox
                        label="Tüm Mezunlara Gönder (Filtresiz)"
                        checked={sendToAll}
                        onChange={(event) => {
                            setSendToAll(event.currentTarget.checked);
                            if (event.currentTarget.checked) {
                                setFacultyId(null);
                                setDepartmentName(null);
                            }
                        }}
                        color="meuBlue"
                        size="md"
                    />

                    <Select
                        label="Fakülte Filtresi"
                        placeholder={sendToAll ? "Tümüne gönderim seçili" : "Önce fakülte seçin"}
                        data={mockFaculties}
                        value={facultyId}
                        onChange={handleFacultyChange}
                        disabled={sendToAll}
                        searchable
                        clearable
                        size="md"
                    />

                    <Select
                        label="Bölüm Filtresi"
                        placeholder={sendToAll ? "Tümüne gönderim seçili" : (facultyId ? "Gönderilecek bölümü seçin" : "Önce fakülte seçmelisiniz")}
                        data={filteredDepartments}
                        value={departmentName}
                        onChange={setDepartmentName}
                        disabled={sendToAll || !facultyId}
                        searchable
                        clearable
                        size="md"
                    />

                    <Button
                        onClick={handleSendMails}
                        loading={loading}
                        color="orange.6"
                        size="md"
                        fullWidth
                        mt="md"
                    >
                        {sendToAll ? "Tümüne Davetiye Gönder" : "Seçili Bölüme Gönder"}
                    </Button>

                    {status && (
                        <Notification
                            icon={status.type === 'success' ? <IconCheck size={18} /> : <IconX size={18} />}
                            color={status.type === 'success' ? 'teal' : 'red'}
                            onClose={() => setStatus(null)}
                            mt="sm"
                        >
                            {status.message}
                        </Notification>
                    )}
                </Stack>
            </Card>
        </Container>
    );
}