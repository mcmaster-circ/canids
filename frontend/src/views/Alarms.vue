<template>
  <div class="column">

    <!-- Title bar -->
    <span style="margin-bottom:20px;">
      <section class="tile is-ancestor" style="margin-bottom: 20px; margin-left: 0; margin-right: 0">
        <article class="tile is-child box" style="padding: 0">
          <h1 style="text-align: center; font-size: 35px; color: black; margin-bottom: 10px">Alarms</h1>
        </article>
      </section>
    </span>

    <!-- date time selectors -->
    <div class="tile is-ancestor">
      <div class="tile is-parent">
        <article class="tile is-child box">
          <b-field label="Start time:">
            <b-datetimepicker v-model="startTime">
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
            <b-datetimepicker v-model="endTime">
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

    <div class="alarm-section">
      <div v-for="alarm in alarms" :key="alarm.asset">
        <div class="alarm-card card">

          <div class="alarm-row">
            <span class="alarm-row-label">
              <b>Asset:</b>
            </span>
            <span class="alarm-row-data">
              {{alarm.asset}}
            </span>
          </div>

          <div class="alarm-row">
            <span class="alarm-row-label">
              <b>Time:</b>
            </span>
            <span class="alarm-row-data">
              {{alarm.time}}
            </span>
          </div>

          <div class="alarm-row">
            <span class="alarm-row-label">
              <b>Source IP:</b>
            </span>
            <span class="alarm-row-data">
              {{alarm.sourceIP}}
            </span>
          </div>

          <div class="alarm-row">
            <span class="alarm-row-label">
              <b>Source Port:</b>
            </span>
            <span class="alarm-row-data">
              {{alarm.sourcePort}}
            </span>
          </div>

          <div class="alarm-row">
            <span class="alarm-row-label">
              <b>Destination IP:</b>
            </span>
            <span class="alarm-row-data">
              {{alarm.destIP}}
            </span>
          </div>

          <div class="alarm-row">
            <span class="alarm-row-label">
              <b>Destination Port:</b>
            </span>
            <span class="alarm-row-data">
              {{alarm.destPort}}
            </span>
          </div>

          <div class="alarm-row">
            <span class="alarm-row-label">
              <b>Alarms:</b>
            </span>
            <span class="alarm-row-data">
              <ul>
                <li class="filter" v-for="(filter, key) in alarm.alarms" :key="key">
                  {{key}}: {{filter}}
                </li>
              </ul>
            </span>
          </div>

        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Alarms',
  data: () => {
    return {
      startTime: new Date(),
      endTime: new Date(),
      alarms: [
        {
          time: "",
          asset: "asset1",
          sourceIP: "",
          sourcePort: 0,
          destIP: "",
          destPort: 0,
          alarms: {
            filter1: true,
            filter2: false
          }
        },
        {
          time: "",
          asset: "asset2",
          sourceIP: "",
          sourcePort: 0,
          destIP: "",
          destPort: 0,
          alarms: {
            filter1: true,
            filter2: false
          }
        },
        {
          time: "",
          asset: "asset3",
          sourceIP: "",
          sourcePort: 0,
          destIP: "",
          destPort: 0,
          alarms: {
            filter1: true,
            filter2: false
          }
        },
        {
          time: "",
          asset: "asset4",
          sourceIP: "",
          sourcePort: 0,
          destIP: "",
          destPort: 0,
          alarms: {
            filter1: true,
            filter2: false
          }
        },
        {
          time: "",
          asset: "asset5",
          sourceIP: "",
          sourcePort: 0,
          destIP: "",
          destPort: 0,
          alarms: {
            filter1: true,
            filter2: false
          }
        }
      ]
    };
  },
  mounted() {
    // default is 30 minutes of data
    this.startTime = new Date()
    this.endTime = new Date()
    this.startTime.setMinutes(this.startTime.getMinutes() - 30)
  },
  methods: {}
}
</script>

<style>
h1 {
  margin: 10px;
  font-size: 24px;
  font-weight: bold;
}

.alarm-section {
  display: flex;
  margin: 5px;
  flex-direction: row;
  flex-wrap: wrap;
}

.alarm-card {
  padding: 10px;
  margin: 5px;
}

.filter {
  margin-left: 20px;
}
</style>
