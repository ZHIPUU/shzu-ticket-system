<template>
  <div class="app-layout">
    <!-- 移动端遮罩 -->
    <div v-if="mobileOpen" class="mobile-overlay" @click="mobileOpen = false" />

    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed, mobileOpen }">
      <div class="brand" @click="$router.push('/tickets')">
        <div class="brand-icon">
          <TicketsPlane :size="20" :stroke-width="2" />
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
          <component :is="item.icon" :size="18" :stroke-width="2" class="menu-icon" />
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
          <component :is="sidebarCollapsed ? ChevronsRight : ChevronsLeft" :size="16" :stroke-width="2" />
        </button>
        <button class="theme-toggle" @click="toggleTheme" :title="isDark ? '切换到亮色' : '切换到暗色'">
          <component :is="isDark ? Sun : Moon" :size="16" :stroke-width="2" />
        </button>
      </div>
    </aside>

    <!-- 主内容区 -->
    <main class="main">
      <header class="topbar">
        <div class="topbar-left">
          <button class="hamburger" @click="mobileOpen = true" aria-label="打开菜单">
            <Menu :size="22" :stroke-width="2" />
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
              <ChevronDown :size="14" :stroke-width="2" class="chev" />
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="settings">
                  <Settings :size="14" :stroke-width="2" class="drop-icon" />
                  个人设置
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <span class="logout-text">
                    <LogOut :size="14" :stroke-width="2" class="drop-icon" />
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
import {
  TicketsPlane,
  Users,
  Settings,
  Menu,
  Sun,
  Moon,
  ChevronsLeft,
  ChevronsRight,
  ChevronDown,
  LogOut,
} from '@lucide/vue'
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
  { path: '/tickets', label: '工单列表', icon: TicketsPlane, badge: 0 },
  { path: '/users', label: '用户管理', icon: Users, adminOnly: true },
  { path: '/settings', label: '个人设置', icon: Settings },
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
.sidebar.collapsed { width: 68px; }

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
  padding: 0 20px;
  cursor: pointer;
  user-select: none;
  border-bottom: 1px solid var(--border-soft);
  height: var(--header-height);
  flex-shrink: 0;
}
.brand-icon {
  width: 34px;
  height: 34px;
  border-radius: var(--radius-md);
  background: var(--gradient-header);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.25);
}
.brand-text { line-height: 1.3; }
.brand-name { font-weight: 600; font-size: 15px; color: var(--text-primary); }
.brand-sub { font-size: 11px; color: var(--text-tertiary); margin-top: 2px; }

.menu {
  flex: 1;
  padding: 12px 12px;
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
  padding: 9px 12px;
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
  height: 18px;
  background: var(--color-primary);
  border-radius: 0 2px 2px 0;
}
.menu-icon { flex-shrink: 0; }
.menu-label { flex: 1; }
.menu-badge { margin-left: auto; }

.sidebar-footer {
  border-top: 1px solid var(--border-soft);
  padding: 10px 12px;
  display: flex;
  align-items: center;
  gap: 6px;
}
.collapse-btn {
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  width: 34px;
  height: 34px;
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

.theme-toggle {
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  width: 34px;
  height: 34px;
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
}

/* ─── 主内容 ─── */
.main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--bg-base);
  min-width: 0;
}

.topbar {
  height: var(--header-height);
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border-soft);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
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
  border-radius: var(--radius-sm);
}
.hamburger:hover { color: var(--text-primary); background: var(--bg-hover); }
@media (max-width: 768px) { .hamburger { display: flex; align-items: center; justify-content: center; } }

.page-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: -0.2px;
}
@media (max-width: 768px) { .page-title { font-size: 14px; } }

.user-chip {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 5px 10px 5px 5px;
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: background var(--transition-fast);
  user-select: none;
}
.user-chip:hover { background: var(--bg-hover); }
@media (max-width: 768px) { .user-chip { padding: 4px; } }

.avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: var(--color-primary);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 12px;
  box-shadow: 0 2px 4px rgba(37, 99, 235, 0.25);
  flex-shrink: 0;
}
.avatar[data-role="admin"] {
  background: linear-gradient(135deg, #EF4444 0%, #F87171 100%);
  box-shadow: 0 2px 4px rgba(239, 68, 68, 0.25);
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
.chev { color: var(--text-tertiary); flex-shrink: 0; }
.drop-icon { margin-right: 8px; vertical-align: middle; }
.logout-text { color: var(--color-danger); display: inline-flex; align-items: center; }

.page-container {
  flex: 1;
  overflow: auto;
  padding: 24px;
}
@media (max-width: 768px) { .page-container { padding: 16px; } }
</style>
