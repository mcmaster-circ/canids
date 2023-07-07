import { Paper } from '@mui/material'

export default (p: any) => {
  return <Paper>{p.class || 'chart'}</Paper>
}
