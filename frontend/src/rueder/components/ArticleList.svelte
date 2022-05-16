<script lang="ts">
    import { createEventDispatcher, getContext } from "svelte"
    import type { Writable } from "svelte/store"
    import InfiniteScroll from "svelte-infinite-scroll"
    import { scrollToTop } from "svelte-scrollto"

    import CloseButton from "./ui/heroicons/outline-x-circle.svelte"
    import OutlineArrowCircleLeft from "./ui/heroicons/outline-arrow-circle-left.svelte"
    import RefreshButton from "./ui/heroicons/outline-refresh.svelte"
    import SmallButton from "./ui/SmallButton.svelte"
    import Spinner from "./ui/Spinner.svelte"

    import ArticleListItem from "./ArticleListItem.svelte"
    import FeedInfo from "./FeedInfo.svelte"

    import type { ArticlePreview, Feed } from "../api/types"
    import type { FeedAPI } from "../api/feed"
    import type { FeedStateStore } from "../stores/articlestate"
    import { contextKey } from "../helpers/constants"

    export let feedStateStore: Writable<FeedStateStore>
    export let selectedFeed: string
    export let hotkeyPressed: string

    const feedAPI: FeedAPI = getContext(contextKey.feedAPI)

    const dispatch = createEventDispatcher()

    // feed info header
    let feedInfo: Promise<Feed>

    // sequence number of last (oldest) article in list.
    // this is used to get the next batch of articles.
    let nextToken: number
    function loadFeed(_feedID: string) {
        feedInfo = feedAPI.GetFeed(selectedFeed)

        if (scrollContainer) {
            scrollToTop({
                container: scrollContainer,
                duration: 0,
            })
        }

        nextToken = 0
        articleList = []
        newBatch = []
        loadMore()
    }
    $: loadFeed(selectedFeed)

    // update articleList when newBatch is written to
    let articleList: ArticlePreview[] = []
    let newBatch: ArticlePreview[] = []
    $: articleList = [...articleList, ...newBatch]
    let loadingArticles = true

    // latest article sequence number for "mark all read"
    let latestFeedSeq: number
    $: if (articleList && articleList.length > 0) {
        latestFeedSeq = articleList[0].feed_seq
    }

    /// feedstate store work
    // when the button is clicked
    function markAllRead() {
        console.log("marking all articles read in", selectedFeed, "until", latestFeedSeq)
        feedStateStore.update((fss) => {
            fss.markAllReadUntil(selectedFeed, latestFeedSeq)
            return fss
        })
    }

    // update article count when feedinfo is updated
    $: updateArticleCount(feedInfo)
    async function updateArticleCount(feedInfo: Promise<Feed>) {
        const fi = await feedInfo

        const oldArticleCount = $feedStateStore.getTotalArticles(fi.id)
        const newArticleCount = fi.article_count
        if (oldArticleCount != newArticleCount) {
            console.log(`feed ${fi.id} article count updated from ${oldArticleCount} to ${newArticleCount}`)
            feedStateStore.update((fss) => {
                fss.setTotalArticles(selectedFeed, newArticleCount)
                return fss
            })
        }
    }

    $: broadcastFeedInfo(feedInfo)
    async function broadcastFeedInfo(feedInfo: Promise<Feed>) {
        const fi = await feedInfo
        dispatch("feedInfo", fi)
    }

    // called initially and when infinitescroll triggers when we scroll down
    async function loadMore() {
        loadingArticles = true
        // this returns the next 40 articles or so
        newBatch = await feedAPI.GetArticles(selectedFeed, nextToken)
        if (newBatch.length) {
            // get sequence number of last article
            nextToken = newBatch.slice(-1)[0].seq
        }
        loadingArticles = false
    }

    let scrollContainer: any

    function closeFeed() {
        dispatch("close")
    }

    let animateRefreshButton = false
    function refreshFeed() {
        animateRefreshButton = true
        setTimeout(() => (animateRefreshButton = false), 500)
        loadFeed(selectedFeed)
    }

    $: if (hotkeyPressed) {
        switch (hotkeyPressed) {
            case "t":
                console.log("scroll to top in article list hotkey pressed")
                if (scrollContainer) {
                    scrollToTop({
                        container: scrollContainer,
                    })
                }
                break
        }
    }
</script>

<ul class="flex flex-col h-full relative overflow-y-scroll rueder-scrollbar pr-2" bind:this={scrollContainer}>
    <div class="flex flex-row space-x-4 md:hidden bg-gray-900">
        <div title="Back to folder list" on:click={closeFeed} class="text-gray-500 cursor-pointer p-1">
            <OutlineArrowCircleLeft size={8} />
        </div>
        <div class="flex-auto" />
        <div
            title="Reload Feed"
            on:click={refreshFeed}
            class="text-gray-500 cursor-pointer p-1"
            class:animate-spin={animateRefreshButton}
        >
            <RefreshButton size={8} />
        </div>
    </div>
    <div class="hidden absolute top-2 right-2 md:flex md:flex-row md:flex-wrap">
        <div
            title="Reload Feed"
            on:click={refreshFeed}
            class="text-gray-500 hover:text-gray-200 cursor-pointer bg-gray-800 rounded-full"
            class:animate-spin={animateRefreshButton}
        >
            <RefreshButton />
        </div>
        <div
            title="Close Feed"
            on:click={closeFeed}
            class="text-gray-500 hover:text-gray-200 cursor-pointer bg-gray-800 rounded-full"
        >
            <CloseButton />
        </div>
    </div>

    <li>
        <FeedInfo {feedInfo} feedID={selectedFeed} />
    </li>
    {#if articleList && articleList.length}
        <li class="flex-none flex flex-row justify-end px-2 py-1">
            <SmallButton addClass="flex-none" on:click={markAllRead}>Mark All Read</SmallButton>
        </li>
    {/if}

    {#each articleList as item}
        <hr class="border-gray-400 mx-2" />
        <ArticleListItem {feedStateStore} feedID={selectedFeed} article={item} on:articleClick />
    {/each}

    {#if loadingArticles}
        <div class="flex-none flex p-4 justify-center">
            <Spinner />
        </div>
    {/if}

    <hr class="border-gray-400 mx-2" />
    <InfiniteScroll hasMore={articleList.length == 0 || newBatch.length > 0} on:loadMore={loadMore} />

    {#if !newBatch.length}
        <li class="italic p-2 text-gray-400">End of articles.</li>
    {/if}
</ul>
