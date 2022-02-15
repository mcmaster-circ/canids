<template>
  <div class="hello">
    <h1 style="font-size:50px;margin-bottom:20px">Groups</h1>
    <section class="info-tiles" style="margin-top:50px">
      <article class="tile is-child box">
        <b-button
          class="is-centered"
          type="is-primary"
          style="float:right"
          @click="addGroup"
          :disabled="role == Roles.SuperUser"
        >Create Group</b-button>
        <b-table :data="groups" ref="table">
          <template slot-scope="props">
            <b-table-column label="Name" field="name" sortable searchable>{{props.row.name}}</b-table-column>
            <b-table-column label="Uuid" field="uuid" sortable searchable>{{props.row.uuid}}</b-table-column>
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
                  @click="editGroup(props.row)"
                >
                  <b-icon icon="pencil"></b-icon>
                  <span style="margin-left:15px;font-size:15px">Edit</span>
                </b-dropdown-item>
                <b-dropdown-item aria-role="listitem" @click="deleteUser(props.row, props.index)">
                  <b-icon icon="delete"></b-icon>
                  <span
                    style="margin-left:15px;font-size:15px "
                    :disabled="role != Roles.SuperUser"
                  >Delete</span>
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
    <b-modal :active.sync="groupModalActive" has-modal-card>
      <AddGroup :groupToEdit="groupToEdit" @addedGroup="addedGroup" @editedGroup="editedGroup"></AddGroup>
    </b-modal>
  </div>
</template>

<script>
import { Role } from "@/_helpers/role.js";
import AddGroup from "./AddGroup";

export default {
  name: "AddVisualization",
  props: {
    type: String
  },
  data() {
    return {
      groups: [],
      groupToEdit: {},
      groupModalActive: false,
      dataToEdit: undefined,
      role: undefined,
      uuid: undefined,
      Roles: Role
    };
  },
  mounted() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      fetch('/api/group/list')
        .then(response => {
          return response.json()
        })
        .then(json => {
          var array = []
          array.push(json.current)
          if (json.others !== null) {
            json.others.forEach(group => {
              array.push(group)
            })
          }
          this.groups = array
        })
    },
    editGroup(row) {
      this.groupToEdit = row
      this.groupModalActive = true
    },
    addGroup() {
      this.groupToEdit = {}
      this.groupModalActive = true
    },
    addedGroup(group) {
      this.$buefy.toast.open({ message: "Group added", position: "is-top", type: "is-success" })
      this.groups.push(group)
      setTimeout(this.fetchData, 1500)
    },
    editedGroup() {
      this.$buefy.toast.open({ message: "Group edited", position: "is-top", type: "is-success" })
      setTimeout(this.fetchData, 1500)
    },
    changeType(type) {
      this.graphType = type;
    },
    deleteUser(group, index) {
      this.$buefy.dialog.confirm({
        title: 'Deleting Group',
        message: `Are you sure you want to delete <b>Group ${group.name}</b>? This action cannot be undone.`,
        confirmText: 'Delete Group',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          const body = {
            uuid: group.uuid
          };
          fetch(process.env.VUE_APP_ENDPOINT + "group/delete", {
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
                this.$buefy.toast.open({ message: "Group Deleted", type: "is-success", position: "is-top" })
                this.groups.splice(index, 1)
                setTimeout(this.fetchData, 1500);
              } else {
                this.$buefy.toast.open({ message: response[1], position: "is-top", type: "is-danger" })
              }
            });
        }
      })
    }
  },
  components: { AddGroup }
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
