import { useCallback, useMemo } from 'react'
import { useRouter } from 'next/router'
import { Typography, Button, MenuItem } from '@mui/material'
import { Logout } from '@mui/icons-material'
import useAuth from '@context/authContext'
import { dashboardRoutes } from '@constants/routes'

interface NavLinkProps {
  name: string
  link: string
  small?: boolean
  handleClose?: () => void
}

export default ({ name, link, small, handleClose }: NavLinkProps) => {
  const { logout } = useAuth()
  const { push, query } = useRouter()

  const handlePressLink = useCallback(
    (link: string) => {
      if (link === 'logout') {
        return logout()
      }
      push(link)
      small && handleClose && handleClose()
    },
    [handleClose, logout, push, small]
  )

  const active = useMemo(() => {
    if (query?.dashboardRoute) {
      return link.includes(query.dashboardRoute as string)
    }
    return link === dashboardRoutes.DASHBOARD
  }, [link, query?.dashboardRoute])

  return small ? (
    <MenuItem
      disabled={link === 'username'}
      onClick={() => handlePressLink(link)}
      sx={{
        mx: 2,
        display: 'flex',
        alignItems: 'center',
        gap: '4px',
      }}
    >
      <Typography
        fontWeight={active ? 800 : 600}
        textAlign="center"
        textTransform={link != 'username' ? 'uppercase' : undefined}
        variant="h6"
      >
        {name}
      </Typography>
      {link === 'logout' && <Logout />}
    </MenuItem>
  ) : (
    <Button
      onClick={() => handlePressLink(link)}
      sx={{
        m: 2,
        px: 0,
        fontSize: '16px',
        gap: '4px',
        display: 'flex',
        alignItems: 'center',
        borderRadius: 0,
        color: 'white',
        borderBottom: active ? '2px solid white' : '2px solid transparent',
        ':hover': { borderBottom: '2px solid white', bgcolor: 'transparent' },
      }}
    >
      {link === 'logout' && <Logout />}
      {name}
    </Button>
  )
}
