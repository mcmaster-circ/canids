import {
  createContext,
  useState,
  useEffect,
  useCallback,
  useMemo,
  ReactNode,
  useContext,
} from 'react'
import { addDays } from 'date-fns'
import { useCookies } from 'react-cookie'
import { login as loginApiCall } from '@api/auth'
import { useRequest } from '@hooks'
import { userProfileCookies, allCookies } from '@constants/cookies'
import useNotification, { NotificationType } from '@context/notificationContext'
import { LoginProps } from '@forms'

interface User {
  jwt?: string | null
  name?: string | null
  email?: string | null
}

interface AuthContextType {
  user?: User
  loading: boolean
  logedIn?: boolean
  login: (d: LoginProps) => void
  logout: () => void
}

const AuthContext = createContext<AuthContextType>({} as AuthContextType)

// Export the provider as we need to wrap the entire app with it
export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const {
    addNotification,
  }: {
    addNotification: (e: any, type?: NotificationType | undefined) => void
  } = useNotification()
  const { makeRequest: loginRequest } = useRequest({
    requestByDefault: false,
    request: loginApiCall,
  })
  const [cookies, setCookie, removeCookie] = useCookies(
    Object.values(allCookies) as string[]
  )

  const [user, setUser] = useState<User>()
  const [logedIn, setLogedIn] = useState<boolean>(false)
  const [loading, setLoading] = useState<boolean>(false)
  const [loadingInitial, setLoadingInitial] = useState<boolean>(true)

  const setUserFields = useCallback((fields: User | any) => {
    const cachedUser: User | any = {}
    userProfileCookies.forEach((f) => (cachedUser[f] = fields[f]))
    setUser(cachedUser)
  }, [])

  const login = useCallback(
    async (d: LoginProps) => {
      setLoading(true)
      try {
        // const res = await loginRequest({ user: { ...d } })
        // TODO: Remove when auth is Setted Up
        const res: any = {
          email: 'example@mail.com',
          name: 'Some User',
          jwt: 'some_jwt',
        }
        if (res) {
          userProfileCookies.forEach((f) =>
            setCookie(f, res[f], { path: '/', expires: addDays(new Date(), 7) })
          )
          setUserFields(res)
          setLogedIn(true)
          addNotification('Successfull Login', 'success')
        }
      } catch (e) {
        addNotification(e)
      }
      setLoading(false)
    },
    [addNotification, setCookie, setUserFields]
  )

  const logout = useCallback(() => {
    setUser(undefined)
    Object.values(allCookies).forEach((f) => removeCookie(f, { path: '/' }))
    addNotification('Successfully logged out', 'success')
    setLogedIn(false)
  }, [addNotification, removeCookie])

  // Check if there is a currently active session
  // when the provider is mounted for the first time.
  useEffect(() => {
    if (!!cookies[allCookies.JWT] && loadingInitial) {
      setUserFields(cookies)
      setLogedIn(true)
      setLoadingInitial(false)
    }
    if (loadingInitial) {
      setLoadingInitial(false)
    }
  }, [cookies, loadingInitial, setUserFields, user])

  const memoedValue = useMemo(
    () => ({
      user,
      loading,
      logedIn,
      login,
      logout,
    }),
    [user, loading, logedIn, login, logout]
  )

  return (
    <AuthContext.Provider value={memoedValue}>
      {!loadingInitial && children}
    </AuthContext.Provider>
  )
}

export default function useAuth() {
  return useContext(AuthContext)
}
