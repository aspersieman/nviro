/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./cmd/static/js/**/*.{html,js}",
    "./cmd/static/templates/**/*.{html,js}",
    "./node_modules/flowbite/**/*.js",
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('flowbite/plugin')
  ]
}
