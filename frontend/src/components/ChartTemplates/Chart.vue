<template>
  <div v-if="errMsg">
    <p>{{ errMsg }}</p>
  </div>
  <div v-else-if="isLine">
    <LineChart :chart-data="chartData" :options="chartOptions" />
  </div>
  <div v-else-if="isBar">
    <BarChart :chart-data="chartData" :options="chartOptions" />
  </div>
  <div v-else-if="isPie">
    <PieChart :chart-data="chartData" :options="chartOptions" />
  </div>
  <div v-else-if="isTable">
    <h1 style="text-align: center; font-size: 25px; font-weight: 400">{{ this.view.name + ' (' + this.view.authorized + ')' }}</h1>
    <b-field label="Page Size">
        <b-input
            v-model="tablePageSize"
            placeholder="Number"
            type="number"
            min="1"
            max="100">
        </b-input>
    </b-field>
    <b-table :data="tableData" :columns="tableColumns"></b-table>
    <b-field label="Page">
      <b-select v-model="tablePageNum" expanded>
        <option v-for="page in tablePages" :value="page" :key="page"> {{ page }} </option>
      </b-select>
    </b-field>
  </div>
</template>

<script>
import BarChart from './BarChart'
import PieChart from './PieChart'
import LineChart from './LineChart'
import moment from 'moment'

