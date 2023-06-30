import { Control, useController } from 'react-hook-form'
import { OutlinedTextFieldProps, TextField } from '@mui/material'

export interface TextFieldProps extends OutlinedTextFieldProps {
  control: Control | any
  helperText?: string
  name: string
  defaultValue?: string
  error?: boolean
}

export default ({
  control,
  helperText,
  name,
  defaultValue,
  error,
  ...rest
}: TextFieldProps) => {
  const { field } = useController({
    name,
    control,
    defaultValue: defaultValue || '',
  })
  return (
    <TextField
      helperText={helperText}
      name={field.name}
      value={field.value}
      onChange={field.onChange}
      fullWidth
      error={error}
      {...rest}
    />
  )
}
