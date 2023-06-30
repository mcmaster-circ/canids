import { useRouter } from 'next/router'
import { authRoutes } from '@constants/routes'

export default () => {
  const { replace } = useRouter()
  replace(authRoutes.LOGIN)

  return null
}
