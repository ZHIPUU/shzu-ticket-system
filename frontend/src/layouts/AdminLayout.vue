<template>
  <div class="layout">
    <!-- ════════ 桌面端侧边栏 ════════ -->
    <aside v-if="!isMobile" class="sidebar" :class="{ mini }">
      <div class="side-brand" @click="$router.push('/tickets')">
        <div class="side-logo"><TicketsPlane :size="20" :stroke-width="2" /></div>
        <transition name="fade-t">
          <div v-if="!mini" class="side-brand-text">
            <div class="sb-name">石小易工单</div>
            <div class="sb-sub">迎新智能体后台</div>
          </div>
        </transition>
      </div>

      <nav class="side-nav">
        <div v-if="!mini" class="nav-caption">工作台</div>
        <router-link
          v-for="item in visibleMenu"
          :key="item.path"
          :to="item.path"
          class="nav-item"
          :class="{ active: isActive(item.path) }"
          :title="mini ? item.label : ''"
        >
          <span class="nav-active-bar" />
          <component :is="item.icon" :size="19" :stroke-width="2" class="nav-icon" />
          <span v-if="!mini" class="nav-label">{{ item.label }}</span>
          <span v-if="!mini && item.badge > 0" class="nav-badge tnum">{{ item.badge > 99 ? '99+' : item.badge }}</span>
        </router-link>
      </nav>

      <div class="side-foot">
        <button class="foot-btn" :title="isDark ? '切换亮色' : '切换暗色'" @click="toggleTheme">
          <component :is="isDark ? Sun : Moon" :size="17" :stroke-width="2" />
          <span v-if="!mini" class="foot-label">{{ isDark ? '亮色模式' : '暗色模式' }}</span>
        </button>
        <button class="foot-btn" :title="mini ? '展开侧栏' : '收起侧栏'" @click="toggleMini">
          <ChevronsLeft v-if="!mini" :size="17" :stroke-width="2" />
          <ChevronsRight v-else :size="17" :stroke-width="2" />
        </button>
      </div>
    </aside>

    <!-- ════════ 主区域 ════════ -->
    <div class="main-col">
      <!-- 顶栏 -->
      <header class="topbar">
        <div class="topbar-left">
          <div v-if="isMobile" class="m-logo" @click="$router.push('/tickets')">
            <TicketsPlane :size="17" :stroke-width="2.2" />
          </div>
          <h2 class="topbar-title">{{ pageTitle }}</h2>
        </div>
        <div class="topbar-right">
          <button v-if="isMobile" class="icon-btn" aria-label="切换主题" @click="toggleTheme">
            <component :is="isDark ? Sun : Moon" :size="18" :stroke-width="2" />
          </button>

          <UiDropdown :items="userMenuItems" placement="end" @select="onUserCommand">
            <template #trigger>
              <div class="user-chip">
                <div class="avatar" :data-role="userStore.user?.role">{{ avatarLetter }}</div>
                <div v-if="!isMobile" class="user-meta">
                  <div class="user-name">{{ userStore.user?.display_name || userStore.user?.username }}</div>
                  <div class="user-role">{{ userStore.user?.role === 'admin' ? '管理员' : '工作人员' }}</div>
                </div>
                <ChevronDown v-if="!isMobile" :size="14" :stroke-width="2" class="user-chev" />
              </div>
            </template>
          </UiDropdown>
        </div>
      </header>

      <!-- 页面内容 -->
      <main class="page-scroll">
        <div class="page-inner" :class="{ 'has-bottomnav': isMobile }">
          <router-view v-slot="{ Component }">
            <transition name="page" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </main>
    </div>

    <!-- ════════ 移动端底部导航 ════════ -->
    <nav v-if="isMobile" class="bottomnav">
      <router-link
        v-for="item in visibleMenu"
        :key="item.path"
        :to="item.path"
        class="bn-item"
        :class="{ active: isActive(item.path) }"
      >
        <span class="bn-icon-wrap">
          <component :is="item.icon" :size="21" :stroke-width="2" />
          <span v-if="item.badge > 0" class="bn-badge tnum">{{ item.badge > 99 ? '99+' : item.badge }}</span>
        </span>
        <span class="bn-label">{{ item.shortLabel || item.label }}</span>
      </router-link>
    </nav>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  TicketsPlane, Users, Settings, Sun, Moon,
  ChevronsLeft, ChevronsRight, ChevronDown, LogOut, UserRound,
} from '@lucide/vue'
import UiDropdown from '../ui/UiDropdown.vue'
import { confirmDialog } from '../ui/confirm'
import { toast } from '../ui/toast'
import { useTheme } from '../composables/useTheme'
import { useIsMobile } from '../composables/useMediaQuery'
import { useUserStore } from '../stores/user'
import { listTickets } from '../api/ticket'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const { isDark, toggle: toggleTheme } = useTheme()
const isMobile = useIsMobile()

