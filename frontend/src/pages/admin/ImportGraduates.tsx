import { useState } from "react";
import Papa from "papaparse";
import { Dropzone, MS_EXCEL_MIME_TYPE } from "@mantine/dropzone";
import {
    Container,
    Title,
    Text,
    Tabs,
    Card,
    Stack,
    Group,
    Button,
    Table,
    Switch,
    TextInput,
    NumberInput,
    ActionIcon,
    Alert,
    ScrollArea,
    Badge,
    Divider,
    Code,
    Select,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";

import { api } from "../../api/client";
import type { StandardizedError } from "../../api/client";
import type { GraduateImportRow, ImportResponse } from "../../types";

const CSV_COLUMNS =
    "obsHashId,firstName,facultyName,departmentName,graduationYear,studentNoHash,email,phone";

export const FACULTY_DEPARTMENTS: Record<string, string[]> = {
    // --- FAKÜLTELER (Lisans) ---
    "Mühendislik Fakültesi": [
        "Bilgisayar Mühendisliği", "Elektrik-Elektronik Mühendisliği", "Makine Mühendisliği",
        "Çevre Mühendisliği", "Gıda Mühendisliği", "Harita Mühendisliği",
        "İnşaat Mühendisliği", "Kimya Mühendisliği", "Metalurji ve Malzeme Mühendisliği", "Jeoloji Mühendisliği"
    ],
    "Tıp Fakültesi": ["Tıp"],
    "Diş Hekimliği Fakültesi": ["Diş Hekimliği"],
    "Eczacılık Fakültesi": ["Eczacılık"],
    "Eğitim Fakültesi": [
        "Sınıf Öğretmenliği", "Okul Öncesi Öğretmenliği", "Rehberlik ve Psikolojik Danışmanlık",
        "İngilizce Öğretmenliği", "İlköğretim Matematik Öğretmenliği", "Türkçe Öğretmenliği",
        "Fen Bilgisi Öğretmenliği", "Özel Eğitim Öğretmenliği"
    ],
    "Fen Fakültesi": ["Matematik", "Fizik", "Kimya", "Biyoloji", "İstatistik"],
    "İnsani ve Toplum Bilimleri Fakültesi": [
        "Psikoloji", "Sosyoloji", "Tarih", "Türk Dili ve Edebiyatı",
        "Mütercim ve Tercümanlık", "Felsefe", "Arkeoloji", "Sanat Tarihi"
    ],
    "İktisadi ve İdari Bilimler Fakültesi": [
        "İşletme", "İktisat", "Kamu Yönetimi", "Uluslararası İlişkiler",
        "Maliye", "Çalışma Ekonomisi ve Endüstri İlişkileri", "Yönetim Bilişim Sistemleri"
    ],
    "Mimarlık Fakültesi": ["Mimarlık", "Şehir ve Bölge Planlama", "İç Mimarlık"],
    "İletişim Fakültesi": ["Gazetecilik", "Radyo, Televizyon ve Sinema", "Halkla İlişkiler ve Reklamcılık"],
    "Turizm Fakültesi": ["Turizm İşletmeciliği", "Gastronomi ve Mutfak Sanatları", "Turizm Rehberliği", "Rekreasyon Yönetimi"],
    "Denizcilik Fakültesi": ["Denizcilik İşletmeleri Yönetimi", "Deniz Ulaştırma İşletme Mühendisliği", "Gemi Makineleri İşletme Mühendisliği"],
    "Sağlık Bilimleri Fakültesi": ["Ebelik", "Fizyoterapi ve Rehabilitasyon", "Sağlık Yönetimi", "Sosyal Hizmet", "Odyoloji", "Ergoterapi"],
    "Hemşirelik Fakültesi": ["Hemşirelik"],
    "Spor Bilimleri Fakültesi": ["Beden Eğitimi ve Spor Öğretmenliği", "Antrenörlük Eğitimi", "Spor Yöneticiliği", "Rekreasyon"],
    "Su Ürünleri Fakültesi": ["Su Ürünleri Mühendisliği"],
    "İslami İlimler Fakültesi": ["İslami İlimler"],
    "Güzel Sanatlar Fakültesi": ["Resim", "Heykel", "Grafik", "Tekstil ve Moda Tasarımı", "Seramik ve Cam", "Müzik"],

    // --- YÜKSEKOKULLAR (Lisans) ---
    "Devlet Konservatuvarı": ["Müzik", "Sahne Sanatları"],
    "Takı Teknolojisi ve Tasarımı Yüksekokulu": ["Takı Teknolojisi ve Tasarımı"],
    "Erdemli Uygulamalı Teknoloji ve İşletmecilik Yüksekokulu": ["Bilişim Sistemleri ve Teknolojileri", "Yönetim Bilişim Sistemleri"],
    "Yabancı Diller Yüksekokulu": ["Yabancı Dil Hazırlık"],

    // --- MESLEK YÜKSEKOKULLARI (Önlisans) ---
    "Teknik Bilimler Meslek Yüksekokulu": [
        "Bilgisayar Programcılığı", "Makine", "Elektrik", "İnşaat Teknolojisi",
        "Elektronik Teknolojisi", "Mekatronik", "Otomotiv Teknolojisi"
    ],
    "Sosyal Bilimler Meslek Yüksekokulu": [
        "Muhasebe ve Vergi Uygulamaları", "Büro Yönetimi ve Yönetici Asistanlığı",
        "Dış Ticaret", "İşletme Yönetimi", "Lojistik", "Pazarlama", "Turizm ve Otel İşletmeciliği"
    ],
    "Sağlık Hizmetleri Meslek Yüksekokulu": [
        "İlk ve Acil Yardım", "Tıbbi Dokümantasyon ve Sekreterlik", "Tıbbi Görüntüleme Teknikleri",
        "Tıbbi Laboratuvar Teknikleri", "Anestezi", "Yaşlı Bakımı", "Çocuk Gelişimi", "Ağız ve Diş Sağlığı"
    ],
    "Deniz ve Ticaret Meslek Yüksekokulu": ["Deniz ve Liman İşletmeciliği", "Dış Ticaret", "Lojistik"],
    "Erdemli Meslek Yüksekokulu": ["Bankacılık ve Sigortacılık", "Bilgisayar Programcılığı", "Muhasebe ve Vergi Uygulamaları", "İşletme Yönetimi"],
    "Mut Meslek Yüksekokulu": ["Bahçe Tarımı", "Bilgisayar Programcılığı", "Muhasebe ve Vergi Uygulamaları"],
    "Silifke Meslek Yüksekokulu": ["Aşçılık", "Organik Tarım", "Turizm ve Otel İşletmeciliği"],
    "Anamur Meslek Yüksekokulu": ["Turizm ve Otel İşletmeciliği", "Büro Yönetimi ve Yönetici Asistanlığı", "Pazarlama"],
    "Gülnar Mustafa Baysan Meslek Yüksekokulu": ["Ormancılık ve Orman Ürünleri", "Maliye"],

    // --- ENSTİTÜLER (Lisansüstü) ---
    "Fen Bilimleri Enstitüsü": ["Yüksek Lisans Programları", "Doktora Programları"],
    "Sosyal Bilimler Enstitüsü": ["Yüksek Lisans Programları", "Doktora Programları"],
    "Sağlık Bilimleri Enstitüsü": ["Yüksek Lisans Programları", "Doktora Programları"],
    "Eğitim Bilimleri Enstitüsü": ["Yüksek Lisans Programları", "Doktora Programları"],
    "Güzel Sanatlar Enstitüsü": ["Yüksek Lisans Programları", "Sanatta Yeterlik Programları"]
};

const FACULTIES = Object.keys(FACULTY_DEPARTMENTS);

function emptyRow(): GraduateImportRow {
    return {
        obsHashId: "",
        firstName: "",
        facultyName: "",
        departmentName: "",
        graduationYear: "",
        studentNoHash: "",
        email: "",
        phone: "",
    };
}

export default function ImportGraduates() {
    const [csvRows, setCsvRows] = useState<GraduateImportRow[]>([]);
    const [csvFileName, setCsvFileName] = useState<string | null>(null);
    const [manualRows, setManualRows] = useState<GraduateImportRow[]>([emptyRow()]);
    const [sendInvites, setSendInvites] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [lastResult, setLastResult] = useState<ImportResponse | null>(null);

    // --- CSV/Excel yükleme ---
    function handleDrop(files: File[]) {
        const file = files[0];
        if (!file) return;
        setCsvFileName(file.name);

        Papa.parse<Record<string, string>>(file, {
            header: true,
            skipEmptyLines: true,
            complete: (result) => {
                const rows: GraduateImportRow[] = result.data.map((r) => ({
                    obsHashId: r.obsHashId?.trim() ?? "",
                    firstName: r.firstName?.trim() ?? "",
                    facultyName: r.facultyName?.trim() ?? "",
                    departmentName: r.departmentName?.trim() ?? "",
                    graduationYear: r.graduationYear ? Number(r.graduationYear) : "",
                    studentNoHash: r.studentNoHash?.trim() ?? "",
                    email: r.email?.trim() ?? "",
                    phone: r.phone?.trim() ?? "",
                }));
                setCsvRows(rows.filter((r) => r.obsHashId && r.firstName));

                if (result.errors.length > 0) {
                    notifications.show({
                        color: "yellow",
                        title: "Dosya kısmen okundu",
                        message: `${result.errors.length} satırda ayrıştırma sorunu oldu, bu satırlar atlandı.`,
                    });
                }
            },
            error: () => {
                notifications.show({ color: "red", title: "Dosya okunamadı", message: "Lütfen geçerli bir CSV dosyası yükleyin." });
            },
        });
    }

    async function submitBatch(graduates: GraduateImportRow[]) {
        if (graduates.length === 0) {
            notifications.show({ color: "red", title: "Kayıt yok", message: "İçe aktarılacak en az bir mezun kaydı olmalı." });
            return;
        }

        // ZORUNLU ALAN KONTROLÜ: Fakülte ve Bölüm de eklendi
        const invalid = graduates.some((g) => !g.obsHashId || !g.firstName || !g.facultyName || !g.departmentName);
        if (invalid) {
            notifications.show({
                color: "red",
                title: "Eksik alan",
                message: "OBS Hash ID, Ad, Fakülte ve Bölüm alanları tüm satırlarda zorunludur.",
            });
            return;
        }

        setSubmitting(true);
        try {
            const payload = {
                sendInvites,
                graduates: graduates.map((g) => ({
                    ...g,
                    graduationYear: g.graduationYear === "" ? 0 : Number(g.graduationYear),
                })),
            };
            const data = await api.post<never, ImportResponse>("/admin/graduates/import", payload);
            setLastResult(data);
            notifications.show({
                color: "green",
                title: "İçe aktarma tamamlandı",
                message: `${data.result.inserted} eklendi, ${data.result.updated} güncellendi, ${data.result.failed} başarısız.`,
            });
        } catch (err) {
            const e = err as StandardizedError;
            notifications.show({ color: "red", title: "İçe aktarma başarısız", message: e.message });
        } finally {
            setSubmitting(false);
        }
    }

    // --- Manuel form yardımcıları ---
    function updateManualRow(index: number, patch: Partial<GraduateImportRow>) {
        setManualRows((rows) =>
            rows.map((r, i) => {
                if (i !== index) return r;
                const updated = { ...r, ...patch };
                // Eğer fakülte değiştirildiyse ve seçilen bölüm yeni fakültede yoksa bölümü sıfırla
                if (patch.facultyName !== undefined && patch.facultyName !== r.facultyName) {
                    const validDepts = FACULTY_DEPARTMENTS[patch.facultyName] || [];
                    if (!validDepts.includes(updated.departmentName)) {
                        updated.departmentName = "";
                    }
                }
                return updated;
            })
        );
    }

    function addManualRow() {
        setManualRows((rows) => [...rows, emptyRow()]);
    }
    function removeManualRow(index: number) {
        setManualRows((rows) => rows.filter((_, i) => i !== index));
    }

    return (
        <Container size={960} py="md">
            <Stack gap="lg">

                {/* Kurumsal Logolu Başlık */}
                <Group align="center" wrap="nowrap" gap="md" mb="xs">
                    <img src="/logo.png" alt="Mersin Üniversitesi Logo" style={{ width: 55, height: 55, objectFit: "contain" }} />
                    <Stack gap={0}>
                        <Text fw={700} size="sm" tt="uppercase" style={{ color: "#F26722" }}>
                            Mersin Üniversitesi
                        </Text>
                        <Title order={2} style={{ color: "#0A192F" }}>Mezun Anketi Yönetimi</Title>
                    </Stack>
                </Group>

                <Text c="dimmed" size="sm" mt="-10px">
                    OBS'den alınan mezun verisini CSV/Excel dosyasıyla toplu, ya da tek tek manuel olarak
                    içe aktarın. Soyad alanı kasıtlı olarak bulunmuyor; e-posta ve telefon backend'e
                    ulaştığı anda şifrelenip öyle saklanıyor.
                </Text>

                <Switch
                    color="orange"
                    label="Kayıt sonrası mezunlara davet e-postası gönder"
                    checked={sendInvites}
                    onChange={(e) => setSendInvites(e.currentTarget.checked)}
                />

                <Tabs defaultValue="file" color="orange">
                    <Tabs.List>
                        <Tabs.Tab value="file">Dosya ile İçe Aktar (CSV/Excel)</Tabs.Tab>
                        <Tabs.Tab value="manual">Manuel Ekle</Tabs.Tab>
                    </Tabs.List>

                    {/* --- Dosya ile içe aktarma --- */}
                    <Tabs.Panel value="file" pt="md">
                        <Card withBorder>
                            <Stack gap="md" p="sm">
                                <Dropzone
                                    onDrop={handleDrop}
                                    accept={["text/csv", ...MS_EXCEL_MIME_TYPE]}
                                    maxFiles={1}
                                >
                                    <Group justify="center" gap="md" mih={120} style={{ pointerEvents: "none" }}>
                                        <Stack gap={2} align="center">
                                            <Text size="sm" fw={500}>
                                                Dosyayı buraya sürükleyin veya seçmek için tıklayın
                                            </Text>
                                            <Text size="xs" c="dimmed">
                                                .csv (Excel'den "CSV olarak kaydet" ile de üretebilirsiniz)
                                            </Text>
                                        </Stack>
                                    </Group>
                                </Dropzone>

                                <Text size="xs" c="dimmed">
                                    Beklenen kolon başlıkları: <Code>{CSV_COLUMNS}</Code>
                                </Text>

                                {csvFileName && (
                                    <Group gap="xs">
                                        <Badge color="orange" variant="light">
                                            {csvFileName}
                                        </Badge>
                                        <Text size="sm" c="dimmed">
                                            {csvRows.length} geçerli satır bulundu
                                        </Text>
                                    </Group>
                                )}

                                {csvRows.length > 0 && (
                                    <ScrollArea h={260}>
                                        <Table striped highlightOnHover stickyHeader>
                                            <Table.Thead>
                                                <Table.Tr>
                                                    <Table.Th>OBS Hash ID</Table.Th>
                                                    <Table.Th>Ad</Table.Th>
                                                    <Table.Th>Fakülte</Table.Th>
                                                    <Table.Th>Bölüm</Table.Th>
                                                    <Table.Th>E-posta</Table.Th>
                                                </Table.Tr>
                                            </Table.Thead>
                                            <Table.Tbody>
                                                {csvRows.slice(0, 50).map((r, i) => (
                                                    <Table.Tr key={i}>
                                                        <Table.Td>{r.obsHashId.slice(0, 12)}...</Table.Td>
                                                        <Table.Td>{r.firstName}</Table.Td>
                                                        <Table.Td>{r.facultyName}</Table.Td>
                                                        <Table.Td>{r.departmentName}</Table.Td>
                                                        <Table.Td>{r.email}</Table.Td>
                                                    </Table.Tr>
                                                ))}
                                            </Table.Tbody>
                                        </Table>
                                    </ScrollArea>
                                )}

                                <Group justify="flex-end">
                                    <Button
                                        style={{ backgroundColor: "#F26722", color: "white" }}
                                        loading={submitting}
                                        disabled={csvRows.length === 0}
                                        onClick={() => submitBatch(csvRows)}
                                    >
                                        {csvRows.length} Mezunu İçe Aktar
                                    </Button>
                                </Group>
                            </Stack>
                        </Card>
                    </Tabs.Panel>

                    {/* --- Manuel giriş --- */}
                    <Tabs.Panel value="manual" pt="md">
                        <Card withBorder>
                            <Stack gap="lg" p="sm">
                                {manualRows.map((row, i) => {
                                    // Seçilen fakülteye ait bölümleri al, yoksa boş liste ver
                                    const availableDepartments = FACULTY_DEPARTMENTS[row.facultyName] || [];

                                    return (
                                        <Stack key={i} gap="xs">
                                            <Group justify="space-between">
                                                <Text fw={600} size="sm" style={{ color: "#0A192F" }}>
                                                    Mezun #{i + 1}
                                                </Text>
                                                {manualRows.length > 1 && (
                                                    <ActionIcon color="red" variant="subtle" onClick={() => removeManualRow(i)}>
                                                        ✕
                                                    </ActionIcon>
                                                )}
                                            </Group>
                                            <Group grow>
                                                <TextInput
                                                    label="OBS Hash ID"
                                                    required
                                                    value={row.obsHashId}
                                                    onChange={(e) => updateManualRow(i, { obsHashId: e.currentTarget.value })}
                                                />
                                                <TextInput
                                                    label="Ad"
                                                    required
                                                    value={row.firstName}
                                                    onChange={(e) => updateManualRow(i, { firstName: e.currentTarget.value })}
                                                />
                                            </Group>
                                            <Group grow>
                                                <Select
                                                    label="Fakülte"
                                                    searchable
                                                    clearable
                                                    required
                                                    withAsterisk
                                                    placeholder="Önce fakülte seçin"
                                                    data={FACULTIES}
                                                    value={row.facultyName}
                                                    onChange={(val) => updateManualRow(i, { facultyName: val || "" })}
                                                />
                                                <Select
                                                    label="Bölüm"
                                                    searchable
                                                    clearable
                                                    required
                                                    withAsterisk
                                                    placeholder={row.facultyName ? "Bölüm seçin" : "Önce fakülte seçmelisiniz"}
                                                    disabled={!row.facultyName}
                                                    data={availableDepartments}
                                                    value={row.departmentName}
                                                    onChange={(val) => updateManualRow(i, { departmentName: val || "" })}
                                                />
                                                <NumberInput
                                                    label="Mezuniyet Yılı"
                                                    value={row.graduationYear}
                                                    onChange={(v) => updateManualRow(i, { graduationYear: v === "" ? "" : Number(v) })}
                                                    min={1980}
                                                    max={2100}
                                                />
                                            </Group>
                                            <Group grow>
                                                <TextInput
                                                    label="Öğrenci No Hash (opsiyonel)"
                                                    value={row.studentNoHash}
                                                    onChange={(e) => updateManualRow(i, { studentNoHash: e.currentTarget.value })}
                                                />
                                                <TextInput
                                                    label="E-posta"
                                                    type="email"
                                                    value={row.email}
                                                    required
                                                    withAsterisk
                                                    onChange={(e) => updateManualRow(i, { email: e.currentTarget.value })}
                                                />
                                                <TextInput
                                                    label="Telefon"
                                                    value={row.phone}
                                                    onChange={(e) => updateManualRow(i, { phone: e.currentTarget.value })}
                                                />
                                            </Group>
                                            {i < manualRows.length - 1 && <Divider mt="xs" />}
                                        </Stack>
                                    );
                                })}

                                <Group justify="space-between">
                                    <Button variant="outline" style={{ borderColor: "#0A192F", color: "#0A192F" }} onClick={addManualRow}>
                                        + Başka Mezun Ekle
                                    </Button>
                                    <Button style={{ backgroundColor: "#F26722", color: "white" }} loading={submitting} onClick={() => submitBatch(manualRows)}>
                                        Kaydet ve İçe Aktar
                                    </Button>
                                </Group>
                            </Stack>
                        </Card>
                    </Tabs.Panel>
                </Tabs>

                {lastResult && (
                    <Alert color="green" title="Son içe aktarma sonucu">
                        <Text size="sm">
                            {lastResult.result.inserted} yeni kayıt, {lastResult.result.updated} güncelleme,{" "}
                            {lastResult.result.failed} hata.
                        </Text>
                        {lastResult.result.errors && lastResult.result.errors.length > 0 && (
                            <Text size="xs" c="dimmed" mt={4}>
                                Hatalar: {lastResult.result.errors.join(", ")}
                            </Text>
                        )}
                        {Object.keys(lastResult.inviteLinks).length > 0 && (
                            <Text size="xs" c="dimmed" mt={4}>
                                {Object.keys(lastResult.inviteLinks).length} mezun için giriş linki üretildi
                                {sendInvites ? " ve e-posta kuyruğuna eklendi." : "."}
                            </Text>
                        )}
                    </Alert>
                )}
            </Stack>
        </Container>
    );
}