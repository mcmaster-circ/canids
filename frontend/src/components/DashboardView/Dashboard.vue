<template>
  <div class="hello" style="min-height:100vh">
    <div class="columns">
      <Sidebar class="column is-one-fifth" :selectedViews="views" @saveViews="saveViews" hidden/>
      <div class="column">
        <span style="margin-bottom:20px;">
          <section class="tile is-ancestor" style="margin-bottom: 20px;">
            <article class="tile is-child box">
              <h1 id="panelTitle">{{ dashboard.name }}</h1>
            </article>
          </section>
        </span>
        <!-- date time selectors -->
        <div class="tile is-ancestor">
          <div class="tile is-parent">
            <article class="tile is-child box">
              <b-field label="Start time:">
                <b-datetimepicker v-model="startTime" :max-datetime="endTime" :datepicker="{ yearsRange }">
                  <template slot="left">
                    <button class="button is-primary"
                        @click="startTime = new Date()">
                        <b-icon icon="clock"></b-icon>
                        <span>Now</span>
                    </button>
                  </template>
                </b-datetimepicker>
              </b-field>
            </article>
          </div>
          <div class="tile is-parent">
            <article class="tile is-child box">
              <b-field label="End Time:">
                <b-datetimepicker v-model="endTime" :min-datetime="startTime" :max-datetime="new Date()" :datepicker="{ yearsRange }">
                  <template slot="left">
                    <button class="button is-primary"
                        @click="endTime = new Date()">
                        <b-icon icon="clock"></b-icon>
                        <span>Now</span>
                    </button>
                  </template>
                </b-datetimepicker>
              </b-field>
            </article>
          </div>
        </div>

        <!-- Chart grid -->
        <div class="tile is-ancestor" style="flex-wrap: wrap">
          <div class="tile is-parent is-6"
            v-for="(view, index) in views"
            :key="view"
            :index="index"
            replace-direction="horizontal"
            style="max-width: 100%;"
          >
            <div class="tile is-child box">
              <Chart :view="view" :start="startTime" :end="endTime" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Chart from "@/components/ChartTemplates/Chart";
import Sidebar from "@/components/sidebar_component/sidebar"

export default {
  name: "Dashboard",
  data() {
    var startTime = new Date()
    var endTime = new Date()
    return {
      startTime: startTime,
      endTime: endTime,
      yearsRange: [-100, 200],
      dashboard: {},
      visualizations: [],
      graphArray: [1, 2, 3, 4, 5, 6, 7],
      dragData: {},
      views: []
    };
  },
  created () {
    this.fetchViews()
  },
  mounted() {
    // default is 30 minutes of data
    this.startTime = new Date()
    this.endTime = new Date()
    this.startTime.setMinutes(this.startTime.getMinutes() - 30)
  },
  methods: {
    fetchViews() {
      // live
      fetch('/api/dashboard/get')
        .then(response => {
          return response.json()
        })
        .then(json => {
          this.dashboard = json
          fetch("/api/view/list")
            .then(resp => {
              return Promise.all([resp.ok, resp.json()]);
            })
            .then(resp => {
              const status = resp[0];
              const data = resp[1];
              if (!status) {
                this.$buefy.snackbar.open(data.message)
              } else {
                this.dashboard.views.forEach(view => {
                  data.views.forEach(jsonView => {
                    if (view === jsonView.uuid) {
                      this.views.push(jsonView)
                    }
                  })
                })
              }
            })
        })
    },
    sortend(e) {
      const { oldIndex, newIndex } = e;
      this.rearrange(this.views, oldIndex, newIndex);
      this.rearrange(this.dashboard.views, oldIndex, newIndex);
      this.save()
    },
    rearrange(array, oldIndex, newIndex) {
      if (oldIndex > newIndex) {
        array.splice(newIndex, 0, array[oldIndex]);
        array.splice(oldIndex + 1, 1);
      } else {
        array.splice(newIndex + 1, 0, array[oldIndex]);
        array.splice(oldIndex, 1);
      }
    },
    save() {
      this.dashboard.sizes = []
      this.dashboard.views.forEach(view => {
        this.dashboard.sizes.push("half")
      })
      fetch('/api/dashboard/update', {
        method: 'post',
        body: JSON.stringify(this.dashboard)
      })
        .then(response => response)
        .then(data => {
          if (data.status === 200) {
            this.$buefy.toast.open({ message: "Dashboard Updated", position: "is-top", type: "is-success" })
          }
        })
    },
    saveViews() {
      this.dashboard.views = []
      this.views.forEach(view => {
        this.dashboard.views.push(view.uuid)
      })
      this.save()
    },
    deleteView(index) {
      this.views.splice(index, 1)
      this.dashboard.views.splice(index, 1)
      this.dashboard.sizes.splice(index, 1)
      this.$buefy.toast.open({ message: "View deleted", position: "is-top", type: "is-success" })
    },
    deleteGraph(graph) {
      for (let i = 0; i < this.graphArray.length; i++) {
        if (this.graphArray[i] === graph) {
          this.graphArray.splice(i, 1);
        }
      }
    },
    getChartData(type) {
      return {
        labels: ["1", "2"],
        datasets: [
          {
            label: "Test",
            backgroundColor: "#2980b9",
            borderColor: "#2980b9",
            data: [10, 20],
            fill: false
          }
        ]
      };
    },
    getChartOptions(type) {
      return {
        responsive: true,
        maintainAspectRatio: false,
        title: {
          display: true,
          text: type,
          fontSize: 30
        },
        legend: {
          display: false
        },
        tooltips: {
          mode: "index",
          intersect: false
        },
        hover: {
          mode: "nearest",
          intersect: true
        },
        scales: {
          xAxes: [
            {
              display: true,
              scaleLabel: {
                display: true,
                labelString: "Time"
              }
            }
          ],
          yAxes: [
            {
              scaleLabel: {
                display: true,
                labelString: "Y-Axis"
              }
            }
          ]
        }
      };
    }
  },
  components: {
    Chart,
    Sidebar
  }
};
</script>

<style scoped>

#panelTitle{
    text-align: center;
    font-size: 35px;
    color: black;
    margin-bottom: 10px;
}

</style>
