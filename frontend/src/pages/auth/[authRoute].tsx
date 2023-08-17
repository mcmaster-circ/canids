import { useEffect } from 'react'
import { useRouter } from 'next/router'
import useAuth from '@context/authContext'
import { Auth } from '@layouts'
import {
  ForgotPasswordForm,
  LoginForm,
  RegisterForm,
  ResetPasswordForm,
  SetupForm,
} from '@forms'
import { authRouteParams, dashboardRoutes, authRoutes } from '@constants/routes'

export default () => {
  const { query, replace } = useRouter()
  const { logedIn, isActive } = useAuth()

  useEffect(() => {
    if (logedIn) {
      replace(dashboardRoutes.DASHBOARD)
    }
  }, [replace, logedIn])

  if (logedIn) {
    return null
  }

  switch (query.authRoute) {
    case authRouteParams.REGISTER:
      return (
        <Auth title="Please register a new account">
          <RegisterForm />
        </Auth>
      )
    case authRouteParams.FORGOT_PASSWORD:
      return (
        <Auth title="Please enter your email to reset the password">
          <ForgotPasswordForm />
        </Auth>
      )
    case authRouteParams.RESET_PASSWORD:
      return (
        <Auth title="Please enter your new password">
          <ResetPasswordForm />
        </Auth>
      )
    case authRouteParams.SETUP:
      isActive().then((data) => {
        if (data) {
          console.log('Data = true')
          replace(authRoutes.LOGIN)
          return null
        }
      })
      return (
        <Auth title="Please enter credentials for initial user">
          <SetupForm />
        </Auth>
      )
    case authRouteParams.LOGIN:
    default: {
      return (
        <Auth title="Authenticate to access CanIDS">
          <LoginForm />
        </Auth>
      )
    }
  }
}

export async function getStaticProps() {
  return { props: {} }
}

export async function getStaticPaths() {
  return {
    paths: [
      { params: { authRoute: authRouteParams.REGISTER } },
      { params: { authRoute: authRouteParams.FORGOT_PASSWORD } },
      { params: { authRoute: authRouteParams.RESET_PASSWORD } },
      { params: { authRoute: authRouteParams.LOGIN } },
      { params: { authRoute: authRouteParams.SETUP } },
    ],
    fallback: false,
  }
}
