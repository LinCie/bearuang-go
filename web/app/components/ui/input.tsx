import * as React from "react";

import { cn } from "~/lib/utils";

const Input = React.forwardRef<HTMLInputElement, React.ComponentProps<"input">>(
  ({ className, type, ...props }, ref) => {
    return (
      <input
        ref={ref}
        type={type}
        className={cn(
          "flex h-12 w-full rounded-xl border border-input bg-card px-4 text-base text-foreground shadow-[0_1px_0_color-mix(in_oklab,var(--foreground)_4%,transparent)] outline-none motion-safe:transition-colors placeholder:text-muted-foreground hover:border-input-hover focus:border-primary focus:ring-0 focus-visible:ring-2 focus-visible:ring-foreground focus-visible:ring-offset-1 focus-visible:ring-offset-card disabled:cursor-not-allowed disabled:bg-muted disabled:opacity-70 aria-invalid:border-destructive aria-invalid:ring-2 aria-invalid:ring-destructive/60",
          className,
        )}
        {...props}
      />
    );
  },
);
Input.displayName = "Input";

export { Input };
