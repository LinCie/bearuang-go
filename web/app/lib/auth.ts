import { HTTPError } from "ky";

import {
  api,
  AUTH_STORAGE_KEY,
  type ApiErrorResponse,
  type ApiResponse,
} from "~/lib/api";

export type Credentials = {
  email: string;
  password: string;
};

export type TokenPair = {
  access_token: string;
  refresh_token: string;
};

export type AuthRequestError = {
  code: string;
  message: string;
};

export async function login(credentials: Credentials) {
  const response = await api
    .post("auth/login", { json: credentials })
    .json<ApiResponse<TokenPair>>();

  return response.data;
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

export function getStoredTokenPair(): TokenPair | null {
  if (typeof window === "undefined") return null;

  try {
    const stored = window.sessionStorage.getItem(AUTH_STORAGE_KEY);
    if (!stored) return null;

    const parsed = JSON.parse(stored) as Partial<TokenPair>;
    if (
      typeof parsed.access_token !== "string" ||
      !parsed.access_token ||
      typeof parsed.refresh_token !== "string" ||
      !parsed.refresh_token
    ) {
      return null;
    }

    return {
      access_token: parsed.access_token,
      refresh_token: parsed.refresh_token,
    };
  } catch {
    return null;
  }
}

export function saveTokenPair(tokens: TokenPair) {
  try {
    window.sessionStorage.setItem(AUTH_STORAGE_KEY, JSON.stringify(tokens));
    return true;
  } catch {
    return false;
  }
}

export function clearTokenPair() {
  try {
    window.sessionStorage.removeItem(AUTH_STORAGE_KEY);
  } catch {
    // Storage may be unavailable in a restricted browser context.
  }
}
