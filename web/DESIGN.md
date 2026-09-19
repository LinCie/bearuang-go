## <!-- SEED: established from the requested direction before implementation; re-run /impeccable document once there is shipped UI to capture the actual tokens and components. -->

name: Bearuang
description: A grounded, warm operational system where bear-like steadiness protects the money and daily work of small businesses.

---

# Design System: Bearuang Go

## Brand layer

**Bearuang = beruang + uang.** The bear is a behavioral metaphor for steadiness, protection, resilience, and warmth. The *uang* half keeps the identity tied to the real work of seeing and managing money, stock, and daily operations clearly.

This identity sits on top of the existing autumnal market-ledger direction. The product should feel like a capable bear at the shop counter: grounded, attentive, and strong enough to make routine work feel safe. The bear is not a mascot that needs to appear on every screen. It is felt through the interface's calm protection of the user's time, attention, and resources.

**Brand expression:** warm but not cute, sturdy but not severe, financial but not flashy. Use a restrained bear cue in a logo, empty state, onboarding moment, or illustration when it genuinely helps orientation. Let layout, legible rupiah values, and dependable states carry the identity everywhere else.

## Overview

**Creative North Star: "The Friendly Market Ledger"**

Bearuang Go should feel like the well-kept ledger at a neighborhood shop: calm, useful, and easy to pick up without instruction. The visual world is simple and welcoming, with an autumn character carried by toasted leaf, clay, dried-straw, moss, and bark-ink tones rather than seasonal decoration. It should feel warm in changing daylight and under the warmer light of an evening counter.

The bearuang layer gives that ledger a point of view: the bear stands for a steady hand, while *uang* makes clarity around value and stock a first-class concern. The star-atlas idea is translated into an operational almanac: information is arranged into clear fields, categories act as the landmarks, and important values are easy to locate at a glance. The metaphors supply order, care, and memory—not visual complexity. Bearuang Go should make routine work feel human while keeping the reliability and restraint expected from a business system that grows toward ERP and POS.

**Key Characteristics:**

- Warm autumn color used as functional notation, not scattered decoration.
- A grounded bearuang identity expressed through care, legibility, and stable interaction—not mascot wallpaper.
- A single friendly, highly legible interface voice with strong number reading.
- Open, quiet layouts that prioritize the current task over dashboard spectacle.
- Soft tactile states and clear grouping instead of heavy borders or floating card piles.
- A visual grammar that can extend from catalog work to future operational modules.

## Colors

Use the implemented autumn palette as Bearuang's grounded material language: warm paper and bark ink establish trust, toasted leaf carries action, and moss and dried-straw explain state. The bearuang identity does not need a separate brown or gold layer; it is already present in the bark, counter, and ledger relationship.

### Primary

- **Toasted Leaf — `#A7472A`:** Primary action, selected navigation, and the clearest active state.

### Secondary

- **Moss — `#4F6546`:** Confirmed, healthy, or available states and occasional secondary emphasis.

### Tertiary

- **Dried Straw — `#D7B56A`:** Pending attention, gentle highlights, and wayfinding. It is not a wealth cue and should never become a gold wash.

### Neutral

- **Warm Paper — `#F7F0E6`:** Default page ground; light, soft, and comfortable for long operational sessions.
- **Ledger White — `#FFFAF3`:** Resting work surfaces and readable overlays.
- **Bark Ink — `#2E2119`:** Main text and structural contrast; dark and warm rather than absolute black.
- **Ashed Paper — `#EDE4D8`:** Muted fills and quiet states.
- **Counter Stroke — `#D8C8B6`:** Light dividers and field boundaries.

### Named Rules

**The One Warm Signal Rule.** One warm action color should lead a screen. Other autumn hues explain state or hierarchy; they do not compete for attention.

**Money Stays Legible.** Rupiah values use clear hierarchy, generous room, and tabular numerals where comparison matters. Color can support a financial state, but never carries the meaning alone.

## Typography

**Display and Body:** Geist Variable, the current interface family. Use its open, compact sans character for Indonesian labels, dense catalog data, and long operational sessions; create hierarchy through scale and weight rather than adding a mascot-like display face.
**Label/Mono:** No separate mono face by default. Use the primary family with tabular numerals for quantities, SKU values, and IDR amounts.

**Character:** Typography should be approachable without becoming playful. A limited weight range, clear numerals, and generous enough line spacing should make catalog names, statuses, and rupiah amounts quick to scan.

### Hierarchy

