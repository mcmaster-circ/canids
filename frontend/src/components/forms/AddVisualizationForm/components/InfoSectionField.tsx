import { Fragment, useCallback, useMemo } from 'react'
import Grid from '@mui/material/Unstable_Grid2'
import {
  FormControl,
  IconButton,
  InputLabel,
  MenuItem,
  Select,
  TextField,
} from '@mui/material'
import { Delete } from '@mui/icons-material'
import {
  FieldsListProps,
  FieldsProps,
  InfoSectionFieldProps,
} from '../constants'

export default ({ i, data, fieldsList, setData }: InfoSectionFieldProps) => {
  const isI: boolean = useMemo(() => typeof i === 'number', [i])

  const indexFieldsList = useMemo(
    () =>
      fieldsList
        ?.find((f: FieldsListProps) => f.index === data.index)
        ?.fields.map((f: FieldsProps) => f.name)
        .sort((a: string, b: string) => {
          return a.localeCompare(b)
        }),
    [data.index, fieldsList]
  )

  const handleDelete = useCallback(
    (i: number) =>
      setData({
        ...data,
        fieldNames: data.fieldNames.filter(
          (_: string, index: number) => index !== i
        ),
        fields: data.fields.filter((_: string, index: number) => index !== i),
      }),
    [data, setData]
  )

  const handleSelect = useCallback(
    (v: string, i?: number) =>
      setData({
        ...data,
        fields: isI
          ? data.fields.map((f: string, index: number) => (index === i ? v : f))
          : [v],
      }),
    [data, isI, setData]
  )

  const handleInput = useCallback(
    (v: string, i?: number) =>
      setData({
        ...data,
        fieldNames: isI
          ? data.fieldNames.map((f: string, index: number) =>
              index === i ? v : f
            )
          : [v],
      }),
    [data, isI, setData]
  )

  return (
    indexFieldsList && (
      <Fragment key={i}>
        {isI && (
          <Grid xs={1}>
            <IconButton onClick={() => handleDelete(i)}>
              <Delete />
            </IconButton>
          </Grid>
        )}
        <Grid xs={isI ? 5 : 6}>
          <FormControl size="small" fullWidth>
            <InputLabel>Data</InputLabel>
            <Select
              value={data.fields[isI ? i : 0] || ''}
              label="Data"
              onChange={(e) => handleSelect(e.target.value, i)}
            >
              {indexFieldsList?.map((f: any) => (
                <MenuItem key={f} value={f}>
                  {f}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
        <Grid xs={6}>
          <TextField
            label="Name"
            size="small"
            value={isI ? data.fieldNames[i] || '' : data.fieldNames[0] || ''}
            onChange={(e) => handleInput(e.target.value, i)}
            fullWidth
          />
        </Grid>
      </Fragment>
    )
  )
}
