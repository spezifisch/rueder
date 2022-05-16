import { get, writable } from "svelte/store"

import jwt_decode from "jwt-decode"

// our session state
export class SessionStore {
    loggedIn = false
    jwtToken: string
}

function jwtDecode() {
    return jwt_decode(get(sessionStore).jwtToken)
}

export function getID(): string {
    const dec = jwtDecode()
    const sub: string = dec["sub"]
    const origin: string = dec["origin"]
    return `${origin}:${sub}`
}

export function getUsername(): string {
    const dec = jwtDecode()
    return dec["sub"]
}

// invalidate session (on http 4xx status or user request)
export function logout() {
    sessionStore.update(
        (store: SessionStore): SessionStore => {
            store.loggedIn = false
            store.jwtToken = null
            return store
        }
    )
}

// store
export const sessionStore = writable<SessionStore>(JSON.parse(localStorage.getItem("session")) || new SessionStore())
sessionStore.subscribe((value) => (localStorage.session = JSON.stringify(value)))
