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
  open: boolean
  handleClose: () => void
  title: string
  children: ReactNode
}

export default ({ children, open, handleClose, title }: Props) => {
  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogTitle
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        Rename{' ' + title}
        <IconButton aria-label="close" onClick={handleClose}>
          <Close />
        </IconButton>
      </DialogTitle>
      <Divider variant="fullWidth" />
      <DialogContent>{children}</DialogContent>
    </Dialog>
  )
}
