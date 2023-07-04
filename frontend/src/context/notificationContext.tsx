import {
  createContext,
  useState,
  useEffect,
  useCallback,
  useMemo,
  ReactNode,
  useContext,
} from 'react'
import { Snackbar, Alert } from '@mui/material'

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
        <Snackbar
          open={!!notification}
          onClose={removeNotification}
          anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
        >
          <Alert
            onClose={removeNotification}
            severity={notificationType}
            variant="filled"
            sx={{ width: '100%' }}
          >
            {notification}
          </Alert>
        </Snackbar>
      )}
    </NotificationContext.Provider>
  )
}

export default function useNotification() {
  return useContext(NotificationContext)
}
