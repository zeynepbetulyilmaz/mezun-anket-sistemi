import { SegmentedControl, Radio, Group, Textarea, NumberInput, Stack, Text } from "@mantine/core";
import type { SurveyQuestion } from "../types";

interface Props {
  question: SurveyQuestion;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export default function QuestionField({ question, value, onChange, error }: Props) {
  const options: string[] = question.optionsJson ? JSON.parse(question.optionsJson) : [];

  return (
    <Stack gap={6}>
      <Text fw={500}>
        {question.text}
        {question.required && (
          <Text component="span" c="red" ml={4}>
            *
          </Text>
        )}
      </Text>

      {question.answerType === "scale_1_10" && (
        <SegmentedControl
          fullWidth
          value={value}
          onChange={onChange}
          data={Array.from({ length: 10 }, (_, i) => String(i + 1))}
          color="meuBlue"
        />
      )}

      {question.answerType === "single_choice" && (
        <Radio.Group value={value} onChange={onChange}>
          <Group mt={4}>
            {options.map((opt) => (
              <Radio key={opt} value={opt} label={opt} color="meuBlue" />
            ))}
          </Group>
        </Radio.Group>
      )}

      {(question.answerType === "text") && (
        <Textarea
          value={value}
          onChange={(e) => onChange(e.currentTarget.value)}
          autosize
          minRows={2}
          placeholder="Cevabınızı yazın..."
        />
      )}

      {(question.answerType === "number" || question.answerType === "duration_months") && (
        <NumberInput
          value={value === "" ? "" : Number(value)}
          onChange={(v) => onChange(v === "" ? "" : String(v))}
          min={0}
          placeholder={question.answerType === "duration_months" ? "Ay cinsinden" : undefined}
        />
      )}

      {error && (
        <Text c="red" size="sm">
          {error}
        </Text>
      )}
    </Stack>
  );
}
