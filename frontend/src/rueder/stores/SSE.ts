import { get, writable } from "svelte/store"
import type { Writable } from "svelte/store"

import { EventSourcePolyfill } from "@spezifisch/event-source-polyfill"

import { sessionStore } from "../stores/session"

export declare type SSEData = string

export class SSEStore {
    endpoint: string
    usePolyfill: boolean
    eventSource: EventSource
    store: Writable<SSEData>

    // usePolyfill parameter can be set to false to use browser's native EventSource implementation
    // which doesn't support setting custom headers (and thus can't send the JWT for authentication)
    constructor(endpoint = "http://127.0.0.1:8083/sse", usePolyfill = true) {
        this.endpoint = endpoint
        this.usePolyfill = usePolyfill
        this.store = writable<SSEData>()
    }

    connect() {
        if (this.usePolyfill) {
            const sessionStoreValue = get(sessionStore)
            if (!sessionStoreValue.loggedIn) {
                throw new Error("tried to connect to SSE without being logged in")
            }

            this.eventSource = new EventSourcePolyfill(this.endpoint, {
                headers: {
                    Authorization: "Bearer " + sessionStoreValue.jwtToken,
                },
            })
        } else {
            this.eventSource = new EventSource(this.endpoint)
        }
        this.eventSource.onmessage = (e) => {
            this.store.update(() => e.data)
        }
    }

    close() {
        if (this.eventSource) {
            this.eventSource.close()
        }
    }
}
