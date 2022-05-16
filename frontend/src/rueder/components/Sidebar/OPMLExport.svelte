<script lang="ts">
    import { createEventDispatcher, getContext } from "svelte"
    import CopyToClipboard from "svelte-copy-to-clipboard"

    import type { FeedAPI } from "../../api/feed"
    import { getUsername } from "../../stores/session"

    import SimpleModal from "../ui/SimpleModal.svelte"
    import Button from "../ui/Button.svelte"
    import Spinner from "../ui/Spinner.svelte"

    import { OPML, Header, Body, Outline } from "../../helpers/opml"
    import { contextKey } from "../../helpers/constants"

    const feedAPI: FeedAPI = getContext(contextKey.feedAPI)

    const username = getUsername()
    const dispatch = createEventDispatcher()

    async function generateOPML(): Promise<string> {
        const folders = await feedAPI.GetFolders()

        let header = {
            title: `${username}'s rueder feeds`,
            ownerName: username,
            dateCreated: new Date(),
            dateModified: new Date(),
        }
        let feeds = []

        for (const folderIdx in folders) {
            const folder = folders[folderIdx]
            for (const feedIdx in folder.feeds) {
                const feed = folder.feeds[feedIdx]

                feeds.push({
                    text: feed.title,
                    title: feed.title,
                    type: "rss",
                    xmlUrl: feed.url,
                    htmlUrl: feed.site_url,
                })
            }
        }

        const head = new Header(header)

        let outlines: Outline[] = []
        for (const feed of feeds) {
            outlines.push(new Outline(feed))
        }

        const body = new Body({ outlines })
        const opml = new OPML({ head, body })
        return opml.toXML().toString()
    }
    const opmlOutput = generateOPML()
</script>

<SimpleModal on:close>
    <div class="bg-white text-black p-4 sm:p-6 sm:pb-4">
        <div class="flex items-center">
            <div
                class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-gray-100 sm:mx-0 sm:h-10 sm:w-10"
            >
                <!-- heroicon: outline/save -->
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-6 w-6 text-gray-800"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4"
                    />
                </svg>
            </div>
            <div class="text-center sm:mt-0 sm:ml-4 sm:text-left">
                <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">OPML Export</h3>
            </div>
        </div>
        <div class="pt-3 p-1 text-sm text-gray-500">
            <p>Download the list of your subscribed feeds.</p>
        </div>

        {#await opmlOutput}
            <Spinner />
        {:then opmlOutput}
            <div class="overflow-auto" style="max-height: 40vh">
                <pre
                    class="bg-gray-200 text-xs whitespace-pre-wrap break-all p-1">
                    {opmlOutput}
                </pre>
            </div>
        {:catch}
            <p class="italic">Something went wrong.</p>
        {/await}
    </div>

    <span slot="buttons">
        {#await opmlOutput}
            <Button mode="dark">...</Button>
        {:then opmlOutput}
            <CopyToClipboard text={opmlOutput} let:copy>
                <Button mode="dark" on:click={copy}>Copy to Clipboard</Button>
            </CopyToClipboard>
        {/await}
        <Button on:click={() => dispatch("close")}>Close</Button>
    </span>
</SimpleModal>
