<script lang="ts">
    // based on: https://svelte-recipes.netlify.app/testing/
    import { setContext, SvelteComponent } from "svelte"

    export let Component: typeof SvelteComponent
    export let context_key: string | symbol
    export let context_value: any

    setContext(context_key, context_value)
</script>

<svelte:component this={Component} {...$$restProps} />
