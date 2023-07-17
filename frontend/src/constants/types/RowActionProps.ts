import { ReactNode } from 'react'

export interface RowActionProps {
  label: string
  key: string
  // eslint-disable-next-line no-unused-vars
  action: (v: any) => void
  icon?: ReactNode
}
