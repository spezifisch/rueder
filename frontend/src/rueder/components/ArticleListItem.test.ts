import { render } from "@testing-library/svelte"
import { tick } from "svelte"

import TestContext from "../../test/TestContext.svelte"
import ArticleListItem from "./ArticleListItem.svelte"

import { ArticlePreview } from "../api/types"
import { contextKey } from "../helpers/constants"
import { ImageProxy } from "../helpers/ImageProxy"
import { localizeTime } from "../helpers/time"

// stub it out because for some reason it gives typing problems
const mockSanitizedValue = "SANITIZED"
jest.mock("sanitize-html", () => ({
    default: jest.fn(() => mockSanitizedValue),
}))

// fake time
jest.useFakeTimers().setSystemTime(new Date("2022-05-15"))

function mockArticlePreview(): ArticlePreview {
    return new ArticlePreview({
        id: "8eebb085-baf9-4284-af1e-ee5858c9113e",
        seq: 1528,
        feed_seq: 526,
        title: "Mock: Test test lorem ipsum",
        time: "2022-04-10T11:09:00Z", // fake today date is set above
        feed_title: "jest online News",
        feed_icon: "https://example.com/icon.png",
        teaser: "banana",
    })
}
const todayTime = "2022-05-15T11:09:00Z" // a date that "is today" when used in the article

const mockImageProxyURL = "https://imageproxy.example.com/foo/"

it("shows feed icon when feed_icon is set", async () => {
    // set up mock image proxy
    const mockImageProxy = new ImageProxy(mockImageProxyURL, false)
    const spyBuildURL = jest.spyOn(mockImageProxy, "buildURL")

    const article = mockArticlePreview()
    const { getByAltText } = render(TestContext, {
        Component: ArticleListItem,
        context_key: contextKey.imageProxy,
        context_value: mockImageProxy,

        article,
    })
    await tick()

    expect(spyBuildURL).toHaveBeenCalledWith("icon", article.feed_icon)

    const icon = getByAltText(article.feed_title)
    expect(icon.tagName).toBe("IMG")
    const imgSrc = icon.getAttribute("src")
    expect(imgSrc).toBe(spyBuildURL.mock.results[0].value)
    expect(imgSrc).toContain(mockImageProxyURL)

    spyBuildURL.mockRestore()
})

it("shows dummy feed icon when feed_icon is missing", () => {
    const article = mockArticlePreview()
    article.feed_icon = undefined
    const { getByAltText } = render(ArticleListItem, {
        article,
    })

    const icon = getByAltText(article.feed_title)
    expect(icon.tagName).toBe("IMG")
    expect(icon.getAttribute("src")).toBe("dummy.png")
})

it("shows feed title, article date and time", () => {
    const article = mockArticlePreview()
    article.feed_icon = undefined // unset so we don't need the imageproxy context
    const { getByText } = render(ArticleListItem, {
        article,
    })

    const title = getByText(article.feed_title)
    expect(title).toBeTruthy()
    expect(title.tagName).toBe("DIV")

    const dateTimeStr = localizeTime(article.time)
    const date = getByText(dateTimeStr.date)
    expect(date).toBeTruthy()
    expect(date.tagName).toBe("SPAN")
    const time = getByText(dateTimeStr.time)
    expect(time).toBeTruthy()
    expect(time.tagName).toBe("SPAN")
})

it("hides article date and shows only time when it's today", () => {
    const article = mockArticlePreview()
    article.feed_icon = undefined
    article.time = todayTime
    const { getByText } = render(ArticleListItem, {
        article,
    })

    const title = getByText(article.feed_title)
    expect(title).toBeTruthy()
    expect(title.tagName).toBe("DIV")

    const dateTimeStr = localizeTime(article.time)
    expect(() => getByText(dateTimeStr.date)).toThrow()

    const time = getByText(dateTimeStr.time)
    expect(time).toBeTruthy()
    expect(time.tagName).toBe("SPAN")
})

it("strips html from article title", () => {
    const article = mockArticlePreview()
    article.feed_icon = undefined
    article.title = "<a href=\"foo\">bar<b>baz</b></a>"
    const { getByRole } = render(ArticleListItem, {
        article,
    })

    const title = getByRole("heading")
    expect(title).toBeTruthy()
    expect(title.innerHTML).toBe(mockSanitizedValue)
    // we don't want or need to test the external module sanitizeHtml here
})

it("shows 'untitled' if article title is empty", () => {
    const article = mockArticlePreview()
    article.feed_icon = undefined
    article.title = ""
    const { getByRole } = render(ArticleListItem, {
        article,
    })

    const title = getByRole("heading")
    expect(title).toBeTruthy()
    expect(title.textContent).toBe("untitled")
})


it("strips html from teaser", () => {
    const article = mockArticlePreview()
    article.feed_icon = undefined
    article.teaser = "<p><a href=\"foo\">bar<b>baz</b></a></p>"
    const { getByTestId } = render(ArticleListItem, {
        article,
    })

    const teaser = getByTestId("teaser")
    expect(teaser).toBeTruthy()
    expect(teaser.innerHTML).toBe(mockSanitizedValue)
})

it("hides teaser if empty", () => {
    const article = mockArticlePreview()
    article.feed_icon = undefined
    article.teaser = ""
    const { getByTestId } = render(ArticleListItem, {
        article,
    })

    expect(() => getByTestId("teaser")).toThrow()
})
