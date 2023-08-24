import { useCallback, useMemo, useState } from 'react'
import { Loader, RowActionsMenu } from '@atoms'
import { useRequest } from '@hooks'
import { Delete, Edit, LockReset } from '@mui/icons-material'
import { Box, Button, Typography } from '@mui/material'
import { DataGrid, GridRenderCellParams } from '@mui/x-data-grid'
import { deleteUser, resetUserPass, userList } from '@api/user'
import {
  defaultAddModalState,
  defaultDeleteModalState,
  defaultEditModalState,
  userColumns,
} from '../constants'
import { AddEditModal, DeleteModal, EditModal } from '@modals'
import { AddUserForm, EditUserForm } from '@forms'

export default () => {
  const [addModal, setAddModal] = useState(defaultAddModalState)
  const [editModal, setEditModal] = useState(defaultEditModalState)
  const [deleteModal, setDeleteModal] = useState(defaultDeleteModalState)
  const { data, loading, makeRequest } = useRequest({
    request: userList,
  })
  const { makeRequest: resetRequest } = useRequest({
    request: resetUserPass,
    requestByDefault: false,
    needSuccess: 'A password reset has been successfully issued for the user',
  })
  const { makeRequest: deleteRequest } = useRequest({
    request: deleteUser,
    requestByDefault: false,
    needSuccess: 'The user account has been successfully deleted',
  })

  const handleCloseAdd = useCallback(() => {
    setAddModal(defaultAddModalState)
    setTimeout(() => makeRequest(), 3000)
  }, [makeRequest])

  const handleCloseEdit = useCallback(() => {
    setEditModal(defaultEditModalState)
    setTimeout(() => makeRequest(), 3000)
  }, [makeRequest])

  const handleCloseDelete = useCallback(() => {
    setDeleteModal(defaultDeleteModalState)
    setTimeout(() => makeRequest(), 3000)
  }, [makeRequest])

  const columns = useMemo(
    () =>
      userColumns(({ row, id }: GridRenderCellParams) => {
        return [
          <RowActionsMenu
            key={id}
            actions={[
              {
                label: 'Edit',
                icon: <Edit />,
                action: () => setEditModal({ open: true, values: row }),
                key: 'edit',
              },
              {
                label: 'Reset Password',
                icon: <LockReset />,
                action: () => resetRequest({ uuid: row.uuid }),
                key: 'reset',
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
    [resetRequest]
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
          Users
        </Typography>
        <Button
          variant="contained"
          onClick={() => setAddModal((s) => ({ ...s, open: true }))}
        >
          Create User
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
        title="User"
        handleClose={handleCloseAdd}
      >
        <AddUserForm
          isUpdate={false}
          values={addModal.values}
          handleClose={handleCloseAdd}
        />
      </AddEditModal>
      <EditModal
        open={editModal.open}
        title="User"
        handleClose={handleCloseEdit}
      >
        <EditUserForm values={editModal.values} handleClose={handleCloseEdit} />
      </EditModal>
      <DeleteModal
        open={deleteModal}
        title="User"
        request={deleteRequest}
        handleClose={handleCloseDelete}
      />
    </>
  )
}
