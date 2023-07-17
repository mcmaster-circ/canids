import {
  IconButton,
  Divider,
  Dialog,
  DialogContent,
  DialogTitle,
  Typography,
  Button,
} from '@mui/material'
import { Close, Error } from '@mui/icons-material'
import { useCallback } from 'react'
import Grid from '@mui/material/Unstable_Grid2/Grid2'

interface ModalTypes {
  open: { open: boolean; label?: string; params?: any }
  handleClose: () => void
  title: string
  request: (v: any) => void
}

export default ({ open, handleClose, request, title }: ModalTypes) => {
  const handleClick = useCallback(async () => {
    await request(open.params)
    handleClose()
  }, [handleClose, open.params, request])

  return (
    <Dialog open={open.open} onClose={handleClose}>
      <DialogTitle
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        Deleting{' ' + title}
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
              Are you sure you want to delete
              <Typography component="span" fontWeight={700}>
                {` ${title} ${open.label}`}
              </Typography>
              ? This action cannot be undone.
            </Typography>
          </Grid>
          <Button
            variant="contained"
            color="error"
            onClick={handleClick}
            sx={{ mt: 2 }}
          >
            Delete{' ' + title}
          </Button>
        </Grid>
      </DialogContent>
    </Dialog>
  )
}
