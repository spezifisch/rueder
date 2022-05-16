module.exports = {
    env: {
        browser: true,
        es2021: true,
    },
    extends: ["eslint:recommended", "plugin:@typescript-eslint/recommended"],
    parser: "@typescript-eslint/parser",
    parserOptions: {
        ecmaVersion: "latest",
        sourceType: "module",
    },
    plugins: ["@typescript-eslint"],
    rules: {
        "no-unused-vars": [
            // https://stackoverflow.com/a/69006568
            "error",
            {
                varsIgnorePattern: "^_",
                argsIgnorePattern: "^_",
            },
        ],
        "@typescript-eslint/no-unused-vars": [
            // as above
            "error",
            {
                varsIgnorePattern: "^_",
                argsIgnorePattern: "^_",
            },
        ],
    },
}
