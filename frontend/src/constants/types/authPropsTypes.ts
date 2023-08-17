export interface LoginProps {
  uuid: string
  password: string
}

export interface RegisterProps {
  name: string
  uuid: string
  password: string
  passwordConfirm: string
}

export interface ForgotProps {
  uuid: string
}

export interface ResetProps {
  password: string
  passwordConfirm: string
}

export interface SetupProps {
  name: string
  uuid: string
  password: string
  passwordConfirm: string
}
