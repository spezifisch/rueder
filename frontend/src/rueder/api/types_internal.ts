import type { Folder } from './types'

export class apiResponse {
    ok: boolean
}

// /feed POST endpoint

export class apiAddFeedRequest {
    url: string
}

export class apiAddFeedResponse {
    feed_id: string
    message: string
    status: string
}

// /folders POST endpoint
export class apiChangeFoldersRequest {
    constructor(folders: Folder[]) {
        this.folders = []

        for (const folder of folders) {
            const f = new lightFolder()
            f.id = folder.id
            f.title = folder.title
            f.feeds = []
            for (const feed of folder.feeds) {
                const lf = new lightFeed()
                lf.id = feed.id
                f.feeds.push(lf)
            }
            this.folders.push(f)
        }
    }

    folders: lightFolder[]
}

export class apiChangeFoldersResponse {
    message: string
    status: string
}

// a version of Folder with only the things needed to update them
export class lightFolder {
    id?: string
    title: string // contains new name (or old name if unchanged)
    feeds: lightFeed[] // contains new feed list
}

export class lightFeed {
    id: string
}

// errors

export class apiErrorResponse {
    code: string
    message: string
}
