import { Button, Typography } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { SectionProps, defaultFormValues, typeButtons } from '../constants'

export default ({ data, setData, initialData }: SectionProps) => {
  return (
    <>
      <Typography variant="h6" textAlign="center" mb={4}>
        Select Service
      </Typography>
      <Grid container justifyContent="center" rowSpacing={2} columnSpacing={4}>
        {typeButtons.map((b) => (
          <Grid key={b.key}>
            <Button
              size="large"
              variant={data.service === b.key ? 'contained' : 'outlined'}
              startIcon={b.icon}
              onClick={() => {
                console.log(b)
                console.log(initialData)

                if (b.key === initialData.service.value) {
                  defaultFormValues.accessURL = initialData.accessURL.value
                  defaultFormValues.apiKey = initialData.apiKey.value
                  defaultFormValues.domain = initialData.domain.value
                  defaultFormValues.fromAddress = initialData.fromAddress.value
                  defaultFormValues.fromName = initialData.fromName.value
                  defaultFormValues.service = initialData.service.value
                  defaultFormValues.url = initialData.url.value
                } else {
                  defaultFormValues.accessURL = ''
                  defaultFormValues.apiKey = ''
                  defaultFormValues.domain = ''
                  defaultFormValues.fromAddress = ''
                  defaultFormValues.fromName = ''
                  defaultFormValues.service = ''
                  defaultFormValues.url = ''
                }
                setData({
                  ...defaultFormValues,
                  service: b.key,
                })
              }}
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
