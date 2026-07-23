<template>
  <div class="app-layout">
    <!-- 移动端遮罩 -->
    <div v-if="mobileOpen" class="mobile-overlay" @click="mobileOpen = false" />

    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed, mobileOpen }">
      <div class="brand" @click="$router.push('/tickets')">
        <div class="brand-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 7l9-4 9 4-9 4-9-4z" />
            <path d="M3 12l9 4 9-4" />
            <path d="M3 17l9 4 9-4" />
          </svg>
        </div>
        <transition name="fade">
          <div v-show="!sidebarCollapsed" class="brand-text">
            <div class="brand-name">石小易工单</div>
            <div class="brand-sub">迎新智能体后台</div>
          </div>
        </transition>
      </div>

      <nav class="menu">
        <div class="menu-section" v-show="!sidebarCollapsed">主菜单</div>
        <router-link
          v-for="item in visibleMenu"
          :key="item.path"
          :to="item.path"
          class="menu-item"
          :class="{ active: isActive(item.path) }"
          :title="sidebarCollapsed ? item.label : ''"
          @click="mobileOpen = false"
        >
          <span class="menu-icon" v-html="item.icon" />
          <transition name="fade">
            <span v-show="!sidebarCollapsed" class="menu-label">{{ item.label }}</span>
          </transition>
          <transition name="fade">
            <el-badge
              v-if="!sidebarCollapsed && item.badge && item.badge > 0"
              :value="item.badge"
              :max="99"
              class="menu-badge"
            />
          </transition>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <button class="collapse-btn" @click="toggleSidebar" :title="sidebarCollapsed ? '展开' : '收起'">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path :d="sidebarCollapsed ? 'M9 6l6 6-6 6' : 'M15 6l-6 6 6 6'" />
          </svg>
        </button>
        <button class="theme-toggle" @click="toggleTheme" :title="isDark ? '切换到亮色' : '切换到暗色'">
          <svg v-if="isDark" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="4" />
            <path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41" />
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
          </svg>
        </button>
      </div>
    </aside>

    <!-- 主内容区 -->
    <main class="main">
      <header class="topbar">
        <div class="topbar-left">
          <button class="hamburger" @click="mobileOpen = true" aria-label="打开菜单">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 12h18M3 6h18M3 18h18" />
            </svg>
          </button>
          <h2 class="page-title">{{ pageTitle }}</h2>
        </div>
        <div class="topbar-right">
          <el-dropdown trigger="click" @command="onUserCommand">
            <div class="user-chip">
              <div class="avatar" :data-role="userStore.user?.role">
                {{ avatarLetter }}
              </div>
              <div class="user-info">
                <div class="user-name">{{ userStore.user?.display_name || userStore.user?.username }}</div>
                <div class="user-role">{{ userStore.user?.role === 'admin' ? '管理员' : '工作人员' }}</div>
              </div>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:14px;height:14px;opacity:0.6">
                <path d="M6 9l6 6 6-6" />
              </svg>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="settings">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:14px;height:14px;margin-right:8px;vertical-align:middle">
                    <circle cx="12" cy="12" r="3"/>
                    <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
                  </svg>
                  个人设置
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <span style="color: #d63031">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:14px;height:14px;margin-right:8px;vertical-align:middle">
                      <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4M16 17l5-5-5-5M21 12H9"/>
                    </svg>
                    退出登录
                  </span>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>

      <div class="page-container">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { useTheme } from '../composables/useTheme'
import { useUserStore } from '../stores/user'
import { listTickets } from '../api/ticket'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const { isDark, toggle: toggleTheme } = useTheme()
const sidebarCollapsed = ref(localStorage.getItem('ticket_sidebar_collapsed') === 'true')
const mobileOpen = ref(false)
const isMobile = ref(window.innerWidth < 768)

const checkMobile = () => { isMobile.value = window.innerWidth < 768 }
onMounted(() => {
  window.addEventListener('resize', checkMobile)
  refreshBadges()
  setInterval(refreshBadges, 60000)
})
onUnmounted(() => window.removeEventListener('resize', checkMobile))

const menuItems = ref([
  {
    path: '/tickets',
    label: '工单列表',
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 3h18v4H3zM3 11h18v4H3zM3 19h18v2H3z"/></svg>',
    badge: 0,
  },
  {
    path: '/users',
    label: '用户管理',
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75"/></svg>',
    adminOnly: true,
  },
  {
    path: '/settings',
    label: '个人设置',
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>',
  },
])

const pageTitle = computed(() => {
  return menuItems.value.find((m) => isActive(m.path))?.label || '工单系统'
})

const avatarLetter = computed(() => {
  const u = userStore.user
  if (!u) return '?'
  return (u.display_name?.[0] || u.username?.[0] || '?').toUpperCase()
})

const visibleMenu = computed(() => {
  return menuItems.value.filter((m) => !m.adminOnly || userStore.isAdmin)
})

const isActive = (path) => {
  if (path === '/tickets') return route.path === '/' || route.path.startsWith('/tickets')
  return route.path === path
}

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
  localStorage.setItem('ticket_sidebar_collapsed', sidebarCollapsed.value)
}

