module.exports = {
    content: ["./index.html", "./src/**/*.html", "./src/**/*.{svelte,js,ts,jsx,tsx}"],
    theme: {
        extend: {},
    },
    variants: {
        extend: {},
        scrollbar: ["rounded"],
    },
    plugins: [require("@tailwindcss/forms"), require("@spezifisch/tailwind-scrollbar")],
}
