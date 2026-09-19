import { useState } from "react";

import { useForm } from "@tanstack/react-form";
import { z } from "zod";
import {
  ArrowRight,
  AlertCircle,
  CheckCircle2,
  Eye,
  EyeOff,
  LoaderCircle,
  LockKeyhole,
  Mail,
} from "lucide-react";
import { Link, useNavigate, useSearchParams } from "react-router";

import { Button } from "~/components/ui/button";
import { Field, FieldContent, FieldError } from "~/components/ui/field";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import {
  getAuthRequestError,
  login,
  register,
  saveTokenPair,
} from "~/lib/auth";

import { AuthShell } from "./auth-shell";

type AuthMode = "login" | "register";

type AuthValues = {
  email: string;
  password: string;
  confirmPassword: string;
};

const emailSchema = z
  .string()
  .trim()
  .min(1, "Email wajib diisi.")
  .email("Masukkan alamat email yang valid.");
const passwordSchema = z.string().min(1, "Kata sandi wajib diisi.");
const confirmPasswordSchema = z.string().min(1, "Konfirmasi kata sandi wajib diisi.");
const loginSchema = z.object({
  email: emailSchema,
  password: passwordSchema,
  confirmPassword: z.string(),
});
const registerSchema = z
  .object({
    email: emailSchema,
    password: passwordSchema,
    confirmPassword: confirmPasswordSchema,
  })
  .refine((values) => values.password === values.confirmPassword, {
    message: "Kata sandi belum sama.",
    path: ["confirmPassword"],
  });

function getFieldError(errors: unknown[]) {
  const [firstError] = errors;

  if (typeof firstError === "string") return firstError;
  if (Array.isArray(firstError)) return getFieldError(firstError);

  if (firstError && typeof firstError === "object" && "message" in firstError) {
    const message = firstError.message;
    return typeof message === "string" ? message : undefined;
  }

  return undefined;
}

