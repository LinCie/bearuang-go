import type { Route } from "./+types/login";

import { AuthForm } from "~/components/auth/auth-form";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Masuk · Bearuang" },
    { name: "description", content: "Masuk ke ruang kerja Bearuang." },
  ];
}

export default function LoginRoute() {
  return <AuthForm mode="login" />;
}
