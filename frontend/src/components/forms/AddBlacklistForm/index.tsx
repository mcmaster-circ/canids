import { useCallback } from 'react'
import { useForm } from 'react-hook-form'
import { Button } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { FormRender } from '@molecules'
import { AddBlacklistProps, UpdateBlacklistProps } from '@constants/types'
import { useRequest } from '@hooks'
import { addBlacklist, updateBlacklist } from '@api/blacklist'
import { defaultValues, addFormConfig, resolver } from './constants'

interface FormProps {
  handleClose: () => void
  isUpdate?: boolean
  values?: UpdateBlacklistProps
}

export default ({ handleClose, isUpdate, values }: FormProps) => {
  const { makeRequest } = useRequest({
    request: addBlacklist,
    requestByDefault: false,
    needSuccess: 'Blacklist Added',
  })
  const { makeRequest: update } = useRequest({
    request: updateBlacklist,
    requestByDefault: false,
    needSuccess: 'Blacklist Updated',
  })

  const onSubmit = useCallback(
    async (data: AddBlacklistProps) => {
      isUpdate
        ? await update({ ...data, uuid: values?.uuid })
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
          {addFormConfig.map((c) => (
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
