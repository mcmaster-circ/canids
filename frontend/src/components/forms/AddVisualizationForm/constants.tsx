import { BarChart, PieChart, TableChart } from '@mui/icons-material'

export interface SectionProps {
  data: any
  setData: (a: any) => void
}

export interface FormProps {
  name: string
  index: string
  class: string
  fields: string[]
  fieldNames: string[]
}
export interface InfoSectionFieldProps {
  i?: any
  data: FormProps
  fieldsList: any[]
  setData: (d: FormProps) => void
}
export interface FieldsProps {
  name: string
  type: string
}
export type FieldsListProps = { index: string; fields: FieldsProps[] }

export const defaultFormValues = {
  name: '',
  index: '',
  class: '',
  fields: [],
  fieldNames: [],
}

export const steps = ['Type', 'Information', 'Name']

export const typeButtons = [
  { label: 'Bar', icon: <BarChart />, key: 'bar' },
  { label: 'Pie', icon: <PieChart />, key: 'pie' },
  { label: 'Table', icon: <TableChart />, key: 'table' },
]
