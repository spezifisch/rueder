<script lang="ts">
    import { createEventDispatcher, getContext } from "svelte"

    import type { FeedAPI } from "../../../api/feed"
    import type { AddFeedResponse } from "../../../api/types"
    import { contextKey } from "../../../helpers/constants"

    const feedAPI: FeedAPI = getContext(contextKey.feedAPI)

    const dispatch = createEventDispatcher()

    // Switch modal content
    enum Card {
        Add = 0, // type in a URL
        Fetch, // wait until feed is fetched
    }
    let card = Card.Add

    // State for the various modal parts
    enum State {
        Initial = 0,
        WaitAdd,
        Error,
    }
    let state = State.Initial

    // disable input box
    let disableInputBox: boolean
    $: disableInputBox = state == State.WaitAdd

    // disable add button
    let disableAdd: boolean
    $: disableAdd = state == State.WaitAdd

    // URL input (bind)
    let url = ""

    // Error message box (bind)
    let errorMsg = ""

    async function addNewFeed() {
        state = State.WaitAdd

        if (!url) {
            state = State.Error
            errorMsg = "Please enter a URL!"
            return
        }

        try {
            const resp: AddFeedResponse = await feedAPI.AddFeed(url)

            if (!resp.ok) {
                // failure adding due to network/server problems or existing feed or invalid url
                errorMsg = resp.message
                state = State.Error
                return
            }

            card = Card.Fetch
            state = State.Initial
            url = ""

            // add new feed to a folder
            dispatch("added", [resp.feed_id])
        } catch (e) {
            // rejected promise: shouldn't happen
            errorMsg = e.message
            state = State.Error
        }
    }

    function backButtonClicked() {
        dispatch("back")
    }

    function cancelButtonClicked() {
        dispatch("cancel")
    }
</script>

<div>
    {#if card == Card.Add}
        <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4 rounded-t-xl">
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
                <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                    <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">Add Feed URL</h3>
                </div>
            </div>

            <div class="pt-3 p-1 text-sm text-gray-500">
                <p>The feed URL can be in RSS, Atom or JSON format. You need to supply the direct link to the feed.</p>
                <div class="relative mt-2">
                    <input
                        bind:value={url}
                        disabled={disableInputBox}
                        type="text"
                        placeholder="URL"
                        class="p-2 rounded text-black text-sm border border-gray-300 outline-none focus:ring-gray-600 w-full {disableInputBox
                            ? 'bg-gray-200'
                            : ''}"
                    />
                    <!-- clear button from: https://github.com/tailwindlabs/tailwindui-issues/issues/504#issuecomment-782291191 -->
                    <button
                        class="absolute inset-y-0 right-0 pr-2 flex items-center"
                        type="button"
                        disabled={disableInputBox}
                    >
                        <!-- heroicon: solid/x-circle -->
                        <svg
                            class="h-5 w-5 text-gray-400 hover:text-gray-600 {url == '' ? 'hidden' : ''}"
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 20 20"
                            fill="currentColor"
                            aria-hidden="true"
                            on:click={() => {
                                url = ""
                            }}
                        >
                            <path
                                fill-rule="evenodd"
                                d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                                clip-rule="evenodd"
                            />
                        </svg>
                    </button>
                </div>
            </div>
            {#if state == State.Error && !disableInputBox}
                <div class="mt-2">
                    <p class="text-red-800">
                        {errorMsg}
                    </p>
                </div>
            {/if}
        </div>
        <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse rounded-b-xl">
            <button
                on:click={addNewFeed}
                disabled={disableAdd}
                type="button"
                class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-green-600 text-base font-medium text-white hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:ml-3 sm:w-auto sm:text-sm"
            >
                Add new feed
            </button>
            <button
                on:click={backButtonClicked}
                type="button"
                class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
            >
                Back
            </button>
            <button
                on:click={cancelButtonClicked}
                type="button"
                class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
            >
                Cancel
            </button>
        </div>
    {:else if card == Card.Fetch}
        <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4 rounded-t-xl">
            <div class="sm:flex sm:items-start">
                <div
                    class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-gray-100 sm:mx-0 sm:h-10 sm:w-10"
                >
                    <!-- Heroicon name: outline/check-circle -->
                    <svg
                        class="h-6 w-6 text-green-600"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                    </svg>
                </div>
                <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                    <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">Feed added</h3>
                    <div class="mt-2">
                        <p class="text-sm text-gray-500">
                            The feed was added to your folders. Please be patient while we fetch it.
                        </p>
                    </div>
                </div>
            </div>
        </div>
        <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse rounded-b-xl">
            <button
                on:click={cancelButtonClicked}
                type="button"
                class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-green-600 text-base font-medium text-white hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:ml-3 sm:w-auto sm:text-sm"
            >
                Close
            </button>
            <button
                on:click={() => (card = Card.Add)}
                type="button"
                class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
            >
                Add another feed
            </button>
        </div>
    {:else}
        invalid card {card}
    {/if}
</div>
