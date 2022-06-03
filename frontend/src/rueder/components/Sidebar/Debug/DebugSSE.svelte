<script lang="ts">
    import { createEventDispatcher, onDestroy, tick } from "svelte"
    import { derived } from "svelte/store"
    import type { Readable } from "svelte/store"

    import { scrollToBottom } from "svelte-scrollto"

    import FullModal from "../../ui/FullModal.svelte"
    import Button from "../../ui/Button.svelte"
    import OutlineChip from "../../ui/heroicons/outline-chip.svelte"

    import { SSEStore } from "../../../stores/SSE"

    const dispatch = createEventDispatcher()

    // SSE connection
    const sse = new SSEStore()
    onDestroy(() => {
        console.log("closing sse")
        sse.close()
    })
    sse.connect()

    // our SSE store always only contains the latest message,
    // create a derived store that concatenated the messages
    let scrollContainer: HTMLDivElement
    const sseLogger: Readable<string> = derived(sse.store, ($store) => {
        if ($sseLogger) {
            return $sseLogger + "\n" + $store
        }
        return $store
    })
    // scroll down after the store has been updated
    sseLogger.subscribe(async () => {
        await handleLogAppend()
    })

    // command box
    let commandInput: HTMLInputElement
    let enteredCommand: string
    async function handleCommand() {
        commandInput.value = ""

        await handleLogAppend()
    }

    async function handleKeypress(e: KeyboardEvent) {
        if (e.key == "Escape") {
            // close
            e.preventDefault()
            dispatch("close")
        } else if (e.key == "Enter") {
            // enter command the same way as clicking "Enter"
            e.preventDefault()
            await handleCommand()
        }
    }

    // called when anything is appended to the log
    async function handleLogAppend() {
        await tick() // need to wait until the appended stuff is rendered so the element height is updated
        scrollToBottom({
            container: scrollContainer,
            offset: 200, // need to add an offset or it stops a bit before the bottom
            duration: 100,
        })
    }
</script>

<FullModal on:close>
    <span slot="header">
        <div class="flex items-center bg-gray-900 p-2">
            <div
                class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-gray-800 sm:mx-0 sm:h-10 sm:w-10"
            >
                <OutlineChip size={6} addClass="text-green-600" />
            </div>
            <div class="text-center sm:mt-0 sm:ml-4 sm:text-left">
                <h3 class="text-lg leading-6 font-medium text-green-600" id="modal-title">Debug Server-Sent Events</h3>
            </div>
        </div>
    </span>

    <div slot="scrolling" class="flex-auto overflow-auto rueder-scrollbar" bind:this={scrollContainer}>
        <p class="text-left p-2 text-white">
            {#if !$sseLogger}
                Connecting to {sse.endpoint}
            {:else}
                <pre>{$sseLogger}</pre>
            {/if}
        </p>
    </div>

    <span slot="buttons">
        <div class="flex flex-row">
            <!-- svelte-ignore a11y-autofocus -->
            <input
                type="text"
                class="flex-auto m-2 bg-gray-800 text-green-600 focus:ring-green-700 focus:border-none"
                placeholder="Enter command"
                autofocus
                bind:value={enteredCommand}
                bind:this={commandInput}
                on:keypress={handleKeypress}
            />
            <div class="flex-none">
                <Button mode="dark" on:click={handleCommand}>Enter</Button>
                <Button mode="error-dark" on:click={() => dispatch("close")}>Close</Button>
            </div>
        </div>
    </span>
</FullModal>
