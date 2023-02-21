<template>
  <nav class="navbar is-dark" id="root">
    <div class="container">
      <!-- Logo -->
      <div class="navbar-brand">
        <div id="logo-container" to="" class="navbar-item brand-text">
          <img id="logo" src="../assets/logo.png">
        </div>
        <!-- Mobile menu -->
        <div class="navbar-burger burger is-centered" data-target="navMenu"
          v-bind:class="{ 'is-active': showNav }"
          v-on:click="toggle">
          <span></span>
          <span></span>
          <span></span>
        </div>
      </div>
      <!-- Menu links -->
      <div id="navMenu" class="navbar-menu" v-bind:class="{ 'is-active': showNav }">
        <div class="navbar-start" v-on:click="condToggle">
          <router-link to="/" class="navbar-item"> Dashboard </router-link>
          <router-link to="/alarms" class="navbar-item"> Alarms </router-link>
          <router-link to="/admin" class="navbar-item" v-if="user.class === 'superuser' || user.class === 'admin'"> Admin </router-link>
          <a href="/logout" class="navbar-item"> Logout </a>
          <div class="navbar-item">Currently logged in as: {{ user.name }}</div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
export default {
  name: 'Navbar',
  data () {
    return {
      showNav: false,
      user: {}
    }
  },
  mounted () {
    fetch('/api/user/info')
      .then(response => {
        return Promise.all([response.ok, response.json()]);
      })
      .then(data => {
        if (data[0] === true) {
          localStorage.setItem('User', JSON.stringify(data[1]));
          this.user = data[1]
        } else {
          window.location.replace('/login')
        }
      })
  },
  methods: {
    toggle () {
      // toggle between displaying and hiding mobile menu
      this.showNav = !this.showNav
    },
    condToggle () {
      // if user clicks on mobile menu, hide the menu
      if (this.showNav) {
        this.toggle()
      }
    }
  }
}
</script>

<style>
#white {
  color: white;
}
#root {
  background: linear-gradient(to right, rgba(169, 21, 26, 255), rgba(221, 123, 50, 1));
}
nav.navbar {
  margin-bottom: 1rem;
}
.navbar-item.brand-text {
  font-weight: 300;
}
.navbar-item, .navbar-link {
  font-size: 15px;
  font-weight: 700;
}
#logo {
  max-height: 4rem;
}

</style>
