<script lang="ts">
    import "./ArticleContent.postcss"

    import { createEventDispatcher, getContext } from "svelte"
    import type { Readable } from "svelte/store"
    import { derived } from "svelte/store"

    import CloseButton from "./ui/heroicons/outline-x-circle.svelte"
    import OutlineArrowCircleLeft from "./ui/heroicons/outline-arrow-circle-left.svelte"
    import OutlinePaperClip from "./ui/heroicons/outline-paper-clip.svelte"
    import Spinner from "./ui/Spinner.svelte"
    import ToastError from "./ui/ToastError.svelte"
    import LabelEditor from "./Label/LabelEditor.svelte"

    import { humanFileSize } from "../api/helpers"
    import { localizeTime, timeAgoString } from "../helpers/time"
    import { cleanupURL } from "../helpers/url"
    import type { Article, Enclosure } from "../api/types"
    import { contextKey } from "../helpers/constants"
    import { ImageProxyType } from "../helpers/ImageProxy"
    import type { ImageProxy } from "../helpers/ImageProxy"
    import { ArticleCleaner } from "../helpers/ArticleCleaner"

    export let article_content: Readable<Promise<Article>>
    export let isFocused: boolean

    const imageProxy: ImageProxy = getContext(contextKey.imageProxy)

    const dispatch = createEventDispatcher()

    const article_cleaner = new ArticleCleaner(imageProxy)

    // get displayed filename for attachments
    function getEnclosureName(url: string): string {
        // source: https://stackoverflow.com/questions/511761/js-function-to-get-filename-from-url/48554885#comment61576914_17143667
        // modified to also split at encoded slashes (%2f) which leads to better results for eg. substack feeds
        url = url.split("#").shift().split("?").shift().replaceAll("%2f", "/").replaceAll("%2F", "/")
        url = url.split("/").pop()
        // resolve url encoded characters (html code is still escaped when showing this name)
        return decodeURIComponent(url)
    }

    // convert length in bytes to human-readable string
    function getEnclosureLength(val: string): string {
        return humanFileSize(val)
    }

    $: dateTimeStr = derived(article_content, async ($article_content) => {
        try {
            return localizeTime((await $article_content).time, false)
        } catch (e) {}
    })

    let allEnclosures: Enclosure[]
    $: updateEnclosures($article_content, article_cleaner.extractedEnclosures)
    async function updateEnclosures(article: Promise<Article>, extractedEnclosures: Enclosure[]) {
        try {
            const ac = await article
            allEnclosures = [...extractedEnclosures]
            if (ac.content && ac.content.enclosures) {
                allEnclosures = [...ac.content.enclosures, ...allEnclosures]
            }
        } catch (e) {}
    }

    function closeArticle() {
        dispatch("close")
    }

    function onLabelEditorFocusChange(event: any) {
        isFocused = event.detail
    }
</script>

