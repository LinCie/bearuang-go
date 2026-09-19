import type { Route } from "./+types/register";

import { AuthForm } from "~/components/auth/auth-form";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Daftar · Bearuang" },
    { name: "description", content: "Buat akun Bearuang." },
  ];
}

export default function RegisterRoute() {
  return <AuthForm mode="register" />;
}
