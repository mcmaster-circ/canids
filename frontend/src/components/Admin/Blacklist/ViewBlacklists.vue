<template>
  <div class="hello">
    <h1 style="font-size:50px;margin-bottom:20px">Blacklists</h1>
    <section class="info-tiles" style="margin-top:50px">
      <article class="tile is-child box">
        <b-button
          class="is-centered"
          type="is-primary"
          style="float:right"
          @click="addBlacklist"
        >Add Blacklist</b-button>
        <b-table :data="blacklists" ref="table">
          <template slot-scope="props">
            <b-table-column label="Name" field="name" sortable searchable>{{props.row.name}}</b-table-column>
            <b-table-column label="Url" field="url" sortable searchable>{{props.row.url}}</b-table-column>
            <b-table-column field="actions">
              <b-dropdown
                aria-role="list"
                class="is-pulled-right"
                position="is-bottom-left"
                style="cursor:pointer"
              >
                <b-icon icon="dots-vertical" slot="trigger"></b-icon>
                <b-dropdown-item
                  aria-role="listitem"
                  @click="editBlacklist(props.row)"
                >
                  <b-icon icon="pencil"></b-icon>
                  <span style="margin-left:15px;font-size:15px">Edit</span>
                </b-dropdown-item>
                <b-dropdown-item aria-role="listitem" @click="deleteBlacklist(props.row, props.index)">
                  <b-icon icon="delete"></b-icon>
                  <span style="margin-left:15px;font-size:15px">Delete</span>
                </b-dropdown-item>
              </b-dropdown>
            </b-table-column>
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
      </article>
    </section>
    <b-modal :active.sync="blacklistModalActive" has-modal-card>
      <AddBlacklist :blacklistToEdit="blacklistToEdit" @addedBlacklist="addedBlacklist" @editedBlacklist="editedBlacklist"></AddBlacklist>
    </b-modal>
  </div>
</template>

<script>
import AddBlacklist from "./AddBlacklist.vue";

export default {
  name: "ViewBlacklists",
  props: {
    type: String
  },
  data() {
    return {
      blacklists: [],
      blacklistToEdit: {},
      blacklistModalActive: false,
      uuid: undefined
    };
  },
  mounted() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      fetch('/api/blacklist/list')
        .then(response => {
          return response.json()
        })
        .then(json => {
          this.blacklists = json.blacklists
        })
    },
    editBlacklist(row) {
      this.blacklistToEdit = row
      this.blacklistModalActive = true
    },
    addBlacklist() {
      this.blacklistToEdit = {}
      this.blacklistModalActive = true
    },
    addedBlacklist(blacklist) {
      this.$buefy.toast.open({ message: "Blacklist added", position: "is-top", type: "is-success" })
      this.blacklists.push(blacklist)
      setTimeout(this.fetchData, 1500)
    },
    editedBlacklist() {
      this.$buefy.toast.open({ message: "Blacklist edited", position: "is-top", type: "is-success" })
      setTimeout(this.fetchData, 1500)
    },
    changeType(type) {
      this.graphType = type;
    },
    deleteBlacklist(blacklist, index) {
      this.$buefy.dialog.confirm({
        title: 'Deleting Blacklist',
        message: `Are you sure you want to delete <b>Blacklist ${blacklist.name}</b>? This action cannot be undone.`,
        confirmText: 'Delete Blacklist',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          const body = {
            uuid: blacklist.uuid
          };
          fetch(process.env.VUE_APP_ENDPOINT + "blacklist/delete", {
            method: "post",
            body: JSON.stringify(body)
          })
            .then(response => {
              if (response.ok) {
                return Promise.all([response.ok, response.json()]);
              } else {
                return Promise.all([response.ok, response.text()]);
              }
            })
            .then(response => {
              if (response[0] === true) {
                this.$buefy.toast.open({ message: "Blacklist Deleted", type: "is-success", position: "is-top" })
                this.blacklists.splice(index, 1)
                setTimeout(this.fetchData, 1500);
              } else {
                this.$buefy.toast.open({ message: response[1], position: "is-top", type: "is-danger" })
              }
            });
        }
      })
    }
  },
  components: { AddBlacklist }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.type-button {
  width: 150px;
}
.button {
  margin-top: 20px;
  margin-bottom: 20px;
}
</style>