const mini = ref(localStorage.getItem('ticket_sidebar_mini') === '1')
const toggleMini = () => {
  mini.value = !mini.value
  localStorage.setItem('ticket_sidebar_mini', mini.value ? '1' : '0')
}

const menuItems = ref([
  { path: '/tickets', label: '工单列表', icon: TicketsPlane, badge: 0 },
  { path: '/users', label: '用户管理', icon: Users, adminOnly: true },
  { path: '/settings', label: '个人设置', shortLabel: '我的', icon: Settings },
])

const visibleMenu = computed(() => menuItems.value.filter((m) => !m.adminOnly || userStore.isAdmin))

const isActive = (path) => {
  if (path === '/tickets') return route.path === '/' || route.path.startsWith('/tickets')
  return route.path === path
}

const pageTitle = computed(() => {
  if (route.name === 'detail') return '工单详情'
  return menuItems.value.find((m) => isActive(m.path))?.label || '工单系统'
})

const avatarLetter = computed(() => {
  const u = userStore.user
  if (!u) return '?'
  return (u.display_name?.[0] || u.username?.[0] || '?').toUpperCase()
})

const userMenuItems = computed(() => [
  { label: '个人设置', value: 'settings', icon: UserRound },
  { label: '退出登录', value: 'logout', icon: LogOut, danger: true, divided: true },
])

const onUserCommand = async (cmd) => {
  if (cmd === 'settings') return router.push('/settings')
  if (cmd === 'logout') {
    const ok = await confirmDialog({
      title: '退出登录',
      message: '确认退出当前账号？',
      confirmText: '退出',
      danger: true,
    })
    if (!ok) return
    userStore.logout()
    toast.success('已退出登录')
    router.push('/login')
  }
}

// 待处理角标轮询
let timer = null
const refreshBadge = async () => {
  if (!userStore.isLoggedIn) return
  try {
    const r = await listTickets({ status: 'pending', page: 1, page_size: 1 })
    menuItems.value[0].badge = r.total
  } catch {}
}
onMounted(() => {
  refreshBadge()
  timer = setInterval(refreshBadge, 60000)
})
onBeforeUnmount(() => clearInterval(timer))
</script>

<style scoped>
.layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

/* ════════ 侧边栏 ════════ */
.sidebar {
  width: var(--sidebar-w);
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  background: var(--bg-surface);
  border-right: 1px solid var(--border-soft);
  transition: width var(--d-base) var(--ease-out);
}
.sidebar.mini { width: var(--sidebar-w-mini); }

.side-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  height: var(--topbar-h);
  padding: 0 20px;
  cursor: pointer;
  user-select: none;
  flex-shrink: 0;
  overflow: hidden;
}
.sidebar.mini .side-brand { padding: 0; justify-content: center; }
.side-logo {
  width: 36px;
  height: 36px;
  border-radius: 11px;
  background: var(--gradient-brand);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: var(--shadow-brand);
  transition: transform var(--d-fast) var(--ease-spring);
}
.side-brand:hover .side-logo { transform: scale(1.06) rotate(-4deg); }
.sb-name { font-size: 15px; font-weight: 700; letter-spacing: -0.3px; color: var(--text-1); line-height: 1.25; white-space: nowrap; }
.sb-sub { font-size: 11px; color: var(--text-3); white-space: nowrap; }

.side-nav {
  flex: 1;
  padding: 14px 12px;
  overflow-y: auto;
  overflow-x: hidden;
}
.nav-caption {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 1px;
  color: var(--text-3);
  padding: 4px 12px 10px;
  text-transform: uppercase;
  white-space: nowrap;
}

.nav-item {
  position: relative;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  margin-bottom: 3px;
  border-radius: var(--r-md);
  color: var(--text-2);
  font-size: 14px;
  cursor: pointer;
  transition: all var(--d-fast) var(--ease-out);
  white-space: nowrap;
}
.sidebar.mini .nav-item { justify-content: center; padding: 10px 0; }
.nav-item:hover { background: var(--bg-hover); color: var(--text-1); }
.nav-item.active {
  background: var(--primary-soft);
  color: var(--primary);
  font-weight: 600;
}
.nav-active-bar {
  position: absolute;
  left: -12px;
  top: 50%;
  width: 3.5px;
  height: 0;
  border-radius: 0 3px 3px 0;
  background: var(--gradient-brand);
  transform: translateY(-50%);
  transition: height var(--d-base) var(--ease-spring);
}
.sidebar.mini .nav-active-bar { display: none; }
.nav-item.active .nav-active-bar { height: 20px; }
.nav-icon { flex-shrink: 0; }
.nav-label { flex: 1; overflow: hidden; text-overflow: ellipsis; }
.nav-badge {
  font-size: 10.5px;
  font-weight: 600;
  line-height: 1;
  padding: 3px 7px;
  border-radius: var(--r-full);
  background: var(--gradient-brand);
  color: #fff;
  box-shadow: var(--shadow-brand);
}

