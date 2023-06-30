import { useEffect, useState } from 'react'
import { useRouter } from 'next/router'
import useAuth from '@context/authContext'
import { Button } from '@mui/material'
import { ModalExample } from '@modals'
import styles from './styles.module.scss'

export default () => {
  const router = useRouter()
  const { user, logout } = useAuth()
  const [openModal, setOpenModal] = useState(false)

  useEffect(() => {
    if (!user?.name) {
      router.push('/auth/login')
    }
  }, [router, user?.name])

  return (
    <>
      <div className={styles.container}>
        <Button onClick={logout}>Log Out</Button>
        <Button onClick={() => setOpenModal(true)}>Open Modal</Button>
        <ModalExample
          open={openModal}
          handleClose={() => setOpenModal(false)}
        />
      </div>
    </>
  )
}
