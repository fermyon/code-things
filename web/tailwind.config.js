const defaultColors = require('tailwindcss/colors');

/*
 * NOTE: shades were generated using https://www.tailwindshades.com
 */
const brandColors = {
  // core brand colors
  // this will effectively replace the default color pallette
  seagreen: {
    DEFAULT: "#34E8BD", // styleguide / seagreen
    50: "#D9FBF3",
    100: "#C7F9ED",
    200: "#A2F4E1",
    300: "#7DF0D5",
    400: "#59ECC9",
    500: "#34E8BD",
    600: "#17CDA1",
    700: "#119A7A",
    800: "#0C6852",
    900: "#06362A",
  },
  oxfordblue: {
    DEFAULT: '#0D203F', // styleguide / oxfordblue
    '50': '#95B5E9',
    '100': '#84A9E6',
    '200': '#6291DF',
    '300': '#407AD8',
    '400': '#2965C6',
    '500': '#2254A4',
    '600': '#1B4283',
    '700': '#143161',
    '800': '#0D203F',
    '900': '#030810'
  },
  // secondary accent colors
  rust: {
    DEFAULT: '#EF946C', // styleguide / rust
    '50': '#F9D7C8',
    '100': '#F7CAB6',
    '200': '#F3AF91',
    '300': '#EF946C',
    '400': '#E96F39',
    '500': '#D45117',
    '600': '#A13D12',
    '700': '#6F2A0C',
    '800': '#3C1707',
    '900': '#090401'
  },
  lavender: {
    DEFAULT: '#BEA7E5', // styleguide / lavender
    '50': '#F8F6FC',
    '100': '#EDE6F8',
    '200': '#D5C6EE',
    '300': '#BEA7E5',
    '400': '#9E7CD8',
    '500': '#7E50CB',
    '600': '#6234B0',
    '700': '#4A2784',
    '800': '#321A59',
    '900': '#1A0E2E'
  },
  colablue: {
    DEFAULT: '#0E8FDD', // styleguide / colablue
    '50': '#A9DBFA',
    '100': '#96D3F8',
    '200': '#6FC3F6',
    '300': '#49B3F3',
    '400': '#23A3F1',
    '500': '#0E8FDD',
    '600': '#0B6DA8',
    '700': '#074B73',
    '800': '#04293F',
    '900': '#01060A'
  },
  // complimentary tones
  darkspace: {
    DEFAULT: '#213762', // styleguide / darkspace
    '50': '#AABDE2',
    '100': '#9BB1DD',
    '200': '#7C99D3',
    '300': '#5E82C9',
    '400': '#406ABE',
    '500': '#36599F',
    '600': '#2B4881',
    '700': '#213762',
    '800': '#131F38',
    '900': '#05080E'
  },
  darkocean: {
    DEFAULT: '#0A455A', // styleguide / darkocean
    '50': '#A1DFF5',
    '100': '#8FD8F3',
    '200': '#6ACCEE',
    '300': '#46BFEA',
    '400': '#21B2E6',
    '500': '#1699C8',
    '600': '#127DA3',
    '700': '#0E617F',
    '800': '#0A455A',
    '900': '#041E28'
  },
  darkolive: {
    DEFAULT: '#1F7A8C', // styleguide / darkolive
    '50': '#C3EAF2',
    '100': '#B2E4EE',
    '200': '#90D8E7',
    '300': '#6FCDDF',
    '400': '#4EC1D8',
    '500': '#2EB4CF',
    '600': '#2697AD',
    '700': '#1F7A8C',
    '800': '#15525E',
    '900': '#0B2A30'
  },
  darkplum: {
    DEFAULT: '#525776', // styleguide / darkplum
    '50': '#CCCFDC',
    '100': '#C0C3D4',
    '200': '#A8ACC3',
    '300': '#9095B2',
    '400': '#787EA1',
    '500': '#63698E',
    '600': '#525776',
    '700': '#3B3F55',
    '800': '#242634',
    '900': '#0D0E13'
  },
  midgreen: {
    DEFAULT: '#1FBCA0', // styleguide / midgreen
    '50': '#C6F6ED',
    '100': '#B4F3E8',
    '200': '#91EDDD',
    '300': '#6EE7D2',
    '400': '#4BE1C7',
    '500': '#28DCBC',
    '600': '#1FBCA0',
    '700': '#178C77',
    '800': '#0F5C4E',
    '900': '#072C25'
  },
  midblue: {
    DEFAULT: '#345995', // styleguide / midblue
    '50': '#C0D0E9',
    '100': '#B1C4E4',
    '200': '#93AED9',
    '300': '#7597CF',
    '400': '#5680C4',
    '500': '#3F6BB3',
    '600': '#345995',
    '700': '#25406B',
    '800': '#172742',
    '900': '#080E18'
  },
  lightplum: {
    DEFAULT: '#D3C3D9', // styleguide / lightplum
    '50': '#EEE8F1',
    '100': '#E5DCE9',
    '200': '#D3C3D9',
    '300': '#BAA1C3',
    '400': '#A17EAD',
    '500': '#865E95',
    '600': '#674973',
    '700': '#483351',
    '800': '#2A1D2E',
    '900': '#0B070C'
  },
  lightgrey: {
    '50': '#FBFAFC',
    '100': '#EFEFF5',
    '200': '#D9DBE8', //lightgrey
    '300': '#B6BFD3',
    '400': '#93A7BE',
    '500': '#7094A9',
    '600': '#55818C',
    '700': '#406869',
    '800': '#2A4642',
    '900': '#15231F'
  },
  lightlavendar: {
    DEFAULT: '#ECE5EE', // styleguide / lightlavendar
    '50': '#F5F1F6',
    '100': '#ECE5EE',
    '200': '#D3C3D8',
    '300': '#BAA1C2',
    '400': '#A27FAB',
    '500': '#876093',
    '600': '#684A71',
    '700': '#49344F',
    '800': '#291D2D',
    '900': '#0A070B'
  },
  lightlemon: {
    DEFAULT: '#F9F7EE', // styleguide / lightlemon
    '50': '#FEFEFD',
    '100': '#F9F7EE',
    '200': '#EAE3C5',
    '300': '#DCD09B',
    '400': '#CDBC72',
    '500': '#BEA948',
    '600': '#998736',
    '700': '#6F6227',
    '800': '#463E19',
    '900': '#1C190A'
  },
};

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx,vue}"],
  theme: {
    colors: {
      // include all of our brand colors
      ...brandColors,
      // pass through default color aliases
      gray: brandColors.lightgrey, //TODO: either refactor the html to change colors or change the gray pallette here
      white: defaultColors.white,
      black: defaultColors.black,
      transparent: defaultColors.transparent,
      current: defaultColors.current,
      inherit: defaultColors.inherit,
      neutral: defaultColors.neutral,
    }
  },
  plugins: [
    require("@tailwindcss/forms"),
    require("prettier-plugin-tailwindcss"),
  ],
};