function AuthForm({ mode }: { mode: AuthMode }) {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [requestError, setRequestError] = useState<string | null>(null);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const isRegister = mode === "register";
  const hasRegistrationSuccess =
    !isRegister && searchParams.get("registered") === "1";
  const submitSchema = isRegister ? registerSchema : loginSchema;

  const form = useForm({
    defaultValues: {
      email: "",
      password: "",
      confirmPassword: "",
    },
    validators: {
      onSubmit: submitSchema,
    },
    onSubmit: async ({ value }) => {
      setRequestError(null);

      try {
        if (isRegister) {
          await register({
            email: value.email.trim(),
            password: value.password,
          });
          navigate("/login?registered=1", { replace: true });
          return;
        }

        const tokens = await login({
          email: value.email.trim(),
          password: value.password,
        });
        if (!saveTokenPair(tokens)) {
          setRequestError(
            "Sesi tidak dapat disimpan di peramban ini. Izinkan penyimpanan lalu coba lagi.",
          );
          return;
        }

        navigate("/dashboard");
      } catch (error) {
        const request = await getAuthRequestError(error);
        const message =
          request.code === "invalid_credentials"
            ? "Email atau kata sandi tidak cocok. Periksa kembali lalu coba lagi."
            : request.code === "email_exists"
              ? "Email ini sudah terdaftar. Coba masuk atau gunakan email lain."
              : request.message;

        setRequestError(message);
      }
    },
  });

  const title = isRegister ? "Buat akun Bearuang" : "Masuk ke Bearuang";
  const description = isRegister
    ? "Mulai dengan email dan kata sandi. Setelah akun dibuat, Anda bisa masuk ke ruang kerja dengan kredensial yang sama."
    : "Masukkan kredensial Anda untuk melanjutkan ke ruang kerja bisnis.";

  return (
    <AuthShell title={title} description={description}>
      <form
        aria-describedby="auth-description"
        aria-labelledby="auth-title"
        noValidate
        className="space-y-6"
        onSubmit={(event) => {
          event.preventDefault();
          void form.handleSubmit();
        }}
      >
        {requestError ? (
          <div
            aria-live="assertive"
            className="flex gap-3 rounded-xl border border-destructive/30 bg-destructive/10 px-4 py-3.5 text-sm leading-5 text-destructive"
            role="alert"
          >
            <AlertCircle aria-hidden="true" className="mt-0.5 size-4 shrink-0" />
            <p>{requestError}</p>
          </div>
        ) : null}

        {hasRegistrationSuccess ? (
          <div
            aria-live="polite"
            className="flex gap-3 rounded-xl border border-success/30 bg-success/10 px-4 py-3.5 text-sm leading-5 text-success"
            role="status"
          >
            <CheckCircle2 aria-hidden="true" className="mt-0.5 size-4 shrink-0" />
            <p>Akun berhasil dibuat. Masuk dengan email dan kata sandi Anda.</p>
          </div>
        ) : null}

        <div className="space-y-5">
          <form.Field
            name="email"
            validators={{
              onBlur: emailSchema,
              onSubmit: emailSchema,
            }}
          >
            {(field) => {
              const error = getFieldError(field.state.meta.errors);
              const fieldId = "email";

              return (
                <Field>
                  <FieldContent>
                    <Label htmlFor={fieldId}>Email</Label>
                    <div className="relative">
                      <Mail
                        aria-hidden="true"
                        className="pointer-events-none absolute left-4 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
                        strokeWidth={1.8}
                      />
                      <Input
                        aria-describedby={error ? `${fieldId}-error` : undefined}
                        aria-invalid={Boolean(error)}
                        autoComplete="email"
                        className="pl-11"
                        id={fieldId}
                        inputMode="email"
                        name={field.name}
                        onBlur={() => field.handleBlur()}
                        onChange={(event) => field.handleChange(event.target.value)}
                        placeholder="nama@perusahaan.id"
                        required
                        type="email"
                        value={field.state.value}
                      />
                    </div>
                    <FieldError id={`${fieldId}-error`}>{error}</FieldError>
                  </FieldContent>
                </Field>
              );
            }}
          </form.Field>

          <form.Field
            name="password"
            validators={{
              onBlur: passwordSchema,
              onSubmit: passwordSchema,
            }}
          >
            {(field) => {
              const error = getFieldError(field.state.meta.errors);
              const fieldId = "password";

              return (
                <Field>
                  <FieldContent>
                    <Label htmlFor={fieldId}>Kata sandi</Label>
                    <div className="relative">
                      <LockKeyhole
                        aria-hidden="true"
                        className="pointer-events-none absolute left-4 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
                        strokeWidth={1.8}
                      />
                      <Input
                        aria-describedby={error ? `${fieldId}-error` : undefined}
                        aria-invalid={Boolean(error)}
                        autoComplete={isRegister ? "new-password" : "current-password"}
                        className="pl-11 pr-14"
                        id={fieldId}
                        name={field.name}
                        onBlur={() => field.handleBlur()}
                        onChange={(event) => field.handleChange(event.target.value)}
                        placeholder="Masukkan kata sandi"
                        required
                        type={showPassword ? "text" : "password"}
                        value={field.state.value}
                      />
                      <button
                        aria-label={showPassword ? "Sembunyikan kata sandi" : "Tampilkan kata sandi"}
                        aria-pressed={showPassword}
                        className="absolute right-1.5 top-1/2 flex size-11 -translate-y-1/2 items-center justify-center rounded-lg text-muted-foreground motion-safe:transition-colors hover:text-foreground focus-visible:ring-2 focus-visible:ring-foreground focus-visible:ring-offset-1 focus-visible:ring-offset-card"
                        onClick={() => setShowPassword((prev) => !prev)}
                        type="button"
                      >
                        {showPassword ? (
                          <EyeOff aria-hidden="true" className="size-4" />
                        ) : (
                          <Eye aria-hidden="true" className="size-4" />
                        )}
                      </button>
                    </div>
                    <FieldError id={`${fieldId}-error`}>{error}</FieldError>
                  </FieldContent>
                </Field>
              );
            }}
          </form.Field>

          {isRegister ? (
            <form.Field
              name="confirmPassword"
              validators={{
                onChangeListenTo: ["password"],
                onChange: ({ value, fieldApi }) => {
                  if (!value) return undefined;
                  const parsed = confirmPasswordSchema.safeParse(value);
                  if (!parsed.success) return parsed.error.issues[0]?.message;
                  if (value !== fieldApi.form.state.values.password) {
                    return "Kata sandi belum sama.";
                  }
                  return undefined;
                },
                onBlur: ({ value, fieldApi }) => {
                  if (!value) return undefined;
                  const parsed = confirmPasswordSchema.safeParse(value);
                  if (!parsed.success) return parsed.error.issues[0]?.message;
                  if (value !== fieldApi.form.state.values.password) {
                    return "Kata sandi belum sama.";
                  }
                  return undefined;
                },
                onSubmit: ({ value, fieldApi }) => {
                  const parsed = confirmPasswordSchema.safeParse(value);
                  if (!parsed.success) return parsed.error.issues[0]?.message;
                  if (value !== fieldApi.form.state.values.password) {
                    return "Kata sandi belum sama.";
                  }
                  return undefined;
                },
              }}
            >
              {(field) => {
                const error = getFieldError(field.state.meta.errors);
                const fieldId = "confirmPassword";

                return (
                  <Field>
                    <FieldContent>
                      <Label htmlFor={fieldId}>Konfirmasi kata sandi</Label>
                      <div className="relative">
                        <LockKeyhole
                          aria-hidden="true"
                          className="pointer-events-none absolute left-4 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
                          strokeWidth={1.8}
                        />
                        <Input
                          aria-describedby={error ? `${fieldId}-error` : undefined}
                          aria-invalid={Boolean(error)}
                          autoComplete="new-password"
                          className="pl-11 pr-14"
                          id={fieldId}
                          name={field.name}
                          onBlur={() => field.handleBlur()}
                          onChange={(event) => field.handleChange(event.target.value)}
                          placeholder="Ulangi kata sandi"
                          required
                          type={showConfirmPassword ? "text" : "password"}
                          value={field.state.value}
                        />
                        <button
                          aria-label={showConfirmPassword ? "Sembunyikan kata sandi" : "Tampilkan kata sandi"}
                          aria-pressed={showConfirmPassword}
                          className="absolute right-1.5 top-1/2 flex size-11 -translate-y-1/2 items-center justify-center rounded-lg text-muted-foreground motion-safe:transition-colors hover:text-foreground focus-visible:ring-2 focus-visible:ring-foreground focus-visible:ring-offset-1 focus-visible:ring-offset-card"
                          onClick={() => setShowConfirmPassword((prev) => !prev)}
                          type="button"
                        >
                          {showConfirmPassword ? (
                            <EyeOff aria-hidden="true" className="size-4" />
                          ) : (
                            <Eye aria-hidden="true" className="size-4" />
                          )}
                        </button>
                      </div>
                      <FieldError id={`${fieldId}-error`}>{error}</FieldError>
                    </FieldContent>
                  </Field>
                );
              }}
            </form.Field>
          ) : null}
        </div>

        <form.Subscribe selector={(state) => state.isSubmitting}>
          {(isSubmitting) => (
            <Button
              className="h-12 w-full gap-2 rounded-xl bg-primary text-base text-primary-foreground shadow-auth-action hover:bg-primary hover:brightness-95 focus-visible:ring-2 focus-visible:ring-foreground focus-visible:ring-offset-2 focus-visible:ring-offset-background"
              disabled={isSubmitting}
              size="lg"
              type="submit"
            >
              {isSubmitting ? (
                <>
                  <LoaderCircle aria-hidden="true" className="size-4 motion-safe:animate-spin" />
                  Memproses…
                </>
              ) : (
                <>
                  {isRegister ? "Buat akun" : "Masuk"}
                  <ArrowRight aria-hidden="true" className="size-4" />
                </>
              )}
            </Button>
          )}
        </form.Subscribe>

        <p className="text-center text-sm leading-6 text-muted-foreground">
          {isRegister ? "Sudah punya akun?" : "Belum punya akun?"}{" "}
          <Link
            className="-mx-1 px-1 font-semibold text-primary underline decoration-accent underline-offset-4 motion-safe:transition-colors hover:text-foreground focus-visible:ring-2 focus-visible:ring-foreground focus-visible:ring-offset-2 focus-visible:ring-offset-background"
            to={isRegister ? "/login" : "/register"}
          >
            {isRegister ? "Masuk" : "Daftar"}
          </Link>
        </p>
      </form>
    </AuthShell>
  );
}

export { AuthForm };