export default {
  name: 'Chart',

  // view is same structure in API, start and end are Date() objects
  props: ['view', 'start', 'end'],

  data () {
    return {
      // how often to poll API for changes (milliseconds)
      updateInterval: 10 * 1000,

      // charting modes
      isLine: false,
      isBar: false,
      isPie: false,
      isTable: false,

      // view name
      name: '',

      // errMsg is blank if no error occurred
      errMsg: '',

      // JSON chart data
      data: {
        fieldNames: [],
        class: '',
        data: null,
        availableRows: null
      },

      // ChartJS
      chartData: null,
      chartOptions: {
        responsive: true
      },

      // data tables
      tableData: [],
      tableColumns: [],

      // table settings
      tablePageSize: 10,
      tablePageNum: 1,
      tablePages: [1, 2, 3, 4],
      tableAvailableRows: 0
    }
  },

  watch: {
    start: function (newVal, oldVal) {
      this.fetchData()
    },
    end: function (newVal, oldVal) {
      this.fetchData()
    },
    tablePageSize: function (newVal, oldVal) {
      this.fetchData()
    },
    tablePageNum: function (newVal, oldVal) {
      this.fetchData()
    }
  },

  mounted () {
    // determine chart type
    switch (this.view.class) {
      case 'line':
        this.isLine = true
        break
      case 'bar':
        this.isBar = true
        break
      case 'pie':
        this.isPie = true
        break
      case 'table':
        this.isTable = true
        break
      default:
        this.errMsg = 'Invalid chart type provided'
    }

    // populate default options and chart name
    this.chartOptions = {
      responsive: true,
      maintainAspectRatio: false,
      title: {
        display: true,
        text: this.view.name + ' (' + this.view.authorized + ')',
        fontSize: 26
      },
      tooltips: {
        mode: "index",
        intersect: false
      }
    }

    // initial data fetch, automatically fetch data on set intervals
    this.fetchData()
    setInterval(this.fetchData, this.updateInterval)
  },

  methods: {
    /**
     * fetchData fetches data from the last 30 minutes. It calls generateChart()
     * upon fetching the data.
     */
    fetchData () {
      if (!this.isTableSizeValid()) return

      this.updateTablePages()

      const thing = parseInt((this.end.getTime() - this.start.getTime()) / 24000)
      const request = process.env.VUE_APP_ENDPOINT + 'data/?view=' + this.view.uuid + '&start=' + this.start.toISOString() + '&end=' + this.end.toISOString() + '&interval=' + thing.toString() + '&maxSize=' + this.tablePageSize.toString() + '&from=' + (this.tablePageSize * (this.tablePageNum - 1)).toString()
      fetch(request)
        .then(resp => {
          return Promise.all([resp.ok, resp.json()])
        })
        .then(resp => {
          const ok = resp[0];
          const data = resp[1];
          if (!ok) {
            this.errMsg = this.view.name + ' Failed to generate visualization: ' + data.message
          } else {
            this.data = data
            this.availableRows = data.availableRows

            this.updateTablePages()

            // success, generate chart
            this.generateChart()
          }
        })
    },

    isTableSizeValid() {
      const tablePageSize = parseInt(this.tablePageSize)
      return !isNaN(tablePageSize) && (tablePageSize >= 1) && (tablePageSize <= 100)
    },

    updateTablePages() {
      this.tablePages = []
      if (this.isTableSizeValid() && (this.availableRows > 0)) {
        this.tablePages = (new Array(Math.floor(this.availableRows / this.tablePageSize))).fill(0).map((x, i) => i + 1)
      }

      if (this.tablePages.length === 0) {
        this.tablePages = [1]
        this.tablePageNum = 1
      }

      const maxPage = this.tablePages[this.tablePages.length - 1]
      if (maxPage < this.tablePageNum) {
        this.tablePageNum = maxPage
      }
    },

    // generate data for the correct chart type
    generateChart () {
      if (this.isLine) {
        this.generateLine()
      } else if (this.isBar) {
        this.generateBar()
      } else if (this.isPie) {
        this.generatePie()
      } else if (this.isTable) {
        this.generateTable()
      }
    },

    /**
     * generateLine generates chartData for a line chart.
     */
    generateLine () {
      const xAxis = this.data.fieldNames[0]
      const yAxis = this.data.fieldNames[1]
      const xData = this.data.data[0]
      const yData = this.data.data[1]

      // if xData contains time, format it
      if (moment(xData[0], moment.ISO_8601, true).isValid()) {
        // format all fields
        for (let i = 0; i < xData.length; i++) {
          xData[i] = moment(xData[i]).format('ddd MMM D h:mm:ssa')
        }
      }

      // if yData contains time, format it
      if (moment(yData[0], moment.ISO_8601, true).isValid()) {
        // format all fields
        for (let i = 0; i < yData.length; i++) {
          yData[i] = moment(yData[i]).format('ddd MMM D h:mm:ssa')
        }
      }

      this.chartOptions.scales = {
        xAxes: [{
          ticks: {
            autoSkip: true,
            maxTicksLimit: 10.1 // this has to be fractional
          },
          scaleLabel: {
            display: true,
            labelString: xAxis
          }
        }],
        yAxes: [{
          ticks: {
            autoSkip: true,
            maxTicksLimit: 10.1 // this has to be fractional
          },
          scaleLabel: {
            display: true,
            labelString: yAxis
          }
        }]
      }
      const colors = ['#e74c3c', '#3498db', '#2ecc71', '#9b59b6', '#e67e22', '#34495e']

      this.chartData = {
        labels: xData,
        datasets: [
          {
            label: yAxis,
            data: yData,
            backgroundColor: colors[0],
            borderColor: colors[0],
            fill: false
          }
        ]
      }
    },

    /**
     * generateBar generates chartData for a bar chart.
     */
    generateBar () {
      const barType = 'Connections'
      const barVar = this.data.data[0]
      const barVal = this.data.data[1]

      // disable top legend
      this.chartOptions.legend = {
        display: false
      }
      // set y-axis
      this.chartOptions.scales = {
        yAxes: [{
          scaleLabel: {
            display: true,
            labelString: barType
          }
        }]
      }
      // populate chart data
      this.chartData = {
        labels: barVar,
        datasets: [
          {
            label: barType,
            data: barVal,
            backgroundColor: ['#e74c3c', '#3498db', '#2ecc71', '#9b59b6', '#e67e22', '#34495e']
          }
        ]
      }
    },

    /**
     * generatePie generates chartData for a pie chart.
     */
    generatePie () {
      const pieVar = this.data.data[0]
      const pieVal = this.data.data[1]

      // populate chart data
      this.chartData = {
        labels: pieVar,
        datasets: [
          {
            data: pieVal,
            backgroundColor: ['#e74c3c', '#3498db', '#2ecc71', '#9b59b6', '#e67e22', '#34495e']
          }
        ]
      }
    },

    /**
     * generateTable generates tableData for a data table.
     */
    generateTable () {
      this.tableData = []
      this.tableColumns = []

      // inject columns
      for (let i = 0; i < this.data.fieldNames.length; i++) {
        this.tableColumns.push({
          field: this.data.fields[i],
          label: this.data.fieldNames[i]
        })
      }

      // if count of columns does not match data, append Count
      if (this.data.fieldNames.length !== this.data.data.length) {
        this.tableColumns.push({
          field: 'count',
          label: 'Count'
        })
      }

      // iterate over all rows in the dataset
      for (let i = 0; i < this.data.data[0].length; i++) {
        // generate row by getting all columns
        const row = {}
        for (let j = 0; j < this.tableColumns.length; j++) {
          row[this.tableColumns[j].field] = this.data.data[j][i]
        }
        this.tableData.push(row)
      }
    }
  },

  components: {
    BarChart,
    LineChart,
    PieChart
  }
}
</script>
