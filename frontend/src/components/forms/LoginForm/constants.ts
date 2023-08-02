import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'

export const defaultValues = {
  uuid: '',
  password: '',
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
  {
    component: 'input',
    name: 'password',
    label: 'Password',
    type: 'password',
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
      password: yup
        .string()
        .required('Password is  a required field')
        .min(2, 'Enter at least 2 symbols'),
    })
    .required()
)
