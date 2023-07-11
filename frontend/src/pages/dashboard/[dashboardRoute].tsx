import { useEffect } from 'react'
import { useRouter } from 'next/router'
import useAuth from '@context/authContext'
import { Dashboard } from '@layouts'
import {
  authRoutes,
  dashboardRoutes,
  dashboardRoutesParams,
} from '@constants/routes'
import { Alarms } from '@organisms'

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
      return <Dashboard>ADMIN</Dashboard>
    default: {
      replace(dashboardRoutes.DASHBOARD)
      return
    }
  }
}
