import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'

export const defaultValues = {
  name: '',
  uuid: '',
  password: '',
  passwordConfirm: '',
}

export const setupFormConfig = [
  {
    component: 'input',
    name: 'name',
    label: 'Name',
    size: 'small',
    column: 12,
    autoComplete: 'false',
    color: 'secondary',
  },
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
  },
  {
    component: 'input',
    name: 'passwordConfirm',
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
        .min(2, 'Enter at least 2 symbols'),
      uuid: yup
        .string()
        .required('Email is  a required field')
        .email('Please enter a valid email'),
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