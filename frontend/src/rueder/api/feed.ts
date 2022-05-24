import { derived } from "svelte/store"
import type { Readable, Writable } from "svelte/store"

import type { AddFeedResponse, ChangeFoldersResponse } from "./types"
import { Article, ArticlePreview, Feed, Folder, Label } from "./types"
import { apiGet } from "./feed_get"
import { apiAddFeed, apiChangeFolders } from "./feed_change"

export class FeedAPI {
    baseURL = "/api/v1"

    constructor(baseURL: string) {
        this.baseURL = baseURL.replace(/\/$/, "") // strip trailing slash, so we're consistent
    }

    getBaseURL(): string {
        return this.baseURL
    }

    GetArticle(id: string): ArticleStore {
        const store = new apiGet<Article>(this.baseURL, "article").getStore(id)
        // convert article to a full Article instance with methods
        return derived(store, async ($store) => new Article(await $store))
    }

    GetArticles(feed_id: string, start?: number): ArticleList {
        return new apiGet<ArticlePreview[]>(this.baseURL, "articles").getPromise(feed_id, start)
    }

    GetFeeds(): Promise<Feed[]> {
        return new apiGet<Feed[]>(this.baseURL, "feeds").getPromise()
    }

    async GetFeed(id: string): Promise<Feed> {
        const feed = await new apiGet<Feed>(this.baseURL, "feed").getPromise(id)
        return new Feed(feed)
    }

    AddFeed(url: string): Promise<AddFeedResponse> {
        return new apiAddFeed(this.baseURL, "feed").add(url)
    }

    GetFolders(): Promise<Folder[]> {
        return new apiGet<Folder[]>(this.baseURL, "folders").getPromise()
    }

    ChangeFolders(folders: Folder[]): Promise<ChangeFoldersResponse> {
        return new apiChangeFolders(this.baseURL, "folders").change(folders)
    }
}

export type ArticleStore = Readable<Promise<Article>>
export type ArticleList = Promise<ArticlePreview[]>
export type LabelStore = Writable<Promise<Label[]>>
