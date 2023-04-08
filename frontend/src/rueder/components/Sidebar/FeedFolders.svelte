<script lang="ts">
    import { createEventDispatcher, getContext, onDestroy, tick } from "svelte"
    import type { Writable } from "svelte/store"

    import { dndzone, TRIGGERS, SOURCES } from "@spezifisch/svelte-dnd-action"
    import { scrollToBottom } from "svelte-scrollto"

    import { v4 as uuidv4 } from "uuid"

    import Spinner from "../ui/Spinner.svelte"
    import ToastError from "../ui/ToastError.svelte"
    import DragHandle from "../ui/icons/DragHandle.svelte"
    import OutlineCircleRight from "../ui/heroicons/outline-circle-right.svelte"
    import RefreshButton from "../ui/heroicons/outline-refresh.svelte"
    import SolidFolder from "../ui/heroicons/solid-folder.svelte"
    import SolidFolderOpen from "../ui/heroicons/solid-folder-open.svelte"
    import SolidPlusCircle from "../ui/heroicons/solid-plus-circle.svelte"
    import SolidTrash from "../ui/heroicons/solid-trash.svelte"

    import type { FeedStateStore } from "../../stores/articlestate"

    import { Feed, Folder } from "../../api/types"
    import type { FeedAPI } from "../../api/feed"
    import type { SSEStore } from "../../stores/SSE"
    import { SSEMessageType } from "../../api/sse"

    import EditFolderTitle from "./Folders/EditFolderTitle.svelte"
    import { contextKey } from "../../helpers/constants"
    import { ImageProxy, ImageProxyType } from "../../helpers/ImageProxy"

    export let feedStateStore: Writable<FeedStateStore>
    export let selectedFeed: string
    export let hotkeyPressed: string
    export let feedInfoUpdate: Feed

    export let editMode: boolean
    export let addNewFeeds: Writable<string[]>

    const feedAPI: FeedAPI = getContext(contextKey.feedAPI)
    const sseEvents: SSEStore = getContext(contextKey.sseEvents)
    let imageProxy: ImageProxy = getContext(contextKey.imageProxy)

    const dispatch = createEventDispatcher()

    const dummyIcon = "dummy.png"

    // set to true when a drag handle is grabbed
    let draggingFolderByHandle = false
    let draggingFeedByHandle = false
    // always disable dnd when we're not in edit mode.
    // if we're in edit mode enable dragging when a drag handle is touched
    $: dragFolderDisabled = !editMode || !draggingFolderByHandle
    $: dragFeedDisabled = !editMode || !draggingFeedByHandle

    enum DragType {
        Folder,
        Feed,
    }
    // drag handle handling for drag&drop of feeds/folders in edit mode.
    // these two functions are based on: https://svelte.dev/repl/4949485c5a8f46e7bdbeb73ed565a9c7?version=3.24.1
    function handleDragStart(type: DragType, e: TouchEvent | MouseEvent) {
        // preventing default to prevent lag on touch devices (because of the browser checking for screen scrolling)
        e.preventDefault()

        if (type == DragType.Folder) {
            draggingFolderByHandle = true
        } else if (type == DragType.Feed) {
            draggingFeedByHandle = true
        }
    }

    function handleDragKeyDown(type: DragType, e: KeyboardEvent) {
        if (e.key == "Enter" || e.key == " ") {
            if (type == DragType.Folder && !draggingFolderByHandle) {
                draggingFolderByHandle = true
            } else if (type == DragType.Feed && !draggingFeedByHandle) {
                draggingFeedByHandle = true
            }
        }
    }

    // folder list
    let loadingFolders = true // show spinner
    let loadingFoldersOk = true // false=show error toast
    let folders: Folder[]
    let feedExists = new Map<string, boolean>() // feed id -> true

    let reloading = false // only used to ensure only one folder request is running at a time
    async function loadFolders() {
        if (reloading) {
            return
        }
        reloading = true
        try {
            // fetch folder list
            let newFolders = await feedAPI.GetFolders()

            // clear states
            feedExists.clear()
            feedExists = new Map<string, boolean>()
            let folderExists = new Map<string, boolean>()

            // cleanup and init states
            for (let i = 0; i < newFolders.length; i++) {
                // ensure uuids are unique
                if (folderExists[newFolders[i].id]) {
                    const oldID = newFolders[i].id
                    newFolders[i].id = uuidv4()
                    console.log("duplicate folder id", oldID, "->", newFolders[i].id)
                }
                folderExists[newFolders[i].id] = true

                if (!newFolders[i].feeds) {
                    // it's "omitempty" on server side. use an empty array for easier templating
                    newFolders[i].feeds = []
                }
                if (newFolders[i].open !== true && newFolders[i].open !== false) {
                    // open folders by default, showing their feeds
                    newFolders[i].open = true
                }

                for (let j = 0; j < newFolders[i].feeds.length; j++) {
                    const feed = newFolders[i].feeds[j]
                    if (feed.id.startsWith("id:dnd-shadow-placeholder")) {
                        // this thing is created temporarily by svelte-dnd-action and should not end up in the state list.
                        // sometimes it does during development though. remove it.
                        console.log("removed placeholder feed in folder")
                        newFolders[i].feeds[j] = null
                        continue
                    }
                    if (feedExists[feed.id]) {
                        console.log("removed duplicate feed", feed.id)
                        newFolders[i].feeds[j] = null
                        continue
                    }

                    // add feed to map for quick lookup whether it exists
                    feedExists[feed.id] = true

                    // add to feed states
                    $feedStateStore.createFeedState(feed.id)
                    $feedStateStore.feeds[feed.id].setTotalArticles(feed.article_count)
                }

                newFolders[i].feeds = newFolders[i].feeds.filter((feed) => feed)
            }

            folders = newFolders
            loadingFolders = false
        } catch (e) {
            console.log("dnd folders failed", e)
            loadingFoldersOk = false
            loadingFolders = false
        }

        setTimeout(() => {
            animateRefreshButton = false
            hideRefreshButton = true
            reloading = false
        }, 500)
    }
    loadFolders()

    $: if (hotkeyPressed) {
        switch (hotkeyPressed) {
            case "u":
                console.log("reload folders hotkey pressed")
                animateRefreshButton = true
                hideRefreshButton = false
                loadFolders()
                break
        }
    }

    async function commitChangedFolders(): Promise<void> {
        dirty = false
        try {
            await feedAPI.ChangeFolders(folders)
            console.log("ChangeFolders done")
        } catch (resp) {
            let msg: string = "" + resp.message
            if (msg && msg.startsWith("JSON.parse:")) {
                msg = "Temporary server issue."
            }
            dispatch("error", "Updating folders failed: " + msg)

            // roll back changes
            await loadFolders()
        }
    }

    $: handleFeedInfoUpdate(feedInfoUpdate)
    function handleFeedInfoUpdate(feedInfoUpdate: Feed) {
        if (!feedInfoUpdate) {
            return
        }

        const feedID = feedInfoUpdate.id
        if (!feedExists[feedID]) {
            // not subscribed to the opened feed
            return
        }

        const folderIdx = folders.findIndex(
            (folder) => folder.feeds && folder.feeds.findIndex((feed) => feed.id == feedID) >= 0
        )
        if (folderIdx < 0) {
            return
        }
        const feedIdx = folders[folderIdx].feeds.findIndex((feed) => feed.id == feedID)
        if (feedIdx < 0) {
            return
        }

        // update folder view
        if (feedInfoUpdate.title) {
            folders[folderIdx].feeds[feedIdx].title = feedInfoUpdate.title
        }
        if (feedInfoUpdate.icon) {
            folders[folderIdx].feeds[feedIdx].icon = feedInfoUpdate.icon
        }
        if (feedInfoUpdate.url) {
            folders[folderIdx].feeds[feedIdx].url = feedInfoUpdate.url
        }
        if (feedInfoUpdate.site_url) {
            folders[folderIdx].feeds[feedIdx].url = feedInfoUpdate.site_url
        }
    }

    function getFolderIndex(folderID: string): number {
        return folders.findIndex((f) => f.id == folderID)
    }

    function handleIconError(e: Event) {
        if ((e.target as HTMLImageElement).src != dummyIcon) {
            ;(e.target as HTMLImageElement).src = dummyIcon
        }
    }

    // dnd folder handlers
    function handleDndConsiderFolders(e: CustomEvent<DndEvent>) {
        folders = [...e.detail.items] as Folder[]

        // Ensure dragging is stopped on drag finish via keyboard (code from drag example)
        const {
            info: { source, trigger },
        } = e.detail
        if (source === SOURCES.KEYBOARD && trigger === TRIGGERS.DRAG_STOPPED) {
            draggingFolderByHandle = false
        }
    }

    async function handleDndFinalizeFolders(e: CustomEvent<DndEvent>) {
        // show changes optimistically
        folders = [...e.detail.items] as Folder[]

        // Ensure dragging is stopped on drag finish via pointer (mouse, touch)
        if (e.detail.info.source === SOURCES.POINTER) {
            draggingFolderByHandle = false
        }

        if (e.detail.info.trigger == TRIGGERS.DROPPED_INTO_ZONE) {
            await commitChangedFolders()
        }
    }

    // dnd feed handlers
    function handleDndConsiderFeeds(folderID: string, e: CustomEvent<DndEvent>) {
        const colIdx = getFolderIndex(folderID)
        folders[colIdx].feeds = e.detail.items as Feed[]
        folders = [...folders]

        // Ensure dragging is stopped on drag finish via keyboard (code from drag example)
        const {
            info: { source, trigger },
        } = e.detail
        if (source === SOURCES.KEYBOARD && trigger === TRIGGERS.DRAG_STOPPED) {
            draggingFeedByHandle = false
        }
    }

    async function handleDndFinalizeFeeds(folderID: string, e: CustomEvent<DndEvent>) {
        // show changes optimistically
        const colIdx = getFolderIndex(folderID)
        folders[colIdx].feeds = e.detail.items as Feed[]
        folders = [...folders]

        // Ensure dragging is stopped on drag finish via pointer (mouse, touch)
        if (e.detail.info.source === SOURCES.POINTER) {
            draggingFeedByHandle = false
        }

        // when a feed was moved to another folder we get 2 events (dropped into zone, dropped into another).
        // when it is moved inside the same folder we get only 1 (dropped into zone).
        // the changes to the DOM are already made by the consider function, so we just trigger on the first event and get all the changes.
        if (e.detail.info.trigger == TRIGGERS.DROPPED_INTO_ZONE) {
            await commitChangedFolders()
        }
    }

    // delete folder
    async function deleteFolder(folderID: string, _e: MouseEvent) {
        const colIdx = getFolderIndex(folderID)
        if (folders[colIdx].feeds.length != 0) {
            console.log("tried to delete non-empty folder", folderID)
            return
        }

        folders = folders.filter((f) => f.id != folderID)

        await commitChangedFolders()
    }

    // delete feed
    async function deleteFeed(folderID: string, feedID: string, _e: MouseEvent) {
        const colIdx = getFolderIndex(folderID)
        folders[colIdx].feeds = folders[colIdx].feeds.filter((f) => f.id != feedID)
        folders = [...folders]

        feedExists[feedID] = false
        await commitChangedFolders()
    }

    // add folder
    let newFolderName: string
    let newFolderCount = 0
    async function addFolder(name: string) {
        newFolderCount++

        if (name == null) {
            name = `Unnamed Folder ${newFolderCount}`
        }
        console.log("adding new folder", name, newFolderCount)

        folders[folders.length] = new Folder({
            id: uuidv4(),
            title: name,
            feeds: [],
        })

        // clear input
        newFolderName = null

        // scroll to bottom of folder list where the new folder is
        if (scrollContainer) {
            await tick()
            scrollToBottom({
                container: scrollContainer,
                offset: 200,
            })
        }

        // commit this empty folder only when leaving edit mode
        dirty = true
    }

    // store update handler for newly subscribed feeds
    const unsubscribe = addNewFeeds.subscribe(async (feedIDs: string[]) => {
        if (feedIDs.length == 0) {
            return // this is triggered once at startup
        }

        // add a new folder if we don't have folders
        if (folders.length == 0) {
            addFolder("Default")
        }

        // add new feeds to last folder
        const targetFolder = folders[folders.length - 1]

        console.log(`FeedFolders adding new feeds to folder ${targetFolder.id}`, feedIDs)

        let newFeeds: Feed[] = []
        for (const feedID of feedIDs) {
            if (feedExists[feedID]) {
                console.log("tried to add already subscribed feed", feedID)
                continue
            }

            const feed = new Feed()
            feed.id = feedID

            newFeeds = [...newFeeds, feed]
            feedExists[feedID] = true
        }

        folders[folders.length - 1].feeds = [...targetFolder.feeds, ...newFeeds]

        // send to backend
        await commitChangedFolders()

        // scroll to bottom of folder list where the new feed is
        if (scrollContainer) {
            await tick()
            scrollToBottom({
                container: scrollContainer,
                offset: 200,
            })
        }

        // open newly added feed if only one was added
        if (!editMode && newFeeds.length == 1) {
            dispatch("feedClick", newFeeds[0].id)
        }
    })
    onDestroy(unsubscribe)

    // SSE event handler
    const sseUnsubscribe = sseEvents.store.subscribe(async (sseEvent) => {
        if (!sseEvent) {
            return
        }
        if (sseEvent.message_type == SSEMessageType.FolderUpdate) {
            // HACKY just assume that only one client can be in edit mode
            if (!editMode) {
                console.log("triggering folder reload via sse")
                await loadFolders()
            }
        }
    })
    onDestroy(sseUnsubscribe)

    // true when folder names are changed which are not yet committed
    let dirty = false

    // commit changes to folders when leaving edit mode
    $: if (!editMode) {
        if (dirty) {
            commitChangedFolders()
        }
    }

    // this shows a spinning refresh button as feedback when folders are reloaded via hotkey
    let hideRefreshButton = true
    let animateRefreshButton = false

    // the element with the scroll bar
    let scrollContainer: any
