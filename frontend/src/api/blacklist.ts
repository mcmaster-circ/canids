import { get, post } from './fetchRequests'
import {
  AddBlacklistProps,
  UpdateBlacklistProps,
  DeleteBlacklistProps,
} from '@constants/types'

const baseUrl = process.env.NEXT_PUBLIC_API_URL

export const getBlacklist = async ({
  params,
}: {
  params: { getNames?: boolean }
}) => {
  const data = await get({ url: baseUrl + '/blacklist/list' })
  return params?.getNames
    ? data?.blacklists?.map((l: any) => l.name)
    : data?.blacklists
}

export const addBlacklist = async ({
  params,
}: {
  params: AddBlacklistProps
}) => {
  const data = await post({ url: baseUrl + '/blacklist/add', body: params })
  return data
}

export const updateBlacklist = async ({
  params,
}: {
  params: UpdateBlacklistProps
}) => {
  const data = await post({ url: baseUrl + '/blacklist/update', body: params })
  return data
}

export const deleteBlacklist = async ({
  params,
}: {
  params: DeleteBlacklistProps
}) => {
  const data = await post({ url: baseUrl + '/blacklist/delete', body: params })
  return data
}
