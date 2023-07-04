import { post } from './fetchRequests'
import { GetAlarmsProps } from '@constants/types'

const baseUrl = process.env.NEXT_PUBLIC_API_URL

export const getAlarms = async ({ params }: { params: GetAlarmsProps }) => {
  const data = await post({ url: baseUrl + '/alarm/data', body: params })
  return data
}
