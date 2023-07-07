import { get, post } from './fetchRequests'
import {
  AddViewProps,
  UpdateViewProps,
  DeleteViewProps,
} from '@constants/types'

const baseUrl = process.env.NEXT_PUBLIC_API_URL

export const getViewList = async () => {
  const data = await get({ url: baseUrl + '/view/list' })
  return data?.views
}

export const addView = async ({ params }: { params: AddViewProps }) => {
  const data = await post({ url: baseUrl + '/view/add', body: params })
  return data
}

export const updateView = async ({ params }: { params: UpdateViewProps }) => {
  const data = await post({ url: baseUrl + '/view/update', body: params })
  return data
}

export const deleteView = async ({ params }: { params: DeleteViewProps }) => {
  const data = await post({ url: baseUrl + '/view/delete', body: params })
  return data
}
