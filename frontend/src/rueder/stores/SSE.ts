import { get, writable } from "svelte/store"
import type { Writable } from "svelte/store"

import { EventSourcePolyfill } from "@spezifisch/event-source-polyfill"

import { sessionStore } from "../stores/session"

export declare type SSEData = string

export class SSEStore {
    endpoint: string
    eventSource: EventSource
    store: Writable<SSEData>

    constructor(endpoint = "http://127.0.0.1:8083/sse") {
        this.endpoint = endpoint
        this.store = writable<SSEData>()
    }

    connect() {
        const sessionStoreValue = get(sessionStore)
        if (!sessionStoreValue.loggedIn) {
            throw new Error("tried to connect to SSE without being logged in")
        }

        this.eventSource = new EventSourcePolyfill(this.endpoint, {
            headers: {
                Authorization: "Bearer " + sessionStoreValue.jwtToken,
            },
        })
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
