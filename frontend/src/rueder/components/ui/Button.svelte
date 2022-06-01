<script lang="ts">
    import { createEventDispatcher } from "svelte"

    const dispatch = createEventDispatcher()

    export let mode: string = "light"
    export let addClass: string = ""
</script>

<!-- based on: https://tailwindcomponents.com/component/button-component-default -->

{#if mode == "primary"}
    <button
        on:click={(e) => dispatch("click", e)}
        type="button"
        class="border border-indigo-500 bg-indigo-500 text-white rounded-md px-4 py-2 m-2 transition duration-500 ease select-none hover:bg-indigo-600 focus:outline-none focus:shadow-outline {addClass}"
    >
        <slot />
    </button>
{:else if mode == "success"}
    <button
        on:click={(e) => dispatch("click", e)}
        type="button"
        class="border border-green-500 bg-green-500 text-white rounded-md px-4 py-2 m-2 transition duration-500 ease select-none hover:bg-green-600 focus:outline-none focus:shadow-outline {addClass}"
    >
        <slot />
    </button>
{:else if mode == "error"}
    <button
        on:click={(e) => dispatch("click", e)}
        type="button"
        class="border border-red-500 bg-red-500 text-white rounded-md px-4 py-2 m-2 transition duration-500 ease select-none hover:bg-red-600 focus:outline-none focus:shadow-outline {addClass}"
    >
        <slot />
    </button>
{:else if mode == "error-dark"}
    <button
        on:click={(e) => dispatch("click", e)}
        type="button"
        class="border border-red-700 bg-red-900 text-white rounded-md px-4 py-2 m-2 transition duration-500 ease select-none hover:bg-red-600 focus:outline-none focus:shadow-outline {addClass}"
    >
        <slot />
    </button>
{:else if mode == "warning"}
    <button
        on:click={(e) => dispatch("click", e)}
        type="button"
        class="border border-yellow-500 bg-yellow-500 text-white rounded-md px-4 py-2 m-2 transition duration-500 ease select-none hover:bg-yellow-600 focus:outline-none focus:shadow-outline {addClass}"
    >
        <slot />
    </button>
{:else if mode == "info"}
    <button
        on:click={(e) => dispatch("click", e)}
        type="button"
        class="border border-teal-500 bg-teal-500 text-white rounded-md px-4 py-2 m-2 transition duration-500 ease select-none hover:bg-teal-600 focus:outline-none focus:shadow-outline {addClass}"
    >
        <slot />
    </button>
{:else if mode == "dark"}
    <button
        on:click={(e) => dispatch("click", e)}
        type="button"
        class="border border-gray-600 bg-gray-600 text-white rounded-md px-4 py-2 m-2 transition duration-500 ease select-none hover:bg-gray-800 focus:outline-none focus:shadow-outline {addClass}"
    >
        <slot />
    </button>
{:else if mode == "light"}
    <button
        on:click={(e) => dispatch("click", e)}
        type="button"
        class="border border-gray-200 bg-gray-200 text-gray-700 rounded-md px-4 py-2 m-2 transition duration-500 ease select-none hover:bg-gray-300 focus:outline-none focus:shadow-outline {addClass}"
    >
        <slot />
    </button>
{/if}
