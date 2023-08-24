import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'

export const defaultValues = {
  name: '',
  class: '',
  activated: '',
}

export const formConfig = [
  {
    component: 'input',
    name: 'name',
    label: 'Name',
    size: 'small',
    column: 12,
    autoComplete: 'false',
  },
  {
    component: 'select',
    name: 'class',
    label: 'Role',
    size: 'small',
    options: [
      { name: 'Standard', value: 'standard' },
      { name: 'Admin', value: 'admin' },
    ],
    column: 6,
    autoComplete: 'false',
  },
  {
    component: 'select',
    name: 'activated',
    label: 'Activated',
    size: 'small',
    options: [
      { name: 'Yes', value: 'true' },
      { name: 'No', value: 'false' },
    ],
    column: 6,
    autoComplete: 'false',
  },
]

export const resolver = yupResolver(
  yup
    .object()
    .shape({
      name: yup.string().required('Name is  a required field'),
      class: yup.string().required('Role is a required field'),
      activated: yup.string().required('Activated is a required field'),
    })
    .required()
)
