import type { ReactNode } from "react";

import { BearMark } from "~/components/brand/bear-mark";

type AuthShellProps = {
  children: ReactNode;
  companion?: ReactNode;
  description: ReactNode;
  title: ReactNode;
};

function AuthShell({ children, companion, description, title }: AuthShellProps) {
  return (
    <main className="min-h-svh overflow-x-hidden bg-background text-foreground selection:bg-accent selection:text-foreground">
      <div className="grid min-h-svh w-full lg:grid-cols-[minmax(420px,0.42fr)_minmax(0,0.58fr)]">
        <aside className="relative flex min-w-0 flex-col overflow-hidden bg-auth-panel px-6 py-6 text-auth-panel-foreground sm:px-10 sm:py-8 lg:min-h-svh lg:px-14 lg:py-12">
          <div className="relative z-10 flex items-center justify-between">
            <div className="group inline-flex items-center gap-3 text-xl font-semibold tracking-[-0.035em]">
              <span className="flex size-9 items-center justify-center rounded-full border border-auth-panel-foreground/55 motion-safe:transition-transform motion-safe:group-hover:-rotate-6">
                <BearMark aria-hidden="true" className="size-4.5" />
              </span>
              <span>bearuang</span>
            </div>
            <span className="text-xs font-semibold tracking-[0.16em] text-auth-panel-foreground uppercase">
              ruang kerja
            </span>
          </div>

          <div className="relative z-10 mt-4 max-w-[30rem] pt-2 sm:mt-6 sm:pt-4 lg:mt-auto lg:pb-10 lg:pt-32">
            <p className="max-w-[12ch] text-2xl font-semibold leading-[1.1] tracking-[-0.035em] text-balance sm:text-3xl lg:text-[clamp(2.75rem,5.4vw,5.1rem)] lg:leading-[0.94]">
              Ruang kerja yang tetap jernih.
            </p>
            <p className="mt-2 hidden text-sm leading-6 text-auth-panel-foreground sm:mt-4 sm:block sm:text-base sm:leading-7 lg:mt-7 lg:text-lg lg:leading-8">
              Bearuang membantu pemilik usaha menata pekerjaan penting dengan cara yang sederhana dan dapat diandalkan.
            </p>
          </div>

          <div className="relative z-10 mt-4 hidden border-t border-auth-panel-foreground/30 pt-4 text-xs font-medium leading-5 text-auth-panel-foreground sm:block lg:mt-auto">
            <span>ERP sederhana untuk usaha yang bertumbuh.</span>
          </div>

          <div aria-hidden="true" className="pointer-events-none absolute inset-x-0 bottom-24 hidden opacity-15 lg:block">
            <div className="h-px bg-auth-panel-foreground" />
            <div className="mt-5 h-px bg-auth-panel-foreground" />
            <div className="mt-5 h-px bg-auth-panel-foreground" />
          </div>
          <div aria-hidden="true" className="pointer-events-none absolute -right-20 top-1/2 size-72 -translate-y-1/2 rounded-full border border-auth-panel-foreground/15" />
          <div aria-hidden="true" className="pointer-events-none absolute -right-10 top-1/2 size-52 -translate-y-1/2 rounded-full border border-auth-panel-foreground/15" />
        </aside>

        <section aria-describedby="auth-description" aria-labelledby="auth-title" className="flex min-w-0 items-center bg-background px-5 py-8 sm:px-10 sm:py-12 lg:min-h-svh lg:border-l lg:border-border lg:px-16 xl:px-24">
          <div className="mx-auto w-full max-w-[31rem]">
            {companion ? (
              <div className="mb-6 flex justify-center sm:mb-8">
                {companion}
              </div>
            ) : null}
            <div className="mb-6 sm:mb-8">
              <h1 id="auth-title" className="text-3xl font-semibold leading-[1.08] tracking-[-0.035em] text-foreground sm:text-5xl sm:leading-[1.04]">
                {title}
              </h1>
              <p id="auth-description" className="mt-3 max-w-[38ch] text-sm leading-6 text-muted-foreground sm:mt-4 sm:text-base sm:leading-7">
                {description}
              </p>
              <div aria-hidden="true" className="mt-6 flex items-center sm:mt-8">
                <div className="h-0.5 w-10 rounded-full bg-primary/75" />
                <div className="h-px flex-1 bg-border" />
              </div>
            </div>
            {children}
          </div>
        </section>
      </div>
    </main>
  );
}

export { AuthShell };
