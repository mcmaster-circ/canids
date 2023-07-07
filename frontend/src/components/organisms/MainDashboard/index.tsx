import { useCallback, useEffect, useState } from 'react'
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
import { subMinutes } from 'date-fns'
import { MoreVert, SsidChart } from '@mui/icons-material'
import { DateTimePicker } from '@mui/x-date-pickers'
import Grid from '@mui/material/Unstable_Grid2'
import { getViewList } from '@api/view'
import { useRequest } from '@hooks'
import { Loader } from '@atoms'
import { ViewListItemProps } from '@constants/types'
import { GRAPH_TYPES_ICONS } from '@constants/graphTypes'
import { getChartsData } from '@api/charts'

export default () => {
  const [start, setStart] = useState(subMinutes(new Date(), 30))
  const [end, setEnd] = useState(new Date())
  const [open, setOpen] = useState(null)

  const { data: viewsList, loading: loadingList } = useRequest({
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

  const handleRequest = useCallback(async () => {
    return await requestChartData({
      views: viewsList.map((v: ViewListItemProps) => v.uuid),
      params: {
        start: start.toISOString(),
        end: end.toISOString(),
        interval: 75,
        maxSize: 10,
        from: 0,
      },
    })
  }, [end, requestChartData, start, viewsList])

  const handleClose = useCallback(() => setOpen(null), [])

  useEffect(() => {
    // const interval =
    //   viewsList?.length &&
    //   chartData?.length &&
    //   setInterval(() => handleRequest(), 10000)
    if (viewsList?.length && !chartData) {
      handleRequest()
    }
    // return () => {
    //   clearInterval(interval)
    // }
  }, [chartData, handleRequest, viewsList])

  return (
    <Grid container spacing={2} p={3} m={0}>
      <Grid xs={12} p={0} pb={3}>
        <Grid container>
          <Typography
            variant="h4"
            fontWeight={700}
            lineHeight={1.6}
            sx={{ pb: { xs: 2, md: 0 } }}
          >
            Main Dashboard
          </Typography>
          <Box
            sx={{
              display: 'flex',
              gap: 4,
              flexGrow: 1,
              justifyContent: 'flex-end',
            }}
          >
            <DateTimePicker
              label="Start Time"
              ampm={false}
              value={start}
              timeSteps={{ hours: 1, minutes: 1 }}
              onChange={(date) => setStart(date as Date)}
              onClose={handleRequest}
              maxDateTime={subMinutes(end, 30)}
              slotProps={{
                actionBar: {
                  actions: ['today', 'accept'],
                },
              }}
            />
            <DateTimePicker
              label="End Time"
              ampm={false}
              value={end}
              timeSteps={{ hours: 1, minutes: 1 }}
              maxDateTime={end}
              minDateTime={start}
              onChange={(date) => setEnd(date as Date)}
              onClose={handleRequest}
            />
          </Box>
        </Grid>
        <Divider sx={{ pt: 2, borderColor: '#000' }} />
      </Grid>
      <Grid xs={12} md={6} lg={3} p={0}>
        <Paper
          elevation={3}
          sx={{ minHeight: 'calc(100vh - 257px)', p: 2, borderRadius: 2 }}
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
                <Divider component="li" sx={{ mx: 2, pt: 1 }} />
              </div>
            ))}
          </List>
        </Paper>
      </Grid>
      {(loadingList || loadingChartData) && <Loader />}
    </Grid>
  )
}
