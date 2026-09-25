import { useEffect, useState } from "react";

import { LogOut } from "lucide-react";
import { useNavigate } from "react-router";

import type { Route } from "./+types/dashboard";
import { BearMark } from "~/components/brand/bear-mark";
import { Button } from "~/components/ui/button";
import { logout, refreshSession } from "~/services/auth";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Ruang kerja · Bearuang" },
    { name: "description", content: "Ruang kerja Bearuang." },
  ];
}

export default function DashboardRoute() {
  const navigate = useNavigate();
  const [isReady, setIsReady] = useState(false);

  useEffect(() => {
    let cancelled = false;

    void refreshSession()
      .then(() => {
        if (!cancelled) setIsReady(true);
      })
      .catch(() => {
        if (!cancelled) navigate("/login", { replace: true });
      });

    return () => {
      cancelled = true;
    };
  }, [navigate]);

  if (!isReady) {
    return (
      <main className="flex min-h-svh items-center justify-center bg-background px-6 text-foreground">
        <p aria-live="polite" className="text-sm text-muted-foreground" role="status">
          Memuat ruang kerja…
        </p>
      </main>
    );
  }

  return (
    <main className="min-h-svh bg-background text-foreground">
      <div className="mx-auto flex min-h-svh w-full max-w-5xl flex-col px-5 py-6 sm:px-10 sm:py-8">
        <header className="flex items-center justify-between border-b border-border pb-6">
          <div className="inline-flex items-center gap-3 text-xl font-semibold tracking-[-0.035em]">
            <span className="flex size-9 items-center justify-center rounded-full border border-border">
              <BearMark aria-hidden="true" className="size-4.5" />
            </span>
            <span>bearuang</span>
          </div>
          <Button
            className="min-h-11 gap-2 px-3"
            onClick={() => {
              void logout()
                .catch(() => undefined)
                .finally(() => navigate("/login", { replace: true }));
            }}
            type="button"
            variant="ghost"
          >
            <LogOut aria-hidden="true" className="size-4" />
            Keluar
          </Button>
        </header>

        <section className="flex flex-1 items-center py-16">
          <div className="max-w-2xl">
            <h1 className="max-w-[16ch] text-4xl font-semibold leading-[1.05] tracking-[-0.035em] text-balance sm:text-6xl sm:leading-[0.98]">
              Ruang kerja Bearuang.
            </h1>
            <p className="mt-5 max-w-[52ch] text-base leading-7 text-muted-foreground sm:text-lg sm:leading-8">
              Anda berhasil masuk. Ruang kerja ini siap menerima modul operasional Anda.
            </p>
          </div>
        </section>
      </div>
    </main>
  );
}
