<template>
  <div class="tile is-ancestor">
    <div class="tile is-child box">
      <h1>All Visualizations</h1>
      <b-table
        :data="allViews"
        :loading="isTableLoading"
        @dblclick="dblclick"
        :row-class="(row, index) => doesContain(row, index)">

        <template slot-scope="props">
            <b-table-column field="name" label="Name" style="cursor:pointer;">
                {{ props.row.name }}
            </b-table-column>

            <b-table-column field="Type" label="Type" style="cursor:pointer;">
                {{ props.row.class }}
            </b-table-column>
        </template>

        <template slot="empty">
          <section class="section">
              <div class="content has-text-grey has-text-centered">
                <p>
                  <b-icon
                    icon="emoticon-sad"
                    size="is-large">
                  </b-icon>
                </p>
                <p>Nothing here.</p>
              </div>
          </section>
        </template>
      </b-table>
      <button style="margin-top:20px;width:100%;" class="button is-primary" type="button" @click="$emit('saveViews')" v-if="role === 'superuser' || role === 'admin'"><b-icon icon="check"></b-icon><span>Save</span></button>
    </div>
  </div>
</template>

<script>

var menuIsHidden = false;
export default {
  props: ['selectedViews'],
  name: "Sidebar",
  components: {
  },
  data() {
    return {
      allViews: [],
      role: ""
    };
  },
  mounted () {
    this.fetchData()
    this.fetchUser()
  },
  methods: {
    fetchUser() {
      fetch('/api/user/info')
        .then(response => {
          return response.json()
        })
        .then(json => {
          this.role = json.class
        })
    },
    fetchData() {
      fetch('/api/view/list')
        .then(response => {
          return response.json()
        })
        .then(json => {
          this.allViews = json.views
        })
    },
    dblclick(row) {
      const contains = { exists: false, index: 0 }
      for (let i = 0; i < this.selectedViews.length; i++) {
        if (this.selectedViews[i].uuid === row.uuid) {
          contains.exists = true
          contains.index = i
        }
      }
      if (contains.exists) {
        this.selectedViews.splice(contains.index, 1)
      } else {
        this.selectedViews.push(row)
      }
    },
    doesContain(row, index) {
      let contains = false
      this.selectedViews.forEach(view => {
        if (row.uuid === view.uuid) {
          contains = true
        }
      });
      if (contains) {
        return 'is-info'
      } else {
        return null
      }
    }
  }
};
</script>

<style>
  tr.is-info {
    background: #be3e24;
    color: #fff;
  }
</style>
