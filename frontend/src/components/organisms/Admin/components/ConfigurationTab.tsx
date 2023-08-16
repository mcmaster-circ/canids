import { useCallback, useEffect, useState } from 'react'
import { Box, Button, Grid, Switch, Typography } from '@mui/material'
import { updateConfiguration, getConfiguration } from '@api/configuration'
import { useRequest } from '@hooks'
import { Loader } from '@atoms'
import { SaveModal, SetupMailModal } from '@modals'
import { defaultSaveModalState, defaultSetupMailModalState } from '../constants'
import { Setting, SetupMailSettings, BooleanSettings } from '@constants/types'
import SetupMailForm from 'src/components/forms/SetupMailForm'

export default () => {
  const [saveModal, setSaveModal] = useState(defaultSaveModalState)
  const [setupMailModal, setSetupMailModal] = useState(
    defaultSetupMailModalState
  )
  // prettier-ignore
  const [booleanSettings, setBooleanSettings] = useState<BooleanSettings>({
    middlewareDisable: { name: 'MIDDLEWARE_DISABLE', label: 'Middleware Disable', value: 'false', prevValue: 'false', isAdvanced: false},
    httpsEnabled: { name: 'HTTPS_ENABLED', label: 'HTTPS Enabled', value: 'false', prevValue: 'false', isAdvanced: false},
    userRegistration: { name: 'USER_REGISTRATION', label: 'User Registration', value: 'false', prevValue: 'false', isAdvanced: false},
    userActivated: { name: 'USER_ACTIVATED', label: 'User Activated', value: 'false', prevValue: 'false', isAdvanced: false},
    debugLogging: { name: 'DEBUG_LOGGING', label: 'Debug Logging', value: 'false', prevValue: 'false', isAdvanced: false},
  })
  // prettier-ignore
  const [setupMailSettings, setSetupMailSettings] = useState<SetupMailSettings>({
    service: { name: 'MAIL_SERVICE', label: 'Mail Service', value: 'NONE', prevValue: 'NONE', isAdvanced: false },
    url: { name: 'MAIL_URL', label: 'Mail URL', value: '', prevValue: '', isAdvanced: false },
    apiKey: { name: 'MAIL_API_KEY', label: 'Mail API Key', value: '', prevValue: '', isAdvanced: false },
    fromAddress: { name: 'MAIL_FROM_ADDRESS', label: 'Mail From Address', value: '', prevValue: '', isAdvanced: false },
    fromName: { name: 'MAIL_FROM_NAME', label: 'Mail From Name', value: '', prevValue: '', isAdvanced: false },
    domain: { name: 'MAIL_DOMAIN', label: 'Mail Domain', value: '', prevValue: '', isAdvanced: false},
  })

  const { data, loading, makeRequest } = useRequest({
    request: getConfiguration,
  })
  const { makeRequest: saveRequest } = useRequest({
    request: updateConfiguration,
    requestByDefault: false,
    needSuccess: 'Successfully saved configuration',
  })

  const handleCloseSave = useCallback(() => {
    setSaveModal(defaultSaveModalState)
    setTimeout(() => makeRequest(), 1500)
  }, [makeRequest])

  const handleCloseSetupMail = useCallback(() => {
    setSetupMailModal(defaultSetupMailModalState)
    setTimeout(() => makeRequest(), 1500)
  }, [makeRequest])

  useEffect(() => {
    console.log(booleanSettings.debugLogging.value)
  }, [booleanSettings.debugLogging.value])

  // prettier-ignore
  useEffect(() => {
    if (data) {
      const d: Setting[] = data

      const service = d.find((s) => s.name === setupMailSettings.service.name)
      const url = d.find((s) => s.name === setupMailSettings.url.name)
      const apiKey = d.find((s) => s.name === setupMailSettings.apiKey.name)
      const fromAddress = d.find((s) => s.name === setupMailSettings.fromAddress.name)
      const fromName = d.find((s) => s.name === setupMailSettings.fromName.name)
      const domain = d.find((s) => s.name === setupMailSettings.domain.name)
      const sms: SetupMailSettings = {
        service: { ...setupMailSettings.service, value: service?.value || 'NONE', prevValue: service?.value || 'NONE'},
        url: { ...setupMailSettings.url, value: url?.value || '', prevValue: url?.value || '' },
        apiKey: { ...setupMailSettings.apiKey, value: apiKey?.value || '', prevValue: apiKey?.value || '' },
        fromAddress: { ...setupMailSettings.fromAddress, value: fromAddress?.value || '', prevValue: fromAddress?.value || '' },
        fromName: { ...setupMailSettings.fromName, value: fromName?.value || '', prevValue: fromName?.value || '' },
        domain: { ...setupMailSettings.domain, value: domain?.value || '', prevValue: domain?.value || '' },
      }
      setSetupMailSettings(sms)

      const middlewareDisable = d.find((s) => s.name === booleanSettings.middlewareDisable.name)
      const httpsEnabled = d.find((s) => s.name === booleanSettings.httpsEnabled.name)
      const userRegistration = d.find((s) => s.name === booleanSettings.userRegistration.name)
      const userActivated = d.find((s) => s.name === booleanSettings.userActivated.name)
      const debugLogging = d.find((s) => s.name === booleanSettings.debugLogging.name)
      const bs: BooleanSettings = {
        middlewareDisable: { ...booleanSettings.middlewareDisable, value: middlewareDisable?.value || 'false', prevValue: middlewareDisable?.value || 'false' },
        httpsEnabled: { ...booleanSettings.httpsEnabled, value: httpsEnabled?.value || 'false', prevValue: httpsEnabled?.value || 'false' },
        userRegistration: { ...booleanSettings.userRegistration, value: userRegistration?.value || 'false', prevValue: userRegistration?.value || 'false' },
        userActivated: { ...booleanSettings.userActivated, value: userActivated?.value || 'false', prevValue: userActivated?.value || 'false' },
        debugLogging: { ...booleanSettings.debugLogging, value: debugLogging?.value || 'false', prevValue: debugLogging?.value || 'false' },
      }
      setBooleanSettings(bs)
    }
  }, [data])

  return (
    <>
      <Box
        sx={{
          display: 'flex',
          flexWrap: 'wrap',
          gap: 2,
          justifyContent: 'space-between',
          alignItems: 'center',
          mb: 3,
        }}
      >
        <Typography variant="h6" fontWeight={700}>
          Configuration
        </Typography>
      </Box>
      <Box
        sx={{
          height: '100%',
          width: '100%',
          display: 'grid',
          gridTemplateColumns: '1fr',
        }}
      >
        <Box mb={3}>
          <Grid container spacing={1} columnSpacing={0.5}>
            <Grid item xs={3}>
              <Typography padding={0.5}>Mail Service</Typography>
            </Grid>
            <Grid item xs={9}>
              <Button
                variant="contained"
                onClick={() => {
                  setSetupMailModal((s) => ({ ...s, open: true }))
                }}
              >
                {setupMailSettings.service.value !== 'NONE'
                  ? setupMailSettings.service.value
                  : 'Setup'}
              </Button>
            </Grid>
          </Grid>
        </Box>
        <Box mb={3}>
          <Grid container spacing={1} columnSpacing={0.5}>
            <Grid item xs={3}>
              <Typography padding={0.5}>User Registration</Typography>
            </Grid>
            <Grid item xs={9}>
              <Switch
                value={booleanSettings.userRegistration.value}
                onChange={(e) =>
                  setBooleanSettings((d: BooleanSettings) => ({
                    ...d,
                    userRegistration: {
                      ...d.userRegistration,
                      value: e.target.checked.toString(),
                    },
                  }))
                }
              />
            </Grid>
            <Grid item xs={3}>
              <Typography padding={0.5}>User Activated</Typography>
            </Grid>
            <Grid item xs={9}>
              <Switch
                value={booleanSettings.userActivated.value}
                onChange={(e) =>
                  setBooleanSettings((d: BooleanSettings) => ({
                    ...d,
                    userActivated: {
                      ...d.userActivated,
                      value: e.target.checked.toString(),
                    },
                  }))
                }
              />
            </Grid>
          </Grid>
        </Box>
        <Box
          sx={{
            border: '3px solid red',
            borderRadius: '5px',
            backgroundColor: '#ffcccb',
            mb: 3,
          }}
        >
          <Grid container spacing={1} columnSpacing={0.5}>
            <Grid item xs={3}>
              <Typography padding={0.5}>Middleware Disable</Typography>
            </Grid>
            <Grid item xs={9}>
              <Switch
                value={booleanSettings.middlewareDisable.value}
                onChange={(e) =>
                  setBooleanSettings((d: BooleanSettings) => ({
                    ...d,
                    middlewareDisable: {
                      ...d.middlewareDisable,
                      value: e.target.checked.toString(),
                    },
                  }))
                }
              />
            </Grid>
            <Grid item xs={3}>
              <Typography padding={0.5}>HTTPS Enabled</Typography>
            </Grid>
            <Grid item xs={9}>
              <Switch
                value={booleanSettings.httpsEnabled.value}
                onChange={(e) =>
                  setBooleanSettings((d: BooleanSettings) => ({
                    ...d,
                    httpsEnabled: {
                      ...d.httpsEnabled,
                      value: e.target.checked.toString(),
                    },
                  }))
                }
              />
            </Grid>
            <Grid item xs={3}>
              <Typography padding={0.5}>Debug Logging</Typography>
            </Grid>
            <Grid item xs={9}>
              <Switch
                value={booleanSettings.debugLogging.value}
                onChange={(e) =>
                  setBooleanSettings((d: BooleanSettings) => ({
                    ...d,
                    debugLogging: {
                      ...d.debugLogging,
                      value: e.target.checked.toString(),
                    },
                  }))
                }
              />
            </Grid>
          </Grid>
        </Box>
        <Button
          variant="contained"
          onClick={() => {
            setSaveModal(() => ({ open: true }))
          }}
        >
          Save
        </Button>
      </Box>
      {loading && <Loader />}
      <SaveModal
        open={saveModal}
        title="Settings"
        data={booleanSettings}
        request={saveRequest}
        handleClose={handleCloseSave}
      />
      <SetupMailModal
        open={setupMailModal.open}
        title="Setup Mail"
        handleClose={handleCloseSetupMail}
      >
        <SetupMailForm
          values={setupMailSettings}
          handleClose={handleCloseSetupMail}
        />
      </SetupMailModal>
    </>
  )
}
