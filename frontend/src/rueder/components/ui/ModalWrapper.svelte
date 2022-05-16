<script lang="ts">
    export let centered = false
</script>

<!-- based on: https://tailwindui.com/components/application-ui/overlays/modals -->
<div class="fixed z-10 inset-0" aria-labelledby="modal-title" role="dialog" aria-modal="true">
    <div class="flex items-end justify-center text-center sm:block h-screen">
        <div class="fixed inset-0 bg-gray-800 bg-opacity-25 transition-opacity" aria-hidden="true" />

        {#if centered}
            <!-- This element is to trick the browser into centering the modal contents. -->
            <span class="hidden sm:inline-block sm:align-middle h-screen" aria-hidden="true">&#8203;</span>
        {/if}

        <div class="inline-block transform align-middle sm:max-w-lg sm:w-full h-screen">
            <div class="flex flex-col h-screen justify-center">
                <div class="flex-none spacer" />
                <div class="flex-auto flex flex-col justify-center dialog-box text-left">
                    <slot />
                </div>
                <div class="flex-none h-8 spacer" />
            </div>
        </div>
    </div>
</div>

<style>
    .spacer {
        height: 8vh;
    }

    .dialog-box {
        max-height: 82vh;
    }
</style>
