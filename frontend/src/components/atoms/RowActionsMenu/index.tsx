import { useCallback, useState } from 'react'
import { MoreVert } from '@mui/icons-material'
import { IconButton, Menu, MenuItem } from '@mui/material'
import { RowActionProps } from '@constants/types'

interface MenuProps {
  actions: RowActionProps[]
}

export default ({ actions }: MenuProps) => {
  const [open, setOpen] = useState(null)

  const handleClose = useCallback(() => setOpen(null), [])
  const handleClick = useCallback(
    (action: (v?: any) => void) => {
      action()
      handleClose()
    },
    [handleClose]
  )

  return (
    <>
      <IconButton
        aria-label="Acations"
        aria-haspopup="true"
        onClick={(e: any) => setOpen(e.target)}
      >
        <MoreVert />
      </IconButton>
      <Menu
        id="basic-menu"
        elevation={1}
        anchorEl={open}
        open={!!open}
        onClose={handleClose}
      >
        {actions.map(({ label, key, action, icon }) => (
          <MenuItem
            key={key}
            onClick={() => handleClick(action)}
            sx={{ gap: 1 }}
          >
            {icon && icon}
            {label}
          </MenuItem>
        ))}
      </Menu>
    </>
  )
}