</script>

{#if loadingFolders}
    <div class="flex w-full justify-center mb-2 mr-2"><Spinner /></div>
{:else if loadingFoldersOk}
    <!-- folder list -->
    <section
        bind:this={scrollContainer}
        class="flex-auto py-2 overflow-y-auto h-auto rueder-scrollbar"
        use:dndzone={{
            type: "folder",
            items: folders,
            dragDisabled: dragFolderDisabled,
            dropTargetStyle: { outline: "none" },
        }}
        on:consider={handleDndConsiderFolders}
        on:finalize={handleDndFinalizeFolders}
    >
        {#each folders as folder (folder.id)}
            <section>
                {#if editMode}
                    <div class="flex flex-row px-2 md:p-0 md:pl-2 bg-gray-800 md:bg-gray-900">
                        <!-- grab handle stuff based on svelte-dnd-action example: https://svelte.dev/repl/4949485c5a8f46e7bdbeb73ed565a9c7?version=3.24.1 -->
                        <div
                            class="flex-none text-white h-8 w-8 p-1"
                            tabindex={draggingFolderByHandle ? -1 : 0}
                            style={draggingFolderByHandle ? "cursor: grabbing" : "cursor: grab"}
                            on:mousedown={(e) => handleDragStart(DragType.Folder, e)}
                            on:touchstart={(e) => handleDragStart(DragType.Folder, e)}
                            on:keydown={(e) => handleDragKeyDown(DragType.Folder, e)}
                            aria-label="drag-handle"
                        >
                            <DragHandle />
                        </div>

                        <!-- folder title input box -->
                        <EditFolderTitle
                            bind:folder
                            bind:dirty
                            on:deleteFolder={(e) => deleteFolder(e.detail.folderID, e.detail.event)}
                        />
                    </div>
                {:else}
                    <!-- folder title display -->
                    <div class="px-1 mt-1 flex flex-row items-center bg-gray-800 md:bg-transparent">
                        <div
                            class="flex-none text-gray-400 h-5 w-5"
                            on:click|stopPropagation={() => (folder.open = !folder.open)}
                        >
                            {#if folder.open}
                                <SolidFolderOpen />
                            {:else}
                                <SolidFolder />
                            {/if}
                        </div>
                        <h4
                            class="flex-auto px-1 cursor-default text-center md:text-left md:text-sm md:font-bold"
                        >
                            {folder.title ? folder.title : "Unnamed Folder"}
                        </h4>
                        <div
                            title="Reload Feed"
                            class="text-gray-200 ml-1 flex-none object-contain h-5 w-5"
                        >
                            <!-- always show wrapping div with same width as the folder icon so that the title is centered -->
                            <div
                                class:animate-spin={animateRefreshButton}
                                class:hidden={hideRefreshButton}
                            >
                                <RefreshButton size={5} />
                            </div>
                        </div>
                    </div>
                {/if}

                <!-- list of feeds in this folder -->
                <ol
                    class="overflow-y-auto"
                    class:pt-6={editMode && (!folder.feeds || !folder.feeds.length)}
                    class:pb-2={editMode}
                    class:hidden={folder.open === false}
                    use:dndzone={{
                        type: "feed",
                        items: folder.feeds,
                        dragDisabled: dragFeedDisabled,
                        dropTargetStyle: { outline: "none" },
                        scrollContainer,
                    }}
                    on:consider={(e) => handleDndConsiderFeeds(folder.id, e)}
                    on:finalize={(e) => handleDndFinalizeFeeds(folder.id, e)}
                >
                    {#each folder.feeds as feed (feed.id)}
                        <li
                            class="py-1 px-2 md:py-0 flex flex-row items-center hover:bg-gray-700 md:text-sm"
                            class:cursor-pointer={editMode || selectedFeed != feed.id}
                            class:bg-gray-900={editMode || selectedFeed != feed.id}
                            class:bg-gray-600={!editMode && selectedFeed == feed.id}
                            on:click={() => {
                                if (!editMode && selectedFeed != feed.id)
                                    dispatch("feedClick", feed.id)
                            }}
                        >
                            <!-- feed icon (and drag handle in edit mode) -->
                            <div class="flex-none flex items-center h-8 md:h-6">
                                {#if editMode}
                                    <!-- grab handle stuff based on svelte-dnd-action example: https://svelte.dev/repl/4949485c5a8f46e7bdbeb73ed565a9c7?version=3.24.1 -->
                                    <div
                                        class="text-white h-8 w-8 p-1"
                                        tabindex={draggingFeedByHandle ? -1 : 0}
                                        style={draggingFeedByHandle
                                            ? "cursor: grabbing"
                                            : "cursor: grab"}
                                        on:mousedown={(e) => handleDragStart(DragType.Feed, e)}
                                        on:touchstart={(e) => handleDragStart(DragType.Feed, e)}
                                        on:keydown={(e) => handleDragKeyDown(DragType.Feed, e)}
                                        aria-label="drag-handle"
                                    >
                                        <DragHandle />
                                    </div>
                                {:else if selectedFeed == feed.id}
                                    <!-- highlighting arrow icon -->
                                    <div
                                        title="Currently opened feed"
                                        class="h-8 w-8 p-1 md:w-6 md:h-6 text-gray-400"
                                        aria-hidden="true"
                                    >
                                        <OutlineCircleRight size={4} />
                                    </div>
                                {:else}
                                    <img
                                        src={feed.icon
                                            ? imageProxy.buildURL(ImageProxyType.Icon, feed.icon)
                                            : dummyIcon}
                                        alt=""
                                        class="object-contain h-8 w-8 p-1 md:w-6 md:h-6"
                                        aria-hidden="true"
                                        on:error={handleIconError}
                                    />
                                {/if}
                            </div>
                            {#if feed.title}
                                <!-- title -->
                                <h5
                                    class="ml-2 mr-2 flex-auto text-gray-100 truncate"
                                    title={feed.title}
                                >
                                    {feed.title}
                                </h5>
                            {:else}
                                <!-- for untitled feeds -->
                                <h5
                                    class="ml-2 mr-2 flex-auto italic truncate"
                                    title="Untitled feed with URL: {feed.url}"
                                >
                                    Feed {feed.id.slice(0, 8)}
                                </h5>
                            {/if}
                            {#if $feedStateStore.getUnreadCount(feed.id) > 0 && !editMode}
                                <!-- unread count bubble -->
                                <span
                                    class="flex-none bg-gray-800 text-gray-200 text-xs rounded-full py-1 md:py-0 px-2 mx-1"
                                    >{$feedStateStore.getUnreadCount(feed.id)}</span
                                >
                            {/if}
                            {#if editMode}
                                <!-- trash icon (in edit mode) -->
                                <div
                                    class="flex-none mr-1 text-red-900 hover:text-red-600 cursor-pointer"
                                    on:click={(e) => deleteFeed(folder.id, feed.id, e)}
                                >
                                    <SolidTrash />
                                </div>
                            {/if}
                        </li>
                    {/each}
                </ol>

                {#if !editMode && (!folder.feeds || !folder.feeds.length)}
                    <p class="italic px-2 text-gray-400">empty</p>
                {/if}
            </section>
        {:else}
            <p class="italic px-2 text-gray-400">No folders</p>
        {/each}
    </section>

    {#if editMode}
        <!-- add folder box -->
        <div class="flex-none p-2 pr-0 text-white border-t border-green-500 bg-gray-800">
            <p class="text-sm">Add Folder</p>
            <div class="flex flex-row items-center">
                <input
                    bind:value={newFolderName}
                    type="text"
                    placeholder="New folder name"
                    class="flex-grow text-sm text-black min-w-0"
                    on:keydown={(e) => {
                        e.key == "Enter" && addFolder(newFolderName)
                    }}
                />
                <div
                    class="flex-none p-2 hover:text-green-500 cursor-pointer"
                    title="Add folder"
                    on:click={() => addFolder(newFolderName)}
                >
                    <SolidPlusCircle />
                </div>
            </div>
        </div>
    {/if}
{:else}
    <ToastError message="Failed fetching feeds" />
{/if}
