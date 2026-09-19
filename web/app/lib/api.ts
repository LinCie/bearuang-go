import ky from "ky";

const apiBaseUrl = (import.meta.env.VITE_API_URL || "/api").replace(/\/+$/, "");

export const AUTH_STORAGE_KEY = "bearuang.auth";

function getStoredAccessToken() {
  if (typeof window === "undefined") return null;

  try {
    const stored = window.sessionStorage.getItem(AUTH_STORAGE_KEY);
    if (!stored) return null;

    const parsed = JSON.parse(stored) as { access_token?: unknown };
    return typeof parsed.access_token === "string" && parsed.access_token
      ? parsed.access_token
      : null;
  } catch {
    return null;
  }
}

export const api = ky.create({
  baseUrl: apiBaseUrl,
  hooks: {
    beforeRequest: [
      ({ request }) => {
        if (request.url.includes("/auth/")) return;

        const accessToken = getStoredAccessToken();
        if (accessToken) {
          request.headers.set("Authorization", `Bearer ${accessToken}`);
        }
      },
    ],
  },
});

export type ApiResponse<T> = {
  data: T;
};

export type ApiErrorResponse = {
  code?: string;
  message?: string;
};
