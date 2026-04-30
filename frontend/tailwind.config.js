/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,jsx,ts,tsx}'],
  theme: {
    extend: {
      colors: {
        relay: {
          50: '#f0f4ff',
          100: '#dde8ff',
          200: '#c3d3ff',
          300: '#9ab4ff',
          400: '#7090ff',
          500: '#4f6ef7',
          600: '#3a56e8',
          700: '#2d45cc',
          800: '#2238a8',
          900: '#1a2a7a',
        },
      },
    },
  },
  plugins: [],
}
