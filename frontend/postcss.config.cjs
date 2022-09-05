const path = require("path")
const tailwindcss = require("tailwindcss")

module.exports = {
    plugins: {
        "tailwindcss/nesting": {},
        tailwindcss: path.resolve(__dirname, "./tailwind.config.js"),
        autoprefixer: {},
    },
}
