import { ReactNode, useEffect, useState } from 'react'
import { useRouter } from 'next/router'
import useAuth from '@context/authContext'
// import { Button } from '@mui/material'
// import { ModalExample } from '@modals'
import { Header } from '@organisms'
import styles from './styles.module.scss'
import { authRoutes } from '@constants/routes'

interface Props {
  children: ReactNode
}

export default ({ children }: Props) => {
  const router = useRouter()
  const { logedIn } = useAuth()
  const [openModal, setOpenModal] = useState(false)

  useEffect(() => {
    if (!logedIn) {
      router.push(authRoutes.LOGIN)
    }
  }, [logedIn, router])

  return (
    <div className={styles.container}>
      <Header />
      {/* <Button onClick={() => setOpenModal(true)}>Open Modal</Button> */}
      {/* <ModalExample open={openModal} handleClose={() => setOpenModal(false)} /> */}
      {children}
    </div>
  )
}
