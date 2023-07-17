import {
  Check,
  Clear,
  PeopleAltRounded,
  PlaylistRemoveRounded,
  QueryStatsRounded,
} from '@mui/icons-material'
import UsersTab from './components/UsersTab'
import VisualizationsTab from './components/VisualizationsTab'
import BlacklistsTab from './components/BlacklistsTab'
import { GridRenderCellParams } from '@mui/x-data-grid'

const tabSx = { justifyContent: 'flex-start', pl: 3, pr: 4 }

export const tabs = [
  {
    icon: <PeopleAltRounded />,
    iconPosition: 'start' as 'start',
    label: 'View Users',
    sx: tabSx,
  },
  {
    icon: <QueryStatsRounded />,
    iconPosition: 'start' as 'start',
    label: 'View Visualizations',
    sx: tabSx,
  },
  {
    icon: <PlaylistRemoveRounded />,
    iconPosition: 'start' as 'start',
    label: 'View Blacklists',
    sx: tabSx,
  },
]

export const tabsPanels = [
  {
    c: <UsersTab />,
  },
  {
    c: <VisualizationsTab />,
  },
  {
    c: <BlacklistsTab />,
  },
]

export const blacklistColumns = (actions: (props: any) => JSX.Element[]) => [
  { field: 'name', headerName: 'Name', flex: 0.475 },
  {
    field: 'url',
    headerName: 'Url',
    flex: 0.475,
  },
  {
    field: 'actions',
    type: 'actions',
    flex: 0.05,
    getActions: actions,
  },
]
export const defaultAddModalState = {
  open: false,
  values: undefined,
  isUpdate: false,
}
export const defaultDeleteModalState: {
  open: boolean
  label?: string
  params?: any
} = {
  open: false,
  label: undefined,
  params: undefined,
}

export const visualizationColumns = [
  { field: 'name', headerName: 'Name', flex: 0.33 },
  {
    field: 'class',
    headerName: 'Class',
    flex: 0.33,
  },
]

export const userColumns = [
  { field: 'name', headerName: 'Name', flex: 0.2 },
  {
    field: 'uuid',
    headerName: 'Email',
    flex: 0.2,
  },
  {
    field: 'class',
    headerName: 'Role',
    flex: 0.2,
  },
  {
    field: 'activated',
    headerName: 'Activated',
    flex: 0.2,
    type: 'boolean',
    renderCell: ({ value }: GridRenderCellParams) =>
      value ? <Check color="success" /> : <Clear color="error" />,
  },
]
