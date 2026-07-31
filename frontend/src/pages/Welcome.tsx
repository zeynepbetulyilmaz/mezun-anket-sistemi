import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Container, Card, Title, Text, Button, Stack, Loader, Center } from "@mantine/core";

import { api } from "../api/client";
import type { GraduateProfile } from "../types";

export default function Welcome() {
  const navigate = useNavigate();
  const [profile, setProfile] = useState<GraduateProfile | null>(null);

  useEffect(() => {
    api.get<never, GraduateProfile>("/me").then(setProfile).catch(() => navigate("/giris"));
  }, [navigate]);

  if (!profile) {
    return (
      <Center h="100vh">
        <Loader color="meuBlue" />
      </Center>
    );
  }

  return (
    <Container size={560} py={80}>
      <Card withBorder>
        <Stack gap="md" p="md">
          <Text c="meuBlue.8" fw={600} size="sm" tt="uppercase">
            Mersin Üniversitesi · Mezun Takip Anketi
          </Text>
          <Title order={2}>
            Hoş geldin {profile.firstName}, {profile.graduationYear} yılı {profile.departmentName}{" "}
            mezunumuz!
          </Title>
          <Text c="dimmed">
            Seni tekrar aramızda görmek çok güzel. Aşağıdaki kısa anket, mezunlarımızın eğitim ve
            kariyer deneyimlerini anlayarak üniversitemizi geliştirmemize yardımcı oluyor. Anket 5
            bölümden oluşuyor ve yaklaşık 8-10 dakika sürüyor; istediğiniz an bırakıp kaldığınız
            yerden devam edebilirsiniz.
          </Text>
          <Button size="md" onClick={() => navigate("/anket")} color="meuBlue">
            Ankete Başla
          </Button>
        </Stack>
      </Card>
    </Container>
  );
}
