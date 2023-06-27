<template>
  <form action>
    <div class="modal-card" style="width: 500px">
      <header class="modal-card-head">
        <p class="modal-card-title">Edit</p>
      </header>
      <section class="modal-card-body">
        <b-field label="Name">
          <b-input type="text" v-model="passedObject.name" placeholder="Name..." required></b-input>
        </b-field>
        <b-field label="Email">
          <b-input type="text" v-model="passedObject.uuid" placeholder="Email..." required></b-input>
        </b-field>
        <b-field label="Role">
          <b-select expanded v-model="passedObject.class" placeholder="Role..." :disabled="Role.Standard == loggedIn.role">
            <option v-for="role in roles" :value="role" :key="role">{{ role }}</option>
          </b-select>
        </b-field>
        <b-field label="Activated">
          <b-select expanded v-model="passedObject.activated" placeholder="Activated..." :disabled="Role.Standard == loggedIn.role">
            <option v-for="type in [true,false]" :value="type" :key="type">{{ type }}</option>
          </b-select>
        </b-field>
      </section>
      <footer class="modal-card-foot" style="display: flex;justify-content: flex-end">
        <b-button class="button" type="button" @click="cancel">Close</b-button>
        <b-button type="is-primary" @click="save">Save</b-button>
      </footer>
    </div>
  </form>
</template>

<script>
import { Role } from '@/_helpers/role.js';
import format from '@/_helpers/snackbar.js';

export default {
  props: {
    passedObject: {}
  },
  data() {
    return {
      panelistList: [],
      roles: ["admin", "standard"],
      uuid: this.passedObject.uuid,
      loggedIn: {
        role: undefined,
        uuid: undefined
      },
      Role: Role
    };
  },
  mounted() {
    this.loggedIn.role = JSON.parse(localStorage.getItem('User')).class;
    this.loggedIn.uuid = JSON.parse(localStorage.getItem('User')).uuid;
  },
  methods: {
    cancel() {
      this.$emit('fetchData', this.user);
      this.$parent.close();
    },
    save() {
      fetch(
        process.env.VUE_APP_ENDPOINT + "user/update?uuid=" + this.uuid,
        {
          method: "post",
          body: JSON.stringify(this.passedObject)
        }
      )
        .then(response => {
          if (response.ok) {
            return Promise.all([response.ok, response.json()]);
          } else {
            return Promise.all([response.ok, response.text()]);
          }
        })
        .then(data => {
          if (data[0] === false) {
            this.$buefy.snackbar.open();
          } else {
            this.$emit('editedUser', this.user)
            this.$parent.close();
          }
        });
    }
  },
  components: {}
};
</script>
