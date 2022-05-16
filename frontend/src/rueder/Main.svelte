<script lang="ts">
    import { onDestroy, setContext } from "svelte"

    import Hamburger from "./components/ui/Hamburger.svelte"
    import ModalWrapper from "./components/ui/ModalWrapper.svelte"
    import ToastError from "./components/ui/ToastError.svelte"

    import ArticleList from "./components/ArticleList.svelte"
    import ArticleContent from "./components/ArticleContent.svelte"
    import LabelArticles from "./components/LabelArticles.svelte"
    import SidebarNavigation from "./components/Sidebar/SidebarNavigation.svelte"
    import SidebarMenu from "./components/Sidebar/SidebarMenu.svelte"

    import { FeedAPI } from "./api/feed"
    import type { ArticleStore } from "./api/feed"
    import type { Feed } from "./api/types"
    import { getUsername } from "./stores/session"
    import { feedStateStore } from "./stores/articlestate"
    import url from "./stores/url"
    import { labelStore } from "./stores/labels"
    import { contextKey } from "./helpers/constants"
    import { ImageProxy } from "./helpers/ImageProxy"

    export let imageProxyBaseURL: string
    export let imageProxyUseTypePrefixes: boolean
    export let imageProxyKey: string
    export let imageProxySalt: string

    // backend communication
    export let baseURL: string
    let feedAPI = new FeedAPI(baseURL)

    // imageproxy wrapper
    let imageProxy = new ImageProxy(imageProxyBaseURL, imageProxyUseTypePrefixes, imageProxyKey, imageProxySalt)

    // context
    setContext(contextKey.feedAPI, feedAPI)
    setContext(contextKey.imageProxy, imageProxy)

    // stores!
    let article_content: ArticleStore

    let selectedFeed: string
    let selectedLabel: string

    // sidebar menu overlay
    let sidebarMenuOpen = false
    let sidebarMenuClosed: boolean
    $: sidebarMenuClosed = !sidebarMenuOpen

    // show an error toast if needed
    let doShowErrorToast = false
    let errorToastMessage: string
    function showErrorToast(message: string) {
        doShowErrorToast = true
        errorToastMessage = message
    }

    if ($url.hash) {
        loadFromHash($url.hash)
    }

    const unsubscribe = url.subscribe((url) => {
        if (!url.hash) {
            return
        }
        loadFromHash(url.hash)
    })
    onDestroy(unsubscribe)

    // parse deeplinks like "#/feed/123-456" and "#/article/123-456"
    function loadFromHash(hash: string) {
        const reUUID = new RegExp("^[a-fA-F0-9-]+$")
        const split = hash.split("/")
        if (split.length < 3) {
            console.log("invalid hash", hash)
            showErrorToast("invalid deep link")
            return
        }

        const kind = split[1]
        let id = split[2]
        if (!["article", "feed", "label"].includes(kind)) {
            console.log("invalid kind from hash:", kind)
            showErrorToast("invalid deep link")
            return
        }

        if ((kind == "article" || kind == "feed") && !reUUID.test(id)) {
            console.log("invalid uuid from hash:", id)
            showErrorToast("invalid deep link")
            return
        }

        if (kind == "label") {
            id = decodeURI(id)
            if (!labelStore.labelExists(id)) {
                console.log("non-existent label name from hash:", id)
                showErrorToast("label does not exist")
                return
            }
        }

        switch (kind) {
            case "article":
                _loadArticle(id)
                break
            case "feed":
                _loadFeed(id)
                break
            case "label":
                _loadLabel(id)
                break
        }
    }

    // set hash for deeplinks when feed/article/label is clicked
    function updateHash(kind?: string, id?: string) {
        const hash = `/${kind}/${id}`
        if ($url.hash != hash) {
            window.location.hash = hash
        }
    }

    // remove hash when everything's closed
    function removeHash() {
        // source: https://stackoverflow.com/a/5298684
        history.pushState("", document.title, window.location.pathname + window.location.search)
    }

    /// Feed Functions
    function _loadFeed(uuid: string) {
        selectedLabel = null
        selectedFeed = uuid
    }

    // called when a feed is clicked in the folder list
    function loadFeed(uuid: string) {
        updateHash("feed", uuid)
        _loadFeed(uuid)
    }

    async function closeFeed() {
        selectedFeed = null

        if (article_content) {
            // show article url if one's open
            updateHash("article", (await $article_content).id)
        } else {
            removeHash()
        }
    }

    /// Article Functions
    function _loadArticle(uuid: string) {
        article_content = feedAPI.GetArticle(uuid)
    }

    // called when an article is clicked in the feed list
    function loadArticle(uuid: string) {
        updateHash("article", uuid)
        _loadArticle(uuid)
    }

    function closeArticle() {
        article_content = null

        if (selectedFeed) {
            // show feed url if one's open
            updateHash("feed", selectedFeed)
        } else if (selectedLabel) {
            updateHash("label", selectedLabel)
        } else {
            removeHash()
        }
    }

    /// Label Functions
    function _loadLabel(name: string) {
        selectedFeed = null
        selectedLabel = name
    }

    // called when a label is clicked in the label list
    function loadLabel(name: string) {
        updateHash("label", name)
        _loadLabel(name)
    }

    async function closeLabel() {
        selectedLabel = null

        if (article_content) {
            // show article url if one's open
            updateHash("article", (await $article_content).id)
        } else {
            removeHash()
        }
    }

    function closeAll() {
        selectedFeed = null
        selectedLabel = null
        article_content = null
        removeHash()
    }

    // for the top bar
    const username = getUsername()

    // send feed info update to folder list when feed view is updated
    let feedInfoUpdate: Feed
    async function broadcastFeedInfo(feedInfo: Feed) {
        feedInfoUpdate = feedInfo
    }

    // keyboard shortcuts
    let isDialogOpen = false
    let articleEditorIsFocused = false
    let hotkeyPressed: string
    function handleKeydown(event: KeyboardEvent) {
        // close everything on Esc
        if (event.key == "Escape") {
            event.preventDefault()
            doShowErrorToast = false
            sidebarMenuOpen = false
            hotkeyPressed = event.key
        }
        if (isDialogOpen || articleEditorIsFocused) {
            // don't intercept keyboard chars when we have input fields
            return
        }

        switch (event.key) {
            case "r": // mark all read
                if (selectedFeed) {
                    event.preventDefault()
                    feedStateStore.update((fss) => {
                        fss.markAllRead(selectedFeed)
                        return fss
                    })
                }
            case "u": // update folder list
            case "t": // scroll to top in article list
                event.preventDefault()
                hotkeyPressed = event.key
                break
        }
    }

    function handleKeyup(_event: KeyboardEvent) {
        if (hotkeyPressed) {
            hotkeyPressed = null
        }
    }

    /// responsive design toggles
    let forceShowSidebar = false
    let forceShowFeed = false
    let forceShowLabelList = false
