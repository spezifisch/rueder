<script lang="ts">
    /* based on MultiSelection.svelte from svelte-select,
    license: https://github.com/rob-balfre/svelte-select/blob/master/LICENSE */

    import { createEventDispatcher } from "svelte"

    import SolidTag from "../ui/heroicons/solid-tag.svelte"
    import LabelColorPicker from "./LabelColorPicker.svelte"

    import type { CustomEvent } from "./types"

    const dispatch = createEventDispatcher()

    export let selectedValue = []
    export let activeSelectedValue = undefined
    export let isDisabled = false
    export let multiFullItemClearable = false
    export let getSelectionLabel = undefined

    function handleClear(i: number, event: MouseEvent) {
        event.stopPropagation()
        dispatch("multiItemClear", { i })
    }

    function colorUpdateEvent(labelName: string, hexColor: string) {
        const data: CustomEvent = {
            updateType: "color",
            labelName,
            hexColor,
        }
        dispatch("customEvent", data)
    }

    function getLabelName(option: any): string {
        return option.label
    }

    function getLabelColor(option: any): string {
        return option.color
    }

    // shows/hides the color picker when a label tag is clicked
    let openPickerIndex: number = undefined
    function onLabelClick(labelIdx: number): void {
        if (openPickerIndex === labelIdx) {
            // close picker
            openPickerIndex = undefined
        } else {
            // open picker for this label
            openPickerIndex = labelIdx
        }
    }

    let color: any
    function onApplyClick(labelIdx: number, option: any): void {
        const newColor = color.hex
        selectedValue[labelIdx].color = newColor
        colorUpdateEvent(getLabelName(option), newColor)

        // close picker for this label
        openPickerIndex = undefined
    }
</script>

{#each selectedValue as value, i}
    <button
        class="multiSelectItem {activeSelectedValue === i ? 'active' : ''} {isDisabled
            ? 'disabled'
            : ''}"
        on:click={(event) => (multiFullItemClearable ? handleClear(i, event) : {})}
    >
        {#if i == openPickerIndex}
            <LabelColorPicker
                bind:color
                startColor={getLabelColor(value)}
                on:applyClick={() => onApplyClick(i, value)}
            />
        {/if}
        <div class="multiSelectItem_label">
            <div class="flex flex-row space-x-1">
                <button
                    class="flex-none"
                    style="color: {getLabelColor(value)}"
                    on:click|stopPropagation={() => onLabelClick(i)}
                >
                    <SolidTag size={5} addClass="inline" />
                </button>
                <div>
                    {@html getSelectionLabel(value)}
                </div>
            </div>
        </div>
        {#if !isDisabled && !multiFullItemClearable}
            <button class="multiSelectItem_clear" on:click={(event) => handleClear(i, event)}>
                <svg
                    width="100%"
                    height="100%"
                    viewBox="-2 -2 50 50"
                    focusable="false"
                    role="presentation"
                >
                    <path
                        d="M34.923,37.251L24,26.328L13.077,37.251L9.436,33.61l10.923-10.923L9.436,11.765l3.641-3.641L24,19.047L34.923,8.124 l3.641,3.641L27.641,22.688L38.564,33.61L34.923,37.251z"
                    />
                </svg>
            </button>
        {/if}
    </button>
{/each}

<style>
    .multiSelectItem {
        background: var(--multiItemBG, #ebedef);
        margin: var(--multiItemMargin, 5px 5px 0 0);
        border-radius: var(--multiItemBorderRadius, 16px);
        height: var(--multiItemHeight, 32px);
        line-height: var(--multiItemHeight, 32px);
        display: flex;
        cursor: default;
        padding: var(--multiItemPadding, 0 10px 0 15px);
        max-width: 100%;
    }

    .multiSelectItem_label {
        margin: var(--multiLabelMargin, 0 5px 0 0);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .multiSelectItem:hover,
    .multiSelectItem.active {
        background-color: var(--multiItemActiveBG, #006fff);
        color: var(--multiItemActiveColor, #fff);
    }

    .multiSelectItem.disabled:hover {
        background: var(--multiItemDisabledHoverBg, #ebedef);
        color: var(--multiItemDisabledHoverColor, #c1c6cc);
    }

    .multiSelectItem_clear {
        border-radius: var(--multiClearRadius, 50%);
        background: var(--multiClearBG, #52616f);
        min-width: var(--multiClearWidth, 16px);
        max-width: var(--multiClearWidth, 16px);
        height: var(--multiClearHeight, 16px);
        position: relative;
        top: var(--multiClearTop, 8px);
        text-align: var(--multiClearTextAlign, center);
        padding: var(--multiClearPadding, 1px);
    }

    .multiSelectItem_clear:hover,
    .active .multiSelectItem_clear {
        background: var(--multiClearHoverBG, #fff);
    }

    .multiSelectItem_clear:hover svg,
    .active .multiSelectItem_clear svg {
        fill: var(--multiClearHoverFill, #006fff);
    }

    .multiSelectItem_clear svg {
        fill: var(--multiClearFill, #ebedef);
        vertical-align: top;
    }
</style>
