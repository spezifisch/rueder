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
        "no-unused-vars": "off", // disabled because it clashes with ts enums
        "@typescript-eslint/no-unused-vars": [
            // https://stackoverflow.com/a/69006568
            "error",
            {
                varsIgnorePattern: "^_",
                argsIgnorePattern: "^_",
            },
        ],
        indent: "off",
        "@typescript-eslint/no-redeclare": "error",
        "no-shadow": "off",
        "@typescript-eslint/no-shadow": "error",
        quotes: "off",
        "@typescript-eslint/quotes": "error",
    },
}
