import { get, writable } from "svelte/store"
import type { Writable } from "svelte/store"

import { EventSourcePolyfill } from "@spezifisch/event-source-polyfill"

import { sessionStore } from "../stores/session"

export declare type SSEData = string
export interface SSEPong {
    ping: string
    msg?: string
    claims?: string
}

declare type Headers = Record<string, string>

export class SSEStore {
    base: string
    endpoint: string
    usePolyfill: boolean
    eventSource: EventSource
    store: Writable<SSEData>

    // usePolyfill parameter can be set to false to use browser's native EventSource implementation
    // which doesn't support setting custom headers (and thus can't send the JWT for authentication)
    constructor(base = "http://127.0.0.1:8083/", usePolyfill = true) {
        this.base = base
        this.endpoint = base + "sse"
        this.usePolyfill = usePolyfill
        this.store = writable<SSEData>()
    }

    private getHeaders(): Headers {
        const sessionStoreValue = get(sessionStore)
        if (!sessionStoreValue.loggedIn) {
            throw new Error("tried to connect to SSE without being logged in")
        }
        return {
            Authorization: "Bearer " + sessionStoreValue.jwtToken,
        }
    }

    connect() {
        if (this.usePolyfill) {
            this.eventSource = new EventSourcePolyfill(this.endpoint, {
                headers: this.getHeaders(),
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

    // make request to default route to check if it's up
    async ping(): Promise<SSEPong> {
        const resp = await fetch(this.base, {
            headers: this.getHeaders(),
        })
        const data = await resp.json()
        return data as SSEPong
    }
}
