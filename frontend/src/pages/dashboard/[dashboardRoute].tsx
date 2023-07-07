import { useEffect } from 'react'
import { useRouter } from 'next/router'
import useAuth from '@context/authContext'
import { Dashboard } from '@layouts'
import {
  authRoutes,
  dashboardRoutes,
  dashboardRoutesParams,
} from '@constants/routes'

export default () => {
  const { query, replace } = useRouter()
  const { logedIn } = useAuth()

  useEffect(() => {
    if (!logedIn) {
      replace(authRoutes.LOGIN)
    }
  }, [replace, logedIn])

  if (!logedIn) {
    return null
  }

  switch (query.dashboardRoute) {
    case dashboardRoutesParams.ALARMS:
      return <Dashboard>ALARMS</Dashboard>
    case dashboardRoutesParams.ADMIN:
      return <Dashboard>ADMIN</Dashboard>
    default: {
      replace(dashboardRoutes.DASHBOARD)
      return null
    }
  }
}
