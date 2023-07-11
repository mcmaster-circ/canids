import { Check, Update } from '@mui/icons-material'
import { Box, Button } from '@mui/material'
import { DateTimePicker, PickersActionBarProps } from '@mui/x-date-pickers'

interface CustomActionBarProps {
  start: Date
  end: Date
  setStart: (d: Date) => void
  setEnd: (d: Date) => void
  handleRequest: ({ st, en }: { st?: Date; en?: Date }) => void
}

export default ({
  start,
  end,
  setStart,
  setEnd,
  handleRequest,
}: CustomActionBarProps) => {
  const renderCustomBar = (
    {
      ownerState: { value },
      onCancel,
      onAccept,
      className,
    }: PickersActionBarProps | any,
    opt: 'start' | 'end'
  ) => (
    <Box
      className={className}
      sx={{ display: 'flex', justifyContent: 'space-between', p: 1 }}
    >
      <Button
        variant="contained"
        onClick={() => {
          if (opt === 'start') {
            setStart(new Date())
            setEnd(new Date())
            handleRequest({ st: new Date(), en: new Date() })
          } else {
            setEnd(new Date())
            handleRequest({ en: new Date() })
          }
          onCancel()
        }}
        startIcon={<Update />}
      >
        Now
      </Button>
      <Button
        variant="outlined"
        onClick={() => {
          opt === 'start'
            ? handleRequest({ st: value })
            : handleRequest({ en: value })
          onAccept()
        }}
        endIcon={<Check />}
      >
        Ok
      </Button>
    </Box>
  )
  return (
    <>
      <DateTimePicker
        key="start"
        label="Start Time"
        ampm={false}
        value={start}
        timeSteps={{ hours: 1, minutes: 1 }}
        closeOnSelect={false}
        onChange={(date) => setStart(date as Date)}
        maxDateTime={end}
        slots={{
          actionBar: (props) => renderCustomBar(props, 'start'),
        }}
        sx={{
          '.MuiInputBase-root': {
            pointerEvents: 'none',
          },
          '& .MuiInputBase-root button': {
            pointerEvents: 'all',
          },
        }}
      />
      <DateTimePicker
        label="End Time"
        ampm={false}
        value={end}
        timeSteps={{ hours: 1, minutes: 1 }}
        disableFuture
        minDateTime={start}
        onChange={(date) => setEnd(date as Date)}
        slots={{
          actionBar: (props) => renderCustomBar(props, 'end'),
        }}
        sx={{
          '.MuiInputBase-root': {
            pointerEvents: 'none',
          },
          '& .MuiInputBase-root button': {
            pointerEvents: 'all',
          },
        }}
      />
    </>
  )
}
