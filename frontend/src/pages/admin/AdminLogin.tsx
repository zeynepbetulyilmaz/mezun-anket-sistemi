import { useState } from "react";
import type { FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import { Container, Card, Stack, Title, TextInput, PasswordInput, Button, Text } from "@mantine/core";

import { api, adminTokenStorage } from "../../api/client";
import type { StandardizedError } from "../../api/client";

export default function AdminLogin() {
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      const data = await api.post<never, { accessToken: string; role: string }>("/admin/login", {
        username,
        password,
      });
      adminTokenStorage.set(data.accessToken);
      navigate("/admin");
    } catch (err) {
      setError((err as StandardizedError).message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <Container size={400} py={100}>
      <Card withBorder>
        <form onSubmit={handleSubmit}>
          <Stack gap="md" p="md">
            <Title order={3}>Yönetim Paneli Girişi</Title>
            <TextInput
              label="Kullanıcı adı"
              value={username}
              onChange={(e) => setUsername(e.currentTarget.value)}
              required
            />
            <PasswordInput
              label="Şifre"
              value={password}
              onChange={(e) => setPassword(e.currentTarget.value)}
              required
            />
            {error && (
              <Text c="red" size="sm">
                {error}
              </Text>
            )}
            <Button type="submit" loading={loading} color="meuBlue">
              Giriş Yap
            </Button>
          </Stack>
        </form>
      </Card>
    </Container>
  );
}
