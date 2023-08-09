import { useEffect } from 'react'
import { useRouter } from 'next/router'
import useAuth from '@context/authContext'
import { Dashboard } from '@layouts'
import {
  authRoutes,
  dashboardRoutes,
  dashboardRoutesParams,
} from '@constants/routes'
import { Admin, Alarms } from '@organisms'

export default () => {
  const {
    query: { dashboardRoute },
    replace,
    isReady,
  } = useRouter()
  const { logedIn } = useAuth()

  useEffect(() => {
    if (!logedIn) {
      replace(authRoutes.LOGIN)
    }
  }, [replace, logedIn])

  if (!logedIn || !isReady) {
    return
  }

  switch (dashboardRoute) {
    case dashboardRoutesParams.ALARMS:
      return (
        <Dashboard>
          <Alarms />
        </Dashboard>
      )
    case dashboardRoutesParams.ADMIN:
      return (
        <Dashboard>
          <Admin />
        </Dashboard>
      )
    default: {
      replace(dashboardRoutes.DASHBOARD)
      return
    }
  }
}

export async function getStaticProps() {
  return { props: {} }
}

export async function getStaticPaths() {
  return {
    paths: [
      { params: { dashboardRoute: dashboardRoutesParams.ALARMS }},
      { params: { dashboardRoute: dashboardRoutesParams.ADMIN }},
    ],
    fallback: false,
  }
}