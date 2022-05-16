import { get } from "svelte/store"

import { sessionStore, logout } from "../stores/session"

export async function authFetch(input: RequestInfo, init?: RequestInit): Promise<Response> {
    // add auth header if logged in
    const store = get(sessionStore)
    if (store.loggedIn) {
        init = init ?? {}
        init.headers = init.headers ?? {}
        init.headers["Authorization"] = "Bearer " + store.jwtToken
    }

    const resp = await fetch(input, init)
    if (resp.status == 401 || resp.status == 403) { // unauthorized/forbidden
        console.log("auth failed, invalidating jwt")
        logout()
    }
    return resp
}
