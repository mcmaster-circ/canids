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

export interface DataViewParams {
  view: string
  start: string
  end: string
  interval: number
  maxSize: number
  from: number
}
