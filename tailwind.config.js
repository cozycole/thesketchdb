/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./ui/**/*.html", "./ui/**/*.js"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/line-clamp'),
    require('tailwind-scrollbar-hide'),
  ],
}
