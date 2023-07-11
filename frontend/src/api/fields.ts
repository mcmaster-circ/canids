import { get } from './fetchRequests'

const baseUrl = process.env.NEXT_PUBLIC_API_URL

export const getFields = async () => {
  const data = await get({ url: baseUrl + '/fields/list' })
  return data
    ?.filter((l: any) => l.index.includes('.alarm'))
    .map((l: any) => l.index)
}
