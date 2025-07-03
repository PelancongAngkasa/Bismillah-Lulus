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
      path: '/partner',
      name: 'PartnerList',
      component: () => import('../views/ViewPartnerUser.vue'),
    },
    {
      path: '/admin/partner/add',
      name: 'Add Partner',
      component: () => import('../views/AddPartner.vue'),
    },
    {
      path: '/admin/log',
      name: 'LogView',
      component: () => import('../views/LogViewPage.vue'),
    },
    {
      path: '/admin/pmode/edit',
      name: 'PmodeEdit',
      component: () => import('../views/PmodePage.vue')
    }
  ],
})

export default router