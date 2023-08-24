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
import {
  login as loginApiCall,
  logout as logoutApiCall,
  isActive as isActiveApiCall,
  setup as setupApiCall,
  resetPassword as resetPasswordApiCall,
} from '@api/auth'
import { userInfo } from '@api/user'
import { useRequest } from '@hooks'
import useNotification, { NotificationType } from '@context/notificationContext'
import { userProfileCookies, allCookies as ac } from '@constants/cookies'
import { LoginProps, ResetProps, SetupProps, UserProps } from '@constants/types'

interface AuthContextType {
  user?: UserProps
  loading: boolean
  logedIn?: boolean
  login: (d: LoginProps) => void
  logout: () => void
  isActive: () => Promise<boolean>
  setup: (d: SetupProps) => void
  resetPassword: (r: ResetProps) => void
}

const AuthContext = createContext<AuthContextType>({} as AuthContextType)

// Export the provider as we need to wrap the entire app with it
export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const {
    addNotification,
  }: {
    addNotification: (e: any, type?: NotificationType | undefined) => void
  } = useNotification()
  const { makeRequest: loginRequest, loading: loginLoading } = useRequest({
    requestByDefault: false,
    request: loginApiCall,
    needSuccess: 'Successfull Login',
  })
  const { makeRequest: setupRequest, loading: setupLoading } = useRequest({
    requestByDefault: false,
    request: setupApiCall,
    needSuccess: 'Successfull Login',
  })
  const { makeRequest: userInfoRequest, loading: userLoading } = useRequest({
    requestByDefault: false,
    request: userInfo,
  })
  const { makeRequest: logoutRequest } = useRequest({
    requestByDefault: false,
    request: logoutApiCall,
    needSuccess: 'Successfull Logout',
  })
  const { makeRequest: resetPasswordRequest } = useRequest({
    requestByDefault: false,
    request: resetPasswordApiCall,
    needSuccess: 'Successfull reset',
  })
  const { makeRequest: isActiveRequest } = useRequest({
    requestByDefault: false,
    request: isActiveApiCall,
  })
  const [cookies, setCookie, removeCookie] = useCookies(
    Object.values(ac) as string[]
  )

  const [user, setUser] = useState<UserProps>()
  const [logedIn, setLogedIn] = useState<boolean>(false)
  const [loadingInitial, setLoadingInitial] = useState<boolean>(true)

  const setUserFields = useCallback((fields: UserProps | any) => {
    const cachedUser: UserProps | any = {}
    userProfileCookies.forEach((f) => (cachedUser[f] = fields[f]))
    setUser(cachedUser)
  }, [])

  const isActive = useCallback(async () => {
    const res: any = await isActiveRequest()
    if (res) {
      return Boolean(res.active)
    }
    return false
  }, [isActiveRequest])

  const login = useCallback(
    async (d: LoginProps) => {
      await loginRequest({ ...d })
      const res: any = await userInfoRequest()
      if (res) {
        userProfileCookies.forEach((f) => setCookie(f, res[f], { path: '/' }))
        setUserFields(res)
        setLogedIn(true)
      }
    },
    [loginRequest, setCookie, setUserFields, userInfoRequest]
  )

  const setup = useCallback(
    async (d: SetupProps) => {
      await setupRequest({ ...d })
      const res: any = await userInfoRequest()
      if (res) {
        userProfileCookies.forEach((f) => setCookie(f, res[f], { path: '/' }))
        setUserFields(res)
        setLogedIn(true)
      }
    },
    [setupRequest, setCookie, setUserFields, userInfoRequest]
  )

  const logout = useCallback(async () => {
    setUser(undefined)
    await logoutRequest()
    Object.values(ac).forEach((f) => removeCookie(f, { path: '/' }))
    setLogedIn(false)
  }, [logoutRequest, removeCookie])

  const resetPassword = useCallback(
    async (r: ResetProps) => {
      await resetPasswordRequest({ ...r })
    },
    [resetPasswordRequest]
  )

  // Check if there is a currently active session
  // when the provider is mounted for the first time.
  useEffect(() => {
    if (!!cookies[ac.ROLE] && loadingInitial) {
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
      loading: loginLoading || userLoading,
      logedIn,
      login,
      logout,
      isActive,
      setup,
      resetPassword,
    }),
    [
      user,
      loginLoading,
      userLoading,
      logedIn,
      login,
      logout,
      isActive,
      setup,
      resetPassword,
    ]
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