- **Display:** Reserved for the page title or the single most important total; never used as decoration.
- **Headline:** Gives each operational view a clear starting point and a calm sense of place.
- **Title:** Names sections, groups, and focused work surfaces with quiet confidence.
- **Body:** Comfortable, legible copy for descriptions, helper text, and empty states; keep longer explanations narrow enough to read without scanning across a dashboard.
- **Label:** Short, sentence-case labels for controls, states, and data columns; use tabular numerals where values must compare.

### Named Rules

**The One-Family Rule.** Create hierarchy through scale, weight, spacing, and color—not by mixing ornamental typefaces. The interface should feel like one considerate voice.

**The Held Ledger Rule.** Every surface should feel like it is protecting the owner's time, attention, money, or stock. Remove decoration that does not help someone understand or act.

## Motif and illustration

The bearuang theme is a semantic layer, not a pattern layer.

- A bear cue may appear as a restrained single-color mark, rounded ear/shoulder silhouette, or paw-like grouping in the brand lockup, onboarding, or an empty state.
- Illustrations should show calm capability—holding, sorting, or watching over a ledger or stock—not cartoon antics or promises of wealth.
- *Uang* is expressed through clear `Rp` formatting, useful totals, and trustworthy state language. Do not use coin piles, banknotes, literal gold, or finance clichés as decoration.
- When no illustration is needed, omit the bear. The interface should still feel Bearuang through its steadiness and care.

## Layout

The layout follows a quiet market-counter rhythm: orientation stays stable, the working area is obvious, and every view has one dominant task. Prefer one strong work surface with a small amount of supporting context over a wall of equal-weight cards. Let data create density; do not add density through decoration.

Use generous outer breathing room, compact but comfortable rows, and a small repeatable spacing rhythm [to be resolved during implementation]. Group related fields visibly and keep primary actions close to the work they affect. Data tables, catalog lists, and forms should read as continuous operational surfaces rather than isolated promotional panels.

On smaller screens, collapse navigation into a simple top-level control and preserve the same reading order: title, context, primary action, then the work surface. Rows may stack, but labels, values, and actions must remain associated. Avoid horizontal overflow for core tasks.

## Elevation & Depth

Use tonal layering rather than dramatic elevation. A warm paper ground, subtly differentiated work surfaces, and restrained borders should establish depth; shadows belong only to temporary overlays, menus, dialogs, and other surfaces that genuinely sit above the page. The default interface should feel settled on a counter, not made of floating cards.

### Named Rules

**The Grounded Surface Rule.** Resting content stays visually grounded. A surface earns a shadow only when it needs to detach from the page or communicate a temporary state.

**The Bear, Not the Mascot Rule.** Bearuang should be recognizable from its behavior—steady, protective, and clear—before any bear illustration is introduced.

## Shapes

Use gently softened corners with a small-to-medium radius strategy [to be resolved during implementation]. Fields, buttons, and compact rows can feel approachable; larger containers should remain composed rather than pill-shaped. Pills are reserved for statuses, filters, and other compact semantic tokens.

Prefer light, warm strokes when separation is necessary, with consistent line weight and no heavy outlines. Controls should feel tactile through clear fill, stroke, and state changes—not through bevels, gloss, or oversized rounding. Keep tables and numeric grids more disciplined than expressive so the data stays trustworthy.

## Do's and Don'ts

### Do:

- **Do** use autumn hues as a clear functional language: action, confirmation, attention, and neutral structure should each have a job.
- **Do** make Bearuang feel steady and protective through clear hierarchy, calm states, and reliable affordances.
- **Do** keep interactive states understandable through text, shape, and placement as well as color.
- **Do** give Indonesian labels and Rp values enough room to remain legible, aligned, and easy to compare.
- **Do** use bear imagery only where it adds orientation, warmth, or reassurance.
- **Do** make the main task visible immediately and keep the number of competing actions low.
- **Do** preserve the same calm grammar as the product grows from catalog management toward ERP and POS work.

### Don't:

- **Don't** turn the autumn direction into leaf illustrations, pumpkins, paper grain, or other seasonal props.
- **Don't** turn Bearuang into a cartoon mascot, paw-print pattern, forest scene, or recurring animal wallpaper.
- **Don't** use coin piles, literal gold, banknote imagery, or wealth clichés to stand in for the *uang* idea.
- **Don't** inherit the starter's monochrome React Router treatment as the product identity.
- **Don't** make every screen a grid of floating cards or a KPI wall.
- **Don't** use cool gray, neon, glossy gradients, glass effects, or absolute-black contrast as the default voice.
- **Don't** use a decorative serif, handwritten face, or color alone to carry meaning.
