interface BaseViewProps {
  name: string
  class: 'line' | 'bar' | 'pie' | 'table'
  index: string
  fields: string[]
  fieldNames: string[]
}

export interface ViewListItemProps extends BaseViewProps {
  uuid: string
}

export interface AddViewProps extends BaseViewProps {}

export interface UpdateViewProps extends AddViewProps {
  uuid: string
}

export interface DeleteViewProps {
  uuid: string
}
