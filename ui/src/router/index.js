import Vue from 'vue'
import Router from 'vue-router'
import SearchBox from '@/components/SearchBox'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'SearchBox',
      component: SearchBox
    }
  ]
})
