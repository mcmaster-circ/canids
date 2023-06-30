import { MenuItem } from '@mui/material'
import Grid from '@mui/material/Unstable_Grid2'
import { FormSelect, FormTextInput } from '@molecules'

interface FormRenderPropsTypes {
  component: 'select' | 'input' | string
  name: string
  label: string
  options?: { name: string; value: string }[] | string[] | any
  column: number
  control: any
  errors: any
  variant?: any
}

export default ({
  component,
  column,
  name,
  label,
  options,
  control,
  errors,
  variant,
  ...rest
}: FormRenderPropsTypes) => {
  switch (component) {
    case 'select':
      return (
        <Grid xs={column} key={name}>
          <FormSelect
            label={label}
            control={control}
            helperText={errors?.[name]?.message}
            error={!!errors?.[name]}
            name={name}
            {...rest}
          >
            {options.map((o: { name: string; value: string } & string) => (
              <MenuItem key={o.value || o} value={o.value || o}>
                {o.name || o}
              </MenuItem>
            ))}
          </FormSelect>
        </Grid>
      )
    case 'input':
    default:
      return (
        <Grid xs={column} key={name}>
          <FormTextInput
            label={label}
            helperText={errors?.[name]?.message}
            error={!!errors?.[name]}
            name={name}
            variant={variant || 'outlined'}
            control={control}
            {...rest}
          />
        </Grid>
      )
  }
}
