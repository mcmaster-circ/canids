import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'

export const defaultValues = {
  uuid: '',
}

export const formConfig = [
  {
    component: 'input',
    name: 'uuid',
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
      uuid: yup.string().required('Name is  a required field'),
    })
    .required()
)
