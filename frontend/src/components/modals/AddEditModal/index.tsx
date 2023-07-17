import { ReactNode } from 'react'
import {
  IconButton,
  Divider,
  Dialog,
  DialogContent,
  DialogTitle,
} from '@mui/material'
import { Close } from '@mui/icons-material'

interface Props {
  open: { open: boolean; isUpdate: boolean; values?: any }
  handleClose: () => void
  title: string
  children: ReactNode
}

export default ({ open, handleClose, title, children }: Props) => {
  return (
    <Dialog open={open.open} onClose={handleClose}>
      <DialogTitle
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        Add{' ' + title}
        <IconButton aria-label="close" onClick={handleClose}>
          <Close />
        </IconButton>
      </DialogTitle>
      <Divider variant="fullWidth" />
      <DialogContent>{children}</DialogContent>
    </Dialog>
  )
}
