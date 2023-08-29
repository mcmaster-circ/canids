export interface Setting {
  name: string
  value: string
  label: string
  prevValue: string
  isAdvanced: boolean | undefined
}

export interface UpdateSetting {
  name: string
  value: string
}

export type Settings = SetupMailSettings & BooleanSettings

export interface SetupMailSettings {
  service: Setting
  url: Setting
  apiKey: Setting
  fromAddress: Setting
  fromName: Setting
  domain: Setting
  accessURL: Setting
}

export interface BooleanSettings {
  middlewareDisable: Setting
  httpsEnabled: Setting
  userRegistration: Setting
  userActivated: Setting
  debugLogging: Setting
}

export interface ListConfigurationProps {
  success: boolean
  configuration: Setting[]
}

export interface UpdateConfigurationProps {
  configuration: UpdateSetting[]
}

export interface UpdateMailSetupProps {
  service: string
  url: string
  apiKey: string
  fromAddress: string
  fromName: string
  domain: string
  accessURL: string
}
