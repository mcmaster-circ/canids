<template>
  <div class="hello" style="min-height:100vh">
    <section class="info-tiles">
      <span style="margin-bottom:20px;">
        <h1 style="font-size:50px; display:inline-block">Dashboard List</h1>
      </span>
      <article class="tile is-child box" style="width:1000px">
        <div>
          <b-table
            :data="dashboardData"
            ref="table"
            detailed
            hoverable
            custom-detail-row
            :opened-detailed="['Board Games']"
            detail-key="name"
            :show-detail-icon="true"
          >
            <template slot-scope="props">
              <b-table-column label="Name" field="name" sortable searchable>{{props.row.name}}</b-table-column>
              <b-table-column field="actions">
                <button
                  class="button is-primary"
                  @click="editDashboard = true; dataToEdit=props.row"
                >
                  <b-icon icon="pencil"></b-icon>
                </button>
              </b-table-column>
            </template>

            <template slot="detail" slot-scope="props">
              <tr v-for="(view, index) in props.row.views" :key="view">
                <td></td>
                <td>{{views[view].name}} - {{props.row.sizes[index]}} size</td>
                <td></td>
              </tr>
            </template>
            <template slot="empty">
              <section class="section">
                <div class="content has-text-grey has-text-centered">
                  <p>
                    <b-icon icon="emoticon-sad" size="is-large"></b-icon>
                  </p>
                  <p>Nothing here.</p>
                </div>
              </section>
            </template>
          </b-table>
        </div>
      </article>
    </section>
    <b-modal :active.sync="editDashboard" has-modal-card>
      <DashboardEdit :dashboardObject="dataToEdit"></DashboardEdit>
    </b-modal>
  </div>
</template>

<script>
import DashboardEdit from "./DashboardEdit";

export default {
  name: "Dashboard List",
  data() {
    return {
      dataToEdit: undefined,
      editDashboard: false,
      isOpen: 0,
      dashboardObject: {
        name: undefined,
        views: [],
        sizes: []
      },
      views: {},
      dashboardData: []
    };
  },
  mounted() {
    this.fetchData();
  },
  methods: {
    fetchData() {
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
          const status = response[0];
          const data = response[1];
          if (!status) {
            this.$buefy.snackbar.open(data);
          }
          for (let i = 0; i < data.views.length; i++) {
            this.views[data.views[i].uuid] = data.views[i].name;
          }
        });
    }
  },
  components: { DashboardEdit }
};
</script>

<style scoped lang="scss">
</style>
