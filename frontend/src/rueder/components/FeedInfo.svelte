<script lang="ts">
    import type { Feed } from "../api/types"
    import { localizeTime, timeAgoString } from "../helpers/time"
    import { cleanupURL } from "../helpers/url"
    import OutlineExclamation from "./ui/heroicons/outline-exclamation.svelte"

    export let feedID: string
    export let feedInfo: Promise<Feed>

    let showErrorDetails = false
</script>

<div class="bg-gray-800 p-2">
    {#await feedInfo}
        &nbsp;
    {:then feed}
        <h2 class="text-lg md:text-xl text-center border-b border-gray-400 text-gray-100">
            {#if feed.title}
                {feed.title}
            {:else}
                Feed {feed.id.slice(0, 8)}
            {/if}
        </h2>
        {#if feed.fetcher_state.getLastSuccess()}
            <p class="md:hidden text-xs pt-2 text-right text-gray-200">
                Last Update: {localizeTime(feed.fetcher_state.getLastSuccess()).toString(true)}
            </p>
        {/if}
        <div class="hidden md:block text-xs font-mono text-gray-400 pt-2">
            <div class="flex flex-row">
                {#if feed.site_url}
                    <p class="truncate">
                        Site: <a
                            href={cleanupURL(feed.site_url)}
                            title="Go to feed homepage"
                            rel="noreferrer noopener nofollow"
                            target="_blank">{cleanupURL(feed.site_url)}</a
                        >
                    </p>
                {/if}
                <div class="flex-auto px-1" />
                <p class="flex-none">ID: <code>{feed.id}</code></p>
            </div>

            <div class="flex flex-row mt-1">
                <p class="truncate">Feed: <code>{feed.url}</code></p>
                {#if feed.fetcher_state.getLastSuccess()}
                    <div class="flex-auto px-1" />
                    <p class="flex-none text-right text-gray-200">
                        Last Update: {localizeTime(feed.fetcher_state.getLastSuccess()).toString(true)}
                    </p>
                {/if}
            </div>
        </div>
    {:catch}
        Failed getting info for Feed {feedID}.
    {/await}
</div>

{#await feedInfo}
    <!---->
{:then feed}
    {#if feed.fetcher_state.working}
        <!-- nothing to report -->
    {:else if feed.fetcher_state.neverFetched()}
        <div class="p-2 italic">
            <p>This feed was added recently and wasn't fetched yet.</p>
            <p>Please be patient and try again.</p>
        </div>
    {:else}
        <!-- feed isn't working -->
        <div class="m-2 border-gray-500 border">
            <!-- title -->
            <div class="flex flex-row bg-gray-800 items-center cursor-pointer" on:click={(_e) => { showErrorDetails = !showErrorDetails }}>
                <OutlineExclamation size={8} addClass="flex-none text-red-600 m-2" />
                <div class="flex-auto p-2 text-center">
                    <p>This feed is having problems.</p>
                    <p class="text-xs text-gray-500" class:hidden={showErrorDetails}>click for details</p>
                </div>
            </div>
            <!-- description -->
            <ul class="flex-auto p-2 text-sm" class:hidden={!showErrorDetails}>
                <li>
                    <strong>Feed URL:</strong>
                    <a
                        class="underline font-mono text-sm break-all"
                        href={feed.url}
                        rel="noopener noreferrer nofollow"
                        target="_blank">{feed.url}</a
                    >
                </li>
                <li>
                    <strong>Last Success:</strong>
                    {#if feed.fetcher_state.getLastSuccess()}
                        {localizeTime(feed.fetcher_state.getLastSuccess())} (<em
                            >{timeAgoString(feed.fetcher_state.getLastSuccess())}</em
                        >)
                    {:else}
                        <em>never</em>
                    {/if}
                </li>
                <li>
                    <strong>Last Error:</strong>
                    {localizeTime(feed.fetcher_state.getLastError())} (<em
                        >{timeAgoString(feed.fetcher_state.getLastError())}</em
                    >)
                </li>
                {#if feed.fetcher_state.message}
                    <li class="overflow-auto">
                        Details:
                        <pre class="text-xs border-l-4 border-gray-500 p-2 break-all whitespace-pre-wrap">{feed
                                .fetcher_state.message}</pre>
                    </li>
                {/if}
            </ul>
        </div>
    {/if}
{:catch}
    Failed getting info for Feed {feedID}.
{/await}
