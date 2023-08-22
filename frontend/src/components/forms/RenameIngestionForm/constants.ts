import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'

export const defaultValues = {
  name: '',
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
]

export const resolver = yupResolver(
  yup
    .object()
    .shape({
      name: yup.string().required('Name is a required field'),
    })
    .required()
)
