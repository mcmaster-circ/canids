import { get } from './fetchRequests'

const baseUrl = process.env.NEXT_PUBLIC_API_URL

export const getFields = async ({ filter }: { filter?: boolean }) => {
  const data = await get({ url: baseUrl + '/fields/list' })
  return filter
    ? data
        ?.filter((l: any) => l.index.includes('.alarm'))
        .map((l: any) => l.index)
    : data
}
