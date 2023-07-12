export interface GetAlarmsProps {
  index: string[]
  source: string[]
  start: string
  end: string
  maxSize: number
  from: number
}

export interface AlarmProps {
  uid: string
  host: string
  timestamp: string
  id_orig_h: string
  id_orig_p: number
  id_orig_h_pos: string[]
  id_resp_h: string
  id_resp_p: number
  id_resp_h_pos: string[]
}
