import { useCallback } from 'react'
import { useForm } from 'react-hook-form'
import { Button } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { FormRender } from '@molecules'
import { EditUserProps, UpdateUserProps } from '@constants/types'
import { useRequest } from '@hooks'
import { updateUser } from '@api/user'
import { formConfig, resolver } from './constants'

interface FormProps {
  handleClose: () => void
  values?: UpdateUserProps
}

export default ({ handleClose, values }: FormProps) => {
  const { makeRequest: update } = useRequest({
    request: updateUser,
    requestByDefault: false,
    needSuccess:
      'The user account has been successfully updated. Changes will be applied within 1 minute',
  })

  const onSubmit = useCallback(
    async (data: EditUserProps) => {
      await update({
        ...data,
        activated: data.activated === 'true' ? true : false,
        uuid: values?.uuid,
      })
      handleClose()
    },
    [handleClose, update, values?.uuid]
  )

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm({
    resolver: resolver,
    defaultValues: values,
  })

  return (
    <>
      <form onSubmit={handleSubmit(onSubmit)}>
        <Grid container spacing={3} justifyContent="center">
          {formConfig.map((c) => (
            <FormRender key={c.name} {...c} errors={errors} control={control} />
          ))}
          <Grid xs={12} sx={{ display: 'flex', justifyContent: 'center' }}>
            <Button variant="contained" color="secondary" type="submit">
              Save
            </Button>
          </Grid>
        </Grid>
      </form>
    </>
  )
}
