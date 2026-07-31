import { Container, Card, Stack, Title, Text, ThemeIcon, Center } from "@mantine/core";

export default function ThankYou() {
  return (
    <Container size={480} py={100}>
      <Card withBorder>
        <Stack align="center" gap="md" p="md">
          <ThemeIcon size={64} radius="xl" color="meuBlue" variant="light">
            <CheckIcon />
          </ThemeIcon>
          <Title order={2} ta="center">
            Teşekkür ederiz!
          </Title>
          <Text c="dimmed" ta="center">
            Anketi tamamladınız. Cevaplarınız üniversitemizin eğitim kalitesini ve mezunlarımızla
            olan bağını güçlendirmesine katkı sağlayacak. Kısa süre içinde e-posta adresinize bir
            teşekkür mesajı ulaşacak.
          </Text>
        </Stack>
      </Card>
    </Container>
  );
}

function CheckIcon() {
  return (
    <Center>
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none">
        <path
          d="M5 13l4 4L19 7"
          stroke="currentColor"
          strokeWidth="2.5"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
      </svg>
    </Center>
  );
}
