// cleanupURL removed common tracking parameters from URLs
export function cleanupURL(url: string): string {
    let newQuery = ""

    const [tmp1, hash] = url.split("#", 2)
    const [path, query] = tmp1.split("?", 2)
    if (query) {
        const items = query.split("&")
        for (const item of items) {
            const [name, val] = item.split("=", 2)

            if (name === undefined || name === "") {
                // eg. substack just has subscription links with ? at the end for no reason or ?&foo=bar
                continue
            }
            if (name == "fbclid" || name == "igshid" || name == "efg") {
                // facebook and instagram
                continue
            }
            if (name.startsWith("utm_")) {
                // used in various sites
                continue
            }
            if (name == "from" && val == "rss") {
                continue
            }
            if (name.startsWith("wt_")) {
                // eg. heise.de
                continue
            }
            if (name == "source" && val && val.includes("rss")) {
                // eg. medium.com
                continue
            }
            if (name == "CMP") {
                // eg. guardian.co.uk
                continue
            }

            if (newQuery) {
                // there is a previous parameter
                newQuery += "&"
            }
            // an empty val is allowed, we need to handle it so it doesn't say "foo=undefined"
            newQuery += name + "=" + (val ?? "")
        }
    }

    // reconstruct url
    let ret = path
    if (newQuery) {
        ret += "?" + newQuery
    }
    if (hash !== undefined) {
        ret += "#" + hash
    }
    return ret
}

// getURLProtocol returns "https" for a URL like "https://example.com",
// returns undefined on error
export function getURLProtocol(url: string): string {
    return url.split(":").shift()
}

// see https://stackoverflow.com/a/31991870
const absoluteNoProtocolRelativeURLRegex = new RegExp("^[a-zA-Z][a-zA-Z0-9+\\.-]*:")

function isAbsoluteURL(url: string): boolean {
    return absoluteNoProtocolRelativeURLRegex.test(url)
}

// getURLBase returns "https://example.com" for a URL like "https://example.com/foo/bar?a=b#c",
// returns null on error
export function getURLBase(url: string): string {
    if (!url) {
        return null
    }
    const [prefix, nothing, hostname, _path] = url.split("/", 4)
    if (nothing) {
        return null
    }
    if (!prefix || prefix == ":" || !hostname) {
        return null
    }
    return `${prefix}//${hostname}`
}

// getURLBaseDir returns "https://example.com/foo/" for a URL like "https://example.com/foo/bar?a=b#c",
// this function doesn't handle protocol links with ":" instead of "://"
// returns null on error
export function getURLBaseDir(url: string): string {
    if (!url) {
        return null
    }
    // remove query params and hash as they might contain slashes
    const [tmp1, _hash] = url.split("#", 2)
    const [fileURL, _query] = tmp1.split("?", 2)

    // remove the part after the last slash
    const parts = fileURL.split("/")
    let minParts = 1
    if (url.startsWith("//") || isAbsoluteURL(url)) {
        minParts = 3
    }
    if (parts.length > minParts) {
        parts.pop()
    }

    const dirURL = parts.join("/")
    return dirURL
}

// resolve relative URLs from RSS articles
// url is the relative URL vom the "a href"
// base is the feed URL to which relative links are relative to
// anchorBase is the article URL which is used for anchor links
export function fixupRelativeURL(url: string, base: string, anchorBase: string): string {
    if (!url) {
        return null
    }
    if (url.startsWith("#")) {
        // anchor link
        if (anchorBase) {
            return `${anchorBase}${url}`
        }
        return url
    }
    // TODO i forgot what this is for
    if (url.startsWith("///")) {
        return url
    }
    if (base) {
        if (url.startsWith("//")) {
            // protocol-relative link
            return getURLProtocol(base) + ":" + url
        }
        if (url.startsWith("/")) {
            // relative to feed's root directory on host
            return getURLBase(base) + url
        }
        if (!absoluteNoProtocolRelativeURLRegex.test(url)) {
            // relative to feed path, protocol-relative links are ruled out above
            return getURLBaseDir(base) + "/" + url
        }
    }
    return url
}

// removes all '/' characters from end of string
export function stripTrailingSlash(url: string): string {
    return url.replace(/\/$/, "")
}
