import type { ComponentProps } from "react";

export type BearReactionState =
  | "idle"
  | "email"
  | "password"
  | "peek"
  | "submitting"
  | "error"
  | "success";

type BearCompanionProps = {
  state: BearReactionState;
  emailLength?: number;
  className?: string;
} & Omit<ComponentProps<"svg">, "children">;

export function BearCompanion({
  state,
  emailLength = 0,
  className = "",
  ...props
}: BearCompanionProps) {
  // Clamp email length and calculate horizontal pupil progress (-2.2px to +2.2px)
  const clampedEmailLength = Math.max(0, Math.min(emailLength, 32));
  const pupilProgress = clampedEmailLength / 32;
  const emailPupilX = -2.2 + pupilProgress * 4.4;

  // Determine pupil offsets based on state
  let pupilTransform = "translate(0px, 0px)";
  if (state === "email") {
    pupilTransform = `translate(${emailPupilX.toFixed(1)}px, 2.5px)`;
  } else if (state === "peek") {
    pupilTransform = "translate(0.5px, 1px)";
  } else if (state === "submitting") {
    pupilTransform = "translate(0px, -0.6px)";
  } else if (state === "error") {
    pupilTransform = "translate(-0.5px, 1px)";
  }

  // Head tilt for micro-expressions
  let headTransform = "rotate(0deg)";
  if (state === "peek") {
    headTransform = "rotate(2.5deg)";
  } else if (state === "error") {
    headTransform = "rotate(-3.5deg)";
  } else if (state === "email") {
    headTransform = "rotate(1deg) translateY(0.5px)";
  } else if (state === "submitting") {
    headTransform = "translateY(0.5px)";
  }

  // Paw transformations
  // Left eye is at (62, 52), Right eye is at (98, 52)
  const isCovering = state === "password" || state === "peek";
  const isPeeking = state === "peek";

  // Left paw: covers left eye when covering or peeking, otherwise rests on counter
  const leftPawTransform = isCovering
    ? "translate(0px, 0px) rotate(0deg)"
    : "translate(-18px, 60px) rotate(-6deg)";

  // Right paw: covers right eye on password, lowers & angles away on peek, rests on counter otherwise
  const rightPawTransform = isPeeking
    ? "translate(18px, 24px) rotate(22deg)"
    : isCovering
      ? "translate(0px, 0px) rotate(0deg)"
      : "translate(18px, 60px) rotate(6deg)";

  return (
    <div className={`relative flex flex-col items-center select-none ${className}`}>
      <svg
        aria-hidden="true"
        className="h-auto w-32 overflow-hidden sm:w-36"
        viewBox="0 0 160 120"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        {...props}
      >
        <defs>
          <clipPath id="bear-viewport-clip">
            <rect x="0" y="0" width="160" height="120" rx="4" />
          </clipPath>
        </defs>

        <g clipPath="url(#bear-viewport-clip)">
          {/* Shoulders & Body */}
          <path
            d="M 24 118 C 26 94, 46 86, 60 86 L 100 86 C 114 86, 134 94, 136 118 Z"
            fill="var(--foreground)"
          />
          {/* Subtle warm chest bib / counter apron */}
          <path
            d="M 64 86 C 64 102, 96 102, 96 86 L 100 118 L 60 118 Z"
            fill="var(--muted)"
            opacity={0.25}
          />

          {/* Head group (animates tilt) */}
          <g
            className="transition-transform duration-300 ease-out motion-reduce:transition-none"
            style={{
              transform: headTransform,
              transformOrigin: "80px 76px",
            }}
          >
            {/* Outer Ears */}
            <circle cx="50" cy="34" r="14" fill="var(--foreground)" />
            <circle cx="110" cy="34" r="14" fill="var(--foreground)" />

            {/* Inner Ears */}
            <circle cx="50" cy="34" r="8" fill="#8E6F59" />
            <circle cx="110" cy="34" r="8" fill="#8E6F59" />

            {/* Head Silhouette */}
            <path
              d="M 52 38 C 52 26, 108 26, 108 38 C 120 45, 124 64, 118 78 C 112 93, 94 95, 80 95 C 66 95, 48 93, 42 78 C 36 64, 40 45, 52 38 Z"
              fill="var(--foreground)"
            />

            {/* Soft Cheeks Blush */}
            <circle cx="49" cy="69" r="5" fill="var(--primary)" opacity={0.16} />
            <circle cx="111" cy="69" r="5" fill="var(--primary)" opacity={0.16} />

            {/* Muzzle (Snout Base) */}
            <rect
              x="60"
              y="59"
              width="40"
              height="26"
              rx="13"
              fill="var(--muted)"
            />

            {/* Nose */}
            <rect
              x="73"
              y="63"
              width="14"
              height="9"
              rx="4.5"
              fill="var(--foreground)"
            />
            {/* Nose Specular Highlight */}
            <ellipse cx="78" cy="65.5" rx="2" ry="1" fill="#FFFFFF" opacity={0.4} />

            {/* Mouth */}
            <path
              d="M 80 72 L 80 76 M 75 76 Q 80 79.5 85 76"
              stroke="var(--foreground)"
              strokeWidth={1.8}
              strokeLinecap="round"
            />

            {/* Eyebrows */}
            {state === "error" ? (
              <>
                <path
                  d="M 58 45 Q 63 43 68 41"
                  stroke="var(--foreground)"
                  strokeWidth={2}
                  strokeLinecap="round"
                />
                <path
                  d="M 92 41 Q 97 43 102 45"
                  stroke="var(--foreground)"
                  strokeWidth={2}
                  strokeLinecap="round"
                />
              </>
            ) : (
              <>
                <path
                  d="M 58 43 Q 63 41 68 43"
                  stroke="var(--foreground)"
                  strokeWidth={2}
                  strokeLinecap="round"
                />
                <path
                  d="M 92 43 Q 97 41 102 43"
                  stroke="var(--foreground)"
                  strokeWidth={2}
                  strokeLinecap="round"
                />
              </>
            )}

            {/* Eyes */}
            {state === "success" ? (
              <>
                {/* Smiling curved eyes */}
                <path
                  d="M 57 54 Q 62 48 67 54"
                  stroke="var(--foreground)"
                  strokeWidth={2.4}
                  strokeLinecap="round"
                />
                <path
                  d="M 93 54 Q 98 48 103 54"
                  stroke="var(--foreground)"
                  strokeWidth={2.4}
                  strokeLinecap="round"
                />
              </>
            ) : (
              <>
                {/* Left Eye */}
                <g className={state === "idle" ? "animate-bear-blink" : ""}>
                  <ellipse
                    cx="62"
                    cy="52"
                    rx="5.5"
                    ry="6.5"
                    fill="#FFFAF3"
                    stroke="var(--foreground)"
                    strokeWidth={1.2}
                  />
                  <g
                    className="transition-transform duration-150 ease-out motion-reduce:transition-none"
                    style={{ transform: pupilTransform }}
                  >
                    <circle cx="62" cy="52" r="3.6" fill="var(--foreground)" />
                    <circle cx="60.8" cy="50.5" r="1.3" fill="#FFFFFF" />
                  </g>
                </g>

                {/* Right Eye */}
                <g className={state === "idle" ? "animate-bear-blink" : ""}>
                  <ellipse
                    cx="98"
                    cy="52"
                    rx="5.5"
                    ry="6.5"
                    fill="#FFFAF3"
                    stroke="var(--foreground)"
                    strokeWidth={1.2}
                  />
                  <g
                    className="transition-transform duration-150 ease-out motion-reduce:transition-none"
                    style={{ transform: pupilTransform }}
                  >
                    <circle cx="98" cy="52" r="3.6" fill="var(--foreground)" />
                    <circle cx="96.8" cy="50.5" r="1.3" fill="#FFFFFF" />
                  </g>
                </g>
              </>
            )}
          </g>

          {/* Left Paw (Group origin centered at left eye 62, 52) */}
          <g
            className="transition-transform duration-300 ease-out motion-reduce:transition-none"
            style={{
              transform: leftPawTransform,
              transformOrigin: "62px 52px",
            }}
          >
            {/* Forearm & Paw Capsule */}
            <rect
              x="48"
              y="38"
              width="28"
              height="80"
              rx="14"
              fill="var(--foreground)"
            />
            {/* Paw Main Pad */}
            <ellipse cx="62" cy="52" rx="8" ry="6" fill="#8E6F59" />
            {/* 3 Toe Pads */}
            <circle cx="56" cy="44" r="2" fill="#8E6F59" />
            <circle cx="62" cy="42" r="2.2" fill="#8E6F59" />
            <circle cx="68" cy="44" r="2" fill="#8E6F59" />
          </g>

          {/* Right Paw (Group origin centered at right eye 98, 52) */}
          <g
            className="transition-transform duration-300 ease-out motion-reduce:transition-none"
            style={{
              transform: rightPawTransform,
              transformOrigin: "98px 52px",
            }}
          >
            {/* Forearm & Paw Capsule */}
            <rect
              x="84"
              y="38"
              width="28"
              height="80"
              rx="14"
              fill="var(--foreground)"
            />
            {/* Paw Main Pad */}
            <ellipse cx="98" cy="52" rx="8" ry="6" fill="#8E6F59" />
            {/* 3 Toe Pads */}
            <circle cx="92" cy="44" r="2" fill="#8E6F59" />
            <circle cx="98" cy="42" r="2.2" fill="#8E6F59" />
            <circle cx="104" cy="44" r="2" fill="#8E6F59" />
          </g>

          {/* Grounded Shop Counter Shelf */}
          <rect
            x="6"
            y="114"
            width="148"
            height="8"
            rx="4"
            fill="var(--border)"
          />
          <line
            x1="10"
            y1="115"
            x2="150"
            y2="115"
            stroke="#FFFAF3"
            strokeWidth={1}
            opacity={0.65}
          />
        </g>
      </svg>
    </div>
  );
}
