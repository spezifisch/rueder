<script lang="ts">
    import { createEventDispatcher } from "svelte"

    import type { Folder } from "../../../api/types"

    import SolidFolder from "../../ui/heroicons/solid-folder.svelte"
    import SolidFolderOpen from "../../ui/heroicons/solid-folder-open.svelte"
    import SolidTrash from "../../ui/heroicons/solid-trash.svelte"

    const dispatch = createEventDispatcher()

    export let folder: Folder
    export let dirty: boolean

    function clickDeleteFolder(e: MouseEvent) {
        dispatch("deleteFolder", {
            folderID: folder.id,
            event: e,
        })
    }
</script>

<div class="flex-auto flex flex-row items-center mt-1 hover:bg-gray-900 min-w-0">
    <div class="flex-none text-gray-200 mr-1 h-5 w-5" on:click|stopPropagation={() => (folder.open = !folder.open)}>
        {#if folder.open}
            <SolidFolderOpen />
        {:else}
            <SolidFolder />
        {/if}
    </div>
    <input
        class="flex-auto bg-transparent border-b border-white border-t-0 border-r-0 border-l-0 min-w-0 p-0"
        type="text"
        bind:value={folder.title}
        on:keypress={() => (dirty = true)}
        placeholder="Unnamed Folder"
    />
    <div class="flex-none h-7 w-7 min-w-0" on:click={clickDeleteFolder}>
        {#if !folder.feeds || !folder.feeds.length}
            <div class="text-red-900 hover:text-red-600 cursor-pointer p-1" on:click={clickDeleteFolder}>
                <SolidTrash />
            </div>
        {/if}
    </div>
</div>
