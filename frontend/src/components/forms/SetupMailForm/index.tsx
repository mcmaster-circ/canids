import { useCallback, useMemo, useState } from 'react'
import { Button, Step, StepLabel, Stepper } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import {
  SetupMailSettings,
  UpdateConfigurationProps,
  UpdateSetting,
} from '@constants/types'
import { useRequest } from '@hooks'
import { FormProps, steps } from './constants'
import { Check, NavigateBefore, NavigateNext } from '@mui/icons-material'
import { ConfigSection, ServiceSection } from './components'
import { updateConfiguration } from '@api/configuration'

interface SetupMailFormProps {
  handleClose: () => void
  isUpdate?: boolean
  values: SetupMailSettings
}

export default ({ handleClose, values }: SetupMailFormProps) => {
  const [activeStep, setActiveStep] = useState(0)
  const [data, setData] = useState<FormProps>({
    service: values?.service.value,
    url: values?.url.value,
    apiKey: values?.apiKey.value,
    fromAddress: values?.fromAddress.value,
    fromName: values?.fromName.value,
    domain: values?.domain.value,
    accessURL: values?.accessURL.value,
  })
  const { makeRequest: saveRequest } = useRequest({
    request: updateConfiguration,
    requestByDefault: false,
    needSuccess: 'Successfully saved configuration',
  })

  const handleSubmit = useCallback(async () => {
    const settings: UpdateSetting[] = [
      {
        name: 'MAIL_SERVICE',
        value: data.service,
      },
      {
        name: 'MAIL_URL',
        value: data.url,
      },
      {
        name: 'MAIL_API_KEY',
        value: data.apiKey,
      },
      {
        name: 'MAIL_FROM_ADDRESS',
        value: data.fromAddress,
      },
      {
        name: 'MAIL_FROM_NAME',
        value: data.fromName,
      },
      {
        name: 'MAIL_DOMAIN',
        value: data.domain,
      },
      {
        name: 'ACCESS_URL',
        value: data.accessURL,
      },
    ]
    const req: UpdateConfigurationProps = {
      configuration: settings,
    }
    await saveRequest(req)
    handleClose()
  }, [data, handleClose, saveRequest])

  const renderSection = () => {
    switch (activeStep) {
      case 0:
        return <ServiceSection data={data} setData={setData} />
      case 1:
        return <ConfigSection data={data} setData={setData} />
      default:
        return null
    }
  }

  const nextDisabled = useMemo(() => {
    switch (activeStep) {
      case 0:
        return !data.service
      case 1:
        if (data.service === 'NONE') return false
        if (!data.apiKey || !data.fromName || !data.fromAddress) return true
        if (data.service === 'MAILGUN' && !data.domain) return true
        if (
          (data.service === 'MAILGUN' ||
            data.service === 'POSTAL' ||
            data.service === 'SPARKPOST') &&
          !data.url
        )
          return true
        return false
      default:
        return false
    }
  }, [
    activeStep,
    data.service,
    data.url,
    data.apiKey,
    data.fromAddress,
    data.fromName,
    data.domain,
  ])

  return (
    <Grid container spacing={3} justifyContent="center">
      <Grid xs={12}>
        <Stepper
          activeStep={activeStep}
          alternativeLabel
          sx={{
            '& .MuiStepLabel-label': { fontSize: '16px' },
            '& .MuiStepConnector-root': { top: 18 },
            minWidth: '500px',
          }}
        >
          {steps.map((label) => (
            <Step key={label}>
              <StepLabel StepIconProps={{ sx: { width: 36, height: 36 } }}>
                {label}
              </StepLabel>
            </Step>
          ))}
        </Stepper>
      </Grid>
      <Grid xs={12} height={300} overflow="scroll">
        {renderSection()}
      </Grid>
      <Grid xs={12}>
        <Grid container justifyContent="space-between">
          {activeStep !== 0 ? (
            <Button
              variant="contained"
              onClick={() => setActiveStep((a) => a - 1)}
              startIcon={<NavigateBefore />}
            >
              Back
            </Button>
          ) : (
            <div />
          )}
          {activeStep !== 1 ? (
            <Button
              variant="contained"
              color="success"
              onClick={() => setActiveStep((a) => a + 1)}
              endIcon={<NavigateNext />}
              disabled={nextDisabled}
            >
              Next
            </Button>
          ) : (
            <Button
              variant="contained"
              onClick={handleSubmit}
              startIcon={<Check />}
              disabled={nextDisabled}
            >
              Submit
            </Button>
          )}
        </Grid>
      </Grid>
    </Grid>
  )
}
