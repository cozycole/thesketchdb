/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./ui/**/*.gohtml", "./ui/**/*.js"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Rubik", "sans-serif"],
      },
    },
    screens: {
      xs: "480px",
      sm: "640px",
      md: "768px",
      lg: "1080px",
      xl: "1280px",
    },
  },
  plugins: [
    require("@tailwindcss/line-clamp"),
    require("tailwind-scrollbar-hide"),
  ],
};
