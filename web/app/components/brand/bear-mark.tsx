import type { ComponentProps } from "react";

export function BearMark({
  className = "size-4",
  ...props
}: ComponentProps<"svg">) {
  return (
    <svg
      aria-hidden="true"
      className={className}
      fill="none"
      stroke="currentColor"
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={1.75}
      viewBox="0 0 24 24"
      {...props}
    >
      {/* Left ear */}
      <path d="M 6.5 8.5 C 4.5 8.5 3.2 7 3.2 5.2 C 3.2 3.4 4.8 2.2 6.6 2.8 C 7.6 3.2 8.3 4.2 8.6 5.3" />
      {/* Right ear */}
      <path d="M 17.5 8.5 C 19.5 8.5 20.8 7 20.8 5.2 C 20.8 3.4 19.2 2.2 17.4 2.8 C 16.4 3.2 15.7 4.2 15.4 5.3" />
      {/* Head contour */}
      <path d="M 6.5 8.5 C 5 11 5 14.2 6.5 16.8 C 8 19.2 9.8 20.2 12 20.2 C 14.2 20.2 16 19.2 17.5 16.8 C 19 14.2 19 11 17.5 8.5 C 16 9 14.2 8.8 12 8.8 C 9.8 8.8 8 9 6.5 8.5 Z" />
      {/* Eyes */}
      <circle cx="9" cy="12.5" r="0.9" fill="currentColor" stroke="none" />
      <circle cx="15" cy="12.5" r="0.9" fill="currentColor" stroke="none" />
      {/* Snout / Nose */}
      <path
        d="M 10.6 15 C 11 14.4 11.5 14.1 12 14.1 C 12.5 14.1 13 14.4 13.4 15 C 13.4 15.8 12.8 16.6 12 16.6 C 11.2 16.6 10.6 15.8 10.6 15 Z"
        fill="currentColor"
        stroke="none"
      />
      <path d="M 12 16.6 L 12 17.8" />
    </svg>
  );
}
