import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'

export const defaultValues = {
  email: '',
  password: '',
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
  {
    component: 'input',
    name: 'password',
    label: 'Password',
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
      email: yup.string().required().min(3),
      password: yup.string().required().min(3),
    })
    .required()
)
