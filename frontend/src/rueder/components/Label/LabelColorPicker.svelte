<script lang="ts">
    import { createEventDispatcher } from "svelte"

    import ChromePicker from "@spezifisch/svelte-color/Chrome.svelte"

    import OutlineCheck from "../ui/heroicons/outline-check.svelte"

    const dispatch = createEventDispatcher()

    export let color: any = undefined
    export let startColor: string
</script>

<div class="relative inline-block" on:click|stopPropagation={() => {}}>
    <div class="absolute z-10 top-5 left-5 flex flex-col bg-gray-100 text-black">
        <h3 class="text-center">Choose Label Color</h3>
        <ChromePicker bind:color {startColor} disableAlpha={true} />
        <button
            class="text-center flex flex-row justify-center items-center hover:bg-gray-300"
            on:click|stopPropagation={() => dispatch("applyClick")}
        >
            <OutlineCheck />
            <span>Apply</span>
        </button>
    </div>
</div>
