import {
  api,
  apiResult,
  type ApiResult,
  type CommonApiErrorCode,
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

export type AuthErrorCode =
  | CommonApiErrorCode
  | "invalid_email"
  | "invalid_password"
  | "email_exists"
  | "invalid_credentials"
  | "invalid_refresh_token";

export function register(
  credentials: Credentials,
): Promise<ApiResult<AuthUser, AuthErrorCode>> {
  return apiResult<AuthUser, AuthErrorCode>(
    api.post("auth/register", { json: credentials }),
  );
}

export function login(
  credentials: Credentials,
): Promise<ApiResult<AuthSession, AuthErrorCode>> {
  return apiResult<AuthSession, AuthErrorCode>(
    api.post("auth/login", { json: credentials }),
  );
}

export function refreshSession(): Promise<ApiResult<AuthSession, AuthErrorCode>> {
  return apiResult<AuthSession, AuthErrorCode>(api.post("auth/refresh"));
}

export function logout(): Promise<ApiResult<AuthSession, AuthErrorCode>> {
  return apiResult<AuthSession, AuthErrorCode>(api.post("auth/logout"));
}
