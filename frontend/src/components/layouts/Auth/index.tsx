import { ReactNode } from 'react'
import Image from 'next/image'
import { Box } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import logo from '@images/logoGrey.png'
import styles from './styles.module.scss'

type Props = {
  children: ReactNode
}

export default ({ children }: Props) => {
  return (
    <Box sx={{ flexGrow: 1, p: 0 }}>
      <Grid container justifyContent="center">
        <Grid xs={4}>
          <div className={styles.logo}>
            <Image src={logo} alt={'Canids'} width={200} priority={true} />
          </div>
          {children}
        </Grid>
      </Grid>
    </Box>
  )
}
