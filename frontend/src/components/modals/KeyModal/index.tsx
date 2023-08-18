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

interface ModalTypes {
  open: { open: boolean; key?: string; title?: string; params?: any }
  handleClose: () => void
}

export default ({ open, handleClose }: ModalTypes) => {
  return (
    <Dialog open={open.open} onClose={handleClose}>
      <DialogTitle
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        Created Ingestion Client{' ' + open.title}
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
              The encryption key for the client you just created is:
              <Typography component="span" fontWeight={700}>
                {open.key}
              </Typography>
              . Copy this key now! Once you close this modal, you cannot see it
              again.
            </Typography>
          </Grid>
          <Button
            variant="contained"
            color="error"
            onClick={handleClose}
            sx={{ mt: 2 }}
          >
            Close
          </Button>
        </Grid>
      </DialogContent>
    </Dialog>
  )
}
