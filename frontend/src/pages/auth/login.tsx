import { useEffect } from 'react'
import { useRouter } from 'next/router'
import useAuth from '@context/authContext'
import { Auth } from '@layouts'
import { LoginForm } from '@forms'
import { dashboardRoutes } from '@constants/routes'

export default () => {
  const { replace } = useRouter()
  const { logedIn } = useAuth()

  useEffect(() => {
    if (logedIn) {
      replace(dashboardRoutes.DASHBOARD)
    }
  }, [replace, logedIn])

  if (logedIn) {
    return null
  }

  return (
    <Auth title="Authenticate to access CanIDS">
      <LoginForm />
    </Auth>
  )
}
