import { ArticleCleaner } from "./ArticleCleaner"

import { ImageProxy } from "./ImageProxy"

const imageProxy = new ImageProxy()
const ac = new ArticleCleaner(imageProxy)

const someBaseURL = "https://example.com/foo/bar/feed.rss"

const contentWithAnchorLinks = `
<ul>
<li><a href="#background" rel="nofollow">A bit of background explanation</a></li>
<li><a href="#currentmembers" rel="nofollow">What this means for creators who had work(s) on Of Elves and Men</a></li>
<li><a href="#questions" rel="nofollow">And what to do if you still have questions</a></li>
</ul>
<h3><a name="background" rel="nofollow" id="background"></a>Background explanation</h3>
<h3><a name="currentmembers" rel="nofollow" id="currentmembers"></a>What does this mean for creators who had work(s) on Of Elves and Men?</h3>
<h3><a name="questions" rel="nofollow" id="questions"></a>If you still have questions...</h3>
`

const contentWithImage = `
<p>
    <strong>Die Entwicklung der Kryptowährungen und ihr Einfluss auf die Gesellschaft</strong>
</p>
<div>
  <a target="_blank" href="https://media.example.com/media/stuff/stuff224-stuff-stuff-2.jpg">

<img alt="Episode image forSTUFF224 Stuff 2" width="200" src="https://example.com/podlove/image/stuff" />

  </a>
  <p>
  In Fortsetzung des Gesprächs ...
  </p>
`

const contentWithVideo = `
<p>I couldn't find an example for this so here's an artificial one:</p>
<video width="320" height="240" controls>
  <source src="https://movies.example.com/movie.mp4" type="video/mp4">
  <source src="https://movies.example.com/movie.ogg" type="video/ogg">
  Ancient browser alarm.
</video>
`

const contentWithVideoRelative = `
<p>I couldn't find an example for this so here's an artificial one:</p>
<video width="320" height="240" controls>
  <source src="movie.mp4" type="video/mp4">
  <source src="movie.ogg" type="video/ogg">
  Ancient browser alarm.
</video>
`

it("parses empty html", () => {
    expect(ac.cleanupHTML({ html: "" })).toBe("")
    expect(ac.cleanupHTML({ html: "", base: "foo" })).toBe("")
    expect(ac.cleanupHTML({ html: "", parseEnclosures: true })).toBe("")
    expect(ac.cleanupHTML({ html: "", base: "foo", parseEnclosures: true })).toBe("")
})

it("parses plain text", () => {
    const text = "foo bar baz\nboo"
    expect(ac.cleanupHTML({ html: text })).toBe(text)
})

it("retains and resolves anchor links", () => {
    // test case for issue #42
    const base = "https://example.com/feed.rss"
    const anchorBase = "https://example.com/article/23?bar=baz"
    const result = ac.cleanupHTML({ html: contentWithAnchorLinks, base, anchorBase })
    expect(result).toContain("A bit of background explanation")
    expect(result).toContain(`href="${anchorBase}#background"`)
})

it("removes A tags without href", () => {
    const result = ac.cleanupHTML({ html: contentWithAnchorLinks })
    expect(result).toContain("<h3>Background explanation</h3>")
    expect(result).not.toContain("id=")
})

it("sets articleContainsImage=false for content without images", () => {
    const _resultWithout = ac.cleanupHTML({ html: contentWithAnchorLinks })
    expect(ac.articleContainsImage).toBeFalsy()
})

it("sets articleContainsImage=true for content with an image", () => {
    const _resultWith = ac.cleanupHTML({ html: contentWithImage })
    expect(ac.articleContainsImage).toBeTruthy()
})

it("doesn't extract enclosures for content without audio/video/picture tags", () => {
    const _resultWithoutImage = ac.cleanupHTML({ html: contentWithAnchorLinks, parseEnclosures: true })
    expect(ac.extractedEnclosures.length).toBe(0)

    const _resultWithImage = ac.cleanupHTML({ html: contentWithImage, parseEnclosures: true })
    // img isn't an enclosure. only picture
    expect(ac.extractedEnclosures.length).toBe(0)
})

it("extracts an enclosure for content with a video tag", () => {
    const _resultWithVideo = ac.cleanupHTML({ html: contentWithVideo, parseEnclosures: true })
    expect(ac.extractedEnclosures.length).toBe(2)
    expect(ac.extractedEnclosures[0].url).toBe("https://movies.example.com/movie.mp4")
    expect(ac.extractedEnclosures[0].type).toBe("video/mp4")
    expect(ac.extractedEnclosures[1].url).toBe("https://movies.example.com/movie.ogg")
    expect(ac.extractedEnclosures[1].type).toBe("video/ogg")
})

it("extracts an enclosure for content with a video tag with relative URLs", () => {
    const _resultWithVideo = ac.cleanupHTML({ html: contentWithVideoRelative, parseEnclosures: true, base: someBaseURL })
    expect(ac.extractedEnclosures.length).toBe(2)
    expect(ac.extractedEnclosures[0].url).toBe("https://example.com/foo/bar/movie.mp4")
    expect(ac.extractedEnclosures[0].type).toBe("video/mp4")
    expect(ac.extractedEnclosures[1].url).toBe("https://example.com/foo/bar/movie.ogg")
    expect(ac.extractedEnclosures[1].type).toBe("video/ogg")
})
