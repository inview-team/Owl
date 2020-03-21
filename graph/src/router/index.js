import Vue from 'vue';
import Router from 'vue-router';
import Graph from '@/components/Graph';
import Settings from '@/components/Settings';
import Logs from '@/components/Logs';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/graph',
      name: 'Graph',
      component: Graph,
    },
    {
      path: '/settings',
      name: 'Settings',
      component: Settings
    },
    {
      path: '/logs',
      name: 'Logs',
      component: Logs
    }
  ],
});
