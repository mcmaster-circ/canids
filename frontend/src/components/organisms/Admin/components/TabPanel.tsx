import { ReactNode } from 'react'
import { Box } from '@mui/material'

interface TabPanelProps {
  children?: ReactNode
  index: number
  value: number
}

export default (props: TabPanelProps) => {
  const { children, value, index, ...other } = props

  return (
    <Box
      role="tabpanel"
      hidden={value !== index}
      {...other}
      sx={{ width: '100%' }}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </Box>
  )
}
