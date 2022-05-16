import * as url from "./url"

// feed URL
const base = "https://feeds.example.com/~ff/feed/yo/?evil=tracking"
// article URL
const anchorBase = "https://example.com/articles/23?utm_evil=tracking"

describe("cleanupURL", () => {
    it("removes heise tracking elements", () => {
        const a = "https://www.example.com/foo/KI-erspaeht-7092096.html?wt_mc=rss.red.ho.ho.atom.beitrag.beitrag"
        const aClean = "https://www.example.com/foo/KI-erspaeht-7092096.html"
        expect(url.cleanupURL(a)).toBe(aClean)
    })

    it("removes slashdot tracking elements", () => {
        const a =
            "https://tech.example.com/story/22/05/14/2026232/should-social?utm_source=rss1.0mainlinkanon&utm_medium=feed"
        const aClean = "https://tech.example.com/story/22/05/14/2026232/should-social"
        expect(url.cleanupURL(a)).toBe(aClean)
    })

    it("removes medium tracking elements", () => {
        const a = "https://example.com/episode-155-how-the-cafef00db4d?source=rss-abcdef010142------2"
        const aClean = "https://example.com/episode-155-how-the-cafef00db4d"
        expect(url.cleanupURL(a)).toBe(aClean)
    })

    it("removes guardian tracking elements", () => {
        const a = "https://www.example.com/search?entry=10000101994&utm_source=guardian-feast&utm_medium=referral&utm_campaign=guardian-feast-recipe"
        const aClean = "https://www.example.com/search?entry=10000101994"
        expect(url.cleanupURL(a)).toBe(aClean)

        const b = "https://www.example.com/culture/series/saved-for-later?CMP=dcba_sfl"
        const bClean = "https://www.example.com/culture/series/saved-for-later"
        expect(url.cleanupURL(b)).toBe(bClean)
    })
})

describe("getURLProtocol", () => {
    it("returns URLs without colon as-is", () => {
        expect(url.getURLProtocol("")).toBe("")
        expect(url.getURLProtocol("x")).toBe("x")
    })

    it("returns the protocol part of an URL", () => {
        // https://en.wikipedia.org/wiki/List_of_URI_schemes
        const a = "aaa://host.example.com:1813;transport=udp;protocol=radius"
        const b = "facetime://+19995551234"
        const c = "git://github.com/user/project-name.git"
        const d = "market://search?q=pub:Publisher_Name"
        const e = "spotify:track:2jCnn1QPQ3E8exampleNsx"
        const f = "mailto:someone@example.com?subject=Hello&body=Hello%20World"
        expect(url.getURLProtocol(a)).toBe("aaa")
        expect(url.getURLProtocol(b)).toBe("facetime")
        expect(url.getURLProtocol(c)).toBe("git")
        expect(url.getURLProtocol(d)).toBe("market")
        expect(url.getURLProtocol(e)).toBe("spotify")
        expect(url.getURLProtocol(f)).toBe("mailto")
    })
})

describe("getURLBase", () => {
    it("returns null for invalid/incomplete/relative http/https URLs", () => {
        expect(url.getURLBase(null)).toBeNull()
        expect(url.getURLBase("")).toBeNull()
        expect(url.getURLBase("http:")).toBeNull()
        expect(url.getURLBase("http:/")).toBeNull()
        expect(url.getURLBase("http://")).toBeNull()
        expect(url.getURLBase("x://")).toBeNull()
        expect(url.getURLBase("http:/x/x")).toBeNull()
        expect(url.getURLBase("/foo")).toBeNull()
        expect(url.getURLBase("//foo")).toBeNull()
        expect(url.getURLBase("://foo")).toBeNull()
        expect(url.getURLBase("#anchor")).toBeNull()
    })

    it("returns the protocol and host part for valid http/https URLs", () => {
        expect(url.getURLBase("http://x")).toBe("http://x")
        expect(url.getURLBase("https://x")).toBe("https://x")
        expect(url.getURLBase("gopher://x")).toBe("gopher://x") // we don't really care but ok
        expect(url.getURLBase("http://x/y")).toBe("http://x")
        expect(url.getURLBase("https://x/y")).toBe("https://x")
        expect(url.getURLBase("https://x/y/yy////aa/https://http://::://")).toBe("https://x")
        expect(url.getURLBase("x://foo")).toBe("x://foo")
        expect(url.getURLBase(base)).toBe("https://feeds.example.com")
        expect(url.getURLBase(anchorBase)).toBe("https://example.com")
    })
})

