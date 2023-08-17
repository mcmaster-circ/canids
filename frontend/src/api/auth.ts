import { get, post } from './fetchRequests'
import {
  LoginProps,
  RegisterProps,
  ForgotProps,
  ResetProps,
  SetupProps,
} from '@constants/types'

const baseUrl = process.env.NEXT_PUBLIC_API_URL

export const login = async ({ params }: { params: LoginProps }) => {
  const data = await post({ url: baseUrl + '/auth/login', body: params })
  return data
}

export const logout = async () => {
  const data = await post({ url: baseUrl + '/auth/logout' })
  return data
}

export const register = async ({ params }: { params: RegisterProps }) => {
  const data = await post({ url: baseUrl + '/auth/registerUser', body: params })
  return data
}

export const forgotPassword = async ({ params }: { params: ForgotProps }) => {
  const data = await post({ url: baseUrl + '/auth/requestReset', body: params })
  return data
}

export const resetPassword = async ({ params }: { params: ResetProps }) => {
  const data = await post({
    url: baseUrl + '/auth/resetPassword',
    body: params,
  })
  return data
}

export const isActive = async () => {
  const data = await get({
    url: baseUrl + '/auth/isActive',
  })
  return data
}

export const setup = async ({ params }: { params: SetupProps }) => {
  const data = await post({ url: baseUrl + '/auth/setup', body: params })
  return data
}
