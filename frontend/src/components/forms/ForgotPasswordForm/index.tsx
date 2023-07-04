import { useCallback } from 'react'
import { useRouter } from 'next/router'
import { useForm } from 'react-hook-form'
import { Button, Divider, Link, Typography } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { defaultValues, loginFormConfig, resolver } from './constants'
import { FormRender } from '@molecules'
import packageJson from 'package.json'
import { authRoutes } from '@constants/routes'
import { ForgotProps } from '@constants/types'

export default () => {
  const { push } = useRouter()

  const onSubmit = useCallback(
    (data: ForgotProps) => {
      console.log(data)
      push(authRoutes.RESET_PASSWORD)
    },
    [push]
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
        <Grid container spacing={3} justifyContent="center">
          {loginFormConfig.map((c) => (
            <FormRender key={c.name} {...c} errors={errors} control={control} />
          ))}
          <Grid xs={12}>
            <Button variant="contained" color="primary" type="submit" fullWidth>
              SUBMIT
            </Button>
          </Grid>
          <Grid
            container
            justifyContent="center"
            columnSpacing={4}
            rowSpacing={1}
            mt={6}
          >
            <Grid>
              <Link
                variant="body1"
                underline="none"
                alignSelf="center"
                href={authRoutes.LOGIN}
              >
                Back to Login
              </Link>
            </Grid>
            <Divider orientation="vertical" flexItem />
            <Grid>
              <Typography variant="body1" alignSelf="center" color="gray">
                v: {packageJson.version}
              </Typography>
            </Grid>
          </Grid>
        </Grid>
      </form>
    </>
  )
}
