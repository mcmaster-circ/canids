import { ReactNode, SyntheticEvent, useState } from 'react'
import { Divider, Tab, Tabs, Typography, Paper } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { tabs, tabsPanels } from './constants'
import { TabPanel } from './components'

export default () => {
  const [value, setValue] = useState(1)

  const handleChange = (_: SyntheticEvent, newValue: number) => {
    setValue(newValue)
  }

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
            Admin Actions
          </Typography>
        </Grid>
        <Divider sx={{ pt: 2, borderColor: '#000' }} />
      </Grid>
      <Grid xs={12} p={0} pt={3}>
        <Paper
          elevation={3}
          sx={{
            flexGrow: 1,
            display: 'flex',
            borderRadius: 2,
          }}
        >
          <Tabs
            orientation="vertical"
            variant="scrollable"
            value={value}
            onChange={handleChange}
            aria-label="Admin Actions"
            sx={{ borderRight: 1, borderColor: 'divider', minWidth: '236px' }}
          >
            {tabs.map((t) => (
              <Tab key={t.label} {...t} />
            ))}
          </Tabs>
          {tabsPanels.map(({ c }: { c: ReactNode }, i: number) => (
            <TabPanel key={'tab-panel' + i} value={value} index={i}>
              {c}
            </TabPanel>
          ))}
        </Paper>
      </Grid>
    </Grid>
  )
}
