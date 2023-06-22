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
      <div class="tile is-parent">
        <article class="tile is-child box">
          <b-field label="Select">
            <b-dropdown
              v-model="selectedIndices"
              multiple
              scrollable
              expanded
              block
              aria-role="list">
              <template #trigger="{ active }">
                <b-button
                  type="is-primary"
                  :icon-right="active ? 'menu-up' : 'menu-down'">
                  {{ indicesText }}
                </b-button>
              </template>
              <b-dropdown-item v-for="index in indices" :key="index" :value="index" aria-role="listitem">
                <span>{{ index }}</span>
              </b-dropdown-item>
            </b-dropdown>
          </b-field>
        </article>
      </div>
      <div class="tile is-parent">
        <article class="tile is-child box">
          <b-field label="Alarm Source Lists">
            <b-dropdown
              v-model="selectedSourceLists"
              multiple
              scrollable
              expanded
              block
              aria-role="list">
              <template #trigger="{ active }">
                <b-button
                  type="is-primary"
                  :icon-right="active ? 'menu-up' : 'menu-down'">
                  {{ sourceListsText }}
                </b-button>
              </template>
              <b-dropdown-item v-for="source in alarmSourceLists" :key="source" :value="source" aria-role="listitem">
                <span>{{ source }}</span>
              </b-dropdown-item>
            </b-dropdown>
          </b-field>
        </article>
      </div>
    </div>

    <!-- alarm cards -->
    <div class="tile is-ancestor">
        <div v-for="alarm in alarms" :key="alarm.uid" class="tile is-parent">
          <article class="tile is-child box">

            <div class="alarm-row">
              <span class="alarm-row-label">
                <b>UID:</b>
              </span>
              <span class="alarm-row-data">
                {{alarm.uid}}
              </span>
            </div>

            <div class="alarm-row">
              <span class="alarm-row-label">
                <b>Host:</b>
              </span>
              <span class="alarm-row-data">
                {{alarm.host}}
              </span>
            </div>

            <div class="alarm-row">
              <span class="alarm-row-label">
                <b>Time:</b>
              </span>
              <span class="alarm-row-data">
                {{alarm.timestamp}}
              </span>
            </div>

            <div class="alarm-row">
              <span class="alarm-row-label">
                <b>Source IP:</b>
              </span>
              <span class="alarm-row-data">
                {{`${alarm.id_orig_h}:${alarm.id_orig_p}`}}
              </span>
            </div>

            <div class="alarm-row">
              <span class="alarm-row-label">
                <b>Dest IP:</b>
              </span>
              <span class="alarm-row-data">
                {{`${alarm.id_resp_h}:${alarm.id_resp_p}`}}
              </span>
            </div>

            <div v-if="alarm.id_orig_h_pos.length>0" class="alarm-row">
              <span class="alarm-row-label">
                <b>Source Alarms:</b>
              </span>
              <span class="alarm-row-data">
                <ul>
                  <li class="filter" v-for="source in alarm.id_orig_h_pos" :key="source">
                    {{ source }}
                  </li>
                </ul>
              </span>
            </div>

            <div v-if="alarm.id_resp_h_pos.length>0" class="alarm-row">
              <span class="alarm-row-label">
                <b>Dest Alarms:</b>
              </span>
              <span class="alarm-row-data">
                <ul>
                  <li class="filter" v-for="source in alarm.id_resp_h_pos" :key="source">
                    {{ source }}
                  </li>
                </ul>
              </span>
            </div>
          </article>
        </div>
    </div>

    <div class="tile is-ancestor" style="justify-content: space-between;">
      <!-- pagination -->
      <div class="tile is-parent is-3" style="justify-content: space-between;">
        <article class="is-child">
          <b-field>
            <b-select v-model="cardPageNum" expanded>
              <option v-for="page in cardPages" :value="page" :key="page"> {{ page }} </option>
            </b-select>
          </b-field>
        </article>
      </div>

      <!-- pagination arrows -->
      <div class="tile is-parent is-3">
        <article class="tile is-child">
          <div class="buttons is-pulled-right">
            <!-- left arrow -->
            <b-button
              type="is-primary"
              icon-left="chevron-left"
              @click="cardPageNum = cardPageNum - 1"
              :disabled="cardPageNum === 1">
            </b-button>

            <!-- right arrow -->
            <b-button
              type="is-primary"
              icon-right="chevron-right"
              @click="cardPageNum = cardPageNum + 1"
              :disabled="cardPageNum === cardPages.length">
            </b-button>
          </div>
        </article>
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
      indices: [],
      selectedIndices: [],
      indicesText: 'None Selected',
      alarmSourceLists: [],
      selectedSourceLists: [],
      sourceListsText: 'None Selected',
      data: {
        fieldNames: [],
        data: null,
        availableRows: null
      },
      alarms: [],
      cardPageNum: 1,
      cardPages: [1, 2, 3, 4]
    }
  },
  created() {
    this.NUMBER_OF_CARDS = 4
    this.fetchIndices()
    this.fetchLists()
    this.fetchData()
  },
  mounted() {
    // default is 30 minutes of data
    this.startTime = new Date()
    this.endTime = new Date()
    this.startTime.setMinutes(this.startTime.getMinutes() - 30)
  },
  watch: {
    startTime: function (val) {
      this.fetchData()
    },
    endTime: function (val) {
      this.fetchData()
    },
    cardPageNum: function (val) {
      this.fetchData()
    },
    selectedIndices: function (val) {
      this.indicesText = this.dropdownText(val, this.indices.length)
      this.fetchData()
    },
    selectedSourceLists: function (val) {
      this.sourceListsText = this.dropdownText(val, this.alarmSourceLists.length)
      this.fetchData()
    }
  },
  methods: {
    fetchIndices() {
      fetch('/api/fields/list')
        .then(response => {
          return response.json()
        })
        .then(json => {
          this.indices = json.filter(i => i.index.includes('log.alarm')).map(i => {
            return i.index
          })
          this.selectedIndices = [...this.indices]
        })
    },
    fetchLists() {
      fetch('/api/blacklist/list')
        .then(response => {
          return response.json()
        })
        .then(json => {
          this.alarmSourceLists = json.blacklists.map(i => {
            return i.name
          })
          this.selectedSourceLists = [...this.alarmSourceLists]
        })
    },
    fetchData() {
      this.updatePages()
      const req = {
        index: this.selectedIndices,
        source: this.selectedSourceLists,
        start: this.startTime.toISOString(),
        end: this.endTime.toISOString(),
        maxSize: this.NUMBER_OF_CARDS,
        from: (this.NUMBER_OF_CARDS * (this.cardPageNum - 1))
      }
      fetch("/api/alarm/data", {
        method: "post",
        body: JSON.stringify(req)
      })
        .then(response => {
          if (response.ok) {
            return response.json()
          } else {
            return response.text()
          }
        })
        .then(data => {
          if (typeof data === 'string' && data.includes('false')) {
            this.$buefy.snackbar.open(data);
          } else {
            this.alarms = data.alarms
            this.availableRows = data.availableRows
            this.updatePages()
          }
        });
    },
    updatePages() {
      this.cardPages = []
      if (this.availableRows > 0) {
        this.cardPages = (new Array(Math.floor(this.availableRows / this.NUMBER_OF_CARDS))).fill(0).map((x, i) => i + 1)
      }

      if (this.cardPages.length === 0) {
        this.cardPages = [1]
        this.cardPageNum = 1
      }

      const maxPage = this.cardPages[this.cardPages.length - 1]
      if (maxPage < this.cardPages) {
        this.cardPageNum = maxPage
      }
    },
    dropdownText(selectedList, maxListLength) {
      if (selectedList.length === 0) {
        return 'None Selected'
      } else if (selectedList.length === 1) {
        return selectedList[0]
      } else if (selectedList.length === maxListLength) {
        return 'All Selected'
      } else {
        return 'Selected (' + selectedList.length + ')'
      }
    }
  }
}
</script>

<style>
h1 {
  margin: 10px;
  font-size: 24px;
  font-weight: bold;
}

.filter {
  margin-left: 20px;
}
</style>
