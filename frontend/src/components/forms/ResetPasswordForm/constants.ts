import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'

export const defaultValues = {
  pass: '',
  passConfirm: '',
}

export const registerFormConfig = [
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
  {
    component: 'input',
    name: 'passConfirm',
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
      pass: yup
        .string()
        .required('Password is a required field')
        .min(3, 'Enter at least 3 symbols'),
      passConfirm: yup
        .string()
        .required('Please confirm your password')
        .oneOf([yup.ref('pass')], 'Your passwords do not match.'),
    })
    .required()
)
