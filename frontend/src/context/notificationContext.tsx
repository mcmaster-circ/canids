import {
  createContext,
  useState,
  useEffect,
  useCallback,
  useMemo,
  ReactNode,
  useContext,
} from 'react'
import { useTheme } from '@mui/material/styles'

interface NotificationContextType {
  notification: any
  addNotification: (e: any, type?: NotificationType) => void
  removeNotification: () => void
}

export type NotificationType = 'error' | 'success' | 'warning' | 'info'

const NotificationContext = createContext<NotificationContextType>(
  {} as NotificationContextType
)

// Export the provider as we need to wrap the entire app with it
export const NotificationProvider = ({ children }: { children: ReactNode }) => {
  const theme: any = useTheme()
  const [notification, setNotification] = useState<any>()
  const [notificationType, setNotificationType] = useState<any>()

  const addNotification = useCallback(
    (e: any, type: NotificationType = 'error') => {
      e && setNotification(e)
      type && setNotificationType(type)
    },
    []
  )

  const removeNotification = useCallback(() => {
    setNotification(undefined)
    setNotificationType(undefined)
  }, [])

  //Remove notification after 5s
  useEffect(() => {
    if (notification) {
      setTimeout(removeNotification, 5000)
    }
  }, [notification, removeNotification, setNotification])

  const memoedValue = useMemo(
    () => ({
      notification,
      addNotification,
      removeNotification,
    }),
    [addNotification, notification, removeNotification]
  )

  return (
    <NotificationContext.Provider value={memoedValue}>
      {children}
      {notification && (
        <div
          style={{
            position: 'fixed',
            bottom: '50px',
            right: '0',
            backgroundColor: theme.palette[notificationType].main,
            color: theme.palette[notificationType].contrastText,
            padding: '16px',
            borderTopLeftRadius: '8px',
            borderBottomLeftRadius: '8px',
            zIndex: 100,
          }}
        >
          {notification}
        </div>
      )}
    </NotificationContext.Provider>
  )
}

export default function useNotification() {
  return useContext(NotificationContext)
}
