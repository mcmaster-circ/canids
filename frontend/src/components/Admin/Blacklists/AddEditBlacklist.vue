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
        <b-button class="button" type="button" @click="cancel">Close</b-button>
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
    cancel() {
      this.$emit('fetchData', this.blacklist);
      this.$parent.close();
    },
    save() {
      if (this.blacklistToEdit.name) {
        fetch("/api/blacklist/update", {
          method: "post",
          body: JSON.stringify(this.blacklist)
        })
          .then(response => response)
          .then(data => {
            if (data.status === 200) {
              this.$parent.close()
              this.$emit('editedBlacklist')
            }
          });
      } else {
        fetch("/api/blacklist/add", {
          method: "post",
          body: JSON.stringify(this.blacklist)
        })
          .then(response => response)
          .then(data => {
            if (data.status === 200) {
              this.$parent.close()
              this.$emit('addedBlacklist', this.blacklist)
            }
          });
      }
    }
  },
  components: {}
};
</script>
