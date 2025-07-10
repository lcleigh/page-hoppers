/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/**/*.{js,ts,jsx,tsx}',
    './pages/**/*.{js,ts,jsx,tsx}',
    './components/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        bubblegum: '#F7A9A8',
        sky: '#A7D8F1',
        lemon: '#FFF5A5',
        powder: '#FFFDF9',
        lavender: '#E9DFF6',
        charcoal: '#333333',
        coolgray: '#6B7280',
        leaf: '#A9E5BB',
        coral: '#FF6B6B',
      },
    },
  },
}; 