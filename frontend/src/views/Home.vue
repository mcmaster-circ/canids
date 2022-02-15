<template>
  <div class="home" style="padding:10px">
    <Dashboard/>
  </div>
</template>

<script>
import Dashboard from '@/components/DashboardView/Dashboard.vue'

export default {
  name: 'Home',
  components: {
    Dashboard
  },
  created () {
    fetch(process.env.VUE_APP_ENDPOINT + 'user/info')
      .then(response => {
        if (response.ok) {
          return Promise.all([response.ok, response.json()])
        } else {
          return Promise.all([response.ok, response.text()])
        }
      })
      .then((response) => {
        const status = response[0]
        const data = response[1]
        if (!status) {
          this.$buefy.snackbar.open(data)
        }
        localStorage.setItem('User', JSON.stringify(data));
      })
  }
}
</script>
