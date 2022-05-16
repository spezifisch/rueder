<script lang="ts">
    import { createEventDispatcher } from "svelte"
    import { writable } from "svelte/store"
    import type { Writable } from "svelte/store"

    import type { Feed } from "../../api/types"

    import ModalWrapper from "../ui/ModalWrapper.svelte"

    import OutlinePlusCircle from "../ui/heroicons/outline-plus-circle.svelte"
    import OutlinePencilAlt from "../ui/heroicons/outline-pencil-alt.svelte"
    import OutlineTag from "../ui/heroicons/outline-tag.svelte"
    import ActiveTriangle from "../ui/icons/ActiveTriangle.svelte"

    import type { FeedStateStore } from "../../stores/articlestate"

    import AddFeed from "./AddFeed/AddFeed.svelte"
    import FeedFolders from "./FeedFolders.svelte"
    import LabelList from "./LabelList.svelte"

    export let feedStateStore: Writable<FeedStateStore>
    export let selectedFeed: string
    export let selectedLabel: string
    export let feedInfoUpdate: Feed
    export let hotkeyPressed: string
    export let show: boolean
    export let isDialogOpen: boolean

    let hide: boolean
    $: hide = !show

    $: isDialogOpen = enableEditMode || doShowAddFeed

    const dispatch = createEventDispatcher()

    enum Mode {
        Folder,
        Label,
    }

    let mode: Mode = Mode.Folder // sidebar mode
    let enableEditMode = false // triggered by edit button

    if (selectedLabel) {
        mode = Mode.Label
    }

    // triggered by add feed button
    let doShowAddFeed = false
    function closeAddFeedDialog() {
        doShowAddFeed = false
    }

    // newFeedsAdded is called when new feeds are added from the "add feed" dialog.
    // it writes the new feed ids to a store. FeedFolder gets these updates and updates the folder list.
    let addNewFeeds = writable<string[]>([])
    function newFeedsAdded(e: any) {
        closeAddFeedDialog()
        dispatch("closeAll")

        const feedIDs: string[] = e.detail
        $addNewFeeds = feedIDs
    }

    $: if (hotkeyPressed) {
        switch (hotkeyPressed) {
            case "Escape":
                doShowAddFeed = false
                enableEditMode = false
                break
        }
    }

    $: disableLabelButton = enableEditMode
    $: disableAddFeedButton = enableEditMode || mode != Mode.Folder
    $: disableEditButton = false

    function handleLabelButton() {
        if (disableLabelButton) {
            return
        }
        if (mode != Mode.Label) {
            mode = Mode.Label
            enableEditMode = false
        } else {
            mode = Mode.Folder
        }
    }

    function handleAddFeedButton() {
        if (disableAddFeedButton) {
            return
        }
        doShowAddFeed = true
    }

    function handleEditButton() {
        if (disableEditButton) {
            return
        }
        enableEditMode = !enableEditMode
    }
</script>

<!-- show "Add Feed" modal -->
{#if doShowAddFeed}
    <ModalWrapper centered={true}>
        <AddFeed on:cancel={closeAddFeedDialog} on:added={newFeedsAdded} />
    </ModalWrapper>
{/if}

<div
    class="flex-auto flex flex-col overflow-y-auto"
    class:edit-mode-border={enableEditMode}
    class:hidden={hide}
>
    {#if mode == Mode.Folder}
        <!-- folders -->
        <FeedFolders
            {feedStateStore}
            {selectedFeed}
            {hotkeyPressed}
            {feedInfoUpdate}
            {addNewFeeds}
            editMode={mode == Mode.Folder && enableEditMode}
            on:feedClick
            on:error
        />
    {:else if mode == Mode.Label}
        <!-- labels -->
        <LabelList {selectedLabel} editMode={mode == Mode.Label && enableEditMode} on:labelClick on:closeLabel />
    {/if}
</div>
<!-- feedlist bottom bar -->
<div
    class="flex-none flex bg-gray-800 default-border"
    class:hidden={hide}
    class:label-active-border={mode == Mode.Label && !enableEditMode}
    class:edit-active-border={enableEditMode}
>
    <div
        class="button button-leftish relative"
        class:disabled-button={disableLabelButton}
        class:text-yellow-500={mode == Mode.Label}
        title="Toggle Label List"
        on:click={handleLabelButton}
    >
        {#if mode == Mode.Label}
            <div class="absolute top-0 text-yellow-500">
                <ActiveTriangle />
            </div>
        {/if}
        <OutlineTag />
    </div>
    <div
        class="button button-leftish"
        class:disabled-button={disableAddFeedButton}
        title="Add Feed"
        on:click={handleAddFeedButton}
    >
        <OutlinePlusCircle />
    </div>
    <div
        class="button relative"
        class:disabled-button={disableEditButton}
        class:text-green-500={enableEditMode}
        title="Toggle Edit {mode == Mode.Label ? 'Labels' : 'Feeds'}"
        on:click={handleEditButton}
    >
        {#if enableEditMode}
            <div class="absolute top-0 text-green-500">
                <ActiveTriangle />
            </div>
        {/if}
        <OutlinePencilAlt />
    </div>
</div>

<style lang="postcss">
    .edit-mode-border {
        @apply border-t border-l border-r border-green-500;
    }

    .button {
        @apply flex-auto p-3 flex justify-center cursor-pointer md:p-1 md:hover:bg-gray-600;
    }

    .button-leftish {
        @apply border-r border-gray-500;
    }

    .disabled-button {
        @apply text-gray-500;
        @apply cursor-not-allowed;
    }

    .default-border {
        @apply border-t-4;
        @apply border-gray-500;
    }

    .edit-active-border {
        @apply border-green-500;
    }

    .label-active-border {
        @apply border-yellow-500;
    }
</style>
