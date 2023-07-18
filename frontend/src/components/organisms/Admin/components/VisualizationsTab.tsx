import { useCallback, useMemo, useState } from 'react'
import { Loader, RowActionsMenu } from '@atoms'
import { useRequest } from '@hooks'
import { Box, Button, Typography } from '@mui/material'
import { DataGrid, GridRenderCellParams } from '@mui/x-data-grid'
import { deleteView, getViewList } from '@api/view'
import {
  defaultAddModalState,
  defaultDeleteModalState,
  visualizationColumns,
} from '../constants'
import { AddEditModal, DeleteModal } from '@modals'
import { AddVisualizationForm } from '@forms'
import { Delete, Edit } from '@mui/icons-material'
import { getChartData } from '@api/charts'

export default () => {
  const [addModal, setAddModal] = useState(defaultAddModalState)
  const { data, loading, makeRequest } = useRequest({
    request: getViewList,
  })
  const { loading: chartLoading, makeRequest: chartRequest } = useRequest({
    request: getChartData,
    requestByDefault: false,
  })
  const [deleteModal, setDeleteModal] = useState(defaultDeleteModalState)
  const { makeRequest: deleteRequest } = useRequest({
    request: deleteView,
    requestByDefault: false,
    needSuccess: 'Successfully deleted view',
  })

  const handleCloseAdd = useCallback(() => {
    setAddModal(defaultAddModalState)
    setTimeout(() => makeRequest(), 1000)
  }, [makeRequest])

  const handleCloseDelete = useCallback(() => {
    setDeleteModal(defaultDeleteModalState)
    setTimeout(() => makeRequest(), 1000)
  }, [makeRequest])

  const handleRequestEdit = useCallback(
    async (row: GridRenderCellParams['row']) => {
      const res = await chartRequest({ uuid: row.uuid })
      setAddModal({ open: true, isUpdate: true, values: { ...row, ...res } })
    },
    [chartRequest]
  )

  const columns = useMemo(
    () =>
      visualizationColumns(({ row, id }: GridRenderCellParams) => {
        return [
          <RowActionsMenu
            key={id}
            actions={[
              {
                label: 'Edit',
                icon: <Edit />,
                action: () => handleRequestEdit(row),
                key: 'edit',
              },
              {
                label: 'Delete',
                icon: <Delete />,
                action: () =>
                  setDeleteModal({
                    open: true,
                    label: row.name,
                    params: { uuid: id },
                  }),
                key: 'delete',
              },
            ]}
          />,
        ]
      }),
    [handleRequestEdit]
  )

  return (
    <>
      <Box
        sx={{
          display: 'flex',
          flexWrap: 'wrap',
          gap: 2,
          justifyContent: 'space-between',
          alignItems: 'center',
          mb: 3,
        }}
      >
        <Typography variant="h6" fontWeight={700}>
          Visualizations
        </Typography>
        <Button
          variant="contained"
          onClick={() => setAddModal((s) => ({ ...s, open: true }))}
        >
          Create Visualization
        </Button>
      </Box>
      <Box
        sx={{
          height: '100%',
          width: '100%',
          display: 'grid',
          gridTemplateColumns: '1fr',
        }}
      >
        {data && (
          <DataGrid
            sx={{
              '.MuiDataGrid-menuIcon': {
                visibility: 'visible',
                width: 'auto',
                mr: 1,
              },
              '.MuiDataGrid-iconButtonContainer': {
                ml: 1,
              },
              '.MuiDataGrid-columnHeaderTitle': {
                fontWeight: 700,
                fontSize: '16px',
              },
            }}
            getRowId={(row) => row.uuid}
            rows={data}
            columns={columns}
            initialState={{
              pagination: {
                paginationModel: {
                  pageSize: 5,
                },
              },
            }}
            pageSizeOptions={[5, 10, 25, 50, 100]}
            disableRowSelectionOnClick
          />
        )}
      </Box>
      {(loading || chartLoading) && <Loader />}
      <AddEditModal
        open={addModal.open}
        title="Visualization"
        handleClose={handleCloseAdd}
      >
        <AddVisualizationForm
          isUpdate={addModal.isUpdate}
          values={addModal.values}
          handleClose={handleCloseAdd}
        />
      </AddEditModal>
      <DeleteModal
        open={deleteModal}
        title="Visualization"
        request={deleteRequest}
        handleClose={handleCloseDelete}
      />
    </>
  )
}
