import { JsonProperty, JsonClassType, JsonIgnore } from "jackson-js"

import { ArticlePreview } from "../api/types"
import { SyncedStore } from "./syncedstore"

export class LabelStore extends SyncedStore<labelData> {
    constructor() {
        super("labels", labelData)
        this.data.initDerived()

        const l = this.getLabels().length
        if (l == 0) {
            this.addLabel("important", "#EF4444")
            this.addLabel("favorite", "#F59E0B")
            this.addLabel("read later", "#10B981")
        }
    }

    sync(): void {
        this.set(this.data)
    }

    getLabels(): Label[] {
        if (!this.data) {
            return []
        }
        return this.data.orderedLabels
    }

    getLabel(name: string): Label {
        if (!this.data) {
            return
        }
        return this.data.labels.get(name)
    }

    labelExists(name: string): boolean {
        if (!this.data) {
            return false
        }
        return this.data.labels.has(name)
    }

    addLabel(name: string, color: string): boolean {
        if (this.data.labels.has(name)) {
            console.error("label exists", name)
            return false
        }
        const l = new Label(name, color)
        this.data.labels.set(name, l)
        this.data.orderedLabels.push(l)
        this.data.sortOrderedLabels()
        this.sync()
        return true
    }

    setLabelName(label: string, newName: string): boolean {
        if (label == newName || !newName || newName.length < 1) {
            console.error("invalid new name", newName)
            return false
        }
        if (!this.data.labels.has(label)) {
            console.error("label doesn't exist", label)
            return false
        }
        if (this.data.labels.has(newName)) {
            console.error("new label name already exists", newName)
            return false
        }
        const oldLabel = this.data.labels.get(label)
        // change article label list
        for (const articleID of oldLabel.articleIDs) {
            this.data.labelNamesByArticleID.get(articleID).add(newName)
            this.data.labelNamesByArticleID.get(articleID).delete(label)
        }
        // update orderedLabels
        const olIdx = this.data.orderedLabels.findIndex((l: Label) => l.name == label)
        if (olIdx != -1) {
            this.data.orderedLabels[olIdx].name = newName
        }
        // change label name
        oldLabel.name = newName // this should be redundant given we replace the name already in orderedLabels
        this.data.labels.set(newName, oldLabel)
        // fix order
        this.data.sortOrderedLabels()
        this.sync()
        return true
    }

    setLabelColor(label: string, color: string): boolean {
        if (!this.data.labels.has(label)) {
            console.error("label doesn't exist", label)
            return false
        }
        this.data.labels.get(label).color = color
        this.sync()
        return true
    }

    removeLabel(label: string): void {
        if (!this.data.labels.has(label)) {
            console.error("label doesn't exist", label)
            return
        }
        for (const articleID of this.data.labels.get(label).articleIDs) {
            this.removeArticleFromLabel(label, articleID)
        }
        this.data.labels.delete(label)
        this.data.orderedLabels = this.data.orderedLabels.filter((l: Label) => l.name != label)
        this.sync()
    }

    labelContains(label: string, articleID: string): boolean {
        if (!this.data.labels.has(label)) {
            console.error("label doesn't exist", label)
            return false
        }
        return this.data.labels.get(label).articleIDs.has(articleID)
    }

    addArticleToLabel(label: string, articlePreview: ArticlePreview): void {
        if (!this.data.labels.has(label)) {
            console.error("label doesn't exist", label)
            return
        }
        const articleID = articlePreview.id
        if (!this.labelContains(label, articleID)) {
            // store AP if needed
            if (!this.data.articles.has(articleID)) {
                this.data.articles.set(articleID, articlePreview)
            }

            // append to label's article list
            this.data.labels.get(label).articleIDs.add(articleID)

            // append to article's label list
            if (!this.data.labelNamesByArticleID.has(articleID)) {
                this.data.labelNamesByArticleID.set(articleID, new Set<string>())
            }
            this.data.labelNamesByArticleID.get(articleID).add(label)

            this.sync()
        }
    }

