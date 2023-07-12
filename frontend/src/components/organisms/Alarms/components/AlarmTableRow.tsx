import { useState } from 'react'
import { format } from 'date-fns'
import { KeyboardArrowDown, KeyboardArrowUp } from '@mui/icons-material'
import Grid from '@mui/material/Unstable_Grid2'
import {
  Collapse,
  IconButton,
  TableCell,
  TableRow,
  Typography,
} from '@mui/material'
import { AlarmProps } from '@constants/types'

export default ({ row }: { row: AlarmProps }) => {
  const [open, setOpen] = useState(false)

  return (
    <>
      <TableRow>
        <TableCell>
          <IconButton
            aria-label="expand row"
            size="small"
            onClick={() => setOpen(!open)}
          >
            {open ? <KeyboardArrowUp /> : <KeyboardArrowDown />}
          </IconButton>
        </TableCell>
        <TableCell align="left">{row.uid}</TableCell>
        <TableCell align="left">{row.host}</TableCell>
        <TableCell align="left">
          {row.id_orig_h + ':' + row.id_orig_p}
        </TableCell>
        <TableCell align="left">
          {row.id_resp_h + ':' + row.id_resp_p}
        </TableCell>
        <TableCell align="right">
          {format(new Date(row.timestamp), 'MM/dd/Y HH:mm')}
        </TableCell>
      </TableRow>
      <TableRow>
        <TableCell
          style={{
            padding: 0,
            ...(!open ? { borderBottom: 'unset' } : {}),
          }}
          colSpan={6}
        >
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Grid container spacing={2} p={2}>
              <Grid p={2}>
                <Typography variant={'body1'} fontWeight={700}>
                  Source Alarms
                </Typography>
                {row.id_orig_h_pos.map((v, i) => (
                  <Typography key={i} variant={'body1'}>
                    {v}
                  </Typography>
                ))}
              </Grid>
              <Grid p={2} pl={4}>
                <Typography variant={'body1'} fontWeight={700}>
                  Dest Alarms
                </Typography>
                {row.id_resp_h_pos.map((v, i) => (
                  <Typography key={i} variant={'body1'}>
                    {v}
                  </Typography>
                ))}
              </Grid>
            </Grid>
          </Collapse>
        </TableCell>
      </TableRow>
    </>
  )
}
