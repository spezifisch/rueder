<script lang="ts">
    import { fade } from "svelte/transition"

    import { sessionStore } from "./stores/session"

    export let page: string
    export let baseURL: string

    let disabled = false
    let username: string
    let password: string

    async function login() {
        disabled = true
        doShowError = false

        const url = `${baseURL}login`
        const settings = {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                username: username,
                password: password,
            }),
        }
        const res = await fetch(url, settings)
        if (res.status == 200) {
            const jwtToken: string = await res.text()
            $sessionStore.loggedIn = true
            $sessionStore.jwtToken = jwtToken
        } else {
            $sessionStore.loggedIn = false
            $sessionStore.jwtToken = null

            let reason = `Login failed: ${res.status}`
            try {
                const data = await res.json()
                if (data.error) {
                    reason = data.error
                }
            } catch (e) {
                // pass
            }
            showError(reason)
        }

        disabled = false
    }

    let doShowError = false
    let errorMsg: string
    function showError(message: string) {
        doShowError = true
        errorMsg = message
    }
</script>

<!-- based on: https://tailwindcomponents.com/component/tela-de-login -->
<div class="h-screen font-sans login">
    <div class="container mx-auto h-full flex flex-1 justify-center items-center">
        <div class="w-full max-w-lg">
            <div class="flex flex-col items-center leading-loose">
                {#if doShowError}
                    <div
                        class="flex flex-row items-center space-x-6 max-w-sm m-4 py-2 px-10 bg-red-700 bg-opacity-50 text-red-100 rounded shadow-xl"
                        in:fade
                    >
                        <p class="flex-auto">{errorMsg}</p>
                        <div
                            class="flex-none hover:text-red-300"
                            title="Close message"
                            on:click={() => (doShowError = false)}
                        >
                            <!-- heroicons: solid/x -->
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-5 w-5"
                                viewBox="0 0 20 20"
                                fill="currentColor"
                            >
                                <path
                                    fill-rule="evenodd"
                                    d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                                    clip-rule="evenodd"
                                />
                            </svg>
                        </div>
                    </div>
                {/if}
                <form
                    class="max-w-sm m-4 p-10 bg-black bg-opacity-50 rounded shadow-xl"
                    on:submit|preventDefault={login}
                >
                    <p class="text-white text-center text-lg font-bold">Sign in to {page}</p>
                    <div class="">
                        <label class="block text-sm text-white" for="username">Username</label>
                        <input
                            class="w-full px-5 py-1 text-gray-700 bg-gray-300 rounded focus:outline-none focus:bg-white"
                            type="text"
                            id="username"
                            autocomplete="username"
                            placeholder="Username"
                            aria-label="username"
                            {disabled}
                            bind:value={username}
                            required
                        />
                    </div>
                    <div class="mt-2">
                        <label class="block text-sm text-white" for="password">Password</label>
                        <input
                            class="w-full px-5 py-1 text-gray-700 bg-gray-300 rounded focus:outline-none focus:bg-white"
                            type="password"
                            id="password"
                            autocomplete="current-password"
                            placeholder="Password"
                            aria-label="password"
                            {disabled}
                            bind:value={password}
                            required
                        />
                    </div>

                    <div class="mt-4 items-center flex justify-between">
                        <button
                            class="px-4 py-1 text-white font-light tracking-wider bg-green-700 hover:bg-green-800 rounded"
                            type="submit"
                            {disabled}>Sign in</button
                        >
                        <!--a
                            class="inline-block right-0 align-baseline font-bold text-sm text-500 text-white hover:text-red-400"
                            href="#">Forgot password?</a
                        -->
                    </div>
                    <!--div class="text-center">
                        <a class="inline-block right-0 align-baseline font-light text-sm text-500 hover:text-red-400">
                            Create an account
                        </a>
                    </div-->
                </form>
            </div>
        </div>
    </div>
</div>

<style lang="postcss">
    .login {
        /* https://unsplash.com/photos/l7dP0O8Dj60 */
        background: url("../assets/images/wherda-arsianto-l7dP0O8Dj60-unsplash.jpg");
        background-color: #bbb;
        @apply bg-no-repeat;
        @apply bg-cover;
    }
</style>
