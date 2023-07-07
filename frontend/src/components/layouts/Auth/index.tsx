import { ReactNode } from 'react'
import Image from 'next/image'
import Grid from '@mui/material/Unstable_Grid2'
import logo from '@images/wideLogo.png'
import styles from './styles.module.scss'
import { Typography } from '@mui/material'

type Props = {
  children: ReactNode
  title: string
}

export default ({ children, title }: Props) => {
  return (
    <Grid
      container
      justifyContent="center"
      alignItems="center"
      sx={{ minHeight: '100vh' }}
    >
      <Grid width="100%" maxWidth={420} m={2} spacing={2}>
        <div className={styles.logo}>
          <Image src={logo} alt={'Canids'} priority={true} />
        </div>
        <Typography
          variant="h4"
          fontWeight={700}
          textAlign="center"
          mb={6}
          color="gray"
        >
          McMaster CanIDS
        </Typography>
        <Typography variant="body2" textAlign="center" mb={2}>
          {title}
        </Typography>
        {children}
      </Grid>
    </Grid>
  )
}
