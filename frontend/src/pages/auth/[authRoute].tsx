import { useEffect } from 'react'
import { useRouter } from 'next/router'
import useAuth from '@context/authContext'
import { Auth } from '@layouts'
import {
  ForgotPasswordForm,
  LoginForm,
  RegisterForm,
  ResetPasswordForm,
} from '@forms'
import { authRouteParams, dashboardRoutes } from '@constants/routes'

export default () => {
  const { query, replace } = useRouter()
  const { logedIn } = useAuth()

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