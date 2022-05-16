import { render } from "@testing-library/svelte"
import { writable } from "svelte/store"
import { tick } from "svelte"

import TestContext from "../../test/TestContext.svelte"
import TestEmpty from "../../test/TestEmpty.svelte"
import ArticleContent from "./ArticleContent.svelte"

import { Article } from "../api/types"
import { contextKey } from "../helpers/constants"
import { ImageProxy } from "../helpers/ImageProxy"

jest.mock("./Label/LabelEditor.svelte", () => ({
    default: TestEmpty,
}))

it("shows loading screen", () => {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    const store = writable(new Promise<Article>(() => {}))
    const { getByTitle } = render(ArticleContent, {
        article_content: store,
        isFocused: false,
    })

    const back = getByTitle("Back to article list")
    expect(back.firstElementChild.tagName).toBe("svg")

    const close = getByTitle("Close Article")
    expect(close.firstElementChild.tagName).toBe("svg")

    const loading = getByTitle("Loading")
    expect(loading.firstElementChild.classList).toContain("loader")
})

it("shows error screen", async () => {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    const store = writable(new Promise<Article>(() => {}))
    const props = {
        article_content: store,
        isFocused: false,
    }
    const { getByTitle, getByText, rerender } = render(ArticleContent, props)
    // mock failed article fetch
    store.set(Promise.reject("test"))
    rerender(props) // for some reason the DOM isn't updated unless we call this
    await tick() // and for some reason the only tick needed is exactly here

    const back = getByTitle("Back to article list")
    expect(back.firstElementChild.tagName).toBe("svg")

    const close = getByTitle("Close Article")
    expect(close.firstElementChild.tagName).toBe("svg")

    expect(getByText("Failed fetching article")).toBeTruthy()
})

it("shows empty article", async () => {
    const store = writable(Promise.resolve(new Article()))
    const props = {
        article_content: store,
        isFocused: false,
    }
    const { getByTitle, getByText } = render(ArticleContent, props)
    await tick() // this is necessary to not get the "Loading" screen

    const back = getByTitle("Back to article list")
    expect(back.firstElementChild.tagName).toBe("svg")

    const close = getByTitle("Close Article")
    expect(close.firstElementChild.tagName).toBe("svg")

    const content = getByText("No content.")
    expect(content).toBeTruthy()
    expect(content.tagName).toBe("P")
})

it("adds imageproxy prefix to article image", async () => {
    // set up mock image proxy
    const mockImageProxyURL = "https://imageproxy.example.com/foo/"
    const mockImageProxy = new ImageProxy(mockImageProxyURL, false)
    const spyBuildURL = jest.spyOn(mockImageProxy, "buildURL")

    // create article
    const mockImageURL = "https://example.com/image.jpg"
    const mockImageTitle = "foo bar"
    const article = new Article({ image: mockImageURL, image_title: mockImageTitle })
    const store = writable(Promise.resolve(article))

    // we need the TestContext wrapper to pass down context variables
    const { getByAltText } = render(TestContext, {
        Component: ArticleContent,
        // put our mock imageproxy in the context
        context_key: contextKey.imageProxy,
        context_value: mockImageProxy,
        // ArticleContent's props:
        article_content: store,
        isFocused: false,
    })
    await tick()

    expect(spyBuildURL).toHaveBeenCalledWith("image", mockImageURL)

    const img = getByAltText(mockImageTitle)
    const imgSrc = img.getAttribute("src")
    expect(imgSrc).toBe(spyBuildURL.mock.results[0].value)
    expect(imgSrc).toContain(mockImageProxyURL)

    spyBuildURL.mockRestore()
})

it("shows article without title and link correctly", async () => {
    const article = new Article({ title: undefined, link: undefined })
    const store = writable(Promise.resolve(article))
    const props = {
        article_content: store,
        isFocused: false,
    }
    const { getByText } = render(ArticleContent, props)
    await tick()

    const em = getByText("untitled")
    expect(em).toBeTruthy()
    expect(em.tagName).toBe("EM")
    expect(em.parentElement.tagName).toBe("H2")
})

it("shows article without title but with link correctly", async () => {
    const mockLink = "https://example.com/foo"
    const article = new Article({ title: undefined, link: mockLink })
    const store = writable(Promise.resolve(article))
    const props = {
        article_content: store,
        isFocused: false,
    }
    const { getByText } = render(ArticleContent, props)
    await tick()

    const em = getByText("untitled")
    expect(em).toBeTruthy()
    expect(em.tagName).toBe("EM")
    expect(em.parentElement.tagName).toBe("A")
    expect(em.parentElement.parentElement.tagName).toBe("H2")
})

it("shows article with title and link correctly", async () => {
    const mockTitle = "Turtle"
    const mockLink = "https://example.com/foo"
    const article = new Article({ title: mockTitle, link: mockLink })
    const store = writable(Promise.resolve(article))
    const props = {
        article_content: store,
        isFocused: false,
    }
    const { getByText } = render(ArticleContent, props)
    await tick()

    const a = getByText(mockTitle)
    expect(a).toBeTruthy()
    expect(a.tagName).toBe("A")
})

it("shows article with title and no link correctly", async () => {
    const mockTitle = "Turtle"
    const article = new Article({ title: mockTitle, link: undefined })
    const store = writable(Promise.resolve(article))
    const props = {
        article_content: store,
        isFocused: false,
    }
    const { getByText } = render(ArticleContent, props)
    await tick()

    const h2 = getByText(mockTitle)
    expect(h2).toBeTruthy()
    expect(h2.tagName).toBe("H2")
})

it("shows article with empty title correctly", async () => {
    const mockTitle = ""
    const mockLink = "https://example.com/foo"
    const article = new Article({ title: mockTitle, link: mockLink })
    const store = writable(Promise.resolve(article))
    const props = {
        article_content: store,
        isFocused: false,
    }
    const { getByTitle } = render(ArticleContent, props)
    await tick()

    const a = getByTitle("Go to external article page")
    expect(a).toBeTruthy()
    expect(a.tagName).toBe("A")
    expect(a.textContent).toBe("untitled")
})
