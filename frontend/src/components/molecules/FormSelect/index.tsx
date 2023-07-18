import React from 'react'
import { Control, useController } from 'react-hook-form'
import {
  SelectProps,
  Select,
  FormControl,
  InputLabel,
  FormHelperText,
} from '@mui/material'

export interface SelectPropsTypes extends SelectProps {
  control: Control | any
  helperText?: string
  name: string
  defaultValue?: string
  error?: boolean
  size?: 'small' | 'medium'
}

const FormTextInput = ({
  helperText,
  control,
  name,
  defaultValue,
  error,
  children,
  size,
  ...rest
}: SelectPropsTypes) => {
  const { field } = useController({
    name,
    control,
    defaultValue: defaultValue,
  })
  return (
    <FormControl fullWidth error={error} size={size}>
      <InputLabel id={`select-input-label-${field.name}`}>
        {rest.label}
      </InputLabel>
      <Select
        labelId={`select-label-${field.name}`}
        id={`select-${field.name}`}
        name={field.name}
        value={field.value}
        onChange={field.onChange}
        {...rest}
      >
        {children}
      </Select>
      {helperText && <FormHelperText>{helperText}</FormHelperText>}
    </FormControl>
  )
}

export default FormTextInput
