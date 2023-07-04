import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'

export const defaultValues = {
  name: '',
  email: '',
  pass: '',
  passConfirm: '',
}

export const registerFormConfig = [
  {
    component: 'input',
    name: 'name',
    label: 'Name',
    size: 'small',
    column: 12,
    autoComplete: 'false',
  },
  {
    component: 'input',
    name: 'email',
    label: 'Email',
    size: 'small',
    column: 12,
    autoComplete: 'false',
  },
  {
    component: 'input',
    name: 'pass',
    label: 'Password',
    type: 'password',
    size: 'small',
    column: 12,
    autoComplete: 'false',
  },
  {
    component: 'input',
    name: 'passConfirm',
    label: 'Confirm Password',
    type: 'password',
    size: 'small',
    column: 12,
    autoComplete: 'false',
  },
]

export const resolver = yupResolver(
  yup
    .object()
    .shape({
      name: yup
        .string()
        .required('Name is  a required field')
        .min(3, 'Enter at least 3 symbols'),
      email: yup
        .string()
        .required('Email is  a required field')
        .email('Please enter a valid email'),
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
