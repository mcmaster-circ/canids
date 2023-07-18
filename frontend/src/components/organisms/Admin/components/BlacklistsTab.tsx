import { useCallback, useMemo, useState } from 'react'
import { Box, Button, Typography } from '@mui/material'
import { DataGrid, GridRenderCellParams } from '@mui/x-data-grid'
import { deleteBlacklist, getBlacklist } from '@api/blacklist'
import { useRequest } from '@hooks'
import { AddBlacklistForm } from '@forms'
import { Loader, RowActionsMenu } from '@atoms'
import { AddEditModal, DeleteModal } from '@modals'
import { Delete, Edit } from '@mui/icons-material'
import {
  blacklistColumns,
  defaultAddModalState,
  defaultDeleteModalState,
} from '../constants'

export default () => {
  const [addModal, setAddModal] = useState(defaultAddModalState)
  const [deleteModal, setDeleteModal] = useState(defaultDeleteModalState)
  const { data, loading, makeRequest } = useRequest({
    request: getBlacklist,
  })
  const { makeRequest: deleteRequest } = useRequest({
    request: deleteBlacklist,
    requestByDefault: false,
    needSuccess: 'Successfully deleted blacklist',
  })

  const handleCloseAdd = useCallback(() => {
    setAddModal(defaultAddModalState)
    setTimeout(() => makeRequest(), 1500)
  }, [makeRequest])

  const handleCloseDelete = useCallback(() => {
    setDeleteModal(defaultDeleteModalState)
    setTimeout(() => makeRequest(), 1500)
  }, [makeRequest])

  const columns = useMemo(
    () =>
      blacklistColumns(({ row, id }: GridRenderCellParams) => {
        return [
          <RowActionsMenu
            key={id}
            actions={[
              {
                label: 'Edit',
                icon: <Edit />,
                action: () =>
                  setAddModal({ open: true, isUpdate: true, values: row }),
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
          Blacklists
        </Typography>
        <Button
          variant="contained"
          onClick={() => setAddModal((s) => ({ ...s, open: true }))}
        >
          Add Blacklist
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
        title="Blacklist"
        handleClose={handleCloseAdd}
      >
        <AddBlacklistForm
          isUpdate={addModal.isUpdate}
          values={addModal.values}
          handleClose={handleCloseAdd}
        />
      </AddEditModal>
      <DeleteModal
        open={deleteModal}
        title="Blacklist"
        request={deleteRequest}
        handleClose={handleCloseDelete}
      />
    </>
  )
}
