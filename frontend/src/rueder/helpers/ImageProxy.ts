import { ImgProxy } from "@spezifisch/imgproxyjs"
import urljoin from "url-join"

export class ImageProxy {
    useProxy: boolean
    useTypePrefixes: boolean
    baseURL: string
    imgproxy: ImgProxy

    constructor(baseURL?: string, useTypePrefixes?: boolean, imageProxyKey?: string, imageProxySalt?: string) {
        this.useProxy = !!baseURL
        this.baseURL = baseURL
        this.useTypePrefixes = !!useTypePrefixes

        const key = imageProxyKey ?? undefined
        const salt = imageProxySalt ?? undefined
        this.imgproxy = new ImgProxy({
            url: "", // join with baseURL later to avoid double slashes
            key,
            salt,
            autoreset: true,
            presetOnly: true,
        })
    }

    buildURL(type: ImageProxyType, imageURL: string): string {
        if (!this.useProxy) {
            return imageURL
        }

        const ipURL = this.imgproxy.preset(type).get(imageURL)
        if (this.useTypePrefixes) {
            return urljoin(this.baseURL, type, ipURL)
        }
        return urljoin(this.baseURL, ipURL)
    }
}

export enum ImageProxyType {
    /* eslint-disable no-unused-vars */
    Icon = "icon", // 256-512px
    Thumbnail = "thumbnail",
    Content = "image", // full
    /* eslint-enable no-unused-vars */
}
