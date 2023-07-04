import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'

export const defaultValues = {
  email: '',
}

export const loginFormConfig = [
  {
    component: 'input',
    name: 'email',
    label: 'Email',
    size: 'small',
    column: 12,
    autoComplete: 'false',
  },
]

export const resolver = yupResolver(
  yup
    .object()
    .shape({
      email: yup
        .string()
        .required('Email is  a required field')
        .email('Please enter a valid email'),
    })
    .required()
)
