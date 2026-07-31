import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { Center, Stack, Loader, Text, Alert, Container, Image } from "@mantine/core";

import { api, tokenStorage } from "../api/client";
import type { StandardizedError } from "../api/client";

export default function TokenLogin() {
  const [params] = useSearchParams();
  const navigate = useNavigate();
  const [error, setError] = useState<string | null>(null);
  const expired = params.get("expired") === "1";

  useEffect(() => {
    const token = params.get("token");
    if (!token) {
      if (!expired) setError("Giriş linki eksik. Lütfen size gönderilen e-postadaki bağlantıyı kullanın.");
      return;
    }

    api
      .post("/auth/token-login", { token })
      .then((data: any) => {
        tokenStorage.set(data.accessToken);
        sessionStorage.setItem("meu_graduate_profile", JSON.stringify(data.graduate));
        navigate("/hosgeldin", { replace: true });
      })
      .catch((err: StandardizedError) => {
        setError(err.message);
      });
  }, [params, navigate, expired]);

  return (
    <Container size={420} py={80}>
      <Center mb="lg">
        {/* Kurum logosu için: /public/meu-logo.png yerleştirip buraya ekleyin */}
        <Text fw={700} size="xl" c="meuBlue.8">
          Mersin Üniversitesi
        </Text>
      </Center>
      <Center>
        <Stack align="center" gap="md">
          {!error && (
            <>
              <Loader color="meuBlue" />
              <Text c="dimmed">Giriş doğrulanıyor...</Text>
            </>
          )}
          {expired && !error && (
            <Alert color="yellow" title="Oturum süresi doldu">
              Lütfen size gönderilen son e-postadaki linki tekrar kullanın.
            </Alert>
          )}
          {error && (
            <Alert color="red" title="Giriş yapılamadı">
              {error}
            </Alert>
          )}
        </Stack>
      </Center>
    </Container>
  );
}
