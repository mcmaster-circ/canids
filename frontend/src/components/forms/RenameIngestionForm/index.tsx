import { useCallback } from 'react'
import { useForm } from 'react-hook-form'
import { Button } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { FormRender } from '@molecules'
import { useRequest } from '@hooks'
import { formConfig, resolver } from './constants'
import {
  RenameClientProps,
  RenameNameHouser,
} from '@constants/types/ingestionPropsTypes'
import { ingestionRename } from '@api/ingestion'

interface FormProps {
  handleClose: () => void
  values: RenameClientProps
}

export default ({ handleClose, values }: FormProps) => {
  const { makeRequest } = useRequest({
    request: ingestionRename,
    requestByDefault: false,
  })

  const onSubmit = useCallback(
    async (name: RenameNameHouser) => {
      console.log(name)

      var req: RenameClientProps = {
        name: name.name,
        uuid: values.uuid,
      }
      console.log(req)
      await makeRequest(req)
      handleClose()
    },
    [handleClose, makeRequest, values]
  )

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm({
    resolver: resolver,
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
              Rename
            </Button>
          </Grid>
        </Grid>
      </form>
    </>
  )
}
