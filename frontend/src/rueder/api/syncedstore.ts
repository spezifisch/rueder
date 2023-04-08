import { stripTrailingSlash } from "../helpers/url"

export class SyncedStoreAPI {
    baseURL = "/api/v1"

    constructor(baseURL: string) {
        this.baseURL = stripTrailingSlash(baseURL)
    }


}
