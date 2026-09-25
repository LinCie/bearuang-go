import { HTTPError } from "ky";

import {
  api,
  type ApiErrorResponse,
  type ApiFieldError,
  type ApiResponse,
} from "~/services/api";

export type Credentials = {
  email: string;
  password: string;
};

export type AuthUser = {
  id: string;
  email: string;
  created_at: string;
  updated_at: string;
};

export type AuthSession = {
  authenticated: boolean;
};

export type AuthRequestError = {
  code: string;
  message: string;
  fields?: ApiFieldError[];
};

export async function register(credentials: Credentials): Promise<AuthUser> {
  const response = await api
    .post("auth/register", { json: credentials })
    .json<ApiResponse<AuthUser>>();

  return response.data;
}

export async function login(credentials: Credentials): Promise<AuthSession> {
  const response = await api
    .post("auth/login", { json: credentials })
    .json<ApiResponse<AuthSession>>();

  return response.data;
}

export async function refreshSession(): Promise<AuthSession> {
  const response = await api.post("auth/refresh").json<ApiResponse<AuthSession>>();

  return response.data;
}

export async function logout(): Promise<AuthSession> {
  const response = await api.post("auth/logout").json<ApiResponse<AuthSession>>();

  return response.data;
}

export async function getAuthRequestError(error: unknown): Promise<AuthRequestError> {
  if (!(error instanceof HTTPError)) {
    return {
      code: "network_error",
      message: "Tidak dapat terhubung ke Bearuang. Periksa koneksi Anda lalu coba lagi.",
    };
  }

  const response = (await error.response.clone().json().catch(() => null)) as ApiErrorResponse | null;

  return {
    code: response?.code ?? "request_failed",
    message: response?.message ?? "Permintaan belum dapat diproses. Coba lagi.",
    fields: response?.fields,
  };
}
