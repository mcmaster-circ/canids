import {
  Bar,
  BarChart,
  Cell,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'

export default ({ chartData }: any) => (
  <ResponsiveContainer width="100%" height="100%">
    <BarChart
      data={chartData}
      margin={{
        top: 5,
        right: 30,
        left: 20,
        bottom: 30,
      }}
    >
      <XAxis dataKey="name" style={{ fontSize: '12px' }} />
      <YAxis
        style={{ fontSize: '12px' }}
        label={{
          value: 'Connections',
          angle: -90,
          position: 'insideleft',
          fontSize: '12px',
          dx: -30,
        }}
      />
      <Tooltip />
      <Bar dataKey="Connections">
        {chartData.map((d: any, i: number) => (
          <Cell key={i} fill={d.fill} />
        ))}
      </Bar>
    </BarChart>
  </ResponsiveContainer>
)
