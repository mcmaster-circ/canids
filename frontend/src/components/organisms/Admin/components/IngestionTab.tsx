import { useCallback, useEffect, useMemo, useState } from 'react'
import { Loader, RowActionsMenu } from '@atoms'
import { useRequest } from '@hooks'
import { Check, Delete, Redo } from '@mui/icons-material'
import { Box, Typography } from '@mui/material'
import { DataGrid, GridRenderCellParams } from '@mui/x-data-grid'
import {
  ingestionApprove,
  ingestionDelete,
  ingestionList,
} from '@api/ingestion'
import {
  defaultDeleteModalState,
  defaultKeyModalState,
  defaultRenameModalState,
  ingestionColumns,
} from '../constants'
import { DeleteModal } from '@modals'
import KeyModal from 'src/components/modals/KeyModal'
import { ApproveClientProps } from '@constants/types/ingestionPropsTypes'
import RenameModal from 'src/components/modals/RenameModal'
import RenameIngestionForm from 'src/components/forms/RenameIngestionForm'

export default () => {
  const [renameModal, setRenameModal] = useState(defaultRenameModalState)
  const [deleteModal, setDeleteModal] = useState(defaultDeleteModalState)
  const [keyModal, setKeyModal] = useState(defaultKeyModalState)
  const { data, loading, makeRequest } = useRequest({
    request: ingestionList,
  })

  useEffect(() => {
    const interval = setInterval(() => makeRequest(), 10000)
    return () => clearInterval(interval)
  }, [makeRequest])

  const { makeRequest: deleteRequest } = useRequest({
    request: ingestionDelete,
    requestByDefault: false,
    needSuccess: 'Successfully deleted',
  })
  const { makeRequest: approveRequest } = useRequest({
    request: ingestionApprove,
    requestByDefault: false,
    needSuccess: 'Successfully approved',
  })

  const handleCloseRename = useCallback(() => {
    setRenameModal(defaultRenameModalState)
    setTimeout(() => makeRequest(), 3000)
  }, [makeRequest])

  const handleCloseDelete = useCallback(() => {
    setDeleteModal(defaultDeleteModalState)
    setTimeout(() => makeRequest(), 3000)
  }, [makeRequest])

  const handleCloseKey = useCallback(() => {
    setKeyModal(defaultKeyModalState)
    setTimeout(() => makeRequest(), 3000)
  }, [makeRequest])

  console.log(data)

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
                    label: row.uuid,
                    params: { uuid: id },
                  }),
                key: 'delete',
              },
              {
                label: 'Approve',
                icon: <Check />,
                action: () => {
                  var props: ApproveClientProps = {
                    uuid: row.uuid,
                  }
                  approveRequest(props)
                  setTimeout(() => makeRequest(), 3000)
                },
                key: 'approve',
              },
              {
                label: 'Rename',
                icon: <Redo />,
                action: () => {
                  setRenameModal({
                    isUpdate: true,
                    values: row,
                    open: true,
                  })
                },
                key: 'rename',
              },
            ]}
          />,
        ]
      }),
    [approveRequest, makeRequest]
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
      <RenameModal
        open={renameModal.open}
        title="Ingestion Client"
        handleClose={handleCloseRename}
      >
        <RenameIngestionForm
          values={renameModal.values}
          handleClose={handleCloseRename}
        />
      </RenameModal>
      <DeleteModal
        open={deleteModal}
        title="Ingestion Client"
        request={deleteRequest}
        handleClose={handleCloseDelete}
      />
      <KeyModal open={keyModal} handleClose={handleCloseKey} />
    </>
  )
}
