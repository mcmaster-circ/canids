import {
  createContext,
  useState,
  useEffect,
  useCallback,
  useMemo,
  ReactNode,
  useContext,
} from 'react'
import { useCookies } from 'react-cookie'
import { login as loginApiCall } from '@api/auth'
import { userInfo } from '@api/user'
import { useRequest } from '@hooks'
import useNotification, { NotificationType } from '@context/notificationContext'
import { userProfileCookies, allCookies as ac } from '@constants/cookies'
import { LoginProps, UserProps } from '@constants/types'

interface AuthContextType {
  user?: UserProps
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
  const { makeRequest: userInfoRequest } = useRequest({
    requestByDefault: false,
    request: userInfo,
  })
  const [cookies, setCookie, removeCookie] = useCookies(
    Object.values(ac) as string[]
  )

  const [user, setUser] = useState<UserProps>()
  const [logedIn, setLogedIn] = useState<boolean>(false)
  const [loading, setLoading] = useState<boolean>(false)
  const [loadingInitial, setLoadingInitial] = useState<boolean>(true)

  const setUserFields = useCallback((fields: UserProps | any) => {
    const cachedUser: UserProps | any = {}
    userProfileCookies.forEach((f) => (cachedUser[f] = fields[f]))
    setUser(cachedUser)
  }, [])

  const login = useCallback(
    async (d: LoginProps) => {
      setLoading(true)
      try {
        // const res = await loginRequest({ user: { ...d } })
        const res = await userInfoRequest()
        console.log(res)
        if (res) {
          userProfileCookies.forEach((f) => setCookie(f, res[f], { path: '/' }))
          setUserFields(res)
          setLogedIn(true)
          addNotification('Successfull Login', 'success')
        }
      } catch (e) {
        addNotification(e)
      }
      setLoading(false)
    },
    [addNotification, setCookie, setUserFields, userInfoRequest]
  )

  const logout = useCallback(() => {
    setUser(undefined)
    Object.values(ac).forEach((f) => removeCookie(f, { path: '/' }))
    addNotification('Successfully logged out', 'success')
    setLogedIn(false)
  }, [addNotification, removeCookie])

  // Check if there is a currently active session
  // when the provider is mounted for the first time.
  useEffect(() => {
    if (!!cookies[ac.STATE] && loadingInitial) {
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
