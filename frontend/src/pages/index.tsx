import { useEffect } from 'react'
import { useRouter } from 'next/router'
import useAuth from '@context/authContext'
import { authRoutes, dashboardRoutes } from '@constants/routes'

export default () => {
  const router = useRouter()
  const { logedIn } = useAuth()

  useEffect(() => {
    if (logedIn) {
      router.replace(dashboardRoutes.DASHBOARD)
    } else router.replace(authRoutes.LOGIN)
  }, [logedIn, router])

  return null
}
