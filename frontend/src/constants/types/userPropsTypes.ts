export interface UserProps {
  name: string
  uuid: string
  exp: number
  iat: number
  class: 'admin' | 'standard'
  activated: boolean
}

export interface AddUserProps {
  name: string
  class: 'admin' | 'standard'
  uuid: string
  activated?: boolean
}

export interface UpdateUserProps {
  name: string
  uuid: string
  userId: string
  class?: 'admin' | 'standard'
  activated?: boolean
}

export interface ResetUserPassProps {
  uuid: string
}

export interface DeleteUserProps {
  uuid: string
}
