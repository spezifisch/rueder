<script lang="ts">
    import { createEventDispatcher, getContext } from "svelte"

    import Dropzone from "svelte-file-dropzone"
    import sxml from "sxml"

    import Button from "../../ui/Button.svelte"
    import Spinner from "../../ui/Spinner.svelte"
    import OutlineCheckCircle from "../../ui/heroicons/outline-check-circle.svelte"
    import OutlineExclamationCircle from "../../ui/heroicons/outline-exclamation-circle.svelte"

    import type { FeedAPI } from "../../../api/feed"
    import type { AddFeedResponse } from "../../../api/types"
    import { Feed } from "../../../api/types"
    import { contextKey } from "../../../helpers/constants"

    const feedAPI: FeedAPI = getContext(contextKey.feedAPI)

    const dispatch = createEventDispatcher()

    function backButtonClicked() {
        dispatch("back")
    }

    function cancelButtonClicked() {
        dispatch("cancel")
    }

    enum Card {
        FileUpload = 0,
        DisplayResult,
    }
    let state = Card.FileUpload

    // count of invalid feed entries in opml
    let ignoredFeeds: number = 0
    // uploaded files
    let files: File[] = []
    // temporary feed objects for selection
    let importedTempFeeds: Feed[] = []
    // actually selected feed urls for import
    let selectedFeedIndices: number[] = []
    // import results
    let importResults: ImportResult[] = []
    class ImportResult {
        url: string
        ok: boolean
        feedID?: string
        errorMsg?: string
    }
    // true while import in progress
    let importResultsDone: boolean = false

    function sxmlSafeGetProperty(o: sxml.XML, prop: string): string | undefined {
        try {
            return o.getProperty(prop)
        } catch (e) {
            return
        }
    }

    // file upload handler. reads all files
    function handleFilesSelect(e: any) {
        clear()
        files = e.detail.acceptedFiles
        for (const file of files) {
            const reader = new FileReader()
            reader.onload = () => {
                try {
                    const content = reader.result.toString()
                    processOPML(content)
                } catch (e) {
                    console.log("caught", e)
                }
            }
            reader.readAsText(file)
        }
    }

    // parses the content of an OPML file and fills importedTempFeeds
    function processOPML(data: string) {
        let xml: sxml.XML = new sxml.XML(data)
        const body = xml.get("body").at(0)

        for (const outline of body.get("outline")) {
            const feedURL = sxmlSafeGetProperty(outline, "xmlUrl")
            if (!feedURL) {
                ignoredFeeds++
                continue
            }
            const title =
                sxmlSafeGetProperty(outline, "title") ??
                sxmlSafeGetProperty(outline, "text") ??
                sxmlSafeGetProperty(outline, "description") ??
                "unknown"

            let f = new Feed({
                url: feedURL,
                title,
            })
            importedTempFeeds = [...importedTempFeeds, f]
        }
    }

    // reset state before parsing new uploaded files
    function clear() {
        ignoredFeeds = 0
        files = []
        importedTempFeeds = []
        selectedFeedIndices = []
        importResults = []
    }

    function selectAllClicked() {
        selectedFeedIndices = [...importedTempFeeds.keys()].map((i) => i)
        importedTempFeeds = [...importedTempFeeds]
    }

    function deselectAllClicked() {
        selectedFeedIndices = []
        importedTempFeeds = importedTempFeeds
    }

    // adds all selected feeds to the backend
    async function importButtonClicked() {
        const urls = selectedFeedIndices.map((i) => importedTempFeeds[i].url)

        // show progress live
        state = Card.DisplayResult

        importResults = []
        importResultsDone = false
        for (const url of urls) {
            let result = new ImportResult()
            result.url = url

            try {
                let resp: AddFeedResponse = await feedAPI.AddFeed(url)

                result.ok = resp.ok

                if (resp.ok) {
                    result.feedID = resp.feed_id
                } else {
                    // failure adding due to network/server problems or existing feed or invalid url
                    result.errorMsg = resp.message
                }
            } catch (e) {
                // rejected promise: shouldn't happen
                result.ok = false
                result.errorMsg = e.message
            }

            importResults = [...importResults, result]
        }

        // deactivate spinner
        importResultsDone = true
    }

    function finishButtonClicked() {
        // add new feeds to a folder
        const addedFeeds = importResults.filter((result) => result.ok).map((result) => result.feedID)
        if (addedFeeds.length > 0) {
            dispatch("added", addedFeeds)
        } else {
            cancelButtonClicked()
        }
    }
</script>

