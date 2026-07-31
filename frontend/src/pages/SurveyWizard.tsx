import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
    Container,
    Stepper,
    Button,
    Group,
    Stack,
    Card,
    Loader,
    Center,
    Text,
    Paper
} from "@mantine/core";
import { notifications } from "@mantine/notifications";

import { api } from "../api/client";
import type { StandardizedError } from "../api/client";
import type { SurveyCategory, SurveyAnswer } from "../types";
import QuestionField from "../components/QuestionField";
import ProgressHeader from "../components/ProgressHeader";

export default function SurveyWizard() {
    const navigate = useNavigate();
    const [categories, setCategories] = useState<SurveyCategory[] | null>(null);
    const [answers, setAnswers] = useState<Record<number, string>>({});
    const [active, setActive] = useState(0);
    const [fieldErrors, setFieldErrors] = useState<Record<number, string>>({});
    const [saving, setSaving] = useState(false);
    const [userFaculty, setUserFaculty] = useState<string | null>(null);

    useEffect(() => {
        Promise.all([
            api.get<never, SurveyCategory[]>("/survey/structure"),
            api.get<never, { response: any; answers: SurveyAnswer[] }>("/survey/response"),
            api.get<never, any>("/me")
        ])
            .then(([cats, respData, userData]) => {
                setCategories(cats);

                const realFaculty = userData?.facultyName || userData?.faculty || respData.response?.graduate?.facultyName || respData.response?.graduate?.faculty || null;
                setUserFaculty(realFaculty);

                const initial: Record<number, string> = {};
                respData.answers.forEach((a) => {
                    initial[a.questionId] = a.valueText;
                });
                setAnswers(initial);
                setActive(Math.min(Math.max(respData.response.currentStep - 1, 0), cats.length > 0 ? cats.length - 1 : 0));
            })
            .catch(() => navigate("/giris"));
    }, [navigate]);

    const totalQuestionCount = useMemo(() => {
        if (!categories) return 0;
        return categories.reduce((sum, c) => {
            const visibleInCat = (c.questions || []).filter((q) =>
                !q.targetFaculty ||
                (userFaculty && q.targetFaculty.toLowerCase().includes(userFaculty.toLowerCase()))
            );
            return sum + visibleInCat.length;
        }, 0);
    }, [categories, userFaculty]);

    const answeredCount = useMemo(
        () => Object.values(answers).filter((v) => v && v !== "").length,
        [answers]
    );

    if (!categories) {
        return (
            <Center h="100vh" bg="#0A192F">
                <Loader color="orange" />
            </Center>
        );
    }

    const currentCategory = categories[active];

    if (!currentCategory) {
        return (
            <Center h="100vh" bg="#0A192F">
                <Text c="white">Henüz anket için soru/kategori eklenmemiş.</Text>
            </Center>
        );
    }

    const isLastStep = active === categories.length - 1;

    const visibleQuestions = (currentCategory.questions || []).filter(
        (q) => !q.targetFaculty ||
            (userFaculty && q.targetFaculty.toLowerCase().includes(userFaculty.toLowerCase()))
    );

    function setAnswer(questionId: number, value: string) {
        setAnswers((prev) => ({ ...prev, [questionId]: value }));
        setFieldErrors((prev) => {
            if (!prev[questionId]) return prev;
            const next = { ...prev };
            delete next[questionId];
            return next;
        });
    }

    function validateCurrentStep(): boolean {
        const errors: Record<number, string> = {};
        for (const q of visibleQuestions) {
            if (q.required && (!answers[q.id] || answers[q.id] === "")) {
                errors[q.id] = "Bu soru zorunludur.";
            }
        }
        setFieldErrors(errors);
        return Object.keys(errors).length === 0;
    }

    async function handleNext() {
        if (!validateCurrentStep()) {
            notifications.show({
                color: "red",
                title: "Eksik cevaplar var",
                message: "Lütfen bu adımdaki zorunlu soruları yanıtlayın.",
            });
            return;
        }

        setSaving(true);
        try {
            const stepAnswers = visibleQuestions.map((q) => ({
                questionId: q.id,
                value: answers[q.id] ?? "",
            }));
            const requiredQuestionIds = visibleQuestions.filter((q) => q.required).map((q) => q.id);

            await api.put(`/survey/response/step/${active + 1}`, {
                answers: stepAnswers,
                requiredQuestionIds,
            });

            if (isLastStep) {
                await api.post("/survey/response/complete");
                navigate("/tesekkurler");
                return;
            }
            setActive((s) => s + 1);
            window.scrollTo({ top: 0, behavior: "smooth" });
        } catch (err) {
            const e = err as StandardizedError;
            notifications.show({ color: "red", title: "Kaydedilemedi", message: e.message });
        } finally {
            setSaving(false);
        }
    }

    function handleBack() {
        setActive((s) => Math.max(0, s - 1));
        window.scrollTo({ top: 0, behavior: "smooth" });
    }

    return (
        <div style={{ position: "relative", minHeight: "100vh", width: "100%", display: "flex", alignItems: "center", justifyContent: "center", padding: "1rem 0" }}>

            <div style={{
                position: "fixed", top: 0, left: 0, right: 0, bottom: 0,
                backgroundImage: "url('/mersin-kampus.jpg')",
                backgroundSize: "cover",
                backgroundPosition: "center",
                zIndex: -2
            }} />

            <div style={{
                position: "fixed", top: 0, left: 0, right: 0, bottom: 0,
                backgroundColor: "rgba(10, 25, 47, 0.85)",
                zIndex: -1
            }} />

            <Container size={800} w="100%" px="md" style={{ zIndex: 1 }}>
                <Paper shadow="xl" radius="md" p="xl" bg="white">
                    <Stack gap="xl">

                        <Group align="center" wrap="nowrap" gap="sm">
                            <img src="/logo.png" alt="Mersin Üniversitesi Logo" style={{ width: 45, height: 45, objectFit: "contain" }} />
                            <Stack gap={0}>
                                <Text fw={700} size="xs" tt="uppercase" style={{ color: "#F26722" }}>
                                    Mersin Üniversitesi · Mezun Anketi
                                </Text>
                                <Text fw={700} fz="xl" style={{ color: "#0A192F", lineHeight: 1.2 }}>
                                    {currentCategory.title}
                                </Text>
                            </Stack>
                        </Group>

                        <ProgressHeader answeredCount={answeredCount} totalCount={totalQuestionCount} />

                        <Stepper
                            active={active}
                            onStepClick={setActive}
                            allowNextStepsSelect={false}
                            color="orange"
                            size="sm"
                            wrap={false}
                        >
                            {categories.map((cat) => (
                                <Stepper.Step key={cat.id} />
                            ))}
                        </Stepper>

                        <Card withBorder bg="#f8f9fa" p="md">
                            <Stack gap="xl">
                                {visibleQuestions.map((q) => (
                                    <QuestionField
                                        key={q.id}
                                        question={q}
                                        value={answers[q.id] ?? ""}
                                        onChange={(v) => setAnswer(q.id, v)}
                                        error={fieldErrors[q.id]}
                                    />
                                ))}
                            </Stack>
                        </Card>

                        <Group justify="space-between">
                            <Button variant="outline" color="gray" onClick={handleBack} disabled={active === 0 || saving} size="sm">
                                Geri
                            </Button>
                            <Button onClick={handleNext} loading={saving} style={{ backgroundColor: "#F26722", color: "white" }} size="sm">
                                {isLastStep ? "Anketi Tamamla" : "İleri"}
                            </Button>
                        </Group>
                    </Stack>
                </Paper>
            </Container>
        </div>
    );
}