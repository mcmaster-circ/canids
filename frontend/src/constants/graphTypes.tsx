import { ShowChart, BarChart, PieChart, TableChart } from '@mui/icons-material'

export const GRAPH_TYPES = {
  LINE: 'line',
  BAR: 'bar',
  PIE: 'pie',
  TABLE: 'table',
}

export const GRAPH_TYPES_ICONS = {
  [GRAPH_TYPES.LINE]: <ShowChart fontSize="large" />,
  [GRAPH_TYPES.BAR]: <BarChart fontSize="large" />,
  [GRAPH_TYPES.PIE]: <PieChart fontSize="large" />,
  [GRAPH_TYPES.TABLE]: <TableChart fontSize="large" />,
}
