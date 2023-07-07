import { createTheme } from '@mui/material/styles'

export default createTheme({
  typography: {
    fontFamily: [
      'Open Sans',
      'Roboto',
      'Arial',
      '-apple-system',
      'sans-serif',
    ].join(','),
  },
  palette: {
    primary: {
      main: '#be3e24',
    },
    secondary: {
      main: '#1976d2',
    },
  },
})