describe("getURLBaseDir", () => {
    it("removes anything after hashes and query strings", () => {
        expect(url.getURLBaseDir("x/?foo")).toBe("x")
        expect(url.getURLBaseDir("x/?foo=bar")).toBe("x")
        expect(url.getURLBaseDir("x/?foo=/bar&x=/y&1=2")).toBe("x")
        expect(url.getURLBaseDir("x/?foo=bar&x=y&1=/2?/yy")).toBe("x")
        expect(url.getURLBaseDir("x/??1=2/?yy")).toBe("x")
        expect(url.getURLBaseDir("x/?/??")).toBe("x")
        expect(url.getURLBaseDir("x/?#y")).toBe("x")
        expect(url.getURLBaseDir("x/?a=1#y")).toBe("x")
        expect(url.getURLBaseDir("x/?a=1#y?1=2")).toBe("x")
        expect(url.getURLBaseDir("x/?foo/bar/baz")).toBe("x")
        expect(url.getURLBaseDir("x/#foo/bar/baz")).toBe("x")
        expect(url.getURLBaseDir("x/?#foo/bar/baz")).toBe("x")
        expect(url.getURLBaseDir("x/#?foo/bar/baz")).toBe("x")
        expect(url.getURLBaseDir("x/?#foo/b?ar/ba///z§$()=")).toBe("x")
        expect(url.getURLBaseDir("x/#?foo/bar/baz")).toBe("x")
        expect(url.getURLBaseDir("https://example.com/x/#?foo/bar/baz")).toBe("https://example.com/x")
    })

    it("removes the file name and returns the path without a trailing slash", () => {
        const a = "https://example.com/dir"
        expect(url.getURLBaseDir(a + "/foo.jpg")).toBe(a)
        expect(url.getURLBaseDir(a + "/foo.jpg?a=b")).toBe(a)
        expect(url.getURLBaseDir(a + "/foo.jpg?a=b/x/y/z")).toBe(a)
        expect(url.getURLBaseDir(a + "/foo.jpg#a/b/c")).toBe(a)
        expect(url.getURLBaseDir(a + "/foo")).toBe(a)
        expect(url.getURLBaseDir(a + "/foo.jpg?/////")).toBe(a)
        expect(url.getURLBaseDir(a + "/")).toBe(a)
        expect(url.getURLBaseDir(a + "/bar/baz/boo/")).toBe(a + "/bar/baz/boo")
        expect(url.getURLBaseDir(a + "/bar/baz/boo")).toBe(a + "/bar/baz")
        // the following behaviour is kinda weird but the URLs will still work if we append a file name
        expect(url.getURLBaseDir(a + "//")).toBe(a + "/")
        expect(url.getURLBaseDir(a + "///")).toBe(a + "//")

        const b = "https://example.com"
        expect(url.getURLBaseDir(b + "/")).toBe(b)
        expect(url.getURLBaseDir(b)).toBe(b)

        const c = "//example.com"
        expect(url.getURLBaseDir(c + "/")).toBe(c)
        expect(url.getURLBaseDir(c)).toBe(c)
    })
})

