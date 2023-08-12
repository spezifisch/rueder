<script lang="ts">
    import { createEventDispatcher, onDestroy } from "svelte"

    import OutlineTag from "../ui/heroicons/outline-tag.svelte"
    import SolidTrash from "../ui/heroicons/solid-trash.svelte"
    import SolidPlusCircle from "../ui/heroicons/solid-plus-circle.svelte"
    import Spinner from "../ui/Spinner.svelte"

    import LabelColorPicker from "../Label/LabelColorPicker.svelte"
    import { labelStore, Label } from "../../stores/labels"

    export let selectedLabel: string
    export let editMode: boolean

    const dispatch = createEventDispatcher()

    let labels: Label[]
    let labelNames: string[]
    let badEdits: boolean[]
    const unsubscribe = labelStore.subscribe(() => {
        labels = labelStore.getLabels()
        labelNames = labels.map((l: Label) => l.name)
        badEdits = labels.map(() => false)
    })
    onDestroy(unsubscribe)

    /// trash can handling (edit mode)
    async function onRemoveClick(labelName: string) {
        if (labelName == selectedLabel) {
            // label currently open?
            dispatch("closeLabel")
        }
        labelStore.removeLabel(labelName)
    }

    /// label color picker handling (when in edit mode)
    let openPickerIndex: number = undefined
    let editColor: any

    function onTagApply(idx: number) {
        if (openPickerIndex == undefined || idx != openPickerIndex) {
            return
        }

        // close picker for this label
        openPickerIndex = undefined

        // update tag color immediately
        labels[idx].color = editColor.hex

        labelStore.setLabelColor(labels[idx].name, editColor.hex)
    }

    function onTagClick(event: MouseEvent, idx: number) {
        if (!editMode) {
            return
        }
        event.stopPropagation()

        if (openPickerIndex == idx) {
            // clicked again, close
            openPickerIndex = undefined
        } else {
            openPickerIndex = idx
        }
    }

    $: if (!editMode) {
        // close color picker when leaving edit mode
        openPickerIndex = undefined
    }

    /// add label
    let newLabelName: string
    function addLabel() {
        if (!validLabelName) {
            return
        }

        const color = "#ffffff"
        labelStore.addLabel(newLabelName, color)

        // clear input
        newLabelName = null
    }

    $: validLabelName = newLabelName && !labelNames.includes(newLabelName)
</script>

<div class="flex-auto py-2">
    <h3 class="font-bold px-2 pb-2 text-center md:text-left">Labels</h3>
    <ul>
        {#if labels == undefined}
            <li class="flex w-full justify-center m-2"><Spinner /></li>
        {:else}
            {#each labels as label, i}
                <li
                    class="flex flex-row py-2 px-2 md:p-0 md:pl-2 hover:bg-gray-700 cursor-pointer items-center md:text-sm"
                    class:bg-gray-600={!editMode && selectedLabel == labelNames[i]}
                    on:click={() => !editMode && dispatch("labelClick", label.name)}
                    >
                    {#if i == openPickerIndex}
                        <div class="block z-10">
                            <LabelColorPicker
                                bind:color={editColor}
                                startColor={label.color}
                                on:applyClick={() => onTagApply(i)}
                                />
                        </div>
                    {/if}
                    <div class="h-6 w-6" style="color: {label.color}" on:click={(e) => onTagClick(e, i)}>
                        <OutlineTag />
                    </div>

                    {#if !editMode}
                        <h4 class="flex-auto pl-1">{label.name}</h4>
                    {:else}
                        <input
                            class="flex-auto bg-transparent border-b border-white border-t-0 border-r-0 border-l-0 min-w-0 p-0 focus:ring-0 mx-1"
                            class:border-red-500={badEdits[i]}
                            class:focus:border-red-500={badEdits[i]}
                            type="text"
                            bind:value={label.name}
                            on:input={() => {
                                         const nameChanged = label.name != labelNames[i]
                                         const alreadyExists = nameChanged && labelNames.includes(label.name)
                                         const isBadEdit = alreadyExists || label.name == ""
                                         badEdits[i] = isBadEdit
                                         }}
                        title={badEdits[i]
                                  ? label.name == ""
                                  ? "Label name cannot be empty"
                                  : "A label with the same name already exists"
                                  : ""}
                        placeholder="Unnamed Label"
                        />
                    {/if}

                    <!-- article count bubble -->
                    {#if label.articleIDs.size}
                        <span class="flex-none bg-gray-800 text-gray-200 text-xs rounded-full py-1 md:py-0 px-2 mx-1"
                        >{label.articleIDs.size}</span
                                                >
                    {/if}

                    {#if editMode}
                        <!-- trash icon (in edit mode) -->
                        <button
                            class="flex-none mr-1 text-red-900 hover:text-red-600 cursor-pointer"
                            title="Delete Label"
                            on:click|stopPropagation={(_e) => onRemoveClick(label.name)}
                            >
                            <SolidTrash />
                        </button>
                    {/if}
                </li>
            {:else}
                <li class="italic pl-2">No labels</li>
            {/each}
        {/if}
    </ul>
</div>

{#if editMode}
    <!-- add label box -->
    <div class="flex-none p-2 pr-0 text-white border-t border-green-500 bg-gray-800">
        <p class="text-sm">Add Label</p>
        <div class="flex flex-row items-center">
            <input
                bind:value={newLabelName}
                type="text"
                placeholder="New label name"
                class="flex-grow text-sm text-black min-w-0"
                on:keydown={(e) => {
                    e.key == "Enter" && addLabel()
                }}
            />
            <button
                class="flex-none p-2 cursor-pointer"
                class:cursor-not-allowed={!validLabelName}
                class:text-gray-500={!validLabelName}
                class:hover:text-green-500={validLabelName}
                title="Add label"
                on:click={() => addLabel()}
            >
                <SolidPlusCircle />
            </button>
        </div>
    </div>
{/if}
