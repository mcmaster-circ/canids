import { ChangeEvent } from 'react'
import {
  Table,
  TableBody,
  TableCell,
  TableCellProps,
  TableContainer,
  TableHead,
  TablePagination,
  TableRow,
} from '@mui/material'
import { columnsConfig } from '../constants'
import { AlarmTableRow } from './'
import { AlarmProps } from '@constants/types'

interface TableChartProps {
  setPage: (a: number) => void
  page: number
  setRowsPerPage: (a: number) => void
  rowsPerPage: number
  rows: AlarmProps[]
  count: number
  handleRequest: ({ p, rpp }: { p?: number; rpp?: number }) => void
}

export default ({
  setPage,
  page,
  setRowsPerPage,
  rowsPerPage,
  rows,
  count,
  handleRequest,
}: TableChartProps) => {
  const handleChangePage = (_: unknown, newPage: number) => {
    handleRequest({ p: newPage })
    setPage(newPage)
  }

  const handleChangeRowsPerPage = (event: ChangeEvent<HTMLInputElement>) => {
    handleRequest({ rpp: +event.target.value, p: 0 })
    setRowsPerPage(+event.target.value)
    setPage(0)
  }

  return (
    <>
      <TableContainer sx={{ maxHeight: 'calc(100vh - 320px)' }}>
        <Table stickyHeader>
          <TableHead>
            <TableRow>
              <TableCell />
              {columnsConfig?.map((column) => (
                <TableCell
                  key={column.label}
                  align={column.align as TableCellProps['align']}
                  style={{ fontWeight: 700, fontSize: '16px' }}
                >
                  {column.label}
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {rows?.map((row, i) => (
              <AlarmTableRow key={'alarmRow' + i} row={row} />
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        rowsPerPageOptions={[5, 10, 20]}
        component="div"
        count={count}
        rowsPerPage={rowsPerPage}
        page={page}
        onPageChange={handleChangePage}
        onRowsPerPageChange={handleChangeRowsPerPage}
        showFirstButton
        showLastButton
      />
    </>
  )
}
