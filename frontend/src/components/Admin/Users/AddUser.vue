<template>
  <form action>
    <div class="modal-card" style="width: 500px">
      <header class="modal-card-head">
        <p class="modal-card-title">Add</p>
      </header>
      <section class="modal-card-body">
        <b-field label="Name">
          <b-input type="text" v-model="user.name" placeholder="Name..." required></b-input>
        </b-field>
        <b-field label="Email">
          <b-input type="text" v-model="user.uuid" placeholder="Email..." required></b-input>
        </b-field>
        <b-field label="Group">
          <b-select
            :disabled="role == Roles.Admin"
            expanded
            v-model="user.group"
            placeholder="Group..."
          >
            <option v-for="group in groups" :value="group.uuid" :key="group.uuid">{{ group.name }}</option>
          </b-select>
        </b-field>
        <b-field label="Class">
          <b-select expanded v-model="user.class" placeholder="Class...">
            <option v-for="role in roles" :value="role" :key="role">{{ role }}</option>
          </b-select>
        </b-field>
        <b-field label="Activated">
          <b-select expanded v-model="user.activated" placeholder="Activated...">
            <option v-for="type in [true,false]" :value="type" :key="type">{{ type }}</option>
          </b-select>
        </b-field>
      </section>
      <footer class="modal-card-foot" style="display: flex;justify-content: flex-end">
        <b-button class="button" type="button" @click="$parent.close()">Close</b-button>
        <b-button type="is-primary" @click="save">Save</b-button>
      </footer>
    </div>
  </form>
</template>

<script>
import { Role } from "@/_helpers/role.js";

export default {
  data() {
    return {
      panelistList: [],
      roles: ["admin", "superuser", "standard"],
      user: {
        name: null,
        uuid: null,
        class: null,
        group: null,
        activated: null
      },
      groups: [],
      role: undefined,
      Roles: Role
    };
  },
  mounted() {
    this.role = JSON.parse(localStorage.getItem("User")).class;
    fetch(process.env.VUE_APP_ENDPOINT + "group/list")
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
        this.groups.push({ uuid: data.current.uuid, name: data.current.name });
        if (data.others) {
          for (let i = 0; i < data.others.length; i++) {
            this.groups.push({
              uuid: data.others[i].uuid,
              name: data.others[i].name
            });
          }
        }
      });
  },
  methods: {
    save() {
      fetch(process.env.VUE_APP_ENDPOINT + "user/add", {
        method: "post",
        body: JSON.stringify(this.user)
      })
        .then(response => {
          if (response.ok) {
            return Promise.all([response.ok, response.json()]);
          } else {
            return Promise.all([response.ok, response.text()]);
          }
        })
        .then(data => {
          if (data[0] === false) {
            this.$buefy.snackbar.open(data[1]);
          } else {
            this.$emit('addedUser', this.user)
            this.$parent.close();
          }
        });
    }
  },
  components: {}
};
</script>