const onUserCommand = (cmd) => {
  if (cmd === 'settings') router.push('/settings')
  if (cmd === 'logout') {
    ElMessageBox.confirm('确认退出登录？', '退出', { type: 'warning' })
      .then(() => {
        userStore.logout()
        ElMessage.success('已退出登录')
        router.push('/login')
      })
      .catch(() => {})
  }
}

const refreshBadges = async () => {
  if (!userStore.isLoggedIn) return
  try {
    const r = await listTickets({ status: 'pending', page: 1, page_size: 1 })
    menuItems.value[0].badge = r.total
  } catch {}
}
</script>

<style scoped>
.app-layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

/* ─── 移动端遮罩 ─── */
.mobile-overlay {
  display: none;
}
@media (max-width: 768px) {
  .mobile-overlay {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.4);
    z-index: 20;
  }
}

/* ─── 侧边栏 ─── */
.sidebar {
  width: var(--sidebar-width);
  background: var(--bg-surface);
  border-right: 1px solid var(--border-soft);
  display: flex;
  flex-direction: column;
  transition: width var(--transition-base), transform var(--transition-base);
  flex-shrink: 0;
  position: relative;
  z-index: 10;
}
.sidebar.collapsed { width: 72px; }

@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    width: var(--sidebar-width);
    transform: translateX(-100%);
    z-index: 30;
    box-shadow: 4px 0 24px rgba(0,0,0,0.15);
  }
  .sidebar.mobileOpen {
    transform: translateX(0);
  }
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px 20px;
  cursor: pointer;
  user-select: none;
  border-bottom: 1px solid var(--border-soft);
  height: var(--header-height);
}
.brand-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: var(--gradient-header);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 184, 148, 0.3);
}
.brand-icon svg { width: 20px; height: 20px; }
.brand-text { line-height: 1.3; }
.brand-name { font-weight: 600; font-size: 15px; color: var(--text-primary); }
.brand-sub { font-size: 11px; color: var(--text-tertiary); margin-top: 2px; }

.menu {
  flex: 1;
  padding: 16px 12px;
  overflow-y: auto;
}
.menu-section {
  font-size: 11px;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 8px 12px;
  margin-bottom: 4px;
  font-weight: 500;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  text-decoration: none;
  cursor: pointer;
  margin-bottom: 2px;
  transition: all var(--transition-fast);
  position: relative;
  font-size: 14px;
}
.menu-item:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}
.menu-item.active {
  background: var(--color-primary-soft);
  color: var(--color-primary);
  font-weight: 500;
}
.menu-item.active::before {
  content: '';
  position: absolute;
  left: -12px;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 20px;
  background: var(--color-primary);
  border-radius: 0 2px 2px 0;
}
.menu-icon { width: 20px; height: 20px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.menu-icon :deep(svg) { width: 18px; height: 18px; }
.menu-label { flex: 1; }
.menu-badge { margin-left: auto; }

.sidebar-footer {
  border-top: 1px solid var(--border-soft);
  padding: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.collapse-btn {
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 50%;
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--text-secondary);
  transition: all var(--transition-fast);
  flex-shrink: 0;
}
.collapse-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}
.collapse-btn svg { width: 16px; height: 16px; }

.theme-toggle {
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  border-radius: 50%;
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--text-secondary);
  transition: all var(--transition-base);
  margin-left: auto;
}
.theme-toggle:hover {
  background: var(--color-primary-soft);
  color: var(--color-primary);
  transform: rotate(20deg);
}
.theme-toggle svg { width: 16px; height: 16px; }

/* ─── 主内容 ─── */
.main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--bg-base);
}

.topbar {
  height: var(--header-height);
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border-soft);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;
  flex-shrink: 0;
}
@media (max-width: 768px) { .topbar { padding: 0 16px; } }

.topbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.hamburger {
  display: none;
  background: transparent;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  padding: 4px;
  border-radius: 6px;
}
.hamburger:hover { color: var(--text-primary); background: var(--bg-hover); }
.hamburger svg { width: 22px; height: 22px; }
@media (max-width: 768px) { .hamburger { display: flex; align-items: center; justify-content: center; } }

.page-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: -0.3px;
}
@media (max-width: 768px) { .page-title { font-size: 14px; } }

.user-chip {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 12px 6px 6px;
  border-radius: 999px;
  cursor: pointer;
  transition: background var(--transition-fast);
  user-select: none;
}
.user-chip:hover { background: var(--bg-hover); }
@media (max-width: 768px) { .user-chip { padding: 4px; } }

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--color-primary);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 13px;
  box-shadow: 0 2px 6px rgba(0, 184, 148, 0.3);
  flex-shrink: 0;
}
.avatar[data-role="admin"] {
  background: linear-gradient(135deg, #d63031 0%, #e17055 100%);
  box-shadow: 0 2px 6px rgba(214, 48, 49, 0.3);
}
.user-info { line-height: 1.2; }
@media (max-width: 768px) { .user-info { display: none; } }
.user-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
}
.user-role {
  font-size: 11px;
  color: var(--text-tertiary);
  margin-top: 1px;
}

.page-container {
  flex: 1;
  overflow: auto;
  padding: 32px;
}
@media (max-width: 768px) { .page-container { padding: 16px; } }
</style>
