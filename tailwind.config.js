/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./ui/html/**/*.{html,js}'],
  theme: {
    extend: {},
  },
  plugins: [require('@tailwindcss/typography'), require('daisyui')],
  daisyui: {
    themes: ['retro'],
  },
}
