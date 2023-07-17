import { ReactNode, useEffect } from 'react'
import { useRouter } from 'next/router'
import useAuth from '@context/authContext'
import { Header } from '@organisms'
import styles from './styles.module.scss'
import { authRoutes } from '@constants/routes'

interface Props {
  children: ReactNode
}

export default ({ children }: Props) => {
  const router = useRouter()
  const { logedIn } = useAuth()

  useEffect(() => {
    if (!logedIn) {
      router.push(authRoutes.LOGIN)
    }
  }, [logedIn, router])

  return (
    <div className={styles.container}>
      <Header />
      {children}
    </div>
  )
}
