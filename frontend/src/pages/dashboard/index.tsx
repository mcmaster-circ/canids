import { useEffect } from 'react'
import { useRouter } from 'next/router'
import useAuth from '@context/authContext'
import { authRoutes } from '@constants/routes'
import { Dashboard } from '@layouts'
import { MainDashboard } from '@organisms'

export default () => {
  const router = useRouter()
  const { logedIn } = useAuth()

  useEffect(() => {
    if (!logedIn) {
      router.replace(authRoutes.LOGIN)
    }
  }, [logedIn, router])

  if (!logedIn) {
    return null
  }

  return (
    <Dashboard>
      <MainDashboard />
    </Dashboard>
  )
}
