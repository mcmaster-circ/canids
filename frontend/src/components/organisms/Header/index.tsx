import { useState, useCallback } from 'react'
import { useCookies } from 'react-cookie'
import Image from 'next/image'
import {
  AppBar,
  Box,
  Toolbar,
  IconButton,
  Typography,
  Menu,
} from '@mui/material'
import { Menu as MenuIcon } from '@mui/icons-material'
import { ROLES } from '@constants/roles'
import { dashboardLinks } from '@constants/routes'
import { allCookies as ac } from '@constants/cookies'
import { NavLink } from '@molecules'
import logo from '@images/logo.png'

function ResponsiveAppBar() {
  const [cookies] = useCookies([ac.CLASS, ac.NAME])
  const [mOpen, setMOpen] = useState<null | HTMLElement>(null)

  const handleClose = useCallback(() => setMOpen(null), [])

  return (
    <AppBar
      position="static"
      sx={{
        background: 'linear-gradient(90deg,#a9151a,#dd7b32)',
        height: 80,
      }}
    >
      <Toolbar disableGutters sx={{ height: '100%', px: 4 }}>
        <Image src={logo} alt={'Canids'} priority={true} height={56} />
        <Box
          sx={{
            ml: 2,
            flexGrow: 1,
            display: { xs: 'none', md: 'flex' },
            height: '100%',
          }}
        >
          {dashboardLinks.map((l) =>
            l.adminRequired &&
            !cookies[ac.CLASS].includes(ROLES.ADMIN) ? null : (
              <NavLink key={l.link} name={l.name} link={l.link} />
            )
          )}
          <NavLink name="Logout" link="logout" />
        </Box>
        <Typography
          sx={{ flexGrow: 0, display: { xs: 'none', md: 'flex' } }}
          component="span"
          variant="h6"
        >
          Logged in as:&nbsp;
          <Typography component="span" sx={{ fontWeight: 600 }} variant="h6">
            {cookies[ac.NAME]}
          </Typography>
        </Typography>
        <Box
          sx={{
            flexGrow: 1,
            justifyContent: 'flex-end',
            display: { xs: 'flex', md: 'none' },
          }}
        >
          <IconButton
            aria-label="account of current user"
            aria-controls="menu-appbar"
            aria-haspopup="true"
            onClick={(e: any) => setMOpen(e.target)}
            color="inherit"
          >
            <MenuIcon fontSize="large" />
          </IconButton>
          <Menu
            anchorEl={mOpen}
            anchorOrigin={{
              vertical: 'bottom',
              horizontal: 'left',
            }}
            keepMounted
            transformOrigin={{
              vertical: 'top',
              horizontal: 'right',
            }}
            open={!!mOpen}
            onClose={handleClose}
            sx={{
              display: { xs: 'block', md: 'none' },
              py: 2,
            }}
          >
            {dashboardLinks.map((l) =>
              l.adminRequired &&
              !cookies[ac.CLASS].includes(ROLES.ADMIN) ? null : (
                <NavLink
                  key={l.link}
                  name={l.name}
                  link={l.link}
                  handleClose={handleClose}
                  small
                />
              )
            )}
            {
              <NavLink
                name="Logout"
                link="logout"
                small
                handleClose={handleClose}
              />
            }
            <NavLink
              name={'Logged in as: ' + cookies[ac.NAME]}
              link="username"
              small
              handleClose={handleClose}
            />
          </Menu>
        </Box>
      </Toolbar>
    </AppBar>
  )
}
export default ResponsiveAppBar
