import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'

export const defaultValues = {
  name: '',
  url: '',
}

export const addFormConfig = [
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
    name: 'url',
    label: 'Url',
    size: 'small',
    column: 12,
    autoComplete: 'false',
  },
]

export const resolver = yupResolver(
  yup
    .object()
    .shape({
      name: yup.string().required('Name is  a required field'),
      url: yup.string().required('Url is  a required field'),
    })
    .required()
)
