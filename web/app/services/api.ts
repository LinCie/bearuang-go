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
  throwHttpErrors: false,
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

export type ApiFieldError = {
  field: string;
  rule: string;
  param?: string;
  message: string;
};

export type ApiErrorResponse<C extends string = string> = {
  code: C;
  message: string;
  fields?: ApiFieldError[];
};

export type ApiResult<T, C extends string = string> =
  | {
      ok: true;
      data: T;
    }
  | {
      ok: false;
      status: number;
      error: ApiErrorResponse<C>;
    };

export type CommonApiErrorCode = "invalid_body" | "internal_error";

export class ApiProtocolError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "ApiProtocolError";
  }
}

export function isApiProtocolError(error: unknown): error is ApiProtocolError {
  return error instanceof ApiProtocolError;
}

export async function apiResult<T, C extends string = string>(
  responsePromise: Promise<Response>,
): Promise<ApiResult<T, C>> {
  const response = await responsePromise;
  let body: unknown;
  try {
    body = await response.json();
  } catch {
    throw new ApiProtocolError("Invalid API response: expected valid JSON");
  }

  if (response.ok) {
    if (!isRecord(body) || !("data" in body)) {
      throw new ApiProtocolError("Invalid API success response: expected a data envelope");
    }

    const envelope = body as ApiResponse<T>;
    return { ok: true, data: envelope.data };
  }

  if (!isApiErrorResponse(body)) {
    throw new ApiProtocolError("Invalid API error response: expected code and message fields");
  }

  const error = body as ApiErrorResponse<C>;
  return { ok: false, status: response.status, error };
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function isApiErrorResponse(value: unknown): value is Record<string, unknown> {
  if (
    !isRecord(value) ||
    typeof value.code !== "string" ||
    typeof value.message !== "string"
  ) {
    return false;
  }

  if (value.fields === undefined) {
    return true;
  }

  return (
    Array.isArray(value.fields) &&
    value.fields.every(
      (field) =>
        isRecord(field) &&
        typeof field.field === "string" &&
        typeof field.rule === "string" &&
        typeof field.message === "string" &&
        (field.param === undefined || typeof field.param === "string"),
    )
  );
}
