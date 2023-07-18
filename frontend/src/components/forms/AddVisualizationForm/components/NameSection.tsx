import { TextField, Typography } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { FormProps, SectionProps } from '../constants'

export default ({ data, setData }: SectionProps) => {
  return (
    <>
      <Typography variant="h6" textAlign="center" mb={4}>
        Enter Visualization Name
      </Typography>
      <Grid container justifyContent="center" spacing={6}>
        <Grid xs={12}>
          <TextField
            label="Name"
            value={data.name}
            onChange={(e) =>
              setData((d: FormProps) => ({
                ...d,
                name: e.target.value,
              }))
            }
            fullWidth
          />
        </Grid>
      </Grid>
    </>
  )
}
