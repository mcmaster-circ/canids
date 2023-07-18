import { Button, Typography } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { SectionProps, defaultFormValues, typeButtons } from '../constants'

export default ({ data, setData }: SectionProps) => {
  return (
    <>
      <Typography variant="h6" textAlign="center" mb={4}>
        Select Graph Type
      </Typography>
      <Grid container justifyContent="center" spacing={6}>
        {typeButtons.map((b) => (
          <Grid key={b.key}>
            <Button
              size="large"
              variant={data.class === b.key ? 'contained' : 'outlined'}
              startIcon={b.icon}
              onClick={() =>
                setData({
                  ...defaultFormValues,
                  class: b.key,
                })
              }
              sx={{ width: 200 }}
            >
              {b.label}
            </Button>
          </Grid>
        ))}
      </Grid>
    </>
  )
}