describe("fixupRelativeURL", () => {
    it("returns null for falsy URLs", () => {
        expect(url.fixupRelativeURL("", "foo", "bar")).toBeNull()
        expect(url.fixupRelativeURL(null, "foo", "bar")).toBeNull()
    })

    it("resolves anchor links with the article URL", () => {
        expect(url.fixupRelativeURL("#foo-bar", base, anchorBase)).toBe(anchorBase + "#foo-bar")
        expect(url.fixupRelativeURL("#foo-bar", "", anchorBase)).toBe(anchorBase + "#foo-bar")
        expect(url.fixupRelativeURL("#", "", anchorBase)).toBe(anchorBase + "#")
        expect(url.fixupRelativeURL("#/", "", anchorBase)).toBe(anchorBase + "#/")
        expect(url.fixupRelativeURL("#//", "", anchorBase)).toBe(anchorBase + "#//")
        expect(url.fixupRelativeURL("#///", "", anchorBase)).toBe(anchorBase + "#///")
        expect(url.fixupRelativeURL("#////", "", anchorBase)).toBe(anchorBase + "#////")
        expect(url.fixupRelativeURL("#/feed/42", "", anchorBase)).toBe(anchorBase + "#/feed/42")
        expect(url.fixupRelativeURL("#../feed/42", "", anchorBase)).toBe(anchorBase + "#../feed/42")
        expect(url.fixupRelativeURL("#http://feed/42", "", anchorBase)).toBe(anchorBase + "#http://feed/42")

        // can't do it without article URL
        expect(url.fixupRelativeURL("#foo-bar", base, "")).toBe("#foo-bar")
    })

    it("returns anchor links as-is when no article URL is set", () => {
        expect(url.fixupRelativeURL("#foo-bar", base, "")).toBe("#foo-bar")
        expect(url.fixupRelativeURL("#foo-bar", "", "")).toBe("#foo-bar")
        expect(url.fixupRelativeURL("#", "", "")).toBe("#")
        expect(url.fixupRelativeURL("#/", "", "")).toBe("#/")
        expect(url.fixupRelativeURL("#//", "", "")).toBe("#//")
        expect(url.fixupRelativeURL("#///", "", "")).toBe("#///")
        expect(url.fixupRelativeURL("#////", "", "")).toBe("#////")
        expect(url.fixupRelativeURL("#/feed/23", "", "")).toBe("#/feed/23")
    })

    it("returns triple-slash links as-is", () => {
        expect(url.fixupRelativeURL("///foobar", base, anchorBase)).toBe("///foobar")
        expect(url.fixupRelativeURL("///foo", base, "")).toBe("///foo")
        expect(url.fixupRelativeURL("///bar", "", "")).toBe("///bar")
        expect(url.fixupRelativeURL("///", "", "")).toBe("///")
    })

    it("resolves protocol-relative links using the feed URL' protocol", () => {
        const cdn = "//cdn.example.com/foo?1=2&a=x#f-fF0"
        expect(url.fixupRelativeURL(cdn, base, anchorBase)).toBe(`https:${cdn}`)
        expect(url.fixupRelativeURL(cdn, base, "")).toBe(`https:${cdn}`)
        const evilBase = "http://ancient.example.com/"
        expect(url.fixupRelativeURL(cdn, evilBase, anchorBase)).toBe(`http:${cdn}`)
        expect(url.fixupRelativeURL(cdn, evilBase, "")).toBe(`http:${cdn}`)
        const coolBase = "gopher://ancient.example.com/"
        expect(url.fixupRelativeURL(cdn, coolBase, anchorBase)).toBe(`gopher:${cdn}`)
        expect(url.fixupRelativeURL(cdn, coolBase, "")).toBe(`gopher:${cdn}`)

        // can't do it when no feed URL is set
        expect(url.fixupRelativeURL(cdn, "", "")).toBe(cdn)
    })

    it("returns root-relative links using the feed URL as base", () => {
        const a = "/assets/food.jpg"
        expect(url.fixupRelativeURL(a, base, anchorBase)).toBe("https://feeds.example.com/assets/food.jpg")
        expect(url.fixupRelativeURL(a, base, "")).toBe("https://feeds.example.com/assets/food.jpg")
        const b = "/assets/food.jpg?a=b&c=d#x"
        expect(url.fixupRelativeURL(b, base, anchorBase)).toBe("https://feeds.example.com" + b)

        // can't do it when no feed URL is set
        expect(url.fixupRelativeURL(a, "", "")).toBe(a)
    })

    it("returns relative links using the feed URL as base", () => {
        const a = "food.jpg"
        expect(url.fixupRelativeURL(a, base, anchorBase)).toBe("https://feeds.example.com/~ff/feed/yo/food.jpg")
        expect(url.fixupRelativeURL(a, base, "")).toBe("https://feeds.example.com/~ff/feed/yo/food.jpg")
        const b = "../food.jpg"
        expect(url.fixupRelativeURL(b, base, anchorBase)).toBe("https://feeds.example.com/~ff/feed/yo/../food.jpg")
        const c = "../food.jpg?bar=baz#foo"
        expect(url.fixupRelativeURL(c, base, anchorBase)).toBe("https://feeds.example.com/~ff/feed/yo/" + c)

        // can't do it when no feed URL is set
        expect(url.fixupRelativeURL(a, "", "")).toBe(a)
    })

    it("returns mailto: links unchanged", () => {
        const a = "mailto:foo@example.com"
        const b = "mailto:someone@example.com?subject=Hello&body=Hello%20World"
        expect(url.fixupRelativeURL(a, base, anchorBase)).toBe(a)
        expect(url.fixupRelativeURL(b, base, anchorBase)).toBe(b)
    })

    it("returns other procotol links unchanged", () => {
        // https://en.wikipedia.org/wiki/List_of_URI_schemes
        const a = "aaa://host.example.com:1813;transport=udp;protocol=radius"
        const b = "facetime://+19995551234"
        const c = "git://github.com/user/project-name.git"
        const d = "market://search?q=pub:Publisher_Name"
        const e = "spotify:track:2jCnn1QPQ3E8exampleNsx"
        expect(url.fixupRelativeURL(a, base, anchorBase)).toBe(a)
        expect(url.fixupRelativeURL(b, base, anchorBase)).toBe(b)
        expect(url.fixupRelativeURL(c, base, anchorBase)).toBe(c)
        expect(url.fixupRelativeURL(d, base, anchorBase)).toBe(d)
        expect(url.fixupRelativeURL(e, base, anchorBase)).toBe(e)
    })
})
