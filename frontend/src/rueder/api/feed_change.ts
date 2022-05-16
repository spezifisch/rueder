export {}

import { authFetch } from "./auth"
import { AddFeedResponse, ChangeFoldersResponse } from "./types"
import type { Folder } from "./types"
import { apiAddFeedRequest, apiAddFeedResponse, apiChangeFoldersRequest, apiChangeFoldersResponse, apiErrorResponse } from "./types_internal"

// for GET endpoints returning objects
class apiChange<ApiRequestType, ApiResponseType> {
    baseURL: string
    apiPath: string // eg. "feed"

    constructor(baseURL: string, apiPath: string) {
        this.baseURL = baseURL
        this.apiPath = apiPath
    }

    getURL(id?: string): string {
        let suffix = ""
        if (id != null) {
            suffix = `/${id}`
        }
        return `${this.baseURL}/${this.apiPath}${suffix}`
    }

    async postData(data: ApiRequestType): Promise<ApiResponseType> {
        // eslint-disable-next-line no-async-promise-executor
        return new Promise<ApiResponseType>(async (resolve, reject) => {
            // build api url
            const url = this.getURL()

            try {
                const settings = {
                    method: "POST",
                    headers: {
                        Accept: "application/json",
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(data),
                }
                const res = await authFetch(url, settings)
                if (res.status == 200) {
                    // parse result
                    const data: ApiResponseType = await res.json()
                    resolve(data)
                } else {
                    // parse error message
                    const data: apiErrorResponse = await res.json()
                    console.log(`got error from ${url}:`, res.status, data.message)
                    reject(data)
                }
            } catch (error) {
                // invalid json returned
                console.log(`failed fetching ${url}:`, error)
                reject(error)
            }
        })
    }
}

export class apiAddFeed extends apiChange<apiAddFeedRequest, apiAddFeedResponse> {
    add(url: string): Promise<AddFeedResponse> {
        // eslint-disable-next-line no-async-promise-executor
        return new Promise<AddFeedResponse>(async (resolve, reject) => {
            const req = new apiAddFeedRequest()
            req.url = url

            try {
                const resp = await this.postData(req)

                const data = new AddFeedResponse()
                data.ok = resp.status == "ok"
                if (!data.ok) {
                    data.message = `Couldn't add feed: ${resp.status}`
                    reject(data)
                    return
                }
                data.feed_id = resp.feed_id
                resolve(data)
            } catch (e) {
                const ret = new AddFeedResponse()
                ret.ok = false
                if (e.code) {
                    ret.message = e.message
                } else {
                    ret.message = "Failed parsing server response: " + e
                }
                reject(ret)
            }
        })
    }
}

export class apiChangeFolders extends apiChange<apiChangeFoldersRequest, apiChangeFoldersResponse> {
    change(folders: Folder[]): Promise<ChangeFoldersResponse> {
        // eslint-disable-next-line no-async-promise-executor
        return new Promise<ChangeFoldersResponse>(async (resolve, reject) => {
            const req = new apiChangeFoldersRequest(folders)
            const resp = new ChangeFoldersResponse();

            console.log("apiChangeFolders.change", req)

            try {
                await super.postData(req)
                resp.ok = true
                resolve(resp)
            } catch (e) {
                resp.ok = false
                resp.message = e.message
                reject(resp)
            }
        })
    }
}