<div class="flex-auto md:p-4 md:pl-2 overflow-y-auto relative bg-gray-800 rueder-scrollbar">
    <div class="flex flex-row space-x-4 md:hidden bg-gray-900 mb-2">
        <div title="Back to article list" on:click={closeArticle} class="text-gray-500 cursor-pointer p-1">
            <OutlineArrowCircleLeft size={8} />
        </div>
    </div>
    <div
        class="hidden md:block absolute top-2 right-4 text-gray-500 hover:text-gray-200 cursor-pointer bg-gray-800 rounded-full"
        title="Close Article"
        on:click={closeArticle}
    >
        <CloseButton />
    </div>

    {#await $article_content}
        <div class="flex w-full justify-center mb-2 mr-2" title="Loading"><Spinner /></div>
    {:then article_content}
        <div class="mx-1 md:mx-0">
            <!-- header -->
            <div class="px-4">
                {#if article_content.content && article_content.content.tags}
                    <ul class="pb-2">
                        {#each article_content.content.tags as tag}
                            {#if tag}
                                <li
                                    class="inline-grid rounded-full py-1 px-2 mr-1 text-sm bg-gray-600 text-gray-100 truncate"
                                    title="Tagged as: {tag}"
                                >
                                    {tag}
                                </li>
                            {/if}
                        {/each}
                    </ul>
                {/if}
                <div class="text-gray-300">
                    <a href="#/feed/{article_content.feed_id}" title="Open feed in rueder"
                        >{article_content.feed_title}</a
                    >
                </div>
                <h2 class="text-3xl">
                    {#if article_content.link}
                        <a
                            href={cleanupURL(article_content.link)}
                            title="Go to external article page"
                            rel="noreferrer noopener"
                            target="_blank"
                        >
                            {#if article_content.title}
                                {article_content.title}
                            {:else}
                                <em>untitled</em>
                            {/if}
                        </a>
                    {:else if article_content.title}
                        {article_content.title}
                    {:else}
                        <em>untitled</em>
                    {/if}
                </h2>
                {#await $dateTimeStr}
                    <!---->
                {:then dateTimeStr}
                    <p class="text-gray-300" title="Posting time according to the feed">
                        {dateTimeStr.date}
                        {dateTimeStr.time}
                        ({timeAgoString(article_content.time)})
                    </p>
                {/await}
                {#if article_content.content && article_content.content.authors}
                    <div class="text-gray-300" title="Author list">
                        by {article_content.content.authors}
                    </div>
                {/if}
            </div>
            <!-- labels -->
            <LabelEditor article={article_content} on:focusChange={onLabelEditorFocusChange} />
            <!-- enclosures -->
            {#if allEnclosures && allEnclosures.length > 0}
                <div class="my-2 p-2 bg-gray-900">
                    <h3 class="font-bold ml-1 mb-1">Media</h3>
                    <ul>
                        {#each allEnclosures as enclosure}
                            <li class="flex flex-row space-x-1">
                                <OutlinePaperClip addClass="flex-none" />
                                {#if enclosure.url}
                                    <a
                                        class="underline truncate flex-shrink"
                                        rel="noreferrer noopener nofollow"
                                        target="_blank"
                                        title="Open attachment link"
                                        href={cleanupURL(enclosure.url)}>{getEnclosureName(enclosure.url)}</a
                                    >
                                {/if}

                                {#if enclosure.length && enclosure.length != "0"}
                                    <span class="flex-none">
                                        {getEnclosureLength(enclosure.length)}
                                    </span>
                                {/if}
                                {#if enclosure.type}
                                    <span class="flex-none">
                                        ({enclosure.type})
                                    </span>
                                {/if}
                            </li>
                        {/each}
                    </ul>
                </div>
            {/if}
            <!-- image -->
            {#if article_content.image && !article_cleaner.articleContainsImage}
                <!-- this shows only if the article content itself doesn't contain any image -->
                <div class="m-2 w-5/12 mx-auto">
                    <a
                        href={cleanupURL(article_content.image)}
                        rel="noreferrer noopener nofollow"
                        target="_blank"
                        title="Open original image link"
                    >
                        <img
                            src={imageProxy.buildURL(ImageProxyType.Content, cleanupURL(article_content.image))}
                            alt={article_content.image_title}
                        />
                    </a>
                </div>
            {/if}
            <!-- text -->
            <div class="article-text p-4 my-2 bg-gray-900">
                {#if article_content.content && article_content.content.text}
                    {@html article_cleaner.cleanupHTML({
                        html: article_content.content.text,
                        base: article_content.feed_url,
                        anchorBase: article_content.link,
                        parseEnclosures: true,
                    })}
                {:else if article_content.teaser}
                    <!-- teaser html was stripped by backend -->
                    {@html article_cleaner.cleanupHTML({
                        html: article_content.teaser,
                        base: article_content.feed_url,
                        anchorBase: article_content.link,
                    })}
                {:else}
                    <p class="italic">No content.</p>
                {/if}
            </div>
        </div>
    {:catch}
        <div class="mr-2"><ToastError message="Failed fetching article" /></div>
    {/await}
</div>
