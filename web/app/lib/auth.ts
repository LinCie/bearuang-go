import { HTTPError } from "ky";

import { api, type ApiErrorResponse, type ApiResponse } from "~/lib/api";

export type Credentials = {
  email: string;
  password: string;
};

export type AuthRequestError = {
  code: string;
  message: string;
};

export async function login(credentials: Credentials) {
  await api.post("auth/login", { json: credentials });
}

export async function register(credentials: Credentials) {
  const response = await api
    .post("auth/register", { json: credentials })
    .json<ApiResponse<{ id: string; email: string }>>();

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
  };
}

export async function refreshSession() {
  await api.post("auth/refresh");
}

export async function logout() {
  await api.post("auth/logout");
}
