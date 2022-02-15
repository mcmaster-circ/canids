<template>
  <form action>
    <div class="modal-card" style="width: 500px">
      <header class="modal-card-head">
        <p class="modal-card-title">Add</p>
      </header>
      <section class="modal-card-body">
        <b-field label="Name">
          <b-input type="text" v-model="group.name" placeholder="Name..." required></b-input>
        </b-field>
        <b-field label="Authorized">
          <b-table
          :data="allAuthorized"
          :hoverable="true"
          @click="group.authorized.includes(row) ? group.authorized.splice(group.authorized.indexOf(row), 1) : group.authorized.push(row)"
          :checked-rows.sync="group.authorized"
          checkable
          :checkbox-position="'left'">

          <template slot-scope="props">
              <b-table-column field="name" label="Name" style="cursor:pointer;">
                  {{ props.row }}
              </b-table-column>
          </template>
          </b-table>
        </b-field>
        <b-field label="Add new Authorized">
          <b-input type="text" v-model="newAuthorized" placeholder="Create new..."></b-input>
        </b-field>
        <b-button type="is-primary" style="width: 100%" @click="createAuthorized"><b-icon icon="plus-circle-outline"></b-icon><span>Add</span></b-button>
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
  props: ['groupToEdit'],
  data() {
    return {
      newAuthorized: "",
      allAuthorized: [],
      group: {
        name: null,
        authorized: []
      }
    };
  },
  mounted() {
    this.fetchData()
    if (this.groupToEdit.name) {
      this.group = this.groupToEdit
    }
  },
  methods: {
    fetchData() {
      fetch('/api/assets/list')
        .then(response => {
          return response.json()
        })
        .then(json => {
          this.allAuthorized = json.assets
        })
    },
    save() {
      if (this.groupToEdit.name) {
        fetch("/api/group/update", {
          method: "post",
          body: JSON.stringify(this.group)
        })
          .then(response => response)
          .then(data => {
            if (data.status === 200) {
              this.$parent.close()
              this.$emit('editedGroup')
            }
          });
      } else {
        fetch("/api/group/add", {
          method: "post",
          body: JSON.stringify(this.group)
        })
          .then(response => response)
          .then(data => {
            if (data.status === 200) {
              this.$parent.close()
              this.$emit('addedGroup', this.group)
            }
          });
      }
    },
    createAuthorized() {
      if (this.newAuthorized !== "") {
        this.allAuthorized.push(this.newAuthorized)
        this.group.authorized.push(this.newAuthorized)
        this.newAuthorized = ""
      }
    }
  },
  components: {}
};
</script>
