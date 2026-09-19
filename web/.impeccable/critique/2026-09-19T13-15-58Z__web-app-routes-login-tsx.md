---
target: web/app/routes/login.tsx
total_score: 31
max_score: 40
na_heuristics: 
p0_count: 0
p1_count: 2
target_identity: "file:/home/IDSP35936/code/bearuang-go/web/web/app/routes/login.tsx"
timestamp: 2026-09-19T13-15-58Z
slug: web-app-routes-login-tsx
---
## Design Health Score

| # | Heuristic | Score | Key Issue |
|---|-----------|:-----:|-----------|
| 1 | Visibility of System Status | 3 | Submitting shows a clear pending spinner (`Memproses…`), but registration completes with an abrupt push to `/dashboard` without session receipt or account creation feedback. |
| 2 | Match Between System and Real World | 4 | Flawless, natural Indonesian microcopy (*"Ruang kerja yang tetap jernih"*, *"Kata sandi belum sama"*); shop ledger / folio metaphor feels grounded and authentic. |
| 3 | User Control and Freedom | 2 | No password reveal toggle (`Eye`/`EyeOff`). No credential recovery or password reset path (*"Lupa kata sandi?"*) leaves locked-out users without an exit. |
| 4 | Consistency and Standards | 3 | Input patterns follow web conventions, but extensive arbitrary hex classes (`#c25b36`, `#b95532`, `#2e2119`, `#f7f0e6`) duplicate and bypass Tailwind CSS variables in `app.css`. |
| 5 | Error Prevention | 3 | Real-time password matching and HTML attribute constraints are effective, but inability to preview masked passwords on mobile creates high risk of typos. |
| 6 | Recognition Rather Than Recall | 4 | Persistent field labels, contextual leading icons (Mail, Lock), and clear localized placeholders (`nama@perusahaan.id`). |
| 7 | Flexibility and Efficiency | 3 | Enter-to-submit and password manager auto-fill work as expected; lacks convenience accelerators (password peek, "ingat saya"). |
| 8 | Aesthetic and Minimalist Design | 3 | Cohesive warm paper and toasted leaf palette; decorative concentric wireframe rings clash with the physical shop ledger identity. |
| 9 | Help Users Recognize, Diagnose, and Recover from Errors | 4 | `getAuthRequestError` gracefully translates raw backend error codes into respectful, actionable Indonesian instructions. |
| 10 | Help and Documentation | 2 | Purpose of the workspace is stated, but zero contextual help, onboarding FAQ, or support links exist if authentication fails. |
| **Total** | | **31/40** | **Good (28–35)** |

---

## Design Specificity Verdict

**LLM Assessment:** The *Folio Checkpoint* direction successfully establishes a distinct identity away from generic SaaS boilerplate. The tactile pairing of warm paper ground (`#f7f0e6`), bark ink (`#2e2119`), and toasted leaf (`#c25b36`/`#b95532`) gives authentication the feel of a quiet, physical business ledger. Microcopy (*"Ruang kerja yang tetap jernih"*, *"ERP sederhana untuk usaha yang bertumbuh"*) speaks directly to Indonesian shopkeepers and small business operators. However, the background ornamentation in `AuthShell`—namely floating concentric wireframe circles (`size-72` and `size-52` borders)—feels borrowed from generic tech landing pages, diluting the physical ledger thesis.

**Deterministic Scan:** The mechanical CLI detector scan returned 0 AST parsing errors (`[]`). Code-level inspection of craft-floor rules identified 2 typographic tracking floor violations (`tracking-[-0.065em]` on L28 and `tracking-[-0.055em]` on L52 exceeding the Impeccable `-0.04em` floor), 1 WCAG AA contrast failure on tertiary text (`text-[#fff8ee]/70` on `#c25b36` yielding ~2.9:1 vs the 4.5:1 requirement), and systematic design-system drift via raw hex strings rather than semantic Tailwind classes. A detector flag on tinted drop-shadow (`rgba(185,85,50,0.18)`) was verified as a false positive, as it functions as an 8px diffuse elevation shadow rather than a neon halo glow.

**Visual Overlays:** Browser automation is not attached to this harness environment; mutable DOM script injection and live browser overlay presentation were skipped per protocol. The critique relied on deterministic code AST inspection and dual-agent evaluation.

---

## Overall Impression

The login and register experience has an exceptionally strong editorial foundation: the Indonesian copy is natural and respectful, the ledger aesthetic is warm and purposeful, and accessibility foundations (ARIA bindings, focus rings, semantic forms) are thoroughly implemented. The primary flaws are ergonomics and mobile responsiveness: a towering desktop-first aside pushes the actual form below the fold on mobile, masked password inputs lack an eye toggle, and secondary text on the terracotta panel fails WCAG AA contrast.

---

## What's Working

1. **Authentic Indonesian Voice and Ledger Identity:** Rather than sounding like translated English software, the typography, tone, and metaphors (*"Ruang kerja yang tetap jernih"*, *"Buat akun"*) convey calm, professional shopkeeping confidence.
2. **Exemplary Accessibility Wiring:** Every form field links programmatic `<Label>` to `<Input>` with matching `id`, dynamically attaches `aria-describedby` to error strings, and fires `aria-invalid` flags appropriately.
3. **Actionable, Respectful Error Translation:** The `getAuthRequestError` module maps technical API codes into empathetic Indonesian guidance (e.g., explaining when an email is already in use and offering to switch to login).

