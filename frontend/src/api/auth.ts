import { post } from './fetchRequests'
import {
  LoginProps,
  RegisterProps,
  ForgotProps,
  ResetProps,
} from '@constants/types'

const baseUrl = process.env.NEXT_PUBLIC_API_URL

export const login = async ({ params }: { params: LoginProps }) => {
  const data = await post({ url: baseUrl + '/login', body: params })
  return data
}

export const register = async ({ params }: { params: RegisterProps }) => {
  const data = await post({ url: baseUrl + '/register', body: params })
  return data
}

export const forgotPassword = async ({ params }: { params: ForgotProps }) => {
  const data = await post({ url: baseUrl + '/requestReset', body: params })
  return data
}

export const resetPassword = async ({ params }: { params: ResetProps }) => {
  const data = await post({ url: baseUrl + '/reset', body: params })
  return data
}
