/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_RUEDER_BASE_URL_LOGIN: string
    readonly VITE_RUEDER_BASE_URL_API: string
    readonly VITE_RUEDER_BASE_URL_IMGPROXY: string
    readonly VITE_IMGPROXY_KEY: string
    readonly VITE_IMGPROXY_SALT: string
}

interface ImportMeta {
    readonly env: ImportMetaEnv
}
