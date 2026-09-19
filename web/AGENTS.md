# AGENTS.md

## Frontend

- The frontend lives in `web/` and uses React 19, React Router 8, TypeScript, Vite, Tailwind CSS 4, Base UI, shadcn, and Lucide icons.
- Run frontend commands from `web/`: `bun run dev`, `bun run build`, and `bun run typecheck`.
- Routes are configured in `web/app/routes.ts`; route modules belong in `web/app/routes/` and use generated `./+types/*` types.
- Reusable UI components belong in `web/app/components/`; shadcn components go in `web/app/components/ui/`. Shared helpers belong in `web/app/lib/`.
- Be a shadcn and Tailwind maximalist: prefer shadcn components and Tailwind utility classes for UI, styling, and layout instead of custom alternatives.
- Use the `~/*` TypeScript alias for imports from `web/app`, Tailwind utility classes for component styling, and the CSS variables in `web/app/app.css` for theme values.
- Keep TypeScript strict, prefer accessible semantic HTML, and use existing UI primitives and Lucide icons before adding new dependencies.
