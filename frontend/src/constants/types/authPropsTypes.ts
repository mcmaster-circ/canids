export interface LoginProps {
  user: string
  pass: string
}

export interface RegisterProps {
  name: string
  email: string
  pass: string
  passConfirm: string
}

export interface ForgotProps {
  email: string
}

export interface ResetProps {
  pass: string
  passConfirm: string
}
