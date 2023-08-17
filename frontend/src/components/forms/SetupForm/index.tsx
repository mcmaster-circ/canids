import { useCallback } from 'react'
import { useForm } from 'react-hook-form'
import { Button, Typography } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import useAuth from '@context/authContext'
import { defaultValues, setupFormConfig, resolver } from './constants'
import { FormRender } from '@molecules'
import packageJson from 'package.json'
import { SetupProps } from '@constants/types'

export default () => {
  const { setup } = useAuth()

  const onSubmit = useCallback(
    (data: SetupProps) => {
      setup(data)
    },
    [setup]
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
          {setupFormConfig.map((c) => (
            <FormRender key={c.name} {...c} errors={errors} control={control} />
          ))}
          <Grid xs={12}>
            <Button
              variant="contained"
              color="secondary"
              type="submit"
              fullWidth
            >
              INITIALIZE
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