.side-foot {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 12px;
  border-top: 1px solid var(--border-soft);
}
.foot-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 36px;
  min-width: 36px;
  padding: 0 10px;
  border: 1px solid var(--border);
  border-radius: var(--r-md);
  background: transparent;
  color: var(--text-2);
  font-size: 13px;
  cursor: pointer;
  transition: all var(--d-fast) var(--ease-out);
  white-space: nowrap;
  overflow: hidden;
}
.foot-btn:first-child { flex: 1; }
.sidebar.mini .foot-btn:first-child { flex: none; }
.foot-btn:hover { background: var(--bg-hover); color: var(--text-1); }
.foot-label { overflow: hidden; }

.fade-t-enter-active { transition: opacity var(--d-base) var(--ease-out); }
.fade-t-enter-from, .fade-t-leave-to { opacity: 0; }
.fade-t-leave-active { transition: opacity var(--d-fast); }

/* ════════ 主列 ════════ */
.main-col {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  overflow: hidden;
}

.topbar {
  height: var(--topbar-h);
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: color-mix(in srgb, var(--bg-surface) 82%, transparent);
  -webkit-backdrop-filter: blur(12px);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border-soft);
  z-index: 10;
}
@supports not (background: color-mix(in srgb, red 50%, blue)) {
  .topbar { background: var(--bg-surface); }
}
@media (max-width: 767px) { .topbar { padding: 0 16px; } }

.topbar-left { display: flex; align-items: center; gap: 12px; min-width: 0; }
.m-logo {
  width: 30px;
  height: 30px;
  border-radius: 9px;
  background: var(--gradient-brand);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: var(--shadow-brand);
}
.topbar-title {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
  letter-spacing: -0.3px;
  color: var(--text-1);
  white-space: nowrap;
}

.topbar-right { display: flex; align-items: center; gap: 10px; }
.icon-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: var(--r-md);
  background: transparent;
  color: var(--text-2);
  cursor: pointer;
  transition: all var(--d-fast) var(--ease-out);
}
.icon-btn:hover { background: var(--bg-hover); color: var(--text-1); }

.user-chip {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 5px 8px 5px 5px;
  border-radius: var(--r-full);
  transition: background var(--d-fast) var(--ease-out);
  user-select: none;
}
.user-chip:hover { background: var(--bg-hover); }

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--gradient-brand);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 700;
  flex-shrink: 0;
  box-shadow: var(--shadow-brand);
}
.avatar[data-role="admin"] {
  background: linear-gradient(135deg, #F59E0B 0%, #F97316 100%);
  box-shadow: 0 4px 12px -4px rgba(249, 115, 22, 0.5);
}
.user-meta { line-height: 1.25; }
.user-name { font-size: 13px; font-weight: 600; color: var(--text-1); }
.user-role { font-size: 11px; color: var(--text-3); }
.user-chev { color: var(--text-3); }

/* ════════ 页面滚动区 ════════ */
.page-scroll {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  scroll-behavior: smooth;
}
.page-inner {
  max-width: var(--content-max);
  margin: 0 auto;
  padding: 26px 28px 40px;
}
@media (max-width: 767px) {
  .page-inner { padding: 18px 16px 24px; }
  .page-inner.has-bottomnav {
    padding-bottom: calc(var(--bottomnav-h) + env(safe-area-inset-bottom, 0px) + 28px);
  }
}

/* ════════ 底部导航（移动端） ════════ */
.bottomnav {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 50;
  height: calc(var(--bottomnav-h) + env(safe-area-inset-bottom, 0px));
  padding-bottom: env(safe-area-inset-bottom, 0px);
  display: flex;
  align-items: stretch;
  background: color-mix(in srgb, var(--bg-surface) 88%, transparent);
  -webkit-backdrop-filter: blur(16px);
  backdrop-filter: blur(16px);
  border-top: 1px solid var(--border-soft);
}
@supports not (background: color-mix(in srgb, red 50%, blue)) {
  .bottomnav { background: var(--bg-surface); }
}

.bn-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 3px;
  color: var(--text-3);
  font-size: 10.5px;
  transition: color var(--d-fast) var(--ease-out);
}
.bn-item.active { color: var(--primary); font-weight: 600; }

.bn-icon-wrap {
  position: relative;
  display: inline-flex;
  padding: 4px 18px;
  border-radius: var(--r-full);
  transition: background var(--d-base) var(--ease-out);
}
.bn-item.active .bn-icon-wrap { background: var(--primary-soft); }

.bn-badge {
  position: absolute;
  top: -2px;
  right: 6px;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  border-radius: var(--r-full);
  background: var(--gradient-brand);
  color: #fff;
  font-size: 9.5px;
  font-weight: 700;
  line-height: 16px;
  text-align: center;
  box-shadow: 0 2px 6px rgba(99, 102, 241, 0.5);
}
.bn-label { line-height: 1; }
</style>
