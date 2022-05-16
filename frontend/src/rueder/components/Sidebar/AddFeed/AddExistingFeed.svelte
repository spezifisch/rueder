<script lang="ts">
    import { createEventDispatcher, getContext } from "svelte"

    import Spinner from "../../ui/Spinner.svelte"
    import OutlineRSSIcon from "../../ui/heroicons/outline-rss.svelte"

    import type { FeedAPI } from "../../../api/feed"
    import type { Feed } from "../../../api/types"
    import { contextKey } from "../../../helpers/constants"
    import { ImageProxy, ImageProxyType } from "../../../helpers/ImageProxy"

    export let selectedFeedIDs: string[] = []

    const feedAPI: FeedAPI = getContext(contextKey.feedAPI)
    let imageProxy: ImageProxy = getContext(contextKey.imageProxy)

    const dispatch = createEventDispatcher()

    const dummyIcon = "dummy.png"

    async function getSubscribedFeedIDs(): Promise<string[]> {
        // well, i guess we get all the folders and extract the list of feed ids
        const folders = await feedAPI.GetFolders()
        let feedIDs: string[] = []

        for (const folder of folders) {
            if (folder.feeds && folder.feeds.length > 0) {
                feedIDs = [...feedIDs, ...folder.feeds.map((feed) => feed.id)]
            }
        }
        return feedIDs
    }

    let feeds = new Promise<Feed[]>(async (set) => {
        // take all the feeds the site has and remove the ones we are already subscribed to
        const allFeeds = await feedAPI.GetFeeds()
        const subscribedFeedIDs = await getSubscribedFeedIDs()
        set(allFeeds.filter((feed) => feed.title && !subscribedFeedIDs.includes(feed.id)))
    })

    // disable add button
    let disableAdd: boolean
    $: disableAdd = selectedFeedIDs.length == 0

    function addExistingFeeds() {
        dispatch("added", selectedFeedIDs)
    }

    function backButtonClicked() {
        dispatch("back")
    }

    function cancelButtonClicked() {
        dispatch("cancel")
    }
</script>

<div class="flex-auto flex flex-col overflow-y-auto bg-white rounded-xl">
    <div class="px-4 pt-5 pb-4 sm:p-6 sm:pb-4 flex-auto flex flex-col overflow-y-auto">
        <!-- header -->
        <div class="flex-none flex items-center">
            <!-- icon bubble -->
            <div
                class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full sm:mx-0 sm:h-10 sm:w-10 bg-gray-100 text-gray-800"
            >
                <OutlineRSSIcon />
            </div>
            <!-- modal title -->
            <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                    Subscribe to existing feeds
                </h3>
            </div>
        </div>

        <!-- description -->
        <div class="flex-none pt-3 p-1 text-sm text-gray-500">
            <p>Subscribe to some other feeds in our database.</p>
        </div>

        <!-- list of available feeds -->
        {#await feeds}
            <div class="flex-none flex justify-center p-4">
                <Spinner />
            </div>
        {:then feeds}
            <div class="flex-auto overflow-y-auto overscroll-contain px-1 text-black">
                {#each feeds as feed (feed.id)}
                    <div class="flex flex-row items-center px-3 hover:bg-gray-200">
                        <div class="flex-none">
                            <img
                                src={feed.icon ? imageProxy.buildURL(ImageProxyType.Icon, feed.icon) : dummyIcon}
                                alt={feed.title}
                                class="object-contain h-5 w-5"
                            />
                        </div>
                        <div class="flex-auto px-2 py-1 text-sm overflow-y-auto">
                            {feed.title}
                            <pre class="text-xs">{feed.url}</pre>
                        </div>
                        <div class="flex-none">
                            <!-- from https://tailwindcomponents.com/component/checkbox -->
                            <label class="inline-flex items-center">
                                <input
                                    type="checkbox"
                                    class="form-checkbox rounded-sm my-2 h-5 w-5 text-green-600 border-gray-300 focus:ring-gray-600"
                                    bind:group={selectedFeedIDs}
                                    value={feed.id}
                                />
                            </label>
                        </div>
                    </div>
                {:else}
                    <p class="text-green-700 text-sm p-2">Looks like you're subscribed to all the feeds we have.</p>
                {/each}
            </div>
        {:catch}
            <div class="text-red-700 text-center italic">
                <p>failed fetching feed list</p>
            </div>
        {/await}
    </div>

    <!-- bottom buttons -->
    <div class="flex-none bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
        <button
            on:click={addExistingFeeds}
            disabled={disableAdd}
            type="button"
            class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 {disableAdd
                ? 'bg-gray-600'
                : 'bg-green-600 hover:bg-green-700'} text-base font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:ml-3 sm:w-auto sm:text-sm"
        >
            Add existing feed{selectedFeedIDs.length == 1 ? "" : "s"}
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
</div>
