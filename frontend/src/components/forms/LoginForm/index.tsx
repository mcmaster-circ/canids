import { useState, useCallback } from 'react'
import { useForm } from 'react-hook-form'
import {
  Button,
  Typography,
  Link,
  Checkbox,
  FormControlLabel,
} from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import useAuth from '@context/authContext'
import { defaultValues, loginFormConfig, resolver } from './constants'
import styles from './styles.module.scss'
import { FormRender } from '@molecules'

export default () => {
  const { login } = useAuth()
  const [isRemember, setIsRemember] = useState<boolean>(false)

  const onSubmit = useCallback(
    (data: { email: string; password: string }) => {
      console.log(data)
      login({ email: 'example@mail.com', password: 'securePassword' })
    },
    [login]
  )

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm({
    resolver: resolver,
    defaultValues,
  })

  return (
    <>
      <form onSubmit={handleSubmit(onSubmit)}>
        <Grid container spacing={4}>
          <Grid xs={12}>
            <Typography variant="h4">Log in</Typography>
          </Grid>
          {loginFormConfig.map((c) => (
            <FormRender key={c.name} {...c} errors={errors} control={control} />
          ))}
          <div className={styles.container}>
            <FormControlLabel
              control={
                <Checkbox
                  checked={isRemember}
                  onClick={() => setIsRemember(!isRemember)}
                />
              }
              label="Remember me"
            />
            <Link variant="body1" underline="none" alignSelf="center">
              Forgot your password ?
            </Link>
          </div>
          <Grid xs={12}>
            <Button variant="contained" color="primary" type="submit" fullWidth>
              Continue
            </Button>
          </Grid>
        </Grid>
      </form>
    </>
  )
}
