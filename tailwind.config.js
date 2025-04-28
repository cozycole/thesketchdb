/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./ui/**/*.html", "./ui/**/*.js"],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Poppins', 'sans-serif'],
      },
    },
  },
  plugins: [
    require('@tailwindcss/line-clamp'),
    require('tailwind-scrollbar-hide'),
  ],
}
