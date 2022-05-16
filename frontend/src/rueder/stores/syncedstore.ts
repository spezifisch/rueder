import { writable } from "svelte/store"
import type { Writable, Updater, Subscriber, Unsubscriber } from "svelte/store"

import { ObjectMapper } from "jackson-js"

export class SyncedStore<T> implements Writable<T> {
    first: boolean
    store: Writable<T>
    unsubscriber: Unsubscriber

    data: T

    // eslint-disable-next-line no-unused-vars
    constructor(private name: string, private constructorT: { new (): T }) {
        // get data from server
        this.data = this.requestData()
        if (!this.data) {
            // create a new object if nothing's stored on the server
            this.data = new constructorT()
        }

        // create store and push to server on data updates
        this.store = writable<T>(this.data)
        this.first = true
        this.unsubscriber = this.store.subscribe((value: T) => this.pushData(value))
    }

    requestData(): T {
        const objectMapper = new ObjectMapper()

        // TODO request from server
        const stored = localStorage.getItem(this.name)

        // parse
        return objectMapper.parse<T>(stored, {
            mainCreator: () => [this.constructorT],
        })
    }

    pushData(value: T) {
        if (this.first) {
            // ignore first call which happens right after subscribing to the store
            this.first = false
            return
        }
        if (!value) {
            return
        }
        //console.log("pushing", value)

        const objectMapper = new ObjectMapper()
        localStorage[this.name] = objectMapper.stringify<T>(value)
        // TODO send to server
    }

    // implement wrappers to make this class itself a store
    subscribe(run: Subscriber<T>): Unsubscriber {
        return this.store.subscribe(run)
    }

    set(value: T): void {
        this.store.set(value)
    }

    update(updater: Updater<T>): void {
        this.store.update(updater)
    }
}
