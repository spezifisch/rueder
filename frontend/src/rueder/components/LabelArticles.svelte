<script lang="ts">
    import { createEventDispatcher } from "svelte"
    import type { Unsubscriber } from "svelte/store"
    import type { ArticlePreview } from "../api/types"

    import { labelStore, Label } from "../stores/labels"
    import ArticleListItem from "./ArticleListItem.svelte"

    import OutlineArrowCircleLeft from "./ui/heroicons/outline-arrow-circle-left.svelte"
    import CloseButton from "./ui/heroicons/outline-x-circle.svelte"
    import OutlineTag from "./ui/heroicons/outline-tag.svelte"

    export let selectedLabel: string

    const dispatch = createEventDispatcher()

    let label: Label
    let articles: ArticlePreview[]

    // update subscription when selected label changes
    $: update(selectedLabel)

    let unsubscribe: Unsubscriber
    function update(_: string) {
        if (unsubscribe) {
            unsubscribe()
        }

        // update label info and article list when store changes (label color or name changed, articles changed)
        unsubscribe = labelStore.subscribe(() => {
            if (!labelStore.labelExists(selectedLabel)) {
                // selected label was deleted
                return
            }
            label = labelStore.getLabel(selectedLabel)
            articles = labelStore.getLabelArticles(selectedLabel)
        })
    }

    function closeLabel() {
        dispatch("close")
    }
</script>

<ul class="flex flex-col h-full relative overflow-y-auto">
    <li>
        <div class="hidden md:flex absolute top-2 right-2 flex-row flex-wrap">
            <div title="Close Label" on:click={closeLabel} class="text-gray-500 hover:text-gray-200 cursor-pointer">
                <CloseButton />
            </div>
        </div>

        <div class="bg-gray-800 py-2 md:px-2 mb-2 flex justify-center items-center space-x-2">
            <div title="Back to folder list" on:click={closeLabel} class="flex-none md:hidden text-gray-500 cursor-pointer p-1">
                <OutlineArrowCircleLeft size={8} />
            </div>
            {#if label}
                <div class="flex-auto flex items-center justify-center space-x-2">
                    <div class="flex-none" style="color: {label.color}">
                        <OutlineTag />
                    </div>
                    <h2 class="text-xl text-gray-100">
                        {label.name}
                    </h2>
                </div>
            {/if}
            <div class="flex-none md:hidden p-1">
                <div class="flex h-8 w-8" />
            </div>
        </div>
    </li>

    {#each articles as article, i (article.id)}
        {#if i != 0}
            <hr class="border-gray-400 mx-2" />
        {/if}
        <ArticleListItem {article} on:articleClick />
    {:else}
        <li class="italic p-2 text-gray-400">No articles.</li>
    {/each}
</ul>
