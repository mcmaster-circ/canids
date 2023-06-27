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
import { format } from "@/_helpers/snackbar.js";

export default {
  data() {
    return {
      panelistList: [],
      roles: ["admin", "standard"],
      user: {
        name: null,
        uuid: null,
        class: null,
        activated: null
      },
      role: undefined,
      Roles: Role
    };
  },
  mounted() {
    this.role = JSON.parse(localStorage.getItem("User")).class;
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
            this.$buefy.snackbar.open(format(data))
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
