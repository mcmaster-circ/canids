export interface ViewListItemProps {
  uuid: string
  name: string
  class: 'line' | 'bar' | 'pie' | 'table'
  index: string
  fields: string[]
  fieldNames: string[]
}

export interface AddViewProps {
  name: string
  class: 'line' | 'bar' | 'pie' | 'table'
  field: string
  fieldName: string
}

export interface UpdateViewProps extends AddViewProps {
  uuid: string
}

export interface DeleteViewProps {
  uuid: string
}
