<script lang="ts">
    import { createEventDispatcher } from "svelte"

    const dispatch = createEventDispatcher()
</script>

<!-- based on: https://tailwindui.com/components/application-ui/overlays/modals -->
<div class="fixed z-10 inset-0" aria-labelledby="modal-title" role="dialog" aria-modal="true">
    <div class="flex items-end justify-center text-center sm:block h-screen">
        <div class="inline-block align-middle w-full h-screen bg-black bg-opacity-80">
            <div class="flex flex-col w-full h-screen">
                <div class="flex-none">
                    <slot name="header" />
                </div>
                <slot name="scrolling" />
                <div class="flex-none bg-gray-900">
                    <slot name="buttons">
                        <button
                            on:click={() => dispatch("close")}
                            type="button"
                            class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-500 shadow-sm px-4 py-2 bg-gray-700 text-base font-medium text-white hover:bg-gray-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
                        >
                            Close
                        </button>
                    </slot>
                </div>
            </div>
        </div>
    </div>
</div>
