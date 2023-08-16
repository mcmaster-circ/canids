import { useEffect } from 'react'
import { useRouter } from 'next/router'
import useAuth from '@context/authContext'
import { authRoutes, dashboardRoutes } from '@constants/routes'

export default () => {
  const router = useRouter()
  const { logedIn, isActive } = useAuth()

  useEffect(() => {
    if (logedIn) {
      router.replace(dashboardRoutes.DASHBOARD)
    } else {
      isActive().then((data) => {
        if (data) {
          console.log('Data = true')
        } else {
          console.log('Data = false')
        }
      })
      router.replace(authRoutes.LOGIN)
    }
  }, [logedIn, router, isActive])

  return null
}