---

## Priority Issues

### [P1] Missing Password Visibility Toggle
- **Why it matters:** On touch devices and mobile keyboards, typing complex passwords into masked inputs frequently introduces typos. Without an eye toggle, users submit blindly, hit authentication errors, and face needless friction.
- **Fix:** Add an eye icon button (`Eye` / `EyeOff` from `lucide-react`) inside the password inputs in `AuthForm` with `type={show ? "text" : "password"}` and `aria-label="Tampilkan kata sandi"`.
- **Suggested command:** `/impeccable polish`

### [P1] Mobile Form Positioned Below the Fold
- **Why it matters:** On mobile viewports (<768px), the terracotta `<aside>` panel occupies `min-h-[23rem]` (over 360px) of vertical space above the form, forcing the actual credential inputs completely below the fold. Users must scroll down just to discover the input fields.
- **Fix:** In `AuthShell`, condense the mobile aside into a compact header banner (`py-4 px-6`, hidden secondary text, smaller heading) on mobile screens, reserving the full-height vertical folio layout for `lg:` viewports.
- **Suggested command:** `/impeccable layout`

### [P2] Low Contrast on Secondary Text in Terracotta Panel
- **Why it matters:** `text-[#fff8ee]/70` and `text-[#fff8ee]/75` over `#c25b36` achieve contrast ratios of ~2.9:1 and ~3.6:1 respectively, failing the WCAG 2.1 AA 4.5:1 contrast requirement for small text (`text-xs`).
- **Fix:** Increase opacity to `text-[#fff8ee]/90` or solid `#fff8ee` with adjusted font weight to ensure effortless legibility.
- **Suggested command:** `/impeccable clarify`

### [P2] Typographic Tracking Floor Violations
- **Why it matters:** `tracking-[-0.065em]` on the hero heading and `tracking-[-0.055em]` on the form heading violate the Impeccable craft floor (-0.04em), causing letter collisions and cramped glyph baselines at larger viewport scales.
- **Fix:** Clamp letter-spacing to `tracking-[-0.035em]` or standard `tracking-tight` (`-0.025em`) for comfortable display breathing room.
- **Suggested command:** `/impeccable typeset`

### [P2] Design System Drift with Arbitrary Hex Classes
- **Why it matters:** `app/app.css` defines an entire semantic color system (`--primary`, `--background`, `--foreground`, `--border`), yet `auth-shell.tsx` and `auth-form.tsx` bypass it with hardcoded arbitrary hex strings (`#c25b36`, `#b95532`, `#2e2119`, `#f7f0e6`). This produces color drift between components.
- **Fix:** Refactor arbitrary hex classes to reference semantic Tailwind utility classes (`bg-primary`, `bg-background`, `text-foreground`, `border-border`).
- **Suggested command:** `/impeccable harden`

---

## Persona Red Flags

### 1. Casey (Distracted Mobile User)
- **Context:** Checking the ledger on an iPhone while handling orders at the counter.
- **Red Flag:** Casey opens `/login`. The screen is consumed by the terracotta card (`min-h-[23rem]`). Seeing only brand copy and no input fields, Casey assumes it is a marketing page. Upon scrolling down to find the inputs, Casey mistypes a password on the software keyboard, has no eye toggle to check the typo, submits, fails, and gives up in frustration.

### 2. Jordan (Confused First-Timer)
- **Context:** New business owner invited to set up the Bearuang workspace.
- **Red Flag:** Jordan completes registration, clicks *"Buat akun"*, and is abruptly navigated to `/dashboard` with no confirmation dialog, no welcome banner, and no instruction on email verification or initial login credentials.

### 3. Sam (Accessibility-Dependent User)
- **Context:** Low-vision user relying on standard zoom and WCAG contrast ratios.
- **Red Flag:** Sam struggles to read the bottom tagline (*"ERP sederhana untuk usaha yang bertumbuh"*) in the terracotta aside because the 70% opacity yields a washed-out 2.9:1 contrast ratio against the reddish background.

---

## Minor Observations

- **Focus Ring Clipping on Mode Switch:** In `auth-form.tsx`, the inline toggle link uses `focus-visible:ring-4 focus-visible:ring-[#b95532]/20`. When active, the focus outline clips against adjacent paragraph text. Adding `focus-visible:rounded px-1` provides proper visual padding.
- **Generic Wireframe Rings:** The concentric circular vectors in the aside (`size-72` and `size-52`) feel like generic SaaS graphics. Replacing them with ledger-rule watermarks or a quiet folio index stamp (`FOLIO NO. 01`) would deepen the authentic ledger theme.
- **Warm Ochre Selection:** The custom selection highlight (`selection:bg-[#d7b56a]`) is a thoughtful craft touch that harmonizes beautifully with the paper-and-ink aesthetic.

---

## Questions to Consider

- What if the mobile header was styled as an authentic ledger spine across the top, keeping both inputs immediately within comfortable thumb reach without scrolling?
- Could registration provide an instant "Membuka Buku Kas..." transition receipt before redirecting, cementing the feeling of opening a fresh business workspace?
- How might we reassure locked-out shop owners with a quiet recovery path (e.g., *"Butuh bantuan masuk? Hubungi admin ruang kerja"*) without inventing unsupported backend routes?
