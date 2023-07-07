import Grid from '@mui/material/Unstable_Grid2'
import { Paper, Typography } from '@mui/material'
import { GRAPH_TYPES, GRAPH_WIDTH_TYPES } from '@constants/graphTypes'
import { useCallback, useMemo } from 'react'
import { BarChart, PieChart } from '@atoms'
import { colors } from './constanst'

export default ({ name, width, data, class: type }: any) => {
  const chartData = useMemo(() => {
    switch (type) {
      case GRAPH_TYPES.BAR:
      case GRAPH_TYPES.PIE:
        return data[0]?.length
          ? data[0].map((v: string, i: number) => ({
              name: v,
              Connections: data[1][i],
              fill: colors[i],
            }))
          : [{ name: '', Connections: 0, fill: '#ffffff' }]
      default:
        return
    }
  }, [data, type])

  console.log(chartData)

  const renderData = useCallback(() => {
    switch (type) {
      case GRAPH_TYPES.BAR:
        return <BarChart chartData={chartData} />
      case GRAPH_TYPES.PIE:
        return <PieChart chartData={chartData} />
      default:
        return null
    }
  }, [chartData, type])

  return (
    <Grid xs={12} xl={width === GRAPH_WIDTH_TYPES.FULL ? 12 : 6}>
      <Paper sx={{ p: 2, borderRadius: 2, height: '35vh', minHeight: '400px' }}>
        <Typography
          variant="h5"
          fontWeight={700}
          color="gray"
          textAlign="center"
          mb={1}
        >
          {name}
        </Typography>
        {chartData && renderData()}
      </Paper>
    </Grid>
  )
}
