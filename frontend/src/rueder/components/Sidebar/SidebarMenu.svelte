<script lang="ts">
    import { getID, getUsername, logout } from "../../stores/session"
    import OutlineLogout from "../ui/heroicons/outline-logout.svelte"
    import OutlineUserCircle from "../ui/heroicons/outline-user-circle.svelte"

    import ModalWrapper from "../ui/ModalWrapper.svelte"
    import SimpleButton from "./Menu/SimpleButton.svelte"
    import OPMLExport from "./OPMLExport.svelte"
    import DebugSSE from "./Debug/DebugSSE.svelte"

    export let show: boolean

    const username = getUsername()
    const userid = getID()

    let showOPMLExport = false
    let showDebugSSE = false
</script>

{#if show}
    <!-- show "OPML Export" modal -->
    {#if showOPMLExport}
        <ModalWrapper centered={true}>
            <OPMLExport on:close={() => (showOPMLExport = false)} />
        </ModalWrapper>
    {/if}
    <!-- show "Debug: SSE" modal -->
    {#if showDebugSSE}
        <DebugSSE on:close={() => (showDebugSSE = false)} />
    {/if}

    <div class="flex-auto text-gray-300">
        <ul class="p-2">
            <li class="flex flex-row flex-wrap items-center border-b-2 border-gray-500 pb-2 mb-2">
                <!-- user box -->
                <OutlineUserCircle />
                <div class="flex-auto px-2 flex flex-col">
                    <span class="font-bold">{username}</span>
                    <span class="text-xs text-gray-400">{userid}</span>
                </div>
                <button
                    class="text-sm text-gray-300 bg-gray-800 border-gray-600 border hover:text-red-600 hover:bg-gray-700 p-2 text-right flex-none flex items-center space-x-1"
                    type="button"
                    on:click={logout}
                >
                    <span>Logout</span>
                    <OutlineLogout size={5} addClass="text-red-800 inline" />
                </button>
            </li>

            <li>
                <SimpleButton on:click={() => (showOPMLExport = true)}>OPML Export</SimpleButton>
            </li>
            <li>
                <SimpleButton on:click={() => (showDebugSSE = true)}>Debug: SSE</SimpleButton>
            </li>
        </ul>
    </div>
{/if}

<style lang="postcss">
    li:not(:first-child) {
        @apply my-2;
    }
</style>
