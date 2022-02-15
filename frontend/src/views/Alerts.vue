<template>
  <div class="home" style="height:100vh">
    <div class="columns"  style="height:100%">
      <div class="column">
        <AlertsComponent />
      </div>
    </div>
  </div>
</template>

<script>
// @ is an alias to /src
import AlertsComponent from '@/components/Alerts/AlertsComponent.vue'

export default {
  name: 'Home',
  components: {
    AlertsComponent
  },
  data () {
    return {
      isActive: true, // eslint-disable-line
    }
  },
  methods: {
  }
}
</script>
