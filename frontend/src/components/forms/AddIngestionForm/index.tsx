import { useCallback } from 'react'
import { useForm } from 'react-hook-form'
import { Button } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { FormRender } from '@molecules'
import { useRequest } from '@hooks'
import { defaultValues, formConfig, resolver } from './constants'
import { AddClientProps } from '@constants/types/ingestionPropsTypes'
import { ingestionAdd } from '@api/ingestion'

interface FormProps {
  handleClose: (uuid: string, key: string) => void
  isUpdate?: boolean
  values?: AddClientProps
}

export default ({ handleClose, isUpdate, values }: FormProps) => {
  const { makeRequest } = useRequest({
    request: ingestionAdd,
    requestByDefault: false,
    needSuccess: 'Ingestion client has been created',
  })

  const onSubmit = useCallback(
    async (data: AddClientProps) => {
      var resp = await makeRequest(data)
      handleClose(data.uuid, resp.key)
    },
    [handleClose, makeRequest]
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
