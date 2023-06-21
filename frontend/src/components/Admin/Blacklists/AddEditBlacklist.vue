<template>
  <form action>
    <div class="modal-card" style="width: 500px">
      <header class="modal-card-head">
        <p class="modal-card-title">Add</p>
      </header>
      <section class="modal-card-body">
        <b-field label="Name">
          <b-input type="text" v-model="blacklist.name" placeholder="Name..." required></b-input>
        </b-field>
        <b-field label="URL">
          <b-input type="text" v-model="blacklist.url" placeholder="URL..." required></b-input>
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
export default {
  props: ['blacklistToEdit'],
  data() {
    return {
      blacklist: {
        name: null,
        url: null
      }
    };
  },
  mounted() {
    if (this.blacklistToEdit.name) {
      this.blacklist = this.blacklistToEdit
    }
  },
  methods: {
    save() {
      if (this.blacklistToEdit.name) {
        fetch("/api/blacklist/update", {
          method: "post",
          body: JSON.stringify(this.blacklist)
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
              this.$parent.close()
              this.$buefy.toast.open({ message: "Blacklist Updated", position: "is-top", type: "is-success" })
              this.$emit('editedBlacklist', this.blacklist)
            }
          });
      } else {
        fetch("/api/blacklist/add", {
          method: "post",
          body: JSON.stringify(this.blacklist)
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
              this.$parent.close()
              this.$buefy.toast.open({ message: "Blacklist Added", position: "is-top", type: "is-success" })
              this.$emit('addedBlacklist', this.blacklist)
            }
          });
      }
    }
  },
  components: {}
};
</script>
