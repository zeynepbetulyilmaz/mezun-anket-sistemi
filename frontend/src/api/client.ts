import axios, { type AxiosError } from "axios";
import type { ApiErrorPayload } from "../types";

const TOKEN_KEY = "meu_survey_access_token";
const ADMIN_TOKEN_KEY = "meu_survey_admin_token";

export const tokenStorage = {
  get: () => localStorage.getItem(TOKEN_KEY),
  set: (t: string) => localStorage.setItem(TOKEN_KEY, t),
  clear: () => localStorage.removeItem(TOKEN_KEY),
};

export const adminTokenStorage = {
  get: () => localStorage.getItem(ADMIN_TOKEN_KEY),
  set: (t: string) => localStorage.setItem(ADMIN_TOKEN_KEY, t),
  clear: () => localStorage.removeItem(ADMIN_TOKEN_KEY),
};

export const api = axios.create({ baseURL: "/api/v1" });

api.interceptors.request.use((config) => {
  const isAdminRoute = config.url?.startsWith("/admin");
  const token = isAdminRoute ? adminTokenStorage.get() : tokenStorage.get();
  if (token) {
    config.headers = config.headers ?? {};
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});


export interface StandardizedError extends Error {
  code: string;
  details?: ApiErrorPayload["details"];
}

api.interceptors.response.use(
  (response) => response.data?.data ?? response.data,
  (error: AxiosError<{ error?: ApiErrorPayload }>) => {
    const status = error.response?.status;
    const payload = error.response?.data?.error;

    const err = new Error(payload?.message ?? "Beklenmeyen bir hata oluştu.") as StandardizedError;
    err.code = payload?.code ?? "UNKNOWN_ERROR";
    err.details = payload?.details;

    switch (status) {
      case 401:

        tokenStorage.clear();
        adminTokenStorage.clear();
        if (!window.location.pathname.startsWith("/admin")) {
          window.location.href = "/giris?expired=1";
        } else {
          window.location.href = "/admin/giris?expired=1";
        }
        break;
      case 403:

        break;
      case 400:
      case 404:
      case 409:
      case 500:
      default:
        break;
    }
    return Promise.reject(err);
  }
);
