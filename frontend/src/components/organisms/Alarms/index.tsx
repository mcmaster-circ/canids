import { useCallback, useEffect, useState } from 'react'
import { Divider, Typography, Paper, Box } from '@mui/material'
import { subMinutes } from 'date-fns'
import Grid from '@mui/material/Unstable_Grid2'
import { useRequest } from '@hooks'
import { TimeRangePicker } from '@molecules'
import { Loader } from '@atoms'
import { getBlacklist } from '@api/blacklist'
import { getAlarms } from '@api/alarms'
import { getFields } from '@api/fields'
import { FilterSelect, AlarmsTable } from './components'

interface ChartRequestProps {
  st?: Date
  en?: Date
  p?: number
  rpp?: number
  l?: string[]
  b?: string[]
}

export default () => {
  const [start, setStart] = useState(subMinutes(new Date(), 30))
  const [end, setEnd] = useState(new Date())
  const [page, setPage] = useState(0)
  const [rowsPerPage, setRowsPerPage] = useState(5)
  const [logs, setLogs] = useState<string[] | undefined>()
  const [black, setBlack] = useState<string[] | undefined>()

  const { data: blacklist, loading: loadingBlacklist } = useRequest({
    request: getBlacklist,
    params: { getNames: true },
  })
  const { data: logsList, loading: loadingLogsList } = useRequest({
    request: getFields,
  })
  const {
    data: alarms,
    loading: loadingAlarms,
    makeRequest: requestAlarms,
  } = useRequest({
    request: getAlarms,
    requestByDefault: false,
  })
  const handleRequest = useCallback(
    async ({ st, en, p, rpp, l, b }: ChartRequestProps = {}) => {
      const s = st || start
      const e = en || end
      const pg = p || page
      const rperp = rpp || rowsPerPage
      return await requestAlarms({
        start: s.toISOString(),
        end: e.toISOString(),
        maxSize: rperp,
        from: pg * rperp,
        index: l || logs,
        source: b || black,
      })
    },
    [black, end, logs, page, requestAlarms, rowsPerPage, start]
  )

  const handleChange = useCallback(
    (v: string[], type?: 'logs') => {
      if (type === 'logs') {
        setLogs(v)
        setPage(0)
        handleRequest({ l: v, p: 0 })
      } else {
        setBlack(v)
        setPage(0)
        handleRequest({ b: v, p: 0 })
      }
    },
    [handleRequest]
  )

  useEffect(() => {
    if (logsList?.length && blacklist?.length && !logs && !black) {
      setLogs(logsList)
      setBlack(blacklist)
      handleRequest({ l: logsList, b: blacklist })
    }
  }, [
    black,
    blacklist,
    blacklist?.length,
    handleRequest,
    logs,
    logsList,
    logsList?.length,
  ])

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
            Alarms
          </Typography>
          <Box
            sx={{
              display: 'flex',
              flexWrap: 'wrap',
              gap: 2,
              flexGrow: 1,
              justifyContent: { xs: 'space-between', lg: 'flex-end' },
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
            {logsList && logs && (
              <FilterSelect
                title="Logs"
                list={logsList}
                value={logs}
                handleChange={(e: any) => handleChange(e.target.value, 'logs')}
              />
            )}
            {blacklist && black && (
              <FilterSelect
                title="Alarm Source Lists"
                list={blacklist}
                value={black}
                handleChange={(e: any) => handleChange(e.target.value)}
              />
            )}
          </Box>
        </Grid>
        <Divider sx={{ pt: 2, borderColor: '#000' }} />
      </Grid>
      <Grid xs={12} p={0} pt={3}>
        <Paper
          elevation={3}
          sx={{
            p: 2,
            borderRadius: 2,
          }}
        >
          <AlarmsTable
            setPage={setPage}
            setRowsPerPage={setRowsPerPage}
            page={page}
            rowsPerPage={rowsPerPage}
            rows={alarms?.alarms}
            count={alarms?.availableRows}
            handleRequest={handleRequest}
          />
        </Paper>
      </Grid>
      {(loadingBlacklist || loadingLogsList || loadingAlarms) && <Loader />}
    </Grid>
  )
}
