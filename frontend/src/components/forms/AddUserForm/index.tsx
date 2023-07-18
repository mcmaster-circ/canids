import { useCallback } from 'react'
import { useForm } from 'react-hook-form'
import { Button } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { FormRender } from '@molecules'
import { AddUserProps, UpdateUserProps } from '@constants/types'
import { useRequest } from '@hooks'
import { addUser, updateUser } from '@api/user'
import { defaultValues, formConfig, resolver } from './constants'

interface FormProps {
  handleClose: () => void
  isUpdate?: boolean
  values?: UpdateUserProps
}

export default ({ handleClose, isUpdate, values }: FormProps) => {
  const { makeRequest } = useRequest({
    request: addUser,
    requestByDefault: false,
    needSuccess:
      'The user account has been successfully created. The user has been emailed to complete account activation',
  })
  const { makeRequest: update } = useRequest({
    request: updateUser,
    requestByDefault: false,
    needSuccess:
      'The user account has been successfully updated. Changes will be applied within 1 minute',
  })

  const onSubmit = useCallback(
    async (data: AddUserProps) => {
      isUpdate
        ? await update({
            ...data,
            activated: data.activated === 'true' ? true : false,
            uuid: values?.uuid,
          })
        : await makeRequest(data)
      handleClose()
    },
    [handleClose, isUpdate, makeRequest, update, values?.uuid]
  )

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm({
    resolver: resolver,
    defaultValues: isUpdate ? values : defaultValues,
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
