import {
  AddUserProps,
  DeleteUserProps,
  ResetUserPassProps,
  UpdateUserProps,
} from '@constants/types'
import { get, post } from './fetchRequests'

const baseUrl = process.env.NEXT_PUBLIC_API_URL

export const userList = async () => {
  const data = await get({ url: baseUrl + '/user/list' })
  return data
}

export const userInfo = async () => {
  const data = await get({ url: baseUrl + '/user/info' })
  return data
}

export const addUser = async ({ params }: { params: AddUserProps }) => {
  const data = await post({ url: baseUrl + '/user/add', body: params })
  return data
}

export const updateUser = async ({ params }: { params: UpdateUserProps }) => {
  const data = await post({
    url: baseUrl + '/user/update',
    body: params,
    params: { uuid: params.userId },
  })
  return data
}

export const resetUserPass = async ({
  params,
}: {
  params: ResetUserPassProps
}) => {
  const data = await post({ url: baseUrl + '/user/resetPass', body: params })
  return data
}

export const deleteUser = async ({ params }: { params: DeleteUserProps }) => {
  const data = await post({ url: baseUrl + '/user/delete', body: params })
  return data
}
