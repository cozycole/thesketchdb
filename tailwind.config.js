/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./ui/**/*.html"],
  safelist : [
    'absolute',
    'relative',
    'text-xs',
    'rounded', 
    'whitespace-nowrap',
    '-translate-x-1/2 ',
    'left-1/2',
    '-top-6',
    'bg-black',
    'text-white',
    'text-xs',
    'px-2',
    'py-1',
    'w-fit',
    'opacity-0',
    'opacity-100',
    'transition-opacity',
    'duration-300',
    'transform'
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/line-clamp'),
  ],
}
