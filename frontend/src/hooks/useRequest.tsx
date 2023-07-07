import { useEffect, useState, useCallback } from 'react'
import { getApiErrorMessage } from '@utils/commonApiProcessing'
import { useCookies } from 'react-cookie'
import useNotification, { NotificationType } from '@context/notificationContext'
import { allCookies as ac } from '@constants/cookies'

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
  const [cookies] = useCookies([ac.STATE])

  const [data, setData] = useState<undefined | any>()
  const [completed, setCompleted] = useState(false)
  const [loading, setloading] = useState(true)

  const makeRequest = useCallback(
    async (p?: any) => {
      let r = undefined
      try {
        r = p
          ? await request({ params: p, token: cookies[ac.STATE] })
          : await request({ token: cookies[ac.STATE] })
        setData(r)
        needSuccess && addNotification('Successful request', 'success')
      } catch (e: any) {
        setData(undefined)
        addNotification(getApiErrorMessage(e))
      }
      setloading(false)
      return r
    },
    [addNotification, cookies, needSuccess, request]
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
