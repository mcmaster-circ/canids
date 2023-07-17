import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'

export const defaultValues = {
  user: '',
  pass: '',
}

export const loginFormConfig = [
  {
    component: 'input',
    name: 'user',
    label: 'Email',
    size: 'small',
    column: 12,
    autoComplete: 'false',
    color: 'secondary',
  },
  {
    component: 'input',
    name: 'pass',
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
      user: yup
        .string()
        .required('Email is  a required field')
        .email('Please enter a valid email'),
      pass: yup
        .string()
        .required('Password is  a required field')
        .min(2, 'Enter at least 2 symbols'),
    })
    .required()
)
