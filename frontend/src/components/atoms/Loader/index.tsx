import { Backdrop, CircularProgress } from '@mui/material'

export default () => (
  <Backdrop
    sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
    open
  >
    <CircularProgress color="inherit" />
  </Backdrop>
)
