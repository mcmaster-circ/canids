export interface AddBlacklistProps {
  name: string
  url: string
}

export interface UpdateBlacklistProps extends AddBlacklistProps {
  uuid: string
}

export interface DeleteBlacklistProps {
  uuid: string
}
