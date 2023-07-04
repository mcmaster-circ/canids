import { get, post } from './fetchRequests'
import { UpdateDashboardProps } from '@constants/types'

const baseUrl = process.env.NEXT_PUBLIC_API_URL

export const getDashboard = async () => {
  const data = await get({ url: baseUrl + '/dashboard/get' })
  return data
}

export const updateDashboard = async ({
  params,
}: {
  params: UpdateDashboardProps
}) => {
  const data = await post({ url: baseUrl + '/dashboard/update', body: params })
  return data
}
