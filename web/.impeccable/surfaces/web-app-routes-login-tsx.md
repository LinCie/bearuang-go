---
version: 1
slug: "web-app-routes-login-tsx"
primary_target: "web/app/routes/login.tsx"
related_targets: ["web/app/routes/register.tsx"]
---

# Auth pages

- **Mode:** Operate
- **Primary target:** `web/app/routes/login.tsx`
- **Related target:** `web/app/routes/register.tsx`
- **Scope:** Separate public login and registration pages using one shared authentication shell.

## Direction contract

**THESIS:** Folio Checkpoint makes authentication feel like opening the next page of a dependable shop ledger, not entering a promotional landing page. It refuses the generic floating auth card by giving the form a quiet operational counterpart.

**OWN-WORLD:** Warm paper is the ground, bark ink carries text, toasted leaf owns the orientation panel and primary action, and moss is reserved for calm supporting states. Fine ledger rules and restrained radius keep the surface tactile but disciplined. The interface uses the existing Geist family and Tailwind utilities only.

**STORY:** A first-time visitor understands that Bearuang is a focused business workspace, chooses Masuk or Buat akun, and completes only the required credential task. Login sends email and password to `/auth/login`; registration sends email and password to `/auth/register` after frontend confirmation. Success continues to `/dashboard` even while that route intentionally remains unavailable.

**FIRST VIEWPORT:** On desktop, a broad toasted-leaf orientation field occupies the left side while a warm-paper form occupies the right; the form heading, fields, action, and route link are visible without scrolling. On mobile, the orientation field stacks above the form, with the action remaining prominent and no horizontal overflow.

**FORM:** Folio Checkpoint is the selected composition from seed `bbd92642`, implemented code-led. The login form has email and password. The registration form adds confirm password, validated only in the frontend. Both pages expose loading, field validation, API error, focus, and disabled states; no password-recovery or email-confirmation action is invented.

**FINISH:** unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md, and every shipping raster carrying its provenance
