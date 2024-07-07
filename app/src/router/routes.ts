import { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', redirect: '/live' },
      {
        path: 'creation',
        children: [
          {
            path: '',
            name: 'CreationList',
            component: () => import('pages/creation/list/CreationListPage.vue')
          },
          {
            path: 'information/:id',
            name: 'CreationInformation',
            component: () =>
              import('pages/creation/information/InformationPage.vue')
          }
        ]
      },
      {
        path: 'account',
        component: () => import('pages/account/AccountPage.vue'),
        children: [
          {
            path: '',
            component: () => import('pages/account/DouYinAccountPage.vue')
          }
        ]
      },
      {
        path: 'live',
        children: [
          {
            path: '',
            name: 'LivePage',
            component: () => import('pages/live/LivePage.vue')
          },
          {
            path: '/information/{id}',
            name: 'LiveInformation',
            component: () => import('pages/live/LiveInformation.vue'),
            children: [
              {
                path: '/dashboard',
                name: 'LiveDashboard',
                component: () => import('pages/live/information/LiveDashboard.vue')
              },
              {
                path: '/content',
                name: 'LiveContent',
                component: () => import('pages/live/information/LiveContent.vue')
              }
            ]
          }
        ]
      }
    ]
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue')
  }
]

export default routes
