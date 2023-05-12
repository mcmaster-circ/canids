<template>
  <div class="home" style="height:100vh">
    <div class="columns"  style="height:100%">
      <div class="column is-one-fifth" style="border-right: 1px solid black;">
        <b-menu style="font-size:22px; margin-top:10px">
          <b-menu-list label="Admin Actions" >
            <b-menu-item :active="isActive" icon="account-multiple" label="View Users"  style="text-align:left;font-size:20px" @click="changeTab('view')"></b-menu-item>
            <b-menu-item icon="account-group" label="View Groups"  style="text-align:left;font-size:20px" @click="changeTab('groups')"></b-menu-item>
            <b-menu-item icon="plus-circle" label="View Visualizations"  style="text-align:left;font-size:20px" @click="changeTab('vis')"></b-menu-item>
            <b-menu-item v-if="role == Roles.SuperUser" icon="alarm-light" label="View Blacklists"  style="text-align:left;font-size:20px" @click="changeTab('blacklists')"></b-menu-item>
          </b-menu-list>
        </b-menu>
      </div>
      <div class="column">
        <ViewVisualizations v-if="selectedTab == 'vis'" type="Test" />
        <ViewUsers  v-if="selectedTab == 'view'"/>
        <ViewGroups  v-if="selectedTab == 'groups'"/>
        <ViewBlacklists  v-if="selectedTab == 'blacklists' && role == Roles.SuperUser"/>
      </div>
    </div>
  </div>
</template>

<script>
// @ is an alias to /src
import { Role } from "@/_helpers/role.js"
import ViewUsers from '@/components/Admin/Users/ViewUsers.vue'
import ViewGroups from '@/components/Admin/Groups/ViewGroups.vue'
import ViewVisualizations from '@/components/Admin/Visualization/ViewVisualizations.vue'
import ViewBlacklists from '@/components/Admin/Blacklist/ViewBlacklists.vue'

export default {
  name: 'Admin',
  components: {
    ViewUsers,
    ViewGroups,
    ViewVisualizations,
    ViewBlacklists
  },
  data () {
    return {
      isActive: true, // eslint-disable-line
      selectedTab: 'view',
      role: undefined,
      Roles: Role
    }
  },
  mounted () {
    this.role = JSON.parse(localStorage.getItem("User")).class;
  },
  methods: {
    changeTab (type) {
      this.selectedTab = type
    }
  }
}
</script>
