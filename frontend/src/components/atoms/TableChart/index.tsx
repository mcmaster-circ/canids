import { useMemo, ChangeEvent } from 'react'
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

type ColumnProps = {
  id: string
  label: string
  align: TableCellProps['align']
}

interface TableChartProps {
  setPage: (a: number) => void
  page: number
  setRowsPerPage: (a: number) => void
  rowsPerPage: number
  fieldNames: string[]
  rows: any[]
  count: number
  handleRequest: ({ p, rpp }: { p?: number; rpp?: number }) => void
}

export default ({
  setPage,
  page,
  setRowsPerPage,
  rowsPerPage,
  fieldNames,
  rows,
  count,
  handleRequest,
}: TableChartProps) => {
  const columns: ColumnProps[] | undefined = useMemo(
    () =>
      fieldNames?.length
        ? fieldNames.map((f: string, i: number) => ({
            id: 'c' + i,
            label: f,
            align: 'left',
          }))
        : undefined,
    [fieldNames]
  )

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
      <TableContainer sx={{ maxHeight: 'calc(100% - 80px)' }}>
        <Table stickyHeader size="small">
          <TableHead>
            <TableRow>
              {columns?.map((column) => (
                <TableCell
                  key={column.id}
                  align={column.align}
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
                  {columns?.map((column) => {
                    const value = row[column.id]
                    return (
                      <TableCell key={column.id} align={column.align}>
                        {typeof value === 'boolean' ? value.toString() : value}
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
        rowsPerPageOptions={[10, 25, 50, 100]}
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
