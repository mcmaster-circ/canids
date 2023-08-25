import { TextField, Typography } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { FormProps, SectionProps } from '../constants'

export default ({ data, setData, initialData }: SectionProps) => {
  return (
    <>
      <Typography variant="h6" textAlign="center" mb={4}>
        Configuration
      </Typography>
      <Grid container justifyContent="center" rowSpacing={2}>
        {['NONE'].includes(data.service) && (
          <Grid xs={12}>
            <Typography variant="body1" textAlign="center">
              No email service selected
            </Typography>
          </Grid>
        )}
        {['MAILGUN', 'POSTAL', 'SPARKPOST'].includes(data.service) && (
          <Grid xs={12}>
            <TextField
              label="URL"
              value={data.url}
              size="small"
              type="url"
              onChange={(e) =>
                setData((d: FormProps) => ({
                  ...d,
                  url: e.target.value,
                }))
              }
              fullWidth
            />
          </Grid>
        )}
        {!['NONE'].includes(data.service) && (
          <>
            <Grid xs={12}>
              <TextField
                label="API Key"
                value={data.apiKey}
                size="small"
                onChange={(e) =>
                  setData((d: FormProps) => ({
                    ...d,
                    apiKey: e.target.value,
                  }))
                }
                fullWidth
              />
            </Grid>
            <Grid xs={12}>
              <TextField
                label="From Address"
                value={data.fromAddress}
                size="small"
                onChange={(e) =>
                  setData((d: FormProps) => ({
                    ...d,
                    fromAddress: e.target.value,
                  }))
                }
                fullWidth
              />
            </Grid>
            <Grid xs={12}>
              <TextField
                label="From Name"
                value={data.fromName}
                size="small"
                onChange={(e) =>
                  setData((d: FormProps) => ({
                    ...d,
                    fromName: e.target.value,
                  }))
                }
                fullWidth
              />
            </Grid>
            <Grid xs={12}>
              <TextField
                label="CanIDS Access URL"
                value={data.accessURL}
                size="small"
                onChange={(e) =>
                  setData((d: FormProps) => ({
                    ...d,
                    accessURL: e.target.value,
                  }))
                }
                fullWidth
              />
            </Grid>
          </>
        )}
        {['MAILGUN'].includes(data.service) && (
          <Grid xs={12}>
            <TextField
              label="Domain"
              value={data.domain}
              size="small"
              onChange={(e) =>
                setData((d: FormProps) => ({
                  ...d,
                  domain: e.target.value,
                }))
              }
              fullWidth
            />
          </Grid>
        )}
      </Grid>
    </>
  )
}
