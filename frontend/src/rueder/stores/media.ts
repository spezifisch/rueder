import watchMedia from "@spezifisch/svelte-media"

const mediaqueries = {
    /* same behavior as in tailwind: https://tailwindcss.com/docs/breakpoints */
    sm: "(min-width: 640px)",
    md: "(min-width: 768px)",
    lg: "(min-width: 1024px)",
    xl: "(min-width: 1280px)",
    "2xl": "(min-width: 1536px)",
    /* examples from https://github.com/cibernox/svelte-media#usage */
    dark: "(prefers-color-scheme: dark)",
    noanimations: "(prefers-reduced-motion: reduce)",
}

export const media = watchMedia(mediaqueries)
