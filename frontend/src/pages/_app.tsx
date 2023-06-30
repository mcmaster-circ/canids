import type { AppProps } from 'next/app'
import Head from 'next/head'
import { AuthProvider } from '@context/authContext'
import { NotificationProvider } from '@context/notificationContext'
import { ThemeProvider } from '@mui/material/styles'
import theme from '@styles/theme'
import '@styles/global.css'

export default ({ Component, pageProps }: AppProps) => {
  return (
    <>
      <Head>
        <title>CANIDS</title>
      </Head>
      <ThemeProvider theme={theme}>
        <NotificationProvider>
          <AuthProvider>
            <Component {...pageProps} />
          </AuthProvider>
        </NotificationProvider>
      </ThemeProvider>
    </>
  )
}
