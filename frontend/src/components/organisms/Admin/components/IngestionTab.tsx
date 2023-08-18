import { useCallback, useMemo, useState } from 'react'
import { Loader, RowActionsMenu } from '@atoms'
import { useRequest } from '@hooks'
import { Delete } from '@mui/icons-material'
import { Box, Button, Typography } from '@mui/material'
import { DataGrid, GridRenderCellParams } from '@mui/x-data-grid'
import { ingestionDelete, ingestionList } from '@api/ingestion'
import {
  defaultAddModalState,
  defaultDeleteModalState,
  ingestionColumns,
} from '../constants'
import { AddEditModal, DeleteModal } from '@modals'
import AddIngestionForm from 'src/components/forms/AddIngestionForm'

export default () => {
  const [addModal, setAddModal] = useState(defaultAddModalState)
  const [deleteModal, setDeleteModal] = useState(defaultDeleteModalState)
  const { data, loading, makeRequest } = useRequest({
    request: ingestionList,
  })

  const { makeRequest: deleteRequest } = useRequest({
    request: ingestionDelete,
    requestByDefault: false,
    needSuccess: 'The user account has been successfully deleted',
  })

  const handleCloseAdd = useCallback(() => {
    setAddModal(defaultAddModalState)
    setTimeout(() => makeRequest(), 3000)
  }, [makeRequest])

  const handleCloseDelete = useCallback(() => {
    setDeleteModal(defaultDeleteModalState)
    setTimeout(() => makeRequest(), 3000)
  }, [makeRequest])

  const columns = useMemo(
    () =>
      ingestionColumns(({ row, id }: GridRenderCellParams) => {
        return [
          <RowActionsMenu
            key={id}
            actions={[
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
    []
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
          Ingestion Clients
        </Typography>
        <Button
          variant="contained"
          onClick={() => setAddModal((s) => ({ ...s, open: true }))}
        >
          Create Client
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
      {loading && <Loader />}
      <AddEditModal
        open={addModal.open}
        title="Ingestion Client"
        handleClose={handleCloseAdd}
      >
        <AddIngestionForm
          isUpdate={addModal.isUpdate}
          values={addModal.values}
          handleClose={handleCloseAdd}
        />
      </AddEditModal>
      <DeleteModal
        open={deleteModal}
        title="Ingestion Client"
        request={deleteRequest}
        handleClose={handleCloseDelete}
      />
    </>
  )
}
