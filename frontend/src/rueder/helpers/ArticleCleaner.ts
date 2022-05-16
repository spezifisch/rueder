export {}

import sanitizeHtml from "sanitize-html"

import { cleanupURL, getURLProtocol, fixupRelativeURL } from "./url"
import { Enclosure } from "../api/types"
import { ImageProxyType } from "./ImageProxy"
import type { ImageProxy } from "./ImageProxy"

export class ArticleCleaner {
    private imageProxy: ImageProxy

    extractedEnclosures = []
    articleContainsImage = false

    constructor(imageProxy: ImageProxy) {
        this.imageProxy = imageProxy
    }

    // cleanupHTML implements a kind of tracking blocker and rubbish cleanup for article content
    cleanupHTML({ html, base = "", anchorBase = "", parseEnclosures = false }): string {
        // mark these variables invalid before they get set after sanitizeHtml() in case anything goes wrong
        this.articleContainsImage = undefined
        this.extractedEnclosures = undefined
        // these get set by closures
        let articleContainsImage = false
        let extractedEnclosures = []

        const imageProxy = this.imageProxy
        const allowedSchemes = ["http", "https"]

        let safe = sanitizeHtml(html, {
            parser: {
                lowerCaseTags: true,
            },
            allowedTags: sanitizeHtml.defaults.allowedTags.concat(["img", "audio", "video", "picture"]),
            allowedAttributes: {
                a: ["href", "name", "rel", "target"],
                img: ["src", "alt", "title", "height", "width"],
                audio: ["controls", "crossorigin", "duration", "loop", "muted"], // no autoplay, preload
                video: ["buffered", "controls", "crossorigin", "duration", "loop", "muted", "height", "width"], // no autoplay, preload
            },
            allowedSchemes,
            allowProtocolRelative: false,
            transformTags: {
                a: function (tagName, attribs) {
                    // don't send referer when clicking, and open in new tab
                    attribs["rel"] = "nofollow noopener noreferrer"
                    attribs["target"] = "_blank"

                    if (attribs["href"]) {
                        // rewrite relative URLs
                        attribs["href"] = fixupRelativeURL(attribs["href"], base, anchorBase)

                        const urlProtocol = getURLProtocol(attribs["href"])
                        if (!allowedSchemes.includes(urlProtocol.toLowerCase())) {
                            // forbidden protocol
                            delete attribs["href"]
                        } else if (blockedLink(attribs["href"])) {
                            // adblock/social block
                            delete attribs["href"]
                        } else {
                            // remove tracking parameters
                            attribs["href"] = cleanupURL(attribs["href"])
                        }
                    }

                    return {
                        tagName,
                        attribs,
                    }
                },
                img: function (tagName, attribs) {
                    // just a trick to hide the article image which most of the time is the same as the first image in the article (eg. CRE)
                    articleContainsImage = true

                    if (attribs["src"]) {
                        if (blockedLink(attribs["src"])) {
                            // ublock it
                            delete attribs["src"]
                        } else {
                            // benevolent content, add imageproxy prefix
                            attribs["src"] = imageProxy.buildURL(
                                ImageProxyType.Content,
                                cleanupURL(attribs["src"])
                            )
                        }
                    }

                    if (!attribs["title"] && attribs["alt"]) {
                        // eg. for xkcd where title isn't set
                        attribs["title"] = attribs["alt"]
                    }

                    if (attribs["height"] == "1" && attribs["width"] == "1") {
                        // golem is nice (or evil) enough to tell us the dimensions of their tracking pixel
                        delete attribs["src"]
                        delete attribs["alt"]
                    }

                    // load only when viewed
                    attribs["loading"] = "lazy"
                    return {
                        tagName,
                        attribs,
                    }
                },
                audio: function (tagName, attribs) {
                    // content is still loaded directly from source
                    // also audio/video tags are removed by the exclusiveFilter below, but in case we want to use these tags someday...
                    attribs["preload"] = "none"
                    return {
                        tagName,
                        attribs,
                    }
                },
                video: function (tagName, attribs) {
                    // content is still loaded directly from source
                    attribs["preload"] = "none"
                    return {
                        tagName,
                        attribs,
                    }
                },
                source: function (tagName, attribs) {
                    // audio/video:
                    // * currently we move all audio/video sources into enclosures and remove the source tag and the parent audio/video tag
                    //
                    // picture:
                    // * picture tags must provide an img child as fallback, therefore we currently just remove the source tags there too
                    // * to parse their sources properly we would need to
                    //   - keep the srcset (and media) attributes
                    //   - parse it and check all urls for adblock and append imageproxy prefixes
                    // * we don't do this because it's pretty finicky and it would leak device specs (screen size) to the feed owner
                    if (parseEnclosures) {
                        const enclosure = new Enclosure()
                        enclosure.type = attribs["type"].split(";").shift()
                        enclosure.url = fixupRelativeURL(attribs["src"], base, null /* no anchors allowed */)

                        if (enclosure.url) {
                            // for sources inside of audio and video tags -> extract enclosures
                            const urlProtocol = getURLProtocol(enclosure.url)
                            if (allowedSchemes.includes(urlProtocol.toLowerCase())) {
                                if (!blockedLink(enclosure.url)) {
                                    extractedEnclosures = [...extractedEnclosures, enclosure]
                                }
                            }
                        }
                    }

                    return {
                        tagName,
                        attribs,
                    }
                },
            },
            exclusiveFilter: function (frame) {
                const tag = frame.tag.toLowerCase()

                // remove audio/video tags after they've been added to the enclosure list and all source tags
                if (["audio", "video", "source"].includes(tag)) {
                    return true
                }
                // remove tags that are useless after the filtering above and generally
                if (tag == "img" && !frame.attribs["src"]) {
                    return true
                }
                if (tag == "a" && !frame.attribs["href"]) {
                    return true
                }
                if (tag == "span" && (!frame.text || !frame.text.trim())) {
                    return true
                }

                return false
            },
            textFilter: function (text: string, tagName: string) {
                if (!tagName) {
                    return text
                }
                const tag = tagName.toLowerCase()

                if (tag == "p" || tag == "span") {
                    text = text.replaceAll("\n", "")
                    text = text.replaceAll("\t", " ")
                    text = text.replaceAll("  ", "")
                }

                return text
            },
        })

        // deliver results to our caller
        this.articleContainsImage = articleContainsImage
        this.extractedEnclosures = extractedEnclosures

        safe = safe.replaceAll("<p></p>", "").replaceAll("<br /><br />", "<br />")
        return safe
    }
}

// social and crap link filter
function blockedLink(src: string): boolean {
    src = src.toLowerCase()
    if (!src.startsWith("http://") && !src.startsWith("https://")) {
        return true
    }
    const domain = src.split("/", 3)[2]
    if (!domain) {
        return true
    }

    // blocked domains including all subdomains
    const blockedDomains: string[] = ["fsdn.com", "facebook.com"]
    for (const dom of blockedDomains) {
        if (domain == dom || domain.endsWith("." + dom)) {
            return true
        }
    }

    if (domain == "twitter.com" || domain.endsWith(".twitter.com")) {
        if (src.includes("status=")) {
            // let's be fuzzy. it's some kind of "share as twitter status" link, eg. on slashdot
            return true
        }
    }

    if (domain == "feeds.feedburner.com") {
        if (src.includes(`${domain}/~ff/`)) {
            // social links, eg. on postillon, mydealz
            return true
        }
    }

    return false
}
