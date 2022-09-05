<script lang="ts">
    import "./global.postcss"
    import { sessionStore } from "./stores/session"
    import Main from "./Main.svelte"
    import Login from "./Login.svelte"

    // these defaults are consistent with backend's development docker-compose file.
    // you can override them by setting the environment variable when building.
    // they MUST have a trailing slash.
    // loginBaseURL points to loginsrv's URL (where "login" is appended)
    const loginBaseURL = import.meta.env.VITE_RUEDER_BASE_URL_LOGIN ?? "http://127.0.0.1:8082/"
    // apiBaseURL points to the rueder backend API (where /folders, /feed, etc. are appended)
    const apiBaseURL = import.meta.env.VITE_RUEDER_BASE_URL_API ?? "http://127.0.0.1:8080/api/v1/"
    // sseBaseURL points to the rueder SSE API (where /sse is appended)
    const sseBaseURL = import.meta.env.VITE_RUEDER_SSE_URL_API ?? "http://127.0.0.1:8083/"

    // base path to imgproxy instance (https://github.com/imgproxy/imgproxy) to proxy favicons, thumbnails
    // and article content images. use empty string disable proxying.
    const imageProxyBaseURL = import.meta.env.VITE_RUEDER_BASE_URL_IMGPROXY ?? "http://127.0.0.1:8086/"
    // set key and salt to the ones configured in imgproxy. specify empty strings to use imgproxy without signatures.
    // IMPORTANT this will end up in the client js bundle and is NOT SECRET. it's only a kind of obfuscation and not very useful.
    const imageProxyKey = import.meta.env.VITE_IMGPROXY_KEY ?? ""
    const imageProxySalt = import.meta.env.VITE_IMGPROXY_SALT ?? ""

    const isDev = import.meta.env.DEV
    const imageProxyUseTypePrefixes = !isDev
</script>

<svelte:head>
    {#if isDev}
        <title>rueder3-dev</title>
    {/if}
</svelte:head>

{#if $sessionStore.loggedIn}
    <Main baseURL={apiBaseURL} {sseBaseURL} {imageProxyBaseURL} {imageProxyUseTypePrefixes} {imageProxyKey} {imageProxySalt} />
{:else}
    <Login page="rueder" baseURL={loginBaseURL} />
{/if}