</script>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup} />

<main class="flex flex-row h-screen text-gray-50 bg-black">
    <!-- show error toast -->
    {#if doShowErrorToast}
        <div
            on:click={() => {
                doShowErrorToast = false
            }}
        >
            <ModalWrapper
                on:click-background={() => {
                    doShowErrorToast = false
                }}
            >
                <ToastError message={errorToastMessage} />
            </ModalWrapper>
        </div>
    {/if}

    <!-- sidebar -->
    <div
        class="flex-none w-screen h-screen flex flex-col bg-gray-900 md:w-64 xl:flex"
        class:hidden={!forceShowSidebar && (selectedFeed || selectedLabel || article_content)}
    >
        <!-- top bar -->
        <div class="flex-none flex flex-row items-center bg-gray-800 border-b-4 border-gray-500">
            <Hamburger bind:open={sidebarMenuOpen} />
            <h1 class="flex-auto px-2 text-right text-md font-light">{username}:rueder</h1>
        </div>
        <!-- menu -->
        <SidebarMenu bind:show={sidebarMenuOpen} />
        <!-- folders, labels, buttons -->
        <SidebarNavigation
            {feedStateStore}
            {selectedFeed}
            {selectedLabel}
            {feedInfoUpdate}
            {hotkeyPressed}
            on:feedClick={(e) => loadFeed(e.detail)}
            on:labelClick={(e) => loadLabel(e.detail)}
            on:closeFeed={closeFeed}
            on:closeArticle={closeArticle}
            on:closeLabel={closeLabel}
            on:closeAll={closeAll}
            on:error={(e) => showErrorToast(e.detail)}
            bind:show={sidebarMenuClosed}
            bind:isDialogOpen
        />
    </div>

    <!-- feed view aka article list -->
    {#if selectedFeed}
        <div
            class="flex-none w-screen md:block bg-gray-700 md:w-2/5 h-full"
            class:hidden={!forceShowFeed && article_content}
        >
            <ArticleList
                {selectedFeed}
                {feedStateStore}
                {hotkeyPressed}
                on:feedInfo={(e) => broadcastFeedInfo(e.detail)}
                on:articleClick={(e) => loadArticle(e.detail)}
                on:close={closeFeed}
            />
        </div>
    {/if}

    <!-- label list -->
    {#if selectedLabel}
        <div
            class="flex-none w-screen md:block bg-gray-700 md:w-2/5 h-full"
            class:hidden={!forceShowLabelList && article_content}
        >
            <LabelArticles {selectedLabel} on:articleClick={(e) => loadArticle(e.detail)} on:close={closeLabel} />
        </div>
    {/if}

    <!-- article content -->
    {#if article_content}
        <ArticleContent {article_content} on:close={closeArticle} bind:isFocused={articleEditorIsFocused} />
    {/if}
</main>
