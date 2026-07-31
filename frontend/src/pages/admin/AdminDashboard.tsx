import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
    Container, Grid, Card, Title, Text, Group, Stack, Loader, Center, SimpleGrid, RingProgress,
} from "@mantine/core";
import {
    PieChart, Pie, Cell, BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid,
} from "recharts";

import { api, adminTokenStorage } from "../../api/client";

interface OverviewStats {
    totalGraduates: number;
    totalResponses: number;
    completedCount: number;
    completionRate: number;
    dropOffByStep: Record<string, number>;
}

interface DistributionItem {
    label: string;
    count: number;
}

const PIE_COLORS = ["#0A192F", "#F26722", "#1C355E", "#F58A55", "#2E528E", "#F8AD88"];
const DONUT_COLORS = ["#F26722", "#0A192F", "#F8AD88", "#1C355E", "#F58A55"];

// Kurşun geçirmez dizi kontrol fonksiyonu: Veri yoksa çökmek yerine boş dizi döner
const safeArray = (data: any) => Array.isArray(data) ? data : [];

export default function AdminDashboard() {
    const navigate = useNavigate();
    const [overview, setOverview] = useState<OverviewStats | null>(null);

    const [careerGoalsDist, setCareerGoalsDist] = useState<DistributionItem[] | null>(null);
    const [workPositionDist, setWorkPositionDist] = useState<DistributionItem[] | null>(null);
    const [jobFindingTimeDist, setJobFindingTimeDist] = useState<DistributionItem[] | null>(null);
    const [eduQualityDist, setEduQualityDist] = useState<DistributionItem[] | null>(null);

    useEffect(() => {
        if (!adminTokenStorage.get()) {
            navigate("/admin/giris");
            return;
        }

        Promise.all([
            api.get<never, OverviewStats>("/admin/stats/overview"),
            api.get<never, DistributionItem[]>("/admin/stats/question/A1"),
            api.get<never, DistributionItem[]>("/admin/stats/question/C20"),
            api.get<never, DistributionItem[]>("/admin/stats/question/C16"),
            api.get<never, DistributionItem[]>("/admin/stats/question/B15"),
        ])
            .then(([ov, careerGoals, workPosition, jobTime, eduQuality]) => {
                setOverview(ov);
                setCareerGoalsDist(careerGoals);
                setWorkPositionDist(workPosition);
                setJobFindingTimeDist(jobTime);
                setEduQualityDist(eduQuality);
            })
            .catch(() => navigate("/admin/giris"));
    }, [navigate]);

    if (!overview) {
        return (
            <Center h="100vh">
                <Loader color="orange" />
            </Center>
        );
    }

    const dropOffData = Object.entries(overview.dropOffByStep || {}).map(([step, count]) => ({
        step: `Adım ${step}`,
        count,
    }));

    // Grafiklere gidecek tüm verileri kurşun geçirmez kalkanımızdan geçiriyoruz
    const safeCareerData = safeArray(careerGoalsDist);
    const safeWorkData = safeArray(workPositionDist);
    const safeJobData = safeArray(jobFindingTimeDist);
    const safeEduData = safeArray(eduQualityDist);

    return (
        <Container size={1100} py={40}>
            <Stack gap="xl">

                <Group align="center" wrap="nowrap" gap="md" mb="xs">
                    <img src="/logo.png" alt="Mersin Üniversitesi Logo" style={{ width: 60, height: 60, objectFit: "contain" }} />
                    <Stack gap={0}>
                        <Text fw={700} size="sm" tt="uppercase" style={{ color: "#F26722" }}>
                            Mersin Üniversitesi
                        </Text>
                        <Title order={2} style={{ color: "#0A192F" }}>Mezun Anketi Yönetim Panosu</Title>
                    </Stack>
                </Group>

                <SimpleGrid cols={{ base: 1, sm: 4 }}>
                    <StatCard label="Toplam Mezun" value={overview.totalGraduates || 0} />
                    <StatCard label="Ankete Başlayan" value={overview.totalResponses || 0} />
                    <StatCard label="Tamamlayan" value={overview.completedCount || 0} />
                    <Card withBorder radius="lg">
                        <Group>
                            <RingProgress
                                size={70}
                                thickness={7}
                                roundCaps
                                sections={[{ value: overview.completionRate || 0, color: "#F26722" }]}
                                label={
                                    <Text size="xs" ta="center" fw={700} style={{ color: "#0A192F" }}>
                                        {(overview.completionRate || 0).toFixed(0)}%
                                    </Text>
                                }
                            />
                            <Text size="sm" c="dimmed">
                                Tamamlanma Oranı
                            </Text>
                        </Group>
                    </Card>
                </SimpleGrid>

                <Grid>
                    <Grid.Col span={{ base: 12, md: 6 }}>
                        <Card withBorder h={360}>
                            <Title order={4} mb="md" style={{ color: "#0A192F" }}>
                                Kariyer Hedefleri
                            </Title>
                            <Text size="xs" c="dimmed" mb="sm">
                                Soru A1: Mezun olurken kariyer planlamanızdaki öncelikli hedefiniz neydi?
                            </Text>
                            <ResponsiveContainer width="100%" height={260}>
                                <PieChart>
                                    <Pie
                                        data={safeCareerData}
                                        dataKey="count"
                                        nameKey="label"
                                        cx="50%"
                                        cy="50%"
                                        outerRadius={80}
                                        label={(entry) => entry.label?.length > 15 ? `${entry.label.substring(0, 15)}...` : entry.label}
                                    >
                                        {safeCareerData.map((_, i) => (
                                            <Cell key={`cell-${i}`} fill={PIE_COLORS[i % PIE_COLORS.length]} />
                                        ))}
                                    </Pie>
                                    <Tooltip formatter={(value: number) => [`${value} Kişi`, "Yanıt"]} />
                                </PieChart>
                            </ResponsiveContainer>
                        </Card>
                    </Grid.Col>

                    <Grid.Col span={{ base: 12, md: 6 }}>
                        <Card withBorder h={360}>
                            <Title order={4} mb="md" style={{ color: "#0A192F" }}>
                                Çalışma Durumu / Pozisyon
                            </Title>
                            <Text size="xs" c="dimmed" mb="sm">
                                Soru C20: Şu anki pozisyonunuz/unvanınız hangi seviyeye daha yakın?
                            </Text>
                            <ResponsiveContainer width="100%" height={260}>
                                <BarChart data={safeWorkData} margin={{ bottom: 20 }}>
                                    <CartesianGrid strokeDasharray="3 3" vertical={false} />
                                    <XAxis
                                        dataKey="label"
                                        tickFormatter={(v) => v?.length > 12 ? `${v.substring(0, 10)}...` : v}
                                        interval={0}
                                        angle={-45}
                                        textAnchor="end"
                                        height={60}
                                        tick={{ fontSize: 12 }}
                                    />
                                    <YAxis allowDecimals={false} />
                                    <Tooltip cursor={{ fill: 'transparent' }} />
                                    <Bar dataKey="count" name="Kişi Sayısı" fill="#0A192F" radius={[4, 4, 0, 0]} barSize={40} />
                                </BarChart>
                            </ResponsiveContainer>
                        </Card>
                    </Grid.Col>

                    <Grid.Col span={{ base: 12, md: 6 }}>
                        <Card withBorder h={360}>
                            <Title order={4} mb="md" style={{ color: "#0A192F" }}>
                                İlk İşi Bulma Süresi
                            </Title>
                            <Text size="xs" c="dimmed" mb="sm">
                                Soru C16: Mezuniyet sonrası ilk işinizi bulma süreniz ne kadar sürdü?
                            </Text>
                            <ResponsiveContainer width="100%" height={260}>
                                <BarChart data={safeJobData} margin={{ bottom: 20 }}>
                                    <CartesianGrid strokeDasharray="3 3" vertical={false} />
                                    <XAxis
                                        dataKey="label"
                                        tickFormatter={(v) => v?.length > 15 ? `${v.substring(0, 12)}...` : v}
                                        interval={0}
                                        angle={-25}
                                        textAnchor="end"
                                        height={50}
                                        tick={{ fontSize: 12 }}
                                    />
                                    <YAxis allowDecimals={false} />
                                    <Tooltip cursor={{ fill: 'transparent' }} />
                                    <Bar dataKey="count" name="Kişi Sayısı" fill="#F26722" radius={[4, 4, 0, 0]} barSize={40} />
                                </BarChart>
                            </ResponsiveContainer>
                        </Card>
                    </Grid.Col>

                    <Grid.Col span={{ base: 12, md: 6 }}>
                        <Card withBorder h={360}>
                            <Title order={4} mb="md" style={{ color: "#0A192F" }}>
                                Genel Eğitim Kalitesi
                            </Title>
                            <Text size="xs" c="dimmed" mb="sm">
                                Soru B15: Üniversitedeki eğitim kalitesini genel olarak nasıl değerlendirirsiniz?
                            </Text>
                            <ResponsiveContainer width="100%" height={260}>
                                <PieChart>
                                    <Pie
                                        data={safeEduData}
                                        dataKey="count"
                                        nameKey="label"
                                        cx="50%"
                                        cy="50%"
                                        innerRadius={50}
                                        outerRadius={80}
                                        label={(entry) => entry.label}
                                    >
                                        {safeEduData.map((_, i) => (
                                            <Cell key={`cell-${i}`} fill={DONUT_COLORS[i % DONUT_COLORS.length]} />
                                        ))}
                                    </Pie>
                                    <Tooltip formatter={(value: number) => [`${value} Kişi`, "Yanıt"]} />
                                </PieChart>
                            </ResponsiveContainer>
                        </Card>
                    </Grid.Col>

                    <Grid.Col span={12}>
                        <Card withBorder h={340}>
                            <Title order={4} mb="xs" style={{ color: "#0A192F" }}>
                                Adım Bazlı Terk (Bounce) Analizi
                            </Title>
                            <Text size="sm" c="dimmed" mb="lg">
                                Anketi yarıda bırakan mezunların hangi adımda kaldığını gösterir.
                            </Text>
                            <ResponsiveContainer width="100%" height={220}>
                                <BarChart data={dropOffData}>
                                    <CartesianGrid strokeDasharray="3 3" vertical={false} />
                                    <XAxis dataKey="step" />
                                    <YAxis allowDecimals={false} />
                                    <Tooltip cursor={{ fill: 'transparent' }} />
                                    <Bar dataKey="count" name="Terk Eden Kişi" fill="#1C355E" radius={[4, 4, 0, 0]} maxBarSize={60} />
                                </BarChart>
                            </ResponsiveContainer>
                        </Card>
                    </Grid.Col>
                </Grid>
            </Stack>
        </Container>
    );
}

function StatCard({ label, value }: { label: string; value: number }) {
    return (
        <Card withBorder radius="lg">
            <Stack gap={2}>
                <Text size="xs" c="dimmed" tt="uppercase">
                    {label}
                </Text>
                <Text size="xl" fw={700} style={{ color: "#F26722" }}>
                    {value}
                </Text>
            </Stack>
        </Card>
    );
}