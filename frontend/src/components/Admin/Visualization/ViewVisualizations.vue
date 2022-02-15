<template>
  <div class="hello">
    <h1 style="font-size:50px;margin-bottom:20px">Visualizations</h1>
    <section class="info-tiles" style="margin-top:50px">
      <article class="tile is-child box">
        <b-button
          class="is-centered"
          type="is-primary"
          style="float:right"
          @click="addVisualization"
        >Create Visualization</b-button>

        <b-table
          :data="visualizations"
          :loading="isTableLoading"
          :hoverable="true">

          <template slot-scope="props">
              <b-table-column field="name" label="Name" style="cursor:pointer;">
                  {{ props.row.name }}
              </b-table-column>

              <b-table-column field="authorized" label="Authorized" style="cursor:pointer;">
                  {{ props.row.authorized }}
              </b-table-column>

              <b-table-column field="class" label="Class" style="cursor:pointer;">
                  {{ props.row.class }}
              </b-table-column>

              <b-table-column field="actions" numeric>
                <b-dropdown aria-role="list" class="is-pulled-right" position="is-bottom-left" style="cursor:pointer">
                  <b-icon icon="dots-vertical" slot="trigger"></b-icon>
                    <b-dropdown-item aria-role="listitem" @click="editVisualization(props.row)">
                      <b-icon icon="pencil"></b-icon>
                      Edit
                    </b-dropdown-item>
                    <b-dropdown-item aria-role="listitem" @click="deleteVisualization(props.row, props.index)">
                      <b-icon icon="delete"></b-icon>
                      Delete
                    </b-dropdown-item>
                </b-dropdown>
              </b-table-column>
          </template>

          <template slot="empty">
              <section class="section">
                  <div class="content has-text-grey has-text-centered">
                      <p>
                        <b-icon
                            icon="emoticon-sad"
                            size="is-large">
                        </b-icon>
                      </p>
                      <p>Nothing here.</p>
                  </div>
              </section>
          </template>
        </b-table>
      </article>
    </section>
    <b-modal :active.sync="visualizationModalActive" has-modal-card trap-focus>
      <AddEditVisualization :visualizationToEdit="visualizationToEdit" @addedVisualization="addedVisualization" @editedVisualization="editedVisualization"></AddEditVisualization>
    </b-modal>
  </div>
</template>

<script>
import AddEditVisualization from "./AddEditVisualization";

export default {
  data() {
    return {
      visualizations: [],
      isTableLoading: false,
      visualizationModalActive: false,
      visualizationToEdit: {}
    }
  },
  mounted() {
    this.fetchData();
  },
  methods: {
    addVisualization () {
      this.visualizationToEdit = {}
      this.visualizationModalActive = true
    },
    fetchData() {
      fetch('/api/view/list')
        .then(response => {
          return response.json()
        })
        .then(json => {
          this.visualizations = json.views
        })
    },
    editVisualization(row) {
      this.visualizationToEdit = row
      this.visualizationModalActive = true
    },
    deleteVisualization(row, index) {
      this.$buefy.dialog.confirm({
        title: 'Deleting Visualization',
        message: `Are you sure you want to delete <b>visualization ${row.name}</b>? This action cannot be undone.`,
        confirmText: 'Delete Visualization',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          var body = {
            uuid: row.uuid
          }
          fetch('/api/view/delete', {
            method: 'post',
            body: JSON.stringify(body)
          })
            .then(response => response)
            .then(data => {
              if (data.status === 200) {
                this.visualizations.splice(index, 1)
                this.$buefy.toast.open({ message: "Visualization Deleted", position: "is-top", type: "is-success" })
                setTimeout(this.fetchData, 1500)
              }
            })
        }
      })
    },
    addedVisualization(view) {
      this.visualizations.push(view)
      setTimeout(this.fetchData, 1500)
    },
    editedVisualization(view) {
      setTimeout(this.fetchData, 1500)
    }
  },
  components: {
    AddEditVisualization
  }
};
</script>

<style scoped lang="scss">
.type-button {
  width: 150px;
}
.button {
  margin-top: 20px;
  margin-bottom: 20px;
}
</style>
