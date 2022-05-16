// from https://www.thecodecampus.de/blog/typescript-delay-with-async-await/
export function sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms))
}

/**
 * Format bytes as human-readable text.
 * Source: https://stackoverflow.com/a/14919494
 *
 * @param bytes Number of bytes.
 * @param si True to use metric (SI) units, aka powers of 1000. False to use
 *           binary (IEC), aka powers of 1024.
 * @param dp Number of decimal places to display.
 *
 * @return Formatted string.
 */
export function humanFileSize(bytes, si = false, dp = 1) {
    const thresh = si ? 1000 : 1024

    if (Math.abs(bytes) < thresh) {
        return bytes + " B"
    }

    const units = si
        ? ["kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"]
        : ["KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"]
    let u = -1
    const r = 10 ** dp

    do {
        bytes /= thresh
        ++u
    } while (Math.round(Math.abs(bytes) * r) / r >= thresh && u < units.length - 1)

    return bytes.toFixed(dp) + " " + units[u]
}
