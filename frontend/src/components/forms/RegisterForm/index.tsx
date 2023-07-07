import { useCallback } from 'react'
import { useForm } from 'react-hook-form'
import { Button, Divider, Link, Typography } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { defaultValues, registerFormConfig, resolver } from './constants'
import { FormRender } from '@molecules'
import packageJson from 'package.json'
import { authRoutes } from '@constants/routes'
import { useRouter } from 'next/router'
import { RegisterProps } from '@constants/types'

export default () => {
  const { push } = useRouter()
  const onSubmit = useCallback(
    (data: RegisterProps) => {
      console.log(data)
      push(authRoutes.LOGIN)
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
          {registerFormConfig.map((c) => (
            <FormRender key={c.name} {...c} errors={errors} control={control} />
          ))}
          <Grid xs={12}>
            <Button
              variant="contained"
              color="secondary"
              type="submit"
              fullWidth
            >
              REGISTER
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
                color="secondary"
              >
                Back to login
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
