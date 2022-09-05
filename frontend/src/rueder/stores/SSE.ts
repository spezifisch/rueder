import { get, writable } from "svelte/store"
import type { Writable } from "svelte/store"

import { EventSourcePolyfill } from "@spezifisch/event-source-polyfill"

import { sessionStore } from "../stores/session"

declare type SSEPayload = Record<string, string>
export class SSEEvent {
    id?: number
    message_type: string
    message_data?: SSEPayload
}

export interface SSEPong {
    ping: string
    msg?: string
    claims?: string
}

declare type Headers = Record<string, string>

export class SSEStore {
    message_number: number
    base: string
    endpoint: string
    usePolyfill: boolean
    eventSource: EventSource
    store: Writable<SSEEvent>

    // usePolyfill parameter can be set to false to use browser's native EventSource implementation
    // which doesn't support setting custom headers (and thus can't send the JWT for authentication)
    constructor(base = "http://127.0.0.1:8083/", usePolyfill = true) {
        this.message_number = 0
        this.base = base
        this.endpoint = base + "sse"
        this.usePolyfill = usePolyfill
        this.store = writable<SSEEvent>()
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

        // SSE message handler
        this.eventSource.onmessage = (evt) => {
            console.log("sse received", evt.data)

            const storeData = new SSEEvent()
            try {
                // try to parse it as json as normal events for us are json
                const jsonData = JSON.parse(evt.data)
                storeData.message_type = jsonData.type
                if (jsonData.data !== undefined) {
                    storeData.message_data = jsonData.data
                }
            } catch (err) {
                // not really needed anymore but some experiments returned non-json
                storeData.message_type = "raw"
                storeData.message_data = {
                    data: evt.data,
                }
            }
            // just for stats we count the messages
            storeData.id = this.message_number++

            this.store.update(() => storeData)
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

export function createSSEStore(baseURL: string): SSEStore {
    const sse = new SSEStore(baseURL)
    sse.connect()
    return sse
}
