import { post, patch } from './fetchRequests'

export type loginProps = {
  params: {
    email: string
    password: string
  }
}

const baseUrl = process.env.NEXT_PUBLIC_API_URL

// No auth required
export const login = async ({ params }: loginProps) => {
  const data = await post({ url: baseUrl + '/login', body: params })
  return data
}

export type ChangePasswordProps = {
  params: {
    user: {
      email: string
      password: string
      new_password: string
    }
  }
  token: string
}

// auth token required
export const changePassword = async ({
  params,
  token,
}: ChangePasswordProps) => {
  const data = await patch({
    url: baseUrl + '/change-password',
    body: params,
    token,
  })
  return data
}
