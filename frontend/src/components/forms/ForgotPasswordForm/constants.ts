import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'

export const defaultValues = {
  uuid: '',
}

export const loginFormConfig = [
  {
    component: 'input',
    name: 'uuid',
    label: 'Email',
    size: 'small',
    column: 12,
    autoComplete: 'false',
    color: 'secondary',
  },
]

export const resolver = yupResolver(
  yup
    .object()
    .shape({
      uuid: yup
        .string()
        .required('Email is  a required field')
        .email('Please enter a valid email'),
    })
    .required()
)
