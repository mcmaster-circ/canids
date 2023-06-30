import { getHeaders, postHeaders } from '@utils/commonApiProcessing'
import { GetProps, PostProps } from '@constants/types'

export const get = async ({ url, params, token }: GetProps) => {
  if (params) {
    url += '?' + new URLSearchParams(params).toString()
  }
  const res = await fetch(url, {
    headers: getHeaders(token),
    method: 'GET',
  })
  const data = await res.json()
  if (res.ok) {
    return data
  }
  throw data
}

export const post = async ({ url, body, token }: PostProps) => {
  const res = await fetch(url, {
    headers: postHeaders(token),
    method: 'POST',
    body: JSON.stringify(body),
  })
  const data = await res.json()
  if (res.ok) {
    return data
  }
  throw data
}

export const put = async ({ url, body, token }: PostProps) => {
  const res = await fetch(url, {
    headers: postHeaders(token),
    method: 'PUT',
    body: JSON.stringify(body),
  })
  const data = await res.json()
  if (res.ok) {
    return data
  }
  throw data
}

export const patch = async ({ url, body, token }: PostProps) => {
  const res = await fetch(url, {
    headers: postHeaders(token),
    method: 'PATCH',
    body: JSON.stringify(body),
  })
  const data = await res.json()
  if (res.ok) {
    return data
  }
  throw data
}
