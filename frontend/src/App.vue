<template>
  <div>
    <Navbar/>
      <div id="app" style="background: #ECF0F3">
        <router-view/>
      </div>
  </div>
</template>

<script>
import Navbar from '@/components/Navbar.vue'

export default {
  methods: {
  },
  components: {
    Navbar
  }
}
</script>

<style lang="scss">

// Import Bulma's core
@import "~bulma/sass/utilities/functions";

$primary: #be3e24;
$primary-invert: findColorInvert($primary);

html, body {
  font-family: 'Open Sans', Verdana, Arial, Helvetica, sans-serif;
  font-size: 18px;
  line-height: 1.5;
  height: 100%;
  background: #ECF0F3;
  background-size: cover;
  background-repeat: repeat;
}

// Import Bulma and Buefy styles
@import "~bulma";
@import "~buefy/src/scss/buefy";

</style>

<style>
@import 'assets/style.css';
</style>
