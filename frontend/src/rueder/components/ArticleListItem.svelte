<script lang="ts">
    import { createEventDispatcher, getContext, onDestroy } from "svelte"
    import type { Writable } from "svelte/store"

    import sanitizeHtml from "sanitize-html"

    import type { ArticlePreview } from "../api/types"
    import { localizeTime, timeAgoString } from "../helpers/time"
    import type { FeedStateStore } from "../stores/articlestate"
    import { labelStore } from "../stores/labels"
    import SolidTag from "./ui/heroicons/solid-tag.svelte"
    import { contextKey } from "../helpers/constants"
    import type { ImageProxy } from "../helpers/ImageProxy"
    import { ImageProxyType } from "../helpers/ImageProxy"

    export let article: ArticlePreview
    export let feedStateStore: Writable<FeedStateStore> = undefined
    export let feedID: string = undefined

    let imageProxy: ImageProxy = getContext(contextKey.imageProxy)

    const dispatch = createEventDispatcher()

    const dummyIcon = "dummy.png"

    // strip all HTML for the teaser. the backend does it already but we want to be sure it's clean
    // so that we can use @html when displaying it. otherwise escaped HTML entities might be displayed to the user.
    function stripAllHTML(dirty: string): string {
        return sanitizeHtml(dirty, {
            allowedTags: [],
            allowedAttributes: {},
        })
    }

    const dateTimeStr = localizeTime(article.time)

    // show dummy icon if loading article icon fails
    function handleIconError(e: Event) {
        if ((e.target as HTMLImageElement).src != dummyIcon) {
            ;(e.target as HTMLImageElement).src = dummyIcon
        }
    }

    function articleClicked() {
        dispatch("articleClick", article.id)
        if (feedID) {
            feedStateStore.update((fss) => {
                fss.markArticleRead(feedID, article.feed_seq)
                return fss
            })
            articleRead = true
        }
    }

    let articleRead: boolean = true
    let articleUnread: boolean
    $: articleUnread = !articleRead
    $: if (feedStateStore) {
        if (feedID) {
            articleRead = $feedStateStore.isArticleRead(feedID, article.feed_seq)
        }
    }

    let articleLabels = []
    const unsubscribe = labelStore.subscribe(() => {
        articleLabels = labelStore.getArticleLabels(article.id)
    })
    onDestroy(unsubscribe)
</script>

<li class="px-2 ml-2 my-1" class:articleRead class:articleUnread on:click={(_e) => articleClicked()}>
    <!-- title box -->
    <div class="flex flex-row items-center cursor-pointer">
        <!-- feed logo -->
        <div class="flex-none">
            <img
                src={article.feed_icon ? imageProxy.buildURL(ImageProxyType.Icon, article.feed_icon) : dummyIcon}
                alt={article.feed_title}
                class="object-contain h-8 w-8 ml-0 m-2"
                loading="lazy"
                on:error={handleIconError}
            />
        </div>
        <!-- right part -->
        <div class="flex-auto flex flex-col min-w-0">
            <!-- article info -->
            <div class="flex text-gray-400">
                <div class="flex-auto px-1 truncate">
                    {article.feed_title}
                </div>
                <div class="flex-shrink px-1">
                    {#each articleLabels as label}
                        <div title="Labelled as: {label.name}" style="color: {label.color}" class="inline-block">
                            <SolidTag />
                        </div>
                    {/each}
                </div>
                <div class="flex-none" title={timeAgoString(article.time)}>
                    {#if dateTimeStr.isToday}
                        <span class="text-gray-300">
                            {dateTimeStr.time}
                        </span>
                    {:else}
                        <!-- highlight the more important part -->
                        <span class="text-gray-300">
                            {dateTimeStr.date}
                        </span>
                        <span>
                            {dateTimeStr.time}
                        </span>
                    {/if}
                </div>
            </div>
            <!-- title -->
            <h2 class="text-lg">{@html article.title ? stripAllHTML(article.title) : "<em>untitled</em>"}</h2>
        </div>
    </div>
    <!-- teaser box -->
    {#if article.teaser}
        <p class="p-2 truncate" data-testid="teaser">
            {@html stripAllHTML(article.teaser)}
        </p>
    {/if}
</li>

<style>
    :root {
        --unreadWidth: 0.5rem;
        --unreadColor: rgb(107, 114, 128);
    }

    li.articleUnread {
        border-left: var(--unreadWidth) solid var(--unreadColor);
        border-right: var(--unreadWidth) solid transparent;
    }

    li.articleRead {
        border-left: var(--unreadWidth) solid transparent;
        border-right: var(--unreadWidth) solid transparent;
    }
</style>
