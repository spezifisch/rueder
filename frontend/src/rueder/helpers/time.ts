import TimeAgo from "javascript-time-ago"
import en from 'javascript-time-ago/locale/en.json'

TimeAgo.addDefaultLocale(en)

const timeAgo = new TimeAgo("en-US")

function isToday(val: Date, now?: Date): boolean {
    if (!now) {
        now = new Date()
    }
    return now.getDate() == val.getDate() && now.getMonth() == val.getMonth() && now.getFullYear() == val.getFullYear()
}

function isThisYear(val: Date, now?: Date): boolean {
    if (!now) {
        now = new Date()
    }
    return now.getFullYear() == val.getFullYear()
}

export function localizeTime(val: string, hideYear = true, now?: Date): DateTimeString {
    const d = new Date(val)

    const timeStr = d.toLocaleTimeString("de-DE", {
        hour12: false,
        hour: "numeric",
        minute: "numeric",
    })

    hideYear = hideYear && isThisYear(d, now)
    const dateStr = d.toLocaleDateString("de-DE", {
        year: hideYear ? undefined : "numeric",
        month: "2-digit",
        day: "2-digit",
    })

    return new DateTimeString({
        date: dateStr,
        time: timeStr,
        isToday: isToday(d, now),
    })
}

export function timeAgoString(val: string): string {
    const date = Date.parse(val)
    if (isNaN(date)) {
        return ""
    }
    // we are sure format() returns a string because options.getTimeToNextUpdate is not set
    return timeAgo.format(date) as string
}

export class DateTimeString {
    date?: string
    time?: string

    isToday: boolean

    constructor(val: object = {}) {
        Object.assign(this, val)
    }

    toString(hideToday = false): string {
        if (hideToday && this.isToday) {
            return this.time
        }
        return `${this.date} ${this.time}`
    }
}
