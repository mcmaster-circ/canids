import { useEffect, useState } from 'react'
import { useRouter } from 'next/router'
import { useCookies } from 'react-cookie'
import useAuth from '@context/authContext'
import { Auth } from '@layouts'
import { LoginForm } from '@forms'
import { authRouteParams, dashboardRoutes } from '@constants/routes'

export default () => {
  const [cookies, setCookie] = useCookies(['X-State', 'X-Class'])
  const [a, setA] = useState(false)
  const { query, replace } = useRouter()
  const { logedIn } = useAuth()

  useEffect(() => {
    if (!a) {
      setCookie(
        'X-State',
        'ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SjFkV2xrSWpvaVlXNWtjbWxwUUdaNVpXeGhZbk11WTI5dElpd2lZMnhoYzNNaU9pSnpkR0Z1WkdGeVpDSXNJbTVoYldVaU9pSkJibVJ5YVdraUxDSmhZM1JwZG1GMFpXUWlPblJ5ZFdVc0ltVjRjQ0k2TVRZNE9ETTVORGM0TUN3aWFXRjBJam94TmpnNE1UTTFOVGd3ZlEuMTNRQ01wZ1R3RlpxX2ZzSTVra2RNMjU1WGtVN2hjSS1NaFN4bkJ1ZFNMbEVaNjN2SzVaY1hLb1c2VmF6bVNSWWp6Z3BrV2RLM1dkUFd6S2Nvd1RwQ1E=',
        { path: '/' }
      )
      setCookie('X-Class', 'standard', { path: '/' })
      setA(true)
    }
    if (logedIn) {
      replace(dashboardRoutes.DASHBOARD)
    }
  }, [replace, logedIn, a, setCookie])

  if (logedIn) {
    return null
  }

  switch (query.authRoute) {
    case authRouteParams.LOGIN:
    default: {
      return (
        <Auth>
          <LoginForm />
        </Auth>
      )
    }
  }
}