    removeArticleFromLabel(label: string, articleID: string): void {
        if (!this.data.labels.has(label)) {
            console.error("label doesn't exist", label)
            return
        }
        if (this.labelContains(label, articleID)) {
            // remove from label's articles list
            this.data.labels.get(label).articleIDs.delete(articleID)
            // do not delete empty labels here. we want to keep them.

            // remove from article's label list
            this.data.labelNamesByArticleID.get(articleID).delete(label)

            // remove article if it's not labelled anymore
            if (this.data.labelNamesByArticleID.get(articleID).size == 0) {
                this.data.labelNamesByArticleID.delete(articleID)
                this.data.articles.delete(articleID)
            }

            this.sync()
        }
    }

    getArticleLabels(articleID: string): Label[] {
        return this.data.getLabelsByArticleID(articleID)
    }

    getArticlePreview(articleID: string): ArticlePreview {
        return this.data.getArticlePreview(articleID)
    }

    getLabelArticles(label: string): ArticlePreview[] {
        if (!this.data.labels.has(label)) {
            console.error("label doesn't exist", label)
            return
        }
        return [...this.data.labels.get(label).articleIDs.values()].map((articleID: string) =>
            this.data.articles.get(articleID)
        )
    }
}

class labelData {
    @JsonProperty()
    @JsonClassType({ type: () => [Map, [String, ArticlePreview]] })
    articles: Map<string, ArticlePreview> // maps article id to ArticlePreview

    @JsonProperty()
    @JsonClassType({ type: () => [Map, [String, Label]] })
    labels: Map<string, Label> // maps label name to Label

    @JsonIgnore()
    labelNamesByArticleID: Map<string, Set<string>> // maps article id to label name list

    @JsonIgnore()
    orderedLabels: Label[] // labels ordered by name

    constructor() {
        this.articles = new Map()
        this.labels = new Map()
        this.labelNamesByArticleID = new Map()
        this.orderedLabels = []
    }

    // initialize derived object values
    initDerived() {
        // for quick lookup of labels by article id
        this.labelNamesByArticleID = new Map()
        for (const label of this.labels.values()) {
            for (const articleID of label.articleIDs.values()) {
                if (this.labelNamesByArticleID.has(articleID)) {
                    this.labelNamesByArticleID.get(articleID).add(label.name)
                } else {
                    this.labelNamesByArticleID.set(articleID, new Set([label.name]))
                }
            }
        }

        // for quick listing of available labels
        this.orderedLabels = [...this.labels.values()]
        this.sortOrderedLabels()
    }

    sortOrderedLabels() {
        this.orderedLabels.sort((a, b) => {
            // order by name ignoring case
            const nameA = a.name.toLowerCase()
            const nameB = b.name.toLowerCase()
            if (nameA < nameB) {
                return -1
            }
            if (nameA > nameB) {
                return 1
            }
            return 0
        })
    }

    getLabelsByArticleID(articleID: string): Label[] {
        if (!this.labelNamesByArticleID.has(articleID)) {
            return []
        }
        return [...this.labelNamesByArticleID.get(articleID).values()].map((labelName: string) =>
            this.labels.get(labelName)
        )
    }

    getArticlePreview(articleID: string): ArticlePreview {
        return this.articles.get(articleID)
    }
}

export class Label {
    @JsonProperty()
    @JsonClassType({ type: () => [String] })
    name: string
    @JsonProperty()
    @JsonClassType({ type: () => [String] })
    color: string
    @JsonProperty()
    @JsonClassType({ type: () => [Set, [String]] })
    articleIDs: Set<string>

    constructor(name: string, color: string) {
        this.name = name
        this.color = color
        this.articleIDs = new Set<string>()
    }
}

export const labelStore = new LabelStore()
