import { resolve } from "path"
import { defineConfig, splitVendorChunkPlugin } from "vite"
import { svelte } from "@sveltejs/vite-plugin-svelte"
import transformPlugin from "vite-plugin-transform"

// for transformPlugin
const alias = {
    "@app": resolve(__dirname, "./src/rueder"),
}

export default defineConfig({
    clearScreen: false,
    build: {
        chunkSizeWarningLimit: 1000,
    },
    plugins: [
        svelte(),
        splitVendorChunkPlugin(),
        transformPlugin({
            alias, // enable replace aliases resolver
            exclude: ["node_modules"],
        }),
    ],
})
