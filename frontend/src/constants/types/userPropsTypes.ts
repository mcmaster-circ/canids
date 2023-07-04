export interface AddUserProps {
  name: string
  class: 'admin' | 'standard'
  uuid: string
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
