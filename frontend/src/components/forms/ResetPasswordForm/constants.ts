import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'

export const defaultValues = {
  password: '',
  passwordConfirm: '',
}

export const registerFormConfig = [
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
  {
    component: 'input',
    name: 'passwordConfirm',
    label: 'Confirm Password',
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
      password: yup
        .string()
        .required('Password is a required field')
        .min(2, 'Enter at least 2 symbols'),
      passwordConfirm: yup
        .string()
        .required('Please confirm your password')
        .oneOf([yup.ref('password')], 'Your passwords do not match.'),
    })
    .required()
)
