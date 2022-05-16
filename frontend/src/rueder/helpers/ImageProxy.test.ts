import { ImageProxy, ImageProxyType } from "./ImageProxy"

const proxyURL = "https://ip.example.com/yadda/"
const imageURL = "https://images.example.com/foo/bar/test.png"

it("pipes URLs through if no baseURL is set", () => {
    const ip = new ImageProxy()

    expect(ip.buildURL(ImageProxyType.Icon, imageURL)).toBe(imageURL)
    expect(ip.buildURL(ImageProxyType.Thumbnail, imageURL)).toBe(imageURL)
    expect(ip.buildURL(ImageProxyType.Content, imageURL)).toBe(imageURL)
    expect(ip.buildURL(null, imageURL)).toBe(imageURL)

    const ip_ex = new ImageProxy(null)
    expect(ip_ex.buildURL(ImageProxyType.Icon, imageURL)).toBe(imageURL)
})

it("applies type prefixes if useTypePrefixes is set", () => {
    const ip = new ImageProxy(proxyURL, true)

    let output = ip.buildURL(ImageProxyType.Icon, imageURL)
    expect(output).toContain(proxyURL + "icon/")
    output = ip.buildURL(ImageProxyType.Thumbnail, imageURL)
    expect(output).toContain(proxyURL + "thumbnail/")
    output = ip.buildURL(ImageProxyType.Content, imageURL)
    expect(output).toContain(proxyURL + "image/")
})

it("doesn't apply type prefixes if useTypePrefixes is not set", () => {
    const ip = new ImageProxy(proxyURL)

    let output = ip.buildURL(ImageProxyType.Icon, imageURL)
    expect(output).not.toContain(proxyURL + "icon/")
    output = ip.buildURL(ImageProxyType.Thumbnail, imageURL)
    expect(output).not.toContain(proxyURL + "thumbnail/")
    output = ip.buildURL(ImageProxyType.Content, imageURL)
    expect(output).not.toContain(proxyURL + "image/")
})

it("applies imageproxy presets", () => {
    const ip = new ImageProxy(proxyURL, false)

    let output = ip.buildURL(ImageProxyType.Icon, imageURL)
    expect(output).toContain("/insecure/icon/")
    output = ip.buildURL(ImageProxyType.Thumbnail, imageURL)
    expect(output).toContain("/insecure/thumbnail/")
    output = ip.buildURL(ImageProxyType.Content, imageURL)
    expect(output).toContain("/insecure/image/")
})

it("applies the obfuscation key if set", () => {
    const ip = new ImageProxy(proxyURL, false, "cool_key", "cool_salt")

    let output = ip.buildURL(ImageProxyType.Icon, imageURL)
    expect(output).not.toContain("/insecure/")
    expect(output).not.toContain(imageURL)
    output = ip.buildURL(ImageProxyType.Thumbnail, imageURL)
    expect(output).not.toContain("/insecure/")
    expect(output).not.toContain(imageURL)
    output = ip.buildURL(ImageProxyType.Content, imageURL)
    expect(output).not.toContain("/insecure/")
    expect(output).not.toContain(imageURL)
})
