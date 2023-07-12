import { TableCellProps } from '@mui/material'

export const columnsConfig: {
  id: string
  label: string
  align: TableCellProps['align']
}[] = [
  {
    id: 'uid',
    label: 'UID',
    align: 'left',
  },
  {
    id: 'host',
    label: 'Host',
    align: 'left',
  },
  {
    id: 'id_orig_h',
    label: 'Source IP',
    align: 'left',
  },
  {
    id: 'id_dest_h',
    label: 'Dest IP',
    align: 'left',
  },
  {
    id: 'timestamp',
    label: 'Time',
    align: 'right',
  },
]
