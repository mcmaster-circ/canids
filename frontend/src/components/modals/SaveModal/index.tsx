import {
  IconButton,
  Divider,
  Dialog,
  DialogContent,
  DialogTitle,
  Typography,
  Button,
} from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { Close, Error } from '@mui/icons-material'
import { useCallback } from 'react'
import {
  UpdateSetting,
  Setting,
  UpdateConfigurationProps,
  BooleanSettings,
} from '@constants/types'

interface ModalTypes {
  data: BooleanSettings
  open: { open: boolean; label?: string; params?: any }
  handleClose: () => void
  title: string
  request: (v: any) => void
}

export default ({ data, open, handleClose, request, title }: ModalTypes) => {
  const handleClick = useCallback(async () => {
    const settings: UpdateSetting[] = Object.keys(data).map((key) => {
      const setting = data[key as keyof BooleanSettings] as Setting
      return {
        name: setting.name,
        value: setting.value,
      }
    })
    const req: UpdateConfigurationProps = {
      configuration: settings,
    }
    await request(req)
    handleClose()
  }, [handleClose, data, request])

  return (
    <Dialog open={open.open} onClose={handleClose}>
      <DialogTitle
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        Saving{' ' + title}
        <IconButton aria-label="close" onClick={handleClose}>
          <Close />
        </IconButton>
      </DialogTitle>
      <Divider variant="fullWidth" />
      <DialogContent>
        <Grid container justifyContent="center">
          <Grid xs={2}>
            <Error sx={{ fontSize: '48px' }} color="error" />
          </Grid>
          <Grid xs={10}>
            <Typography>
              Are you sure you want to make these changes?
            </Typography>
            <Grid
              container
              justifyContent="center"
              rowSpacing={2}
              columnSpacing={4}
            >
              {Object.keys(data).map((key) => {
                const setting = data[key as keyof BooleanSettings] as Setting
                if (setting.prevValue !== setting.value) {
                  return (
                    <Grid key={key}>
                      <Typography>
                        {setting.name} {setting.prevValue} {'->'}{' '}
                        {setting.value}
                      </Typography>
                    </Grid>
                  )
                }
              })}
            </Grid>
          </Grid>
          <Button
            variant="contained"
            color="error"
            onClick={handleClick}
            sx={{ mt: 2 }}
          >
            Confirm
          </Button>
        </Grid>
      </DialogContent>
    </Dialog>
  )
}
