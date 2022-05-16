import sxml from "sxml"

export class OPML {
    head: Header
    body: Body

    constructor(values: object = {}) {
        Object.assign(this, values)
    }

    toXML(): sxml.XML {
        const ret = new sxml.XML()
        ret.setTag("opml")
        ret.setProperty("version", "2.0")
        ret.push(this.head.toXML())
        ret.push(this.body.toXML())
        return ret
    }
}

export class Header {
    title?: string
    ownerName?: string
    dateCreated?: string | Date
    dateModified?: string | Date

    constructor(values: object = {}) {
        Object.assign(this, values)
    }

    toXML(): sxml.XML {
        const ret = new sxml.XML()
        ret.setTag("head")

        for (const prop in this) {
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const val: any = this[prop]
            if (val !== undefined) {
                const t = new sxml.XML()
                t.setTag(prop.toString())

                if (prop.startsWith("date") && val instanceof Date) {
                    t.setValue(val.toUTCString())
                } else {
                    t.setValue(val)
                }

                ret.push(t)
            }
        }

        return ret
    }
}

export class Body {
    outlines?: Outline[]

    constructor(values: object = {}) {
        this.outlines = []

        Object.assign(this, values)
    }

    toXML(): sxml.XML {
        const ret = new sxml.XML()
        ret.setTag("body")

        for (const i in this.outlines) {
            const outline = this.outlines[i].toXML()
            ret.push(outline)
        }

        return ret
    }
}

export class Outline {
    text?: string
    title?: string
    type?: string
    xmlUrl?: string
    htmlUrl?: string

    constructor(values: object = {}) {
        Object.assign(this, values)
    }

    toXML(): sxml.XML {
        const ret: sxml.XML = new sxml.XML()
        ret.setTag("outline")

        for (const prop in this) {
            if (this[prop] !== undefined) {
                ret.setProperty(prop, "" + this[prop])
            }
        }

        return ret
    }
}
