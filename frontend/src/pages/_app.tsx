import type { AppProps } from 'next/app'
import Head from 'next/head'
import { AuthProvider } from '@context/authContext'
import { NotificationProvider } from '@context/notificationContext'
import { ThemeProvider } from '@mui/material/styles'
import { LocalizationProvider } from '@mui/x-date-pickers'
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns'
import theme from '@styles/theme'
import '@styles/global.css'

export default ({ Component, pageProps }: AppProps) => {
  return (
    <>
      <Head>
        <title>McMaster CanIDS</title>
      </Head>
      <ThemeProvider theme={theme}>
        <LocalizationProvider dateAdapter={AdapterDateFns}>
          <NotificationProvider>
            <AuthProvider>
              <Component {...pageProps} />
            </AuthProvider>
          </NotificationProvider>
        </LocalizationProvider>
      </ThemeProvider>
    </>
  )
}
