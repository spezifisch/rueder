import { resolve } from "path"
import { defineConfig, splitVendorChunkPlugin } from "vite"
import { svelte } from "@sveltejs/vite-plugin-svelte"
import legacy from "@vitejs/plugin-legacy"
import transformPlugin from "vite-plugin-transform"

// for transformPlugin
const alias = {
    "@app": resolve(__dirname, "./src/rueder"),
}

export default defineConfig({
    clearScreen: false,
    build: {
        chunkSizeWarningLimit: 1000,
        minify: "terser",
    },
    plugins: [
        legacy({
            targets: ["defaults", "not IE 11"],
        }),
        svelte(),
        splitVendorChunkPlugin(),
        transformPlugin({
            alias, // enable replace aliases resolver
            exclude: ["node_modules"],
        }),
    ],
    resolve: {
        // these polyfills are a workaround for https://github.com/vitejs/vite/issues/9200
        // remove them when the issue gets resolved
        alias: {
            fs: "browser-fs-access",
            path: "path-browserify",
            url: "url",
        },
    },
    server: {
        port: 3000,
    },
})
