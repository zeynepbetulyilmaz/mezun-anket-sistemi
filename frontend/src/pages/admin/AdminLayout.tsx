import { AppShell, Group, Text, Button, Burger, Collapse, Stack } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { useLocation, useNavigate, Outlet } from "react-router-dom";

import { adminTokenStorage } from "../../api/client";

const NAV_ITEMS = [
  { label: "Panel", path: "/admin" },
  { label: "Mezun Ekle", path: "/admin/mezun-ekle" },
];

export default function AdminLayout() {
  const [opened, { toggle, close }] = useDisclosure();
  const navigate = useNavigate();
  const location = useLocation();

  function handleLogout() {
    adminTokenStorage.clear();
    navigate("/admin/giris");
  }

  function goTo(path: string) {
    navigate(path);
    close();
  }

  return (
    <AppShell header={{ height: 60 }} padding="md">
      <AppShell.Header>
        <Group h={60} px="md" justify="space-between">
          <Group gap="sm">
            <Burger opened={opened} onClick={toggle} hiddenFrom="sm" size="sm" />
            <Text fw={700} c="meuBlue.8">
              Mersin Üniversitesi · Mezun Anketi Yönetimi
            </Text>
          </Group>

          <Group gap="xs" visibleFrom="sm">
            {NAV_ITEMS.map((item) => (
              <Button
                key={item.path}
                variant={location.pathname === item.path ? "filled" : "subtle"}
                color="meuBlue"
                size="sm"
                onClick={() => navigate(item.path)}
              >
                {item.label}
              </Button>
            ))}
            <Button variant="default" size="sm" onClick={handleLogout}>
              Çıkış Yap
            </Button>
          </Group>
        </Group>

        <Collapse in={opened} hiddenFrom="sm">
          <Stack gap={4} p="sm" pt={0}>
            {NAV_ITEMS.map((item) => (
              <Button
                key={item.path}
                variant={location.pathname === item.path ? "filled" : "subtle"}
                color="meuBlue"
                size="sm"
                fullWidth
                onClick={() => goTo(item.path)}
              >
                {item.label}
              </Button>
            ))}
            <Button variant="default" size="sm" fullWidth onClick={handleLogout}>
              Çıkış Yap
            </Button>
          </Stack>
        </Collapse>
      </AppShell.Header>

      <AppShell.Main bg="gray.0">
        <Outlet />
      </AppShell.Main>
    </AppShell>
  );
}
