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
        <b-field label="Group">
          <b-select
            :disabled="Role.Standard == loggedIn.role || Role.Admin == loggedIn.role"
            expanded
            v-model="passedObject.group"
            placeholder="Group..."
          >
            <option v-for="group in groups" :value="group.uuid" :key="group.uuid">{{ group.name }}</option>
          </b-select>
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
        <b-button class="button" type="button" @click="$parent.close()">Close</b-button>
        <b-button type="is-primary" @click="save">Save</b-button>
      </footer>
    </div>
  </form>
</template>

<script>
import { Role } from '@/_helpers/role.js';

export default {
  props: {
    passedObject: {}
  },
  data() {
    return {
      groups: [],
      panelistList: [],
      roles: ["admin", "superuser", "standard"],
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
        for (let i = 0; i < data.others.length; i++) {
          this.groups.push({
            uuid: data.others[i].uuid,
            name: data.others[i].name
          });
        }
      });
  },
  methods: {
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
            this.$buefy.snackbar.open(data[1]);
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
