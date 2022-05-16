import { JsonProperty, JsonClassType } from "jackson-js"

export class Article {
    id: string

    feed_title: string
    feed_id: string
    feed_url: string

    title: string
    time: string
    link: string

    thumbnail?: string
    image?: string
    image_title?: string

    teaser?: string
    content: ArticleContent

    constructor(values: object = {}) {
        Object.assign(this, values)
        if (this.content) {
            this.content = new ArticleContent(this.content)
        }
    }

    toArticlePreview(): ArticlePreview {
        return new ArticlePreview({
            id: this.id,
            seq: -1,
            feed_seq: -1,
            title: this.title,
            time: this.time,
            feed_title: this.feed_title,
            feed_icon: undefined,
            teaser: this.teaser,
        })
    }
}

export class ArticleContent {
    authors?: string[]
    tags?: string[]
    enclosures?: Enclosure[]
    text?: string

    constructor(values: object = {}) {
        Object.assign(this, values)

        // handle some common formatting issues
        this.authors = cleanupArray(this.authors)
        this.tags = cleanupArray(this.tags)

        // this.text is still unsafe to use after this.
        // iframes and such are removed by the backend but tracking pixels, tracking links, ads etc. may remain.
    }
}

function cleanupArray(val: Array<string>): Array<string> {
    if (val && val.length) {
        // remove whitespaces
        val = val.map((elem) => elem.trim())
        // remove empty strings
        val = val.filter((elem) => elem)

        // remove list if nothing is left
        if (!val.length) {
            val = null
        }
    }
    return val
}

export class Enclosure {
    length?: string
    type?: string
    url?: string
}

export class ArticlePreview {
    @JsonProperty()
    @JsonClassType({ type: () => [String] })
    id: string
    @JsonProperty()
    @JsonClassType({ type: () => [Number] })
    seq: number
    @JsonProperty()
    @JsonClassType({ type: () => [Number] })
    feed_seq: number

    @JsonProperty()
    @JsonClassType({ type: () => [String] })
    title: string
    @JsonProperty()
    @JsonClassType({ type: () => [String] })
    time: string
    @JsonProperty()
    @JsonClassType({ type: () => [String] })
    feed_title: string
    @JsonProperty()
    @JsonClassType({ type: () => [String] })
    feed_icon: string
    @JsonProperty()
    @JsonClassType({ type: () => [String] })
    teaser: string

    constructor(values: object = {}) {
        Object.assign(this, values)
    }
}

export class Feed {
    id: string

    title?: string
    icon?: string
    url?: string
    site_url?: string
    article_count?: number

    fetcher_state?: FetcherState

    articles?: ArticlePreview[]

    constructor(values: object = {}) {
        Object.assign(this, values)

        if (this.fetcher_state) {
            this.fetcher_state = new FetcherState(this.fetcher_state)
        }
    }
}

export class FetcherState {
    working: boolean
    last_success?: string
    last_error?: string
    message?: string

    constructor(values: object = {}) {
        Object.assign(this, values)
    }

    neverFetched(): boolean {
        return isNullTimestamp(this.last_success) && isNullTimestamp(this.last_error)
    }

    getLastSuccess(): string {
        if (isNullTimestamp(this.last_success)) {
            return
        }
        return this.last_success
    }

    getLastError(): string {
        if (isNullTimestamp(this.last_error)) {
            return
        }
        return this.last_error
    }
}

function isNullTimestamp(val: string): boolean {
    return !val || val.startsWith("0001-01-01T")
}

export class Folder {
    id: string
    type?: string
    open?: boolean

    title: string
    feeds: Feed[]

    constructor(values: object = {}) {
        Object.assign(this, values)
    }
}

export class Label {
    id: string

    title: string
    color: string
    articles: ArticlePreview[]

    constructor(values: object = {}) {
        Object.assign(this, values)
    }
}

export class BaseResponse {
    ok: boolean
    message: string // only if !ok
}

export class AddFeedResponse extends BaseResponse {
    feed_id: string // only if ok
}

export class ChangeFoldersResponse extends BaseResponse {}
