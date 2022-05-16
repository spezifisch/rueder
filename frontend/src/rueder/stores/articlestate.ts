import { writable } from "svelte/store"
import { validate as validate_uuid } from "uuid"

export class FeedStateStore {
    feeds: Map<string, FeedState>

    constructor(values: object = {}) {
        this.feeds = new Map<string, FeedState>()
        Object.assign(this, values)

        for (const feedID in this.feeds) {
            if (!validate_uuid(feedID)) {
                console.log("removed invalid uuid from feed state", feedID)
                delete this.feeds[feedID]
                continue
            }
            this.feeds[feedID] = new FeedState(this.feeds[feedID])
        }

        console.log("created FeedStateStore:", this.feeds)
    }

    createFeedState(feedID: string): boolean {
        if (!validate_uuid(feedID)) {
            return false
        }

        if (!this.feeds[feedID]) {
            console.log("created feed state for", feedID)
            this.feeds[feedID] = new FeedState()
        }
        return true
    }

    getUnreadCount(feedID: string): number {
        if (!this.createFeedState(feedID)) {
            return 0
        }
        return this.feeds[feedID].getUnreadCount()
    }

    setTotalArticles(feedID: string, count: number): number {
        if (!this.createFeedState(feedID)) {
            return 0
        }
        return this.feeds[feedID].setTotalArticles(count)
    }

    getTotalArticles(feedID: string): number {
        if (!this.createFeedState(feedID)) {
            return 0
        }
        return this.feeds[feedID].getTotalArticles()
    }

    isArticleRead(feedID: string, feedSeq: number): boolean {
        if (!this.createFeedState(feedID)) {
            return false
        }
        return this.feeds[feedID].isArticleRead(feedSeq)
    }

    markArticleRead(feedID: string, feedSeq: number): void {
        if (!this.createFeedState(feedID)) {
            return
        }
        this.feeds[feedID].markArticleRead(feedSeq)
    }

    markAllReadUntil(feedID: string, feedSeq: number) {
        if (!this.createFeedState(feedID)) {
            return
        }
        this.feeds[feedID].markAllReadUntil(feedSeq)
    }

    markAllRead(feedID: string) {
        if (!this.createFeedState(feedID)) {
            return
        }
        const feedSeq = this.feeds[feedID].totalArticles
        if (!feedSeq) {
            return
        }
        this.feeds[feedID].markAllReadUntil(feedSeq)
    }
}

export class FeedState {
    totalArticles: number // total number of articles in feed
    readArticles: number[] // feed sequence number of read articles which are after readAllUntilFeedSeq
    readAllUntil: number // feed_seq number up until which we have read all articles in this feed

    constructor(values: object = {}) {
        this.totalArticles = 0
        this.readArticles = []
        this.readAllUntil = 0
        Object.assign(this, values)
    }

    setTotalArticles(count: number) {
        this.totalArticles = count
    }

    getTotalArticles(): number {
        return this.totalArticles
    }

    getUnreadCount(): number {
        return this.totalArticles - this.readAllUntil - this.readArticles.length
    }

    isArticleRead(feedSeq: number): boolean {
        if (feedSeq <= this.readAllUntil) {
            return true
        }
        return this.readArticles.includes(feedSeq)
    }

    markArticleRead(feedSeq: number): void {
        if (feedSeq <= this.readAllUntil) {
            // do nothing, article was already read
        } else if (this.readAllUntil + 1 == feedSeq) {
            this.readAllUntil++
        } else if (!this.readArticles.includes(feedSeq)) {
            this.readArticles.push(feedSeq)
        }

        // TODO defrag readArticles
    }

    markAllReadUntil(feedSeq: number) {
        this.readAllUntil = feedSeq
        this.readArticles = this.readArticles.filter((it) => it > feedSeq)
    }
}

function requestFeedState(): FeedStateStore {
    const stored = localStorage.getItem("feedstate")
    return new FeedStateStore(JSON.parse(stored))
}

function pushFeedState(value: FeedStateStore) {
    localStorage.feedstate = JSON.stringify(value)
}

// the store itself
export const feedStateStore = writable<FeedStateStore>(requestFeedState())
feedStateStore.subscribe(pushFeedState)
