import {
  IconButton,
  Divider,
  Dialog,
  DialogContent,
  DialogTitle,
} from '@mui/material'
import { Close } from '@mui/icons-material'

interface ModalTypes {
  open: boolean
  handleClose: () => void
}

export default ({ open, handleClose }: ModalTypes) => {
  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogTitle
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        Dialog Title
        <IconButton aria-label="close" onClick={handleClose}>
          <Close />
        </IconButton>
      </DialogTitle>
      <Divider variant="fullWidth" />
      <DialogContent>Dialog Content</DialogContent>
    </Dialog>
  )
}
