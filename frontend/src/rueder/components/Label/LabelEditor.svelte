<script lang="ts">
    import { onDestroy } from "svelte"
    import Select from "@spezifisch/svelte-select"
    import { difference as _difference } from "lodash"

    import type { Article } from "../../api/types"
    import { labelStore, Label } from "../../stores/labels"
    import { getSolidTag } from "../ui/heroicons/solid-tag"

    import CustomMultiSelection from "./CustomMultiSelection.svelte"
    import type { CustomEvent } from "./types"

    export let article: Article

    class Item {
        label: string
        value: string
        color: string
    }

    let items: Item[] = null
    let selectedValue: Item[] = null

    function labelToItem(label: Label): Item {
        return { label: label.name, value: label.name, color: label.color }
    }

    function getItemByLabel(label: Label): Item {
        const idx = items.findIndex((i: Item) => i.label == label.name)
        if (idx >= 0) {
            return items[idx]
        }
    }

    // called when the label store changes (ie labels are added or colors are changed somewhere else)
    const unsubscribe = labelStore.subscribe(() => {
        items = labelStore.getLabels().map((l: Label) => labelToItem(l))

        // set active labels for this article
        selectedValue = labelStore.getArticleLabels(article.id).map((l: Label) => getItemByLabel(l))
        if (!selectedValue.length) {
            selectedValue = null
        }
    })
    onDestroy(unsubscribe)

    // item content in dropdown list
    function getOptionLabel(option: any, filterText: string): string {
        const text = option.isCreator ? `Create \"${filterText}\"` : option.label
        const color = option.isCreator ? "#ffffff" : option.color
        const tag = getSolidTag(5, "inline")
        return `
        <div class="flex flex-row space-x-1">
            <div class="flex-none" style="color: ${color}">${tag}</div>
            <div>${text}</div>
        </div>`
    }

    // called when the user changes the label color with the color picker
    function handleLabelUpdate(event: any) {
        const update: CustomEvent = event.detail
        if (update.updateType == "color") {
            console.log("update label color of", update.labelName, update.hexColor)

            const ok = labelStore.setLabelColor(update.labelName, update.hexColor)
            if (ok) {
                const itemIdx = items.findIndex((i: Item) => i.label == update.labelName)
                if (itemIdx >= 0) {
                    items[itemIdx].color = update.hexColor
                }
            }
        } else if (update.updateType == "name") {
            console.log("update label name of", update.labelName, update.newName)

            const ok = labelStore.setLabelName(update.labelName, update.newName)
            if (ok) {
                const itemIdx = items.findIndex((i: Item) => i.label == update.labelName)
                if (itemIdx >= 0) {
                    items[itemIdx].label = update.labelName
                    items[itemIdx].value = update.labelName
                    // TODO replace selectedValue names with this
                }
            }
        }
    }

    // called when the user adds or removes labels from this article
    function handleSelect(event: any) {
        const update: Item[] = event.detail ?? []

        const existingLabelNames = labelStore.getLabels().map((l: Label) => l.name)
        const newLabelSet = update.map((i: Item) => i.label)
        const oldLabelSet = labelStore.getArticleLabels(article.id).map((l: Label) => l.name)

        for (const addLabel of _difference(newLabelSet, oldLabelSet)) {
            if (!existingLabelNames.includes(addLabel)) {
                // labels doesn't exist yet, create it
                const color = "#ffffff"
                labelStore.addLabel(addLabel, color)
                items.push({ label: addLabel, value: addLabel, color })
            }

            labelStore.addArticleToLabel(addLabel, article.toArticlePreview())
        }

        for (const removeLabel of _difference(oldLabelSet, newLabelSet)) {
            labelStore.removeArticleFromLabel(removeLabel, article.id)
        }
    }
</script>

<div class="my-2 themed">
    <Select
        MultiSelection={CustomMultiSelection}
        {items}
        {selectedValue}
        isMulti={true}
        isClearable={false}
        isCreatable={true}
        {getOptionLabel}
        on:focusChange
        on:customEvent={handleLabelUpdate}
        on:select={handleSelect}
        placeholder="Add Labels"
        noOptionsMessage="No more labels"
    />
</div>

<style>
    .themed {
        /* hide the fact that we're a selection thing by default */
        --background: transparent;
        --border: none;

        --multiSelectPadding: 0 16px;

        /* when dropdown is closed */
        --multiItemBG: #334155;
        --multiItemActiveBG: #334155;

        /* when dropdown is open */
        --listBackground: #334155;
        --itemHoverBG: #64748b;

        /* when a new label is added by typing */
        --inputColor: white;
    }
</style>
