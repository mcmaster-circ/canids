import {
  AddClientProps,
  DeleteClientProps,
} from '@constants/types/ingestionPropsTypes'

import { get, post } from './fetchRequests'

const baseUrl = process.env.NEXT_PUBLIC_API_URL

export const ingestionList = async () => {
  const data = await get({ url: baseUrl + '/ingestion/list' })
  return data?.clients
}

export const ingestionAdd = async ({ params }: { params: AddClientProps }) => {
  const data = await post({ url: baseUrl + '/ingestion/create', body: params })
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
