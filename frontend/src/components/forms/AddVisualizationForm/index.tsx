import { useCallback, useMemo, useState } from 'react'
import { Button, Step, StepLabel, Stepper } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { UpdateViewProps } from '@constants/types'
import { useRequest } from '@hooks'
import { FormProps, defaultFormValues, steps } from './constants'
import { addView, updateView } from '@api/view'
import { Check, NavigateBefore, NavigateNext } from '@mui/icons-material'
import { InfoSection, NameSection, TypesSection } from './components'

interface AddFormProps {
  handleClose: () => void
  isUpdate?: boolean
  values: UpdateViewProps
}

export default ({ handleClose, isUpdate, values }: AddFormProps) => {
  const [activeStep, setActiveStep] = useState(0)
  const [data, setData] = useState<FormProps>(
    isUpdate
      ? {
          name: values?.name,
          class: values?.class,
          index: values?.index,
          fields: values?.fields,
          fieldNames: values?.fieldNames,
        }
      : defaultFormValues
  )
  const { makeRequest } = useRequest({
    request: addView,
    requestByDefault: false,
    needSuccess: 'View successfully created',
  })
  const { makeRequest: update } = useRequest({
    request: updateView,
    requestByDefault: false,
    needSuccess: 'Successfully updated view',
  })

  const handleSubmit = useCallback(async () => {
    isUpdate
      ? await update({
          ...data,
          uuid: values?.uuid,
        })
      : await makeRequest(data)
    handleClose()
  }, [data, handleClose, isUpdate, makeRequest, update, values?.uuid])

  const renderSection = () => {
    switch (activeStep) {
      case 0:
        return <TypesSection data={data} setData={setData} />
      case 1:
        return <InfoSection data={data} setData={setData} />
      case 2:
        return <NameSection data={data} setData={setData} />
      default:
        return null
    }
  }

  const nextDisabled = useMemo(() => {
    switch (activeStep) {
      case 0:
        return !data.class
      case 1:
        return !data.fields.length || !data.fieldNames.length || !data.index
      case 2:
        return !data.name
      default:
        return false
    }
  }, [
    activeStep,
    data.class,
    data.fieldNames.length,
    data.fields.length,
    data.index,
    data.name,
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
          {activeStep !== 2 ? (
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
