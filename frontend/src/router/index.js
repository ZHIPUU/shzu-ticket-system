import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'
import AdminLayout from '../layouts/AdminLayout.vue'
import TicketList from '../views/TicketList.vue'
import TicketDetail from '../views/TicketDetail.vue'
import Settings from '../views/Settings.vue'
import Login from '../views/Login.vue'
import UserManagement from '../views/UserManagement.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', name: 'login', component: Login, meta: { public: true } },
    {
      path: '/',
      component: AdminLayout,
      children: [
        { path: '', redirect: '/tickets' },
        { path: 'tickets', name: 'list', component: TicketList },
        { path: 'tickets/:id', name: 'detail', component: TicketDetail, props: true },
        { path: 'users', name: 'users', component: UserManagement, meta: { adminOnly: true } },
        { path: 'settings', name: 'settings', component: Settings },
      ],
    },
  ],
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  // 公开页面
  if (to.meta.public) {
    if (to.name === 'login' && userStore.isLoggedIn) {
      return next('/tickets')
    }
    return next()
  }

  // 需要登录
  if (!userStore.isLoggedIn) {
    return next({ path: '/login', query: { redirect: to.fullPath } })
  }

  // admin 专属
  if (to.meta.adminOnly && !userStore.isAdmin) {
    return next('/tickets')
  }

  // 首次登录强制改密码（除了 settings 和 login）
  if (userStore.mustChangePassword && to.name !== 'settings') {
    return next('/settings')
  }

  next()
})

export default router
