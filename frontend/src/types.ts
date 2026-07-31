export type AnswerType =
  | "scale_1_10"
  | "single_choice"
  | "multi_choice"
  | "text"
  | "number"
  | "duration_months";

export interface SurveyQuestion {
  id: number;
  categoryId: number;
  order: number;
  code: string;
  text: string;
  answerType: AnswerType;
  optionsJson?: string; // JSON.parse edilecek seçenek dizisi (choice tipleri için)
    required: boolean;
    targetFaculty?: string | null;
    targetDepartment?: string | null;
    dependsOnQuestionId?: number | null;
    dependsOnAnswer?: string | null;
}

export interface SurveyCategory {
  id: number;
  order: number;
  title: string;
  questions: SurveyQuestion[];
}

export interface SurveyAnswer {
  id: number;
  responseId: number;
  questionId: number;
  valueText: string;
}

export interface SurveyResponse {
  id: number;
  graduateId: number;
  status: "in_progress" | "completed";
  currentStep: number;
  startedAt: string;
  completedAt?: string | null;
}

export interface GraduateProfile {
  firstName: string;
  facultyName: string;
  departmentName: string;
  graduationYear: number;
}

export interface ApiErrorPayload {
  code: string;
  message: string;
  details?: { field: string; message: string }[];
}

// --- Mezun içe aktarma (OBS import) ---

export interface GraduateImportRow {
  obsHashId: string;
  firstName: string;
  facultyName: string;
  departmentName: string;
  graduationYear: number | "";
  studentNoHash: string;
  email: string;
  phone: string;
}

export interface ImportResult {
  inserted: number;
  updated: number;
  failed: number;
  errors?: string[];
}

export interface ImportResponse {
  result: ImportResult;
  inviteLinks: Record<string, string>;
}
