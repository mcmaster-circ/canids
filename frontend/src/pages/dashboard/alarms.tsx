import { useEffect } from 'react'
import { useRouter } from 'next/router'
import useAuth from '@context/authContext'
import { Dashboard } from '@layouts'
import { authRoutes } from '@constants/routes'
import { Alarms } from '@organisms'

export default () => {
  const {
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

  return (
    <Dashboard>
      <Alarms />
    </Dashboard>
  )
}
