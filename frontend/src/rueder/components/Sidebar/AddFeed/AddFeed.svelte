<script lang="ts">
    import { createEventDispatcher } from "svelte"

    import Button from "../../ui/Button.svelte"

    import AddExistingFeed from "./AddExistingFeed.svelte"
    import AddFeedURL from "./AddFeedURL.svelte"
    import ImportFeeds from "./ImportFeeds.svelte"

    const dispatch = createEventDispatcher()

    // Switch modal content
    enum Mode {
        Chooser = 0,
        AddFeedURL,
        AddExistingFeed,
        Import,
    }
    let mode = Mode.Chooser

    function cancelButtonClicked() {
        dispatch("cancel")
    }
</script>

{#if mode == Mode.Chooser}
    <div class="bg-white rounded-xl">
        <div class="p-4 sm:p-6 sm:pb-4">
            <div class="flex items-center">
                <div
                    class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-gray-100 sm:mx-0 sm:h-10 sm:w-10"
                >
                    <!-- Heroicon name: outline/rss -->
                    <svg
                        class="h-6 w-6 text-gray-800"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M6 5c7.18 0 13 5.82 13 13M6 11a7 7 0 017 7m-6 0a1 1 0 11-2 0 1 1 0 012 0z"
                        />
                    </svg>
                </div>
                <div class="text-center sm:mt-0 sm:ml-4 sm:text-left">
                    <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">Add Feed</h3>
                </div>
            </div>
            <div class="pt-3 p-1 text-sm text-gray-500">
                <p>Add a new subscription to your last folder.</p>
            </div>
            <div class="flex flex-col items-center">
                <Button mode="dark" on:click={() => (mode = Mode.AddFeedURL)} addClass="w-full">Add URL</Button>
                <Button mode="dark" on:click={() => (mode = Mode.AddExistingFeed)} addClass="w-full"
                    >Subscribe to existing feed</Button
                >
                <Button mode="light" on:click={() => (mode = Mode.Import)} addClass="w-full">Import</Button>
            </div>
        </div>
        <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse rounded-b-xl">
            <button
                on:click={cancelButtonClicked}
                type="button"
                class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
            >
                Cancel
            </button>
        </div>
    </div>
{:else if mode == Mode.AddExistingFeed}
    <AddExistingFeed on:cancel on:added on:back={() => (mode = Mode.Chooser)} />
{:else if mode == Mode.AddFeedURL}
    <AddFeedURL on:cancel on:added on:back={() => (mode = Mode.Chooser)} />
{:else if mode == Mode.Import}
    <ImportFeeds on:cancel on:added on:back={() => (mode = Mode.Chooser)} />
{:else}
    invalid mode {mode}
{/if}
