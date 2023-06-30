import { useEffect, useState, useCallback } from 'react'
import { getApiErrorMessage } from '@utils/commonApiProcessing'
import { useCookies } from 'react-cookie'
import useNotification, { NotificationType } from '@context/notificationContext'

interface RequestProps {
  params?: any
  token?: string
}

interface UseRequestProps {
  request: (p: RequestProps | any) => Promise<any>
  requestByDefault?: boolean
  params?: any
  needSuccess?: boolean
}

const useRequest = ({
  request,
  requestByDefault = true,
  params,
  needSuccess,
}: UseRequestProps) => {
  const {
    addNotification,
  }: {
    addNotification: (e: any, type?: NotificationType | undefined) => void
  } = useNotification()
  const [cookies] = useCookies(['jwt'])

  const [data, setData] = useState<undefined | any>()
  const [completed, setCompleted] = useState(false)
  const [loading, setloading] = useState(false)

  const makeRequest = useCallback(
    async (p?: any) => {
      let r = undefined
      setloading(true)
      try {
        r = p
          ? await request({ params: p, token: cookies.jwt })
          : await request({ token: cookies.jwt })
        setData(r)
        needSuccess && addNotification('Successful request', 'success')
      } catch (e: any) {
        setData(undefined)
        addNotification(getApiErrorMessage(e))
      }
      setloading(false)
      return r
    },
    [addNotification, cookies.jwt, needSuccess, request]
  )

  useEffect(() => {
    if (!completed && requestByDefault) {
      setCompleted(true)
      makeRequest(params)
    }
  }, [completed, makeRequest, params, requestByDefault])

  return { data, loading, makeRequest }
}

export default useRequest
