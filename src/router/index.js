import { createRouter, createWebHistory } from 'vue-router'
import MailList from '../views/MailBoxView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: MailList,
    },
    {
      path: '/compose',
      name: 'compose',
      component: () => import('../views/ComposeMail.vue'),
    },
    {
      path: '/view/:id',
      name: 'view',
      component: () => import('../views/ViewMail.vue'),
    },
    {
      path: '/partner/add',
      name: 'Add Partner',
      component: () => import('../views/AddPartner.vue'),
    },
    {
      path: '/log',
      name: 'LogView',
      component: () => import('../views/LogViewPage.vue'),
    }
  ],
})

export default router