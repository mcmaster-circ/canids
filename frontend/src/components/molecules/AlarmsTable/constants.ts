import { TableCellProps } from '@mui/material'
import { format } from 'date-fns'

export const columnsConfig: {
  id: string
  label: string
  align: TableCellProps['align']
  // eslint-disable-next-line no-unused-vars
  format?: (v: string) => string
}[] = [
  {
    id: 'uid',
    label: 'UID',
    align: 'left',
  },
  {
    id: 'host',
    label: 'Host',
    align: 'center',
  },
  {
    id: 'timestamp',
    label: 'Time',
    align: 'right',
    format: (v: string) => format(new Date(v), 'MM/dd/Y HH:mm'),
  },
]
