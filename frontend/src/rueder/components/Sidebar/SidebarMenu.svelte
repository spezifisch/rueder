<script lang="ts">
    import { getID, getUsername, logout } from "../../stores/session"

    import ModalWrapper from "../ui/ModalWrapper.svelte"
    import SimpleButton from "./Menu/SimpleButton.svelte"
    import OPMLExport from "./OPMLExport.svelte"

    export let show: boolean

    const username = getUsername()
    const userid = getID()

    let doShowOPMLExport = false
    function opmlExport() {
        doShowOPMLExport = true
    }
</script>

{#if show}
    <!-- show "OPML Export" modal -->
    {#if doShowOPMLExport}
        <ModalWrapper centered={true}>
            <OPMLExport on:close={() => (doShowOPMLExport = false)} />
        </ModalWrapper>
    {/if}

    <div class="flex-auto text-gray-300">
        <ul class="p-2">
            <li class="flex flex-row flex-wrap items-center border-b-2 border-gray-400 pb-2 mb-1">
                <!-- user box -->
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-6 w-6 flex-none"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M5.121 17.804A13.937 13.937 0 0112 16c2.5 0 4.847.655 6.879 1.804M15 10a3 3 0 11-6 0 3 3 0 016 0zm6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                </svg>
                <div class="flex-auto px-2 flex flex-col">
                    <span class="font-bold">{username}</span>
                    <span class="text-xs text-gray-400">{userid}</span>
                </div>
                <button
                    class="text-sm text-gray-300 bg-gray-800 hover:text-red-600 hover:bg-gray-700 p-2 text-right flex-none flex items-center space-x-1"
                    type="button"
                    on:click={logout}
                >
                    <span>Logout</span>
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-5 w-5 text-red-800 inline"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                        />
                    </svg>
                </button>
            </li>

            <SimpleButton on:click={opmlExport}>OPML Export</SimpleButton>
        </ul>
    </div>
{/if}
