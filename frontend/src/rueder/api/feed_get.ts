export {}

import { writable } from "svelte/store"
import type { Writable } from "svelte/store"

import { authFetch } from "./auth"

// generic store for GET endpoints returning objects
export class apiGet<Type> {
    baseURL: string
    apiPath: string // eg. "article" or "folders"

    constructor(baseURL: string, apiPath: string) {
        this.baseURL = baseURL
        this.apiPath = apiPath
    }

    getURL(id?: string, start?: number): string {
        let suffix = ""
        if (id != null) {
            suffix = `/${id}`
        }
        if (start != undefined && start > 0) {
            suffix += `?start=${start}`
        }
        return `${this.baseURL}/${this.apiPath}${suffix}`
    }

    getPromise(id?: string, start?: number): Promise<Type> {
        // eslint-disable-next-line no-async-promise-executor
        return new Promise<Type>(async (resolve, reject) => {
            // build api url
            const url = this.getURL(id, start)
            if (id == "") {
                // against user error
                reject(new Error(`id is empty for ${url}`))
                return
            }

            try {
                const res = await authFetch(url)
                if (res.status == 200) {
                    // parse result
                    const data: Type = await res.json()
                    resolve(data)
                } else {
                    // parse error message
                    const reason = await res.json()
                    reject(new Error(`got error from ${url}: ${res.status} ${reason.message}`))
                }
            } catch (error) {
                // invalid json returned
                reject(new Error(`failed fetching ${url}: ${error}`))
            }
        })
    }

    getStore(id?: string): Writable<Promise<Type>> {
        const store = writable(
            // eslint-disable-next-line @typescript-eslint/no-empty-function
            new Promise<Type>(() => {})
        )

        // the idea here is that the store updates when the promise is resolved/rejected
        this.getPromise(id)
            .then((data: Type) => {
                store.set(Promise.resolve(data))
            })
            .catch((err: Error) => {
                store.set(Promise.reject(err))
            })

        return store
    }
}
