import type { ComponentProps, ReactNode } from "react";

import { cn } from "~/lib/utils";

function Field({ className, ...props }: ComponentProps<"div">) {
  return <div className={cn("space-y-2", className)} {...props} />;
}

function FieldContent({ className, ...props }: ComponentProps<"div">) {
  return <div className={cn("space-y-1.5", className)} {...props} />;
}

function FieldError({ children, className, ...props }: { children?: ReactNode } & ComponentProps<"p">) {
  if (!children) return null;

  return (
    <p
      role="alert"
      className={cn("text-sm leading-5 text-destructive", className)}
      {...props}
    >
      {children}
    </p>
  );
}

export { Field, FieldContent, FieldError };
