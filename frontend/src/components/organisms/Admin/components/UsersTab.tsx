import { Loader } from '@atoms'
import { useRequest } from '@hooks'
import { Box, Button, Typography } from '@mui/material'
import { DataGrid } from '@mui/x-data-grid'
import { userList } from '@api/user'
import { userColumns } from '../constants'

export default () => {
  const { data, loading } = useRequest({
    request: userList,
  })

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
        <Button variant="contained">Create User</Button>
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
            columns={userColumns}
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
    </>
  )
}
