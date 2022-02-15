import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/admin',
    name: 'Admin',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/Admin.vue')
  },
  // {
  //   path: '/alarms', TODO remove or uncomment this
  //   name: 'Alarms',
  //   component: () => import(/* webpackChunkName: "about" */ '../views/Alarms.vue')
  // },
  {
    path: '/dashboard/edit',
    name: 'DashboardEdit',
    component: () => import(/* webpackChunkName: "about" */ '../components/DashboardView/DashboardEdit.vue')
  },
  {
    path: '/dashboard/list',
    name: 'DashboardList',
    component: () => import(/* webpackChunkName: "about" */ '../components/DashboardView/DashboardList.vue')
  }
]

const router = new VueRouter({
  routes
})

export default router
