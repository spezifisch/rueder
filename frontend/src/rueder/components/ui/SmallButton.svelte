<script lang="ts">
    import { createEventDispatcher } from "svelte"

    const dispatch = createEventDispatcher()

    export let addClass: string = ""
</script>

<!-- based on: https://tailwindcomponents.com/component/button-component-default -->

<button
    on:click={(e) => dispatch("click", e)}
    type="button"
    class="border border-gray-600 bg-gray-600 text-white rounded-md px-2 py-1 transition duration-500 ease select-none hover:bg-gray-800 focus:outline-none focus:shadow-outline {addClass}"
>
    <slot />
</button>
