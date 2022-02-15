<template>
  <div class="modal-card" style="min-height:100vh">
    <header class="modal-card-head">
      <p class="modal-card-title">Edit</p>
    </header>
    <section class="modal-card-body">
      <div>
        <b-field label="Groups">
          <b-select disabled expanded v-model="dashboardObject.group" placeholder="Group...">
            <option v-for="group in groups" :value="group.uuid" :key="group.uuid">{{ group.name }}</option>
          </b-select>
        </b-field>
        <b-field label="Name">
          <b-input v-model="dashboardObject.name"></b-input>
        </b-field>
        <b-field label="Views">
          <b-select expanded multiple v-model="dashboardObject.views" placeholder="Views...">
            <option
              @click="addSizes"
              v-for="view in views"
              :value="view.uuid"
              :key="view.uuid"
            >{{view.name}}</option>
          </b-select>
        </b-field>
        <b-field label="View Sizes">
          <div v-for="(view,index) of dashboardObject.views" :key="view">
            {{views[view].name}}
            <b-select expanded v-model="dashboardObject.sizes[index]" placeholder="Select a size">
              <option value="half">Half</option>
              <option value="full">Full</option>
            </b-select>
          </div>
        </b-field>
      </div>
    </section>
    <footer class="modal-card-foot" style="display: flex;justify-content: flex-end">
      <b-button class="button" type="button" @click="$parent.close()">Close</b-button>
      <b-button type="is-primary" style="background-color:#712844" @click="save">Save</b-button>
    </footer>
  </div>
</template>

<script>
export default {
  name: "Dashboard Edit",
  props: {
    dashboardObject: {
      group: undefined,
      name: undefined,
      views: [],
      sizes: []
    }
  },
  data() {
    return {
      isOpen: 0,
      views: {},
      groups: {},
      dashboardData: []
    };
  },
  methods: {
    fetchData() {
      // Fetch Groups
      fetch(process.env.VUE_APP_ENDPOINT + "group/list")
        .then(response => {
          if (response.ok) {
            return Promise.all([response.ok, response.json()]);
          } else {
            return Promise.all([response.ok, response.text()]);
          }
        })
        .then(response => {
          this.groups = {};
          const status = response[0];
          const data = response[1];
          if (!status) {
            this.$buefy.snackbar.open(data);
          }
          this.groups[data.current.uuid] = data.current.name;

          for (let i = 0; i < data.others.length; i++) {
            this.groups[data.others[i].uuid] = data.others[i].name;
          }
        });
      // Fetch Views
      fetch(process.env.VUE_APP_ENDPOINT + "views/list")
        .then(response => {
          if (response.ok) {
            return Promise.all([response.ok, response.json()]);
          } else {
            return Promise.all([response.ok, response.text()]);
          }
        })
        .then(response => {
          this.groups = {};
          const status = response[0];
          const data = response[1];
          if (!status) {
            this.$buefy.snackbar.open(data);
          }
          for (let i = 0; i < data.views.length; i++) {
            this.views[data.views[i].uuid] = data.views[i].name;
          }
        });
    },
    addSizes() {
      this.dashboardObject.sizes = new Array(
        this.dashboardObject.views.length
      ).fill("half");
    },
    save() {
      fetch(process.env.VUE_APP_ENDPOINT + "dashboard/update", {
        method: "post",
        body: JSON.stringify(this.dashboardObject)
      })
        .then(response => {
          if (response.ok) {
            return Promise.all([response.ok, response.json()]);
          } else {
            return Promise.all([response.ok, response.text()]);
          }
        })
        .then(response => {
          const status = response[0];
          const data = response[1];
          if (!status) {
            this.$buefy.snackbar.open(data);
          }
          this.$parent.close();
        });
    }
  },
  components: {}
};
</script>

<style scoped lang="scss">
</style>
