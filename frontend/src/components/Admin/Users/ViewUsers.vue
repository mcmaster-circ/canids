<template>
  <div class="hello">
    <h1 style="font-size:50px;margin-bottom:20px">Users</h1>
    <section class="info-tiles" style="margin-top:50px">
      <article class="tile is-child box">
        <b-button
          class="is-centered"
          type="is-primary"
          style="float:right"
          @click="openAddModal"
          :disabled="role == Roles.Standard"
        >Create User</b-button>
        <b-table :data="users" ref="table">
          <template slot-scope="props">
            <b-table-column label="Name" field="name" sortable searchable>{{props.row.name}}</b-table-column>
            <b-table-column label="Email" field="uuid" sortable searchable>{{props.row.uuid}}</b-table-column>
            <b-table-column label="Role" field="class" sortable searchable>{{props.row.class}}</b-table-column>
            <b-table-column
              label="Activated"
              field="activated"
              sortable
              searchable
            >{{props.row.activated}}</b-table-column>

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
                  @click="dataToEdit = props.row;editUserData = true;"
                >
                  <b-icon icon="pencil"></b-icon>
                  <span
                    style="margin-left:15px;font-size:15px"
                    :disabled="role == Roles.Standard && props.row.uuid!=uuid"
                  >Edit</span>
                </b-dropdown-item>
                <b-dropdown-item
                  aria-role="listitem"
                  @click="resetUser(props.row.uuid)"
                  :disabled="role == Roles.Standard"
                >
                  <b-icon icon="lock-reset"></b-icon>
                  <span style="margin-left:15px;font-size:15px">Reset Password</span>
                </b-dropdown-item>
                <b-dropdown-item
                  aria-role="listitem"
                  @click="deleteUser(props.row, props.index)"
                  :disabled="role == Roles.Standard"
                >
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
    <b-modal :active.sync="editUserData" has-modal-card>
      <EditUser :passedObject="dataToEdit" @editedUser="editedUser"></EditUser>
    </b-modal>
    <b-modal :active.sync="addUserModal" has-modal-card>
      <AddUser @addedUser="addedUser" @editedUser="editedUser"></AddUser>
    </b-modal>
  </div>
</template>

<script>
import { Role } from "@/_helpers/role.js";
import EditUser from "./EditUser";
import AddUser from "./AddUser";

export default {
  name: "ViewUsers",
  props: {
    type: String
  },
  data() {
    return {
      users: [],
      editUserData: false,
      addUserModal: false,
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
      this.role = JSON.parse(localStorage.getItem("User")).class;
      this.role = JSON.parse(localStorage.getItem("User")).uuid;
      fetch(process.env.VUE_APP_ENDPOINT + "user/list")
        .then(response => {
          if (response.ok) {
            return response.json()
          } else {
            this.$buefy.toast.open({ message: response.text(), position: "is-top", type: "is-danger" })
            return Promise.all([response.ok, response.text()]);
          }
        })
        .then(json => {
          this.users = json.users;
        });
    },
    openAddModal() {
      const modal = this.$buefy.modal.open({
        parent: this,
        component: AddUser,
        hasModalCard: true
      });
      modal.$on("close", () => {
        this.fetchData();
      });
    },
    changeType(type) {
      this.graphType = type;
    },
    resetUser(email) {
      const user = {
        uuid: email
      };
      fetch(process.env.VUE_APP_ENDPOINT + "user/resetPass", {
        method: "post",
        body: JSON.stringify(user)
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
          this.$buefy.snackbar.open(data);
          setTimeout(this.fetchData(), 3000);
        });
    },
    deleteUser(row, index) {
      this.$buefy.dialog.confirm({
        title: 'Deleting User',
        message: `Are you sure you want to delete <b>user ${row.name}</b>? This action cannot be undone.`,
        confirmText: 'Delete User',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          const user = {
            uuid: row.uuid
          };
          fetch(process.env.VUE_APP_ENDPOINT + "user/delete", {
            method: "post",
            body: JSON.stringify(user)
          })
            .then(response => {
              if (response.ok) {
                return Promise.all([response.ok, response.json()]);
              } else if (response.status === 403) {
                this.$buefy.snackbar.open({ message: "User cannot delete their own account.", position: "is-top", type: "is-danger" })
                return Promise.all([response.ok, response.text()]);
              } else {
                this.$buefy.snackbar.open({ message: "An error occurred.", position: "is-top", type: "is-danger" })
                return Promise.all([response.ok, response.text()]);
              }
            })
            .then(([ok, data]) => {
              if (ok && data.success) {
                this.$buefy.toast.open({ message: "User Deleted", position: "is-top", type: "is-success" })
                this.users.splice(index, 1)
              }
              setTimeout(this.fetchData, 1500);
            });
        }
      })
    },
    addedUser(user) {
      this.$buefy.toast.open({ message: "User added", position: "is-top", type: "is-success" })
      this.users.push(user)
      setTimeout(this.fetchData, 1500)
    },
    editedUser() {
      this.$buefy.toast.open({ message: "User edited", position: "is-top", type: "is-success" })
      setTimeout(this.fetchData, 1500)
    }
  },
  components: { EditUser, AddUser }
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
