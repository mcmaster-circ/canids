import { get, post } from './fetchRequests'
import {
  ListConfigurationProps,
  UpdateConfigurationProps,
} from '@constants/types'

const baseUrl = process.env.NEXT_PUBLIC_API_URL

export const getConfiguration = async ({
  params,
}: {
  params: { getNames?: boolean }
}) => {
  const data: ListConfigurationProps = await get({
    url: baseUrl + '/configuration/list',
  })
  return params?.getNames
    ? data?.configuration?.map((l: any) => l.name)
    : data?.configuration
}

export const updateConfiguration = async ({
  params,
}: {
  params: UpdateConfigurationProps
}) => {
  const data = await post({
    url: baseUrl + '/configuration/update',
    body: params,
  })
  return data
}
