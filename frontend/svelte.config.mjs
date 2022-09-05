import preprocess from "svelte-preprocess"
import { dirname, join } from "path"
import { fileURLToPath } from "url"

const __dirname = dirname(fileURLToPath(import.meta.url))

export default {
    preprocess: [
        preprocess({
            postcss: {
                configFilePath: join(__dirname, "postcss.config.cjs"),
            },
            typescript: true,
        }),
    ],
}