<div class="bg-white px-2 pt-5 pb-4 sm:p-4 sm:pb-4 rounded-t-xl">
    <div class="flex items-center mb-3">
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
            <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                {#if state == Card.FileUpload}
                    Import Feeds
                {:else if state == Card.DisplayResult}
                    Feed Import Results
                {/if}
            </h3>
        </div>
    </div>

    {#if state == Card.FileUpload}
        <div class="p-1 text-sm text-gray-500">
            <p>
                You can upload an OPML file here to import feeds from. Most other feed services offer an OPML export
                option.
            </p>
            <p>All contained feeds will be imported into your last folder.</p>
            <p>You can edit the feed URLs before confirming the import.</p>
        </div>

        <div class="text-black">
            {#if files.length == 0}
                <div class="mt-3">
                    <Dropzone
                        on:drop={handleFilesSelect}
                        multiple={false}
                        accept=".opml,.xml"
                        disableDefaultStyles={true}
                        containerClasses="bg-gray-400 flex flex-col align-center text-center p-5 border-2 rounded-xl border-gray-700 border-dashed text-gray-900 outline-none"
                    >
                        <p>Click here or drag & drop to import an OPML file</p>
                    </Dropzone>
                </div>
            {:else}
                <div class="flex flex-row justify-between items-center">
                    <div>
                        {#each files as item}
                            <h2 class="text-lg">Feeds from {item.name}</h2>
                        {/each}
                    </div>
                    <div>
                        <Button on:click={clear} mode="error" addClass="text-md">Clear imported data</Button>
                    </div>
                </div>

                <div class="overflow-y-auto px-1 text-black" style="max-height: 40vh">
                    <div class="text-right">
                        {#if importedTempFeeds.length == selectedFeedIndices.length}
                            <p class="pt-2 px-2 underline cursor-pointer" on:click={deselectAllClicked}>Deselect all</p>
                        {:else}
                            <p class="pt-2 px-2 underline cursor-pointer" on:click={selectAllClicked}>Select all</p>
                        {/if}
                    </div>
                    {#each importedTempFeeds as feed, i}
                        <div class="flex flex-row items-center px-3 hover:bg-gray-200">
                            <div class="flex-auto px-2 py-1 text-sm overflow-y-auto">
                                {feed.title}
                                <input
                                    class="text-xs p-0 m-0 w-full border-none font-mono"
                                    type="text"
                                    bind:value={feed.url}
                                />
                            </div>
                            <div class="flex-none">
                                <!-- from https://tailwindcomponents.com/component/checkbox -->
                                <label class="inline-flex items-center">
                                    <input
                                        type="checkbox"
                                        class="form-checkbox rounded-sm my-2 h-5 w-5 text-green-600 border-gray-300 focus:ring-gray-600"
                                        bind:group={selectedFeedIndices}
                                        value={i}
                                    />
                                </label>
                            </div>
                        </div>
                    {:else}
                        <div class="text-sm p-2">
                            <p class="text-green-700">There are no feeds in this file.</p>
                            {#if ignoredFeeds > 0}
                                <p class="p-2">{ignoredFeeds} invalid entries were ignored.</p>
                            {/if}
                        </div>
                    {/each}
                </div>
            {/if}
        </div>
    {:else if state == Card.DisplayResult}
        <div class="overflow-y-auto px-1 text-black" style="max-height: 40vh">
            {#each importResults as result, i}
                <div class="flex flex-row items-center max-w-full px-3 text-sm hover:bg-gray-200">
                    <div class="flex-none px-2 py-1">
                        <p>{i + 1}.</p>
                    </div>
                    <div class="flex-auto flex flex-col py-1 text-sm overflow-y-auto overflow-x-auto">
                        <input class="text-xs p-0 m-0 w-full border-none font-mono" type="text" value={result.url} />
                        {#if result.ok}
                            <div class="font-mono text-xs">
                                Feed ID: {result.feedID}
                            </div>
                        {:else}
                            <div>
                                <p class="text-red-800 text-xs">{result.errorMsg}</p>
                            </div>
                        {/if}
                    </div>
                    <div class="flex-none px-2 py-1">
                        {#if result.ok}
                            <div class="text-green-700">
                                <OutlineCheckCircle />
                            </div>
                        {:else}
                            <div class="text-red-700">
                                <OutlineExclamationCircle />
                            </div>
                        {/if}
                    </div>
                </div>
            {/each}
            {#if !importResultsDone}
                <div class="flex justify-center min-w-0">
                    <Spinner />
                </div>
            {/if}
        </div>
    {:else}
        <p>invalid card {state}</p>
    {/if}
</div>
<div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse rounded-b-xl">
    {#if state == Card.FileUpload}
        <button
            on:click={importButtonClicked}
            type="button"
            disabled={selectedFeedIndices.length == 0}
            class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 text-base font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:ml-3 sm:w-auto sm:text-sm {selectedFeedIndices.length >
            0
                ? 'bg-green-600 hover:bg-green-700'
                : 'bg-gray-600'}"
        >
            Import selected Feed{selectedFeedIndices.length > 1 ? "s" : ""}
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
    {:else if state == Card.DisplayResult}
        <button
            on:click={finishButtonClicked}
            type="button"
            class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 text-base font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:ml-3 sm:w-auto sm:text-sm bg-green-600 hover:bg-green-700"
        >
            Finish
        </button>
    {/if}
</div>
