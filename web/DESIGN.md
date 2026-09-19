<!-- SEED: established from the requested direction before implementation; re-run /impeccable document once there is shipped UI to capture the actual tokens and components. -->
---
name: Bearuang Go
description: A simple, warm autumnal design direction for dependable SME operations.
---

# Design System: Bearuang Go

## Overview

**Creative North Star: "The Friendly Market Ledger"**

Bearuang Go should feel like the well-kept ledger at a neighborhood shop: calm, useful, and easy to pick up without instruction. The visual world is simple and welcoming, with an autumn character carried by toasted leaf, clay, dried-straw, moss, and bark-ink tones rather than seasonal decoration. It should feel warm in changing daylight and under the warmer light of an evening counter.

The star-atlas idea is translated into an operational almanac: information is arranged into clear fields, categories act as the landmarks, and important values are easy to locate at a glance. The metaphor supplies order and memory, not visual complexity. Bearuang Go should make routine work feel human while keeping the reliability and restraint expected from a business system that grows toward ERP and POS.

**Key Characteristics:**
- Warm autumn color used as functional notation, not scattered decoration.
- A single friendly, highly legible interface voice with strong number reading.
- Open, quiet layouts that prioritize the current task over dashboard spectacle.
- Soft tactile states and clear grouping instead of heavy borders or floating card piles.
- A visual grammar that can extend from catalog work to future operational modules.

## Colors

Use a restrained autumn palette: a warm paper ground and deep bark-like ink establish trust, while toasted orange carries action. Muted moss and dried-straw tones explain healthy and pending states without turning the interface into a color chart. Exact values are **to be resolved during implementation**.

### Primary
- **Toasted Leaf:** [to be resolved during implementation]. Use for the primary action, selected navigation, and the clearest active state.

### Secondary
- **Moss:** [to be resolved during implementation]. Use for confirmed, healthy, or available states and occasional secondary emphasis.

### Tertiary
- **Dried Straw:** [to be resolved during implementation]. Use sparingly for pending attention, gentle highlights, and wayfinding—not for decoration.

### Neutral
- **Warm Paper:** [to be resolved during implementation]. The default page ground; light, soft, and comfortable for long operational sessions.
- **Bark Ink:** [to be resolved during implementation]. The main text and structural contrast; dark and warm rather than absolute black.
- **Ashed Paper:** [to be resolved during implementation]. Muted text, quiet fills, dividers, and disabled states.

### Named Rules

**The One Warm Signal Rule.** One warm action color should lead a screen. Other autumn hues explain state or hierarchy; they do not compete for attention.

## Typography

**Display Font:** A sturdy humanist sans with open apertures and soft, friendly terminals [to be resolved during implementation].
**Body Font:** The same family, tuned for dense labels, Indonesian copy, and long operational sessions [to be resolved during implementation].
**Label/Mono Font:** No separate display or code face by default; use the primary family with tabular numerals for quantities and IDR values [to be resolved during implementation].

**Character:** Typography should be approachable without becoming playful. A limited weight range, clear numerals, and generous enough line spacing should make catalog names, statuses, and rupiah amounts quick to scan.

### Hierarchy
- **Display:** [to be resolved during implementation]. Reserved for the page title or the single most important total; never used as decoration.
- **Headline:** [to be resolved during implementation]. Gives each operational view a clear starting point and a calm sense of place.
- **Title:** [to be resolved during implementation]. Names sections, groups, and focused work surfaces with quiet confidence.
- **Body:** [to be resolved during implementation]. Comfortable, legible copy for descriptions, helper text, and empty states; keep longer explanations narrow enough to read without scanning across a dashboard.
- **Label:** [to be resolved during implementation]. Short, sentence-case labels for controls, states, and data columns; use tabular numerals where values must compare.

### Named Rules

**The One-Family Rule.** Create hierarchy through scale, weight, spacing, and color—not by mixing ornamental typefaces. The interface should feel like one considerate voice.

## Layout

The layout follows a quiet market-counter rhythm: orientation stays stable, the working area is obvious, and every view has one dominant task. Prefer one strong work surface with a small amount of supporting context over a wall of equal-weight cards. Let data create density; do not add density through decoration.

Use generous outer breathing room, compact but comfortable rows, and a small repeatable spacing rhythm [to be resolved during implementation]. Group related fields visibly and keep primary actions close to the work they affect. Data tables, catalog lists, and forms should read as continuous operational surfaces rather than isolated promotional panels.

On smaller screens, collapse navigation into a simple top-level control and preserve the same reading order: title, context, primary action, then the work surface. Rows may stack, but labels, values, and actions must remain associated. Avoid horizontal overflow for core tasks.

## Elevation & Depth

Use tonal layering rather than dramatic elevation. A warm paper ground, subtly differentiated work surfaces, and restrained borders should establish depth; shadows belong only to temporary overlays, menus, dialogs, and other surfaces that genuinely sit above the page. The default interface should feel settled on a counter, not made of floating cards.

### Named Rules

**The Grounded Surface Rule.** Resting content stays visually grounded. A surface earns a shadow only when it needs to detach from the page or communicate a temporary state.

## Shapes

Use gently softened corners with a small-to-medium radius strategy [to be resolved during implementation]. Fields, buttons, and compact rows can feel approachable; larger containers should remain composed rather than pill-shaped. Pills are reserved for statuses, filters, and other compact semantic tokens.

Prefer light, warm strokes when separation is necessary, with consistent line weight and no heavy outlines. Controls should feel tactile through clear fill, stroke, and state changes—not through bevels, gloss, or oversized rounding. Keep tables and numeric grids more disciplined than expressive so the data stays trustworthy.

## Do's and Don'ts

### Do:
- **Do** use autumn hues as a clear functional language: action, confirmation, attention, and neutral structure should each have a job.
- **Do** keep interactive states understandable through text, shape, and placement as well as color.
- **Do** give Indonesian labels and Rp values enough room to remain legible, aligned, and easy to compare.
- **Do** make the main task visible immediately and keep the number of competing actions low.
- **Do** preserve the same calm grammar as the product grows from catalog management toward ERP and POS work.

### Don't:
- **Don't** turn the autumn direction into leaf illustrations, pumpkins, paper grain, or other seasonal props.
- **Don't** inherit the starter's monochrome React Router treatment as the product identity.
- **Don't** make every screen a grid of floating cards or a KPI wall.
- **Don't** use cool gray, neon, glossy gradients, glass effects, or absolute-black contrast as the default voice.
- **Don't** use a decorative serif, handwritten face, or color alone to carry meaning.
