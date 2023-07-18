export interface UserProps {
  name: string
  uuid: string
  exp?: number
  iat?: number
  class: 'admin' | 'standard'
  activated: boolean
}

export interface AddUserProps {
  name: string
  class: 'admin' | 'standard' | string
  uuid: string
  activated: string
}

export interface UpdateUserProps extends AddUserProps {
  userId: string
  updatePermission: true
}

export interface ResetUserPassProps {
  uuid: string
}

export interface DeleteUserProps {
  uuid: string
}
