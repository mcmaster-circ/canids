import { useCallback, useEffect, useState } from 'react'
import { Divider, Typography, Paper, Box } from '@mui/material'
import { subMinutes } from 'date-fns'
import Grid from '@mui/material/Unstable_Grid2'
import { useRequest } from '@hooks'
import { TimeRangePicker, FilterSelect } from '@molecules'
import { Loader } from '@atoms'
import { getBlacklist } from '@api/blacklist'
import { getAlarms } from '@api/alarms'
import { getFields } from '@api/fields'
import AlarmsTable from 'src/components/molecules/AlarmsTable'

interface ChartRequestProps {
  st?: Date
  en?: Date
  p?: number
  rpp?: number
  l?: string[]
  b?: string[]
}

const alarms = [
  {
    uid: 'CQNEhi1Wy50SyemtOb',
    host: '',
    timestamp: '2021-03-15T16:59:59Z',
    id_orig_h: '192.168.2.81',
    id_orig_p: 57404,
    id_orig_h_pos: ['firehol_level1'],
    id_resp_h: '192.168.2.88',
    id_resp_p: 53,
    id_resp_h_pos: ['firehol_level1'],
  },
  {
    uid: 'C3C7i1D7MkoA0agcj',
    host: '',
    timestamp: '2021-03-15T16:59:59Z',
    id_orig_h: '192.168.2.81',
    id_orig_p: 11112,
    id_orig_h_pos: ['firehol_level1'],
    id_resp_h: '192.168.2.88',
    id_resp_p: 53,
    id_resp_h_pos: ['firehol_level1'],
  },
  {
    uid: 'CsXR5PDwDTmgqR5Ok',
    host: '',
    timestamp: '2021-03-15T16:59:59Z',
    id_orig_h: '192.168.2.81',
    id_orig_p: 58973,
    id_orig_h_pos: ['firehol_level1'],
    id_resp_h: '192.168.2.88',
    id_resp_p: 53,
    id_resp_h_pos: ['firehol_level1'],
  },
  {
    uid: 'CgMkQs1eegT2HZYUkg',
    host: '',
    timestamp: '2021-03-15T16:59:59Z',
    id_orig_h: '192.168.2.81',
    id_orig_p: 58667,
    id_orig_h_pos: ['firehol_level1'],
    id_resp_h: '192.168.2.88',
    id_resp_p: 53,
    id_resp_h_pos: ['firehol_level1'],
  },
  {
    uid: 'CgMkQs1eegT2HZYUk1',
    host: '',
    timestamp: '2021-03-15T16:59:59Z',
    id_orig_h: '192.168.2.81',
    id_orig_p: 58667,
    id_orig_h_pos: ['firehol_level1'],
    id_resp_h: '192.168.2.88',
    id_resp_p: 53,
    id_resp_h_pos: ['firehol_level1'],
  },
]

export default () => {
  const [start, setStart] = useState(subMinutes(new Date(), 30))
  const [end, setEnd] = useState(new Date())
  const [page, setPage] = useState(0)
  const [rowsPerPage, setRowsPerPage] = useState(5)
  const [logs, setLogs] = useState<string[] | undefined>()
  const [black, setBlack] = useState<string[] | undefined>()

  const { data: blacklist, loading: loadingBlacklist } = useRequest({
    request: getBlacklist,
  })
  const { data: logsList, loading: loadingLogsList } = useRequest({
    request: getFields,
  })
  const {
    // data: alarms,
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
        handleRequest({ l: v })
      } else {
        setBlack(v)
        handleRequest({ b: v })
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

  console.log(logs)

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
            rows={alarms}
            count={5}
            handleRequest={handleRequest}
          />
        </Paper>
      </Grid>
      {(loadingBlacklist || loadingLogsList || loadingAlarms) && <Loader />}
    </Grid>
  )
}
