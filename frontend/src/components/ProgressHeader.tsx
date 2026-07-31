import { Progress, Group, Text, Stack } from "@mantine/core";

interface Props {
  answeredCount: number;
  totalCount: number;
}

export default function ProgressHeader({ answeredCount, totalCount }: Props) {
  const pct = totalCount === 0 ? 0 : Math.round((answeredCount / totalCount) * 100);
  return (
    <Stack gap={4} mb="lg">
      <Group justify="space-between">
        <Text size="sm" c="dimmed">
          İlerleme
        </Text>
        <Text size="sm" c="dimmed">
          {answeredCount} / {totalCount} soru
        </Text>
      </Group>
      <Progress value={pct} color="meuBlue" size="lg" radius="xl" />
    </Stack>
  );
}
