/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        './templates/**/*.{html,go}', // Path to your Go HTML template files
        './*.go',                     // Path to your Go files, if you use template in Go files
    ],
    theme: {
        extend: {},
    },
    plugins: [],
}
