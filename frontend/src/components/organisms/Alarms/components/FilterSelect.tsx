import {
  FormControl,
  InputLabel,
  MenuItem,
  OutlinedInput,
  Select,
} from '@mui/material'

export default ({ list, value, handleChange, title }: any) => (
  <FormControl sx={{ width: 250 }}>
    <InputLabel id="multiple-logs">{title}</InputLabel>
    <Select
      labelId="multiple-logs"
      id="multiple-logs-name"
      multiple
      value={value}
      onChange={handleChange}
      input={<OutlinedInput label={title} />}
      renderValue={(selected: string[]) =>
        selected.length === value.length ? 'All Selected' : selected.join(', ')
      }
    >
      {list.map((l: string) => (
        <MenuItem key={l} value={l}>
          {l}
        </MenuItem>
      ))}
    </Select>
  </FormControl>
)
