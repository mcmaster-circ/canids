import { Fragment, useCallback, useMemo } from 'react'
import Grid from '@mui/material/Unstable_Grid2'
import {
  Button,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Typography,
} from '@mui/material'
import { useRequest } from '@hooks'
import { getFields } from '@api/fields'
import { FieldsListProps, FormProps, SectionProps } from '../constants'
import { Loader } from '@atoms'
import { Add } from '@mui/icons-material'
import { InfoSectionField } from './'

export default ({ data, setData }: SectionProps) => {
  const { data: fieldsList, loading: loadingFields } = useRequest({
    request: getFields,
  })

  const indexList = useMemo(
    () => fieldsList?.map((f: FieldsListProps) => f.index),
    [fieldsList]
  )

  const handleAdd = useCallback(
    () =>
      setData({
        ...data,
        fieldNames: [...data.fieldNames, ''],
        fields: [...data.fields, ''],
      }),
    [data, setData]
  )

  return (
    <>
      <Typography variant="h6" textAlign="center" mb={4}>
        Select Graph Information
      </Typography>
      <Grid container justifyContent="space-around" spacing={4}>
        <Grid xs={12}>
          {indexList && (
            <FormControl size="small" fullWidth>
              <InputLabel>Index</InputLabel>
              <Select
                value={data.index}
                label="Index"
                onChange={(e) =>
                  setData((d: FormProps) => ({
                    ...d,
                    index: e.target.value,
                    fields: [],
                    fieldNames: [],
                  }))
                }
              >
                {indexList?.map((f: string) => (
                  <MenuItem key={f} value={f}>
                    {f}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          )}
        </Grid>
        {data.index && (
          <>
            {data.class === 'table' ? (
              <>
                {data.fields.map((_: string, i: number) => (
                  <InfoSectionField
                    key={i}
                    i={i}
                    data={data}
                    setData={setData}
                    fieldsList={fieldsList}
                  />
                ))}
                <Grid
                  xs={12}
                  sx={{ display: 'flex', justifyContent: 'center' }}
                >
                  <Button
                    variant="contained"
                    sx={{ width: 200 }}
                    startIcon={<Add />}
                    onClick={handleAdd}
                  >
                    Add Field
                  </Button>
                </Grid>
              </>
            ) : (
              <InfoSectionField
                data={data}
                setData={setData}
                fieldsList={fieldsList}
              />
            )}
          </>
        )}
      </Grid>
      {loadingFields && <Loader />}
    </>
  )
}
