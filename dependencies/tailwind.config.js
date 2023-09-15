/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../**/*.{html,js}"],
  theme: {
    extend: {
      keyframes: {
        appear: {
          '0%': { opacity: '0' },
          '100%': {opacity:'1'}
        }
      },
      animation:{
        appear: 'appear 0.5s ease-in-out'
      },
      colors:{
        'main-light': '#fdf3e7',
        'main-dark-green': '#82d7cb',
        'main-light-green': '#beebe5',
        'main-brown': '#7a6a5c'
      }
    },
  },
  plugins: [],
}

