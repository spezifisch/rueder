module.exports = {
    transform: {
        "^.+\\.svelte$": [
            "svelte-jester",
            {
                preprocess: "svelte.config.mjs", // it doesn't find our .mjs extension by default
            },
        ],
        // process js files with ts, too, because jest doesn't understand ESM
        "^.+\\.[tj]s$": "ts-jest",
    },
    // process 3rd party svelte and ESM files with ts
    transformIgnorePatterns: ["/node_modules/(?!(svelte|@.*/svelte|url-join)).+\\.js$"],
    moduleFileExtensions: ["js", "ts", "svelte"],
    // handle css imports
    moduleNameMapper: {
        "\\.(post)?css$": "identity-obj-proxy",
    },
    // keep localstorage by default
    resetMocks: false,
    setupFiles: ["jest-localstorage-mock"],
    //
    testEnvironment: "jsdom",
    slowTestThreshold: 30,
}
