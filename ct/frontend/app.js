import Vue from 'vue'
import VueRouter from 'vue-router'
import PortfolioManager from './vue/portfolio-manager.vue'
import Login from './vue/login.vue'

const routes = [
    { path: '/', name: 'main', component: PortfolioManager },
    { path: '/login', name: 'login', component: Login }
]
Vue.use(VueRouter)
const router = new VueRouter({
    //mode: 'history', //make sure traefik know how to handle /login
    routes
})

const app = new Vue({
    router
}).$mount('#app')