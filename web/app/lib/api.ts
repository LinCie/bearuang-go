import ky from "ky";

const apiBaseUrl = `${(import.meta.env.VITE_API_URL || "/api").replace(/\/+$/, "")}/`;

const refreshApi = ky.create({
  baseUrl: apiBaseUrl,
  credentials: "include",
  retry: 0,
});

export const api = ky.create({
  baseUrl: apiBaseUrl,
  credentials: "include",
  retry: {
    limit: 1,
  },
  hooks: {
    afterResponse: [
      async ({ request, response, retryCount }) => {
        const isAuthRoute = new URL(request.url).pathname.includes("/auth/");
        const isAuthFailure = response.status === 401 || response.status === 403;

        if (retryCount > 0 || !isAuthFailure || isAuthRoute) {
          return;
        }

        try {
          await refreshApi.post("auth/refresh");
        } catch {
          return;
        }

        return ky.retry();
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
