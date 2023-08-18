import {
  Check,
  Clear,
  Computer,
  PeopleAltRounded,
  PlaylistRemoveRounded,
  QueryStatsRounded,
  Settings,
} from '@mui/icons-material'
import UsersTab from './components/UsersTab'
import VisualizationsTab from './components/VisualizationsTab'
import BlacklistsTab from './components/BlacklistsTab'
import ConfigurationTab from './components/ConfigurationTab'
import IngestionTab from './components/IngestionTab'
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
  {
    icon: <Settings />,
    iconPosition: 'start' as 'start',
    label: 'Configuration',
    sx: tabSx,
  },
  {
    icon: <Computer />,
    iconPosition: 'start' as 'start',
    label: 'Ingestion',
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
  {
    c: <ConfigurationTab />,
  },
  {
    c: <IngestionTab />,
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

export const visualizationColumns = (
  actions: (props: any) => JSX.Element[]
) => [
  { field: 'name', headerName: 'Name', flex: 0.49 },
  {
    field: 'class',
    headerName: 'Class',
    flex: 0.49,
  },
  {
    field: 'actions',
    type: 'actions',
    flex: 0.02,
    getActions: actions,
  },
]

export const userColumns = (actions: (props: any) => JSX.Element[]) => [
  { field: 'name', headerName: 'Name', flex: 0.27 },
  {
    field: 'uuid',
    headerName: 'Email',
    flex: 0.27,
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
  {
    field: 'actions',
    type: 'actions',
    flex: 0.02,
    getActions: actions,
  },
]

export const ingestionColumns = (actions: (props: any) => JSX.Element[]) => [
  { field: 'uuid', headerName: 'Name', flex: 0.27 },
  {
    field: 'actions',
    type: 'actions',
    flex: 0.02,
    getActions: actions,
  },
]

export const defaultAddModalState: {
  isUpdate: boolean
  open: boolean
  values: any
} = {
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
export const defaultSaveModalState: {
  open: boolean
} = {
  open: false,
}
export const defaultSetupMailModalState: {
  open: boolean
  values: any
} = {
  open: false,
  values: undefined,
}
export const defaultKeyModalState: {
  open: boolean
  key?: string
  title?: string
  params?: any
} = {
  open: false,
  key: undefined,
  title: undefined,
  params: undefined,
}
