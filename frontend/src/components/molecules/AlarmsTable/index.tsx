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
import { columnsConfig } from './constants'

interface TableChartProps {
  setPage: (a: number) => void
  page: number
  setRowsPerPage: (a: number) => void
  rowsPerPage: number
  rows: any[]
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
      <TableContainer sx={{ maxHeight: 440 }}>
        <Table stickyHeader>
          <TableHead>
            <TableRow>
              {columnsConfig?.map((column) => (
                <TableCell
                  key={column.id}
                  align={column.align as TableCellProps['align']}
                  style={{ fontWeight: 700, fontSize: '16px' }}
                >
                  {column.label}
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {rows.map((row, i) => {
              return (
                <TableRow hover role="checkbox" tabIndex={-1} key={i}>
                  {columnsConfig?.map((column) => {
                    const value = row[column.id]
                    return (
                      <TableCell key={column.id} align={column.align}>
                        {column.format && typeof value === 'string'
                          ? column.format(value)
                          : value}
                      </TableCell>
                    )
                  })}
                </TableRow>
              )
            })}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        rowsPerPageOptions={[5, 10, 25]}
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
