import { useCallback, useEffect, useMemo, useState } from 'react'
import {
  Divider,
  Typography,
  Paper,
  Box,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  IconButton,
  Menu,
  MenuItem,
} from '@mui/material'
import { subMinutes, differenceInDays } from 'date-fns'
import { MoreVert, SsidChart } from '@mui/icons-material'
import Grid from '@mui/material/Unstable_Grid2'
import { getViewList } from '@api/view'
import { getChartsData } from '@api/charts'
import { getDashboard } from '@api/dashboard'
import { useRequest } from '@hooks'
import { GRAPH_TYPES_ICONS } from '@constants/graphTypes'
import { ViewListItemProps } from '@constants/types'
import { TimeRangePicker, ChartCard } from '@molecules'
import { Loader } from '@atoms'

interface ChartRequestProps {
  st?: Date
  en?: Date
  p?: number
  rpp?: number
}

export default () => {
  const [start, setStart] = useState(subMinutes(new Date(), 30))
  const [end, setEnd] = useState(new Date())
  const [page, setPage] = useState(0)
  const [rowsPerPage, setRowsPerPage] = useState(10)
  const [open, setOpen] = useState(null)

  const { data: dashboard, loading: loadingDashboard } = useRequest({
    request: getDashboard,
  })
  const { data: views, loading: loadingList } = useRequest({
    request: getViewList,
  })
  const {
    data: chartData,
    loading: loadingChartData,
    makeRequest: requestChartData,
  } = useRequest({
    request: getChartsData,
    requestByDefault: false,
  })

  const viewsList = useMemo(() => {
    if (dashboard?.views?.length && views?.length) {
      return dashboard.views.map((v: string, i: number) => ({
        size: dashboard.sizes[i],
        ...views.find((view: any) => view.uuid === v),
      }))
    }
  }, [dashboard?.sizes, dashboard?.views, views])

  const handleRequest = useCallback(
    async ({ st, en, p, rpp }: ChartRequestProps = {}) => {
      const s = st || start
      const e = en || end
      const pg = p || page
      const rperp = rpp || rowsPerPage
      return await requestChartData({
        views: viewsList,
        params: {
          start: s.toISOString(),
          end: e.toISOString(),
          interval: differenceInDays(e, s),
          maxSize: rperp,
          from: pg * rperp,
        },
      })
    },
    [end, page, requestChartData, rowsPerPage, start, viewsList]
  )

  const handleClose = useCallback(() => setOpen(null), [])

  useEffect(() => {
    const interval =
      viewsList?.length &&
      chartData?.length &&
      setInterval(() => handleRequest(), 10000)
    if (viewsList?.length && !chartData) {
      handleRequest()
    }
    return () => {
      clearInterval(interval)
    }
  }, [chartData, handleRequest, viewsList])

  return (
    <Grid container spacing={2} p={3} m={0}>
      <Grid xs={12} p={0}>
        <Grid container>
          <Typography
            variant="h4"
            fontWeight={700}
            lineHeight={1.6}
            sx={{ pb: { xs: 2, md: 0 } }}
          >
            {dashboard?.name || 'Main Dashboard'}
          </Typography>
          <Box
            sx={{
              display: 'flex',
              gap: 4,
              flexGrow: 1,
              justifyContent: 'flex-end',
              alignItems: 'center',
            }}
          >
            <TimeRangePicker
              start={start}
              end={end}
              setStart={setStart}
              setEnd={setEnd}
              handleRequest={handleRequest}
            />
          </Box>
        </Grid>
        <Divider sx={{ pt: 2, borderColor: '#000' }} />
      </Grid>
      <Grid xs={12} lg={4} xl={3} p={0} pt={3}>
        <Paper
          elevation={3}
          sx={{
            minHeight: 'calc(100vh - 257px)',
            height: 'calc(100% - 32px)',
            p: 2,
            borderRadius: 2,
          }}
        >
          <Typography variant="h6" fontWeight={600}>
            All Visualizations
          </Typography>
          <List>
            {viewsList?.map((v: ViewListItemProps) => (
              <div key={v.uuid}>
                <ListItem sx={{ p: 0 }}>
                  <ListItemIcon>
                    {GRAPH_TYPES_ICONS[v.class] || (
                      <SsidChart fontSize="large" />
                    )}
                  </ListItemIcon>
                  <ListItemText
                    primary={v.name}
                    secondary={'Type: ' + v.class}
                  />
                  <IconButton
                    aria-label="More"
                    aria-haspopup="true"
                    onClick={(e: any) => setOpen(e.target)}
                  >
                    <MoreVert />
                  </IconButton>
                  <Menu
                    id="basic-menu"
                    elevation={1}
                    anchorEl={open}
                    open={!!open}
                    onClose={handleClose}
                  >
                    <MenuItem onClick={handleClose}>Edit</MenuItem>
                    <MenuItem onClick={handleClose}>Delete</MenuItem>
                  </Menu>
                </ListItem>
                <Divider component="li" sx={{ mx: 2 }} />
              </div>
            ))}
          </List>
        </Paper>
      </Grid>
      <Grid xs={12} lg={8} xl={9} p={0}>
        <Grid container spacing={2} p={0} pt={3} pl={{ xs: 0, lg: 3 }}>
          {chartData?.map((c: any) => (
            <ChartCard
              key={c.uuid}
              {...c}
              setPage={setPage}
              page={page}
              setRowsPerPage={setRowsPerPage}
              rowsPerPage={rowsPerPage}
              handleRequest={handleRequest}
            />
          ))}
        </Grid>
      </Grid>
      {(loadingList || loadingChartData || loadingDashboard) && <Loader />}
    </Grid>
  )
}
