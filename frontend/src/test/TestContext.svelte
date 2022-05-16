<script lang="ts">
    // based on: https://svelte-recipes.netlify.app/testing/
    import { setContext, SvelteComponent } from "svelte"

    export let Component: SvelteComponent
    export let context_key: string
    export let context_value: any

    setContext(context_key, context_value)
</script>

<svelte:component this={Component} {...$$restProps} />
