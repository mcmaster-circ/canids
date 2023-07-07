import {
  Cell,
  Legend,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
} from 'recharts'

export default ({ chartData }: any) => (
  <ResponsiveContainer width="100%" height="100%">
    <PieChart>
      <Legend layout="horizontal" verticalAlign="top" />
      <Tooltip />
      <Pie
        data={chartData}
        cx="50%"
        cy="45%"
        outerRadius={120}
        dataKey="Connections"
        label
      >
        {chartData.map((d: any, i: number) => (
          <Cell key={i} fill={d.fill} />
        ))}
      </Pie>
    </PieChart>
  </ResponsiveContainer>
)
