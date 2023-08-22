import {
  ApproveClientProps,
  DeleteClientProps,
  RenameClientProps,
} from '@constants/types/ingestionPropsTypes'

import { get, post } from './fetchRequests'

const baseUrl = process.env.NEXT_PUBLIC_API_URL

export const ingestionList = async () => {
  const data = await get({ url: baseUrl + '/ingestion/list' })
  return data?.clients
}

export const ingestionApprove = async ({
  params,
}: {
  params: ApproveClientProps
}) => {
  const data = await post({ url: baseUrl + '/ingestion/approve', body: params })
  return data
}

export const ingestionDelete = async ({
  params,
}: {
  params: DeleteClientProps
}) => {
  const data = await post({ url: baseUrl + '/ingestion/delete', body: params })
  return data
}

export const ingestionRename = async ({
  params,
}: {
  params: RenameClientProps
}) => {
  const data = await post({ url: baseUrl + '/ingestion/rename', body: params })
  return data
}
