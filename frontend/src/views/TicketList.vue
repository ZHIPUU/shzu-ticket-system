<template>
  <div class="ticket-list">
    <!-- ═══ 页头 ═══ -->
    <header class="list-head anim-fade-up">
      <div class="head-text">
        <h1>工单列表</h1>
        <p>管理用户提问工单 · 答复 · 归档 · 导出</p>
      </div>
      <div class="head-actions">
        <UiButton variant="ghost" :icon-only="isMobile" @click="openExport">
          <template #icon><Download :size="16" :stroke-width="2" /></template>
          <span v-if="!isMobile">导出</span>
        </UiButton>
        <UiButton variant="primary" :icon-only="isMobile" @click="openCreate">
          <template #icon><Plus :size="16" :stroke-width="2.4" /></template>
          <span v-if="!isMobile">模拟提交</span>
        </UiButton>
      </div>
    </header>

    <!-- ═══ 统计 ═══ -->
    <div class="stats-grid">
      <StatCard v-for="(s, i) in statItems" :key="s.label" v-bind="s" class="anim-stagger" :style="{ '--i': i }" />
    </div>

    <!-- ═══ 工具栏 ═══ -->
    <div class="toolbar card anim-stagger" style="--i: 4">
      <div class="toolbar-main">
        <UiSegmented v-model="filters.status" :options="statusOptions" @change="onFilterChange" />
        <div class="toolbar-right">
          <UiInput
            v-model="filters.q"
            class="search-input"
            placeholder="工单号 / 问题关键字"
            clearable
            @enter="onFilterChange"
            @clear="onFilterChange"
            @update:model-value="onSearchInput"
          >
            <template #prefix><Search :size="15" :stroke-width="2" /></template>
          </UiInput>
          <UiButton
            variant="ghost"
            :icon-only="isMobile"
            :class="{ 'filter-on': filterOpen || hasAdvFilter }"
            @click="filterOpen = !filterOpen"
          >
            <template #icon><Funnel :size="15" :stroke-width="2" /></template>
            <span v-if="!isMobile">筛选</span>
          </UiButton>
          <UiButton v-if="isMobile" variant="ghost" icon-only @click="toggleSelectMode">
            <template #icon>
              <SquareCheckBig v-if="!selectMode" :size="15" :stroke-width="2" />
              <X v-else :size="15" :stroke-width="2" />
            </template>
          </UiButton>
        </div>
      </div>

      <!-- 高级筛选（可折叠） -->
      <div class="adv-wrap" :class="{ open: filterOpen }">
        <div class="adv-inner">
          <div class="adv-field">
            <label>归档状态</label>
            <UiSelect v-model="filters.archived" :options="archivedOptions" @change="onFilterChange" />
          </div>
          <div class="adv-field">
            <label>分类</label>
            <UiInput v-model="filters.category" placeholder="如：宿舍 / 招办" clearable @enter="onFilterChange" @clear="onFilterChange" />
          </div>
          <div class="adv-field">
            <label>开始日期</label>
            <UiInput v-model="filters.startDate" type="date" @update:model-value="onFilterChange" />
          </div>
          <div class="adv-field">
            <label>结束日期</label>
            <UiInput v-model="filters.endDate" type="date" @update:model-value="onFilterChange" />
          </div>
          <div class="adv-field adv-actions">
            <UiButton variant="ghost" size="sm" @click="reset">
              <template #icon><RotateCcw :size="13" :stroke-width="2" /></template>
              重置全部
            </UiButton>
          </div>
        </div>
      </div>
    </div>

    <!-- ═══ 列表 ═══ -->
    <div class="list-card card anim-stagger" style="--i: 5">
      <!-- 加载骨架 -->
      <div v-if="loading" class="list-loading"><UiSkeleton type="rows" :count="6" /></div>

      <!-- 空态 -->
      <UiEmpty
        v-else-if="!rows.length"
        icon="Inbox"
        title="暂无工单"
        :description="hasAnyFilter ? '当前筛选条件下没有工单，试试调整筛选' : '还没有用户提交工单'"
      >
        <UiButton v-if="hasAnyFilter" variant="ghost" size="sm" @click="reset">清除筛选条件</UiButton>
      </UiEmpty>

      <!-- 桌面端表格 -->
      <template v-else-if="!isMobile">
        <div class="t-head">
          <div class="t-cell c-check">
            <UiCheckbox :model-value="allSelected" :indeterminate="partialSelected" @update:model-value="toggleAll" />
          </div>
          <div class="t-cell c-id">工单号</div>
          <div class="t-cell c-q">用户问题</div>
          <div class="t-cell c-status">状态</div>
          <div class="t-cell c-cat">分类</div>
          <div class="t-cell c-time">创建时间</div>
          <div class="t-cell c-arrow" />
        </div>
        <div
          v-for="(row, i) in rows"
          :key="row.ticket_id"
          class="t-row anim-stagger"
          :style="{ '--i': i }"
          :class="{ selected: isSelected(row.ticket_id) }"
          @click="goDetail(row.ticket_id)"
        >
          <div class="t-cell c-check" @click.stop>
            <UiCheckbox :model-value="isSelected(row.ticket_id)" @update:model-value="toggleOne(row.ticket_id)" />
          </div>
          <div class="t-cell c-id">
            <span class="tid mono">{{ row.ticket_id }}</span>
          </div>
          <div class="t-cell c-q">
            <span class="q-text clamp-2">{{ row.question }}</span>
          </div>
          <div class="t-cell c-status">
            <StatusPill :status="row.status" />
            <span v-if="row.archived" class="archived-mark" title="已归档"><Archive :size="13" :stroke-width="2.2" /></span>
          </div>
          <div class="t-cell c-cat"><CategoryTag :category="row.category" /></div>
          <div class="t-cell c-time"><TimeText :value="row.created_at" /></div>
          <div class="t-cell c-arrow"><ChevronRight :size="15" :stroke-width="2" /></div>
        </div>
      </template>

      <!-- 移动端卡片 -->
      <template v-else>
        <div
          v-for="(row, i) in rows"
          :key="row.ticket_id"
          class="m-card anim-stagger"
          :style="{ '--i': i }"
          :class="{ selected: isSelected(row.ticket_id), 'select-mode': selectMode }"
          @click="onCardTap(row.ticket_id)"
        >
          <div v-if="selectMode" class="m-check" @click.stop="toggleOne(row.ticket_id)">
            <UiCheckbox :model-value="isSelected(row.ticket_id)" />
          </div>
          <div class="m-body">
            <div class="m-top">
              <StatusPill :status="row.status" />
              <span v-if="row.archived" class="archived-mark"><Archive :size="12" :stroke-width="2.2" />已归档</span>
              <TimeText :value="row.created_at" mode="relative" class="m-time" />
            </div>
            <div class="m-q clamp-2">{{ row.question }}</div>
            <div class="m-bottom">
              <span class="tid mono">{{ row.ticket_id }}</span>
              <CategoryTag :category="row.category" />
              <ChevronRight v-if="!selectMode" :size="15" :stroke-width="2" class="m-chev" />
            </div>
          </div>
        </div>
      </template>

      <!-- 分页 -->
      <div v-if="!loading && rows.length" class="list-pager">
        <UiPagination
          v-model:page="filters.page"
          v-model:page-size="filters.page_size"
          :total="total"
          :compact="isMobile"
          @change="reload"
        />
      </div>
    </div>

    <!-- ═══ 浮动批量操作条 ═══ -->
    <Transition name="rise">
      <div v-if="selected.length > 0" class="bulk-bar card">
        <span class="bulk-count">已选 <strong class="tnum">{{ selected.length }}</strong> 项</span>
        <span class="bulk-divider" />
        <UiButton variant="text" size="sm" @click="onBulkArchive(true)">
          <template #icon><Archive :size="14" :stroke-width="2" /></template>
          归档
        </UiButton>
        <UiButton variant="text" size="sm" @click="onBulkArchive(false)">
          <template #icon><ArchiveRestore :size="14" :stroke-width="2" /></template>
          取消归档
        </UiButton>
        <UiButton variant="text" size="sm" class="bulk-danger" @click="onBulkDelete">
          <template #icon><Trash2 :size="14" :stroke-width="2" /></template>
          删除
        </UiButton>
        <span class="bulk-divider" />
        <button class="bulk-close" aria-label="退出批量操作" @click="exitBulk">
          <X :size="15" :stroke-width="2.2" />
        </button>
      </div>
    </Transition>

    <!-- ═══ 模拟提交 ═══ -->
    <UiModal v-model="createVisible" title="模拟智能体提交工单" width="540px">
      <div class="form-stack">
        <div class="form-item">
          <label class="form-label">问题 <em>*</em></label>
          <UiInput
            v-model="createForm.question"
            type="textarea"
            :rows="3"
            :maxlength="500"
            show-word-limit
            placeholder="例如：石河子大学2026年计算机学院宿舍分配在哪？"
          />
        </div>
        <div class="form-item">
          <label class="form-label">用户 ID <em>*</em></label>
          <UiInput v-model="createForm.user_id" :maxlength="128" placeholder="HiAgent 平台 session_id" />
        </div>
        <div class="form-item">
          <label class="form-label">RAG 结果</label>
          <UiInput v-model="createForm.rag_result" type="textarea" :rows="2" placeholder="留空表示完全无结果" />
        </div>
      </div>
      <template #footer>
        <UiButton variant="ghost" @click="createVisible = false">取消</UiButton>
        <UiButton variant="primary" :loading="creating" @click="doCreate">提交</UiButton>
      </template>
    </UiModal>

    <!-- ═══ 导出 ═══ -->
    <UiModal v-model="exportVisible" title="导出工单" width="440px">
      <div class="form-stack">
        <div class="form-item">
          <label class="form-label">导出格式</label>
          <UiSegmented v-model="exportFormat" :options="[{ label: 'CSV 表格', value: 'csv' }, { label: 'JSON 数据', value: 'json' }]" />
        </div>
        <p class="export-note">
          <Info :size="14" :stroke-width="2" />
          按当前筛选条件导出全部工单（最多 10000 条），不影响原数据。
        </p>
      </div>
      <template #footer>
        <UiButton variant="ghost" @click="exportVisible = false">取消</UiButton>
        <UiButton variant="primary" :loading="exporting" @click="doExport">
          <template #icon><Download :size="14" :stroke-width="2" /></template>
          开始导出
        </UiButton>
      </template>
    </UiModal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import {
  Download, Plus, Search, Funnel, RotateCcw, Archive, ArchiveRestore,
  Trash2, ChevronRight, X, SquareCheckBig, Info,
} from '@lucide/vue'
import {
  listTickets, submitTicket, batchDeleteTickets, archiveTicket, exportTickets,
} from '../api/ticket'
import UiButton from '../ui/UiButton.vue'
import UiInput from '../ui/UiInput.vue'
import UiSelect from '../ui/UiSelect.vue'
import UiModal from '../ui/UiModal.vue'
import UiSegmented from '../ui/UiSegmented.vue'
import UiCheckbox from '../ui/UiCheckbox.vue'
import UiPagination from '../ui/UiPagination.vue'
import UiSkeleton from '../ui/UiSkeleton.vue'
import UiEmpty from '../ui/UiEmpty.vue'
import StatusPill from '../ui/StatusPill.vue'
import CategoryTag from '../ui/CategoryTag.vue'
import TimeText from '../ui/TimeText.vue'
import StatCard from '../components/StatCard.vue'
import { toast } from '../ui/toast'
import { confirmDialog } from '../ui/confirm'
import { useIsMobile } from '../composables/useMediaQuery'

const router = useRouter()
const isMobile = useIsMobile()

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const filterOpen = ref(false)
const selectMode = ref(false)
const selected = ref([])

const createVisible = ref(false)
const creating = ref(false)
const exportVisible = ref(false)
const exportFormat = ref('csv')
const exporting = ref(false)

const filters = reactive({
  status: '',
  archived: '',
  category: '',
  q: '',
  startDate: '',
  endDate: '',
  page: 1,
  page_size: 20,
})

const createForm = reactive({ question: '', user_id: '', rag_result: '' })

const statusOptions = [
  { label: '全部', value: '' },
  { label: '待处理', value: 'pending' },
  { label: '处理中', value: 'processing' },
  { label: '已答复', value: 'answered' },
  { label: '已关闭', value: 'closed' },
]

const archivedOptions = [
  { label: '全部', value: '' },
  { label: '未归档', value: 'false' },
  { label: '已归档', value: 'true' },
]

const counts = computed(() => ({
  pending: rows.value.filter((r) => r.status === 'pending').length,
  answered: rows.value.filter((r) => r.status === 'answered').length,
}))

const statItems = computed(() => [
  { icon: 'Layers', label: '本页工单', value: rows.value.length, meta: '当前筛选下', accent: '#6366F1' },
  { icon: 'Timer', label: '待处理', value: counts.value.pending, meta: '需尽快响应', accent: '#D98309' },
  { icon: 'MessageSquareCheck', label: '已答复', value: counts.value.answered, meta: '等待复问', accent: '#0EA96E' },
  { icon: 'Database', label: '全部工单', value: total.value, meta: '筛选范围内总数', accent: '#8B5CF6' },
])

const hasAdvFilter = computed(() => !!(filters.archived || filters.category || filters.startDate || filters.endDate))
const hasAnyFilter = computed(
  () => hasAdvFilter.value || !!filters.status || !!filters.q
)

const buildQuery = () => ({
  status: filters.status || undefined,
  archived: filters.archived || undefined,
  category: filters.category?.trim() || undefined,
  q: filters.q?.trim() || undefined,
  start_date: filters.startDate || undefined,
  end_date: filters.endDate || undefined,
  page: filters.page,
  page_size: filters.page_size,
})

const reload = async () => {
  loading.value = true
  try {
    const data = await listTickets(buildQuery())
    rows.value = data.items
    total.value = data.total
    selected.value = []
  } catch (e) {
    /* 拦截器已提示 */
  } finally {
    loading.value = false
  }
}

const onFilterChange = () => {
  filters.page = 1
  reload()
}

// 搜索防抖
let searchTimer = null
const onSearchInput = () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(onFilterChange, 450)
}
onBeforeUnmount(() => clearTimeout(searchTimer))

const reset = () => {
  Object.assign(filters, { status: '', archived: '', category: '', q: '', startDate: '', endDate: '', page: 1 })
  reload()
}

// ─── 选择 ───
const isSelected = (id) => selected.value.includes(id)
const allSelected = computed(() => rows.value.length > 0 && selected.value.length === rows.value.length)
const partialSelected = computed(() => selected.value.length > 0 && !allSelected.value)

const toggleOne = (id) => {
  const i = selected.value.indexOf(id)
  if (i === -1) selected.value.push(id)
  else selected.value.splice(i, 1)
}
const toggleAll = () => {
  selected.value = allSelected.value ? [] : rows.value.map((r) => r.ticket_id)
}
const toggleSelectMode = () => {
  selectMode.value = !selectMode.value
  if (!selectMode.value) selected.value = []
}
const exitBulk = () => {
  selected.value = []
  selectMode.value = false
}

const goDetail = (id) => router.push(`/tickets/${id}`)
const onCardTap = (id) => {
  if (selectMode.value) toggleOne(id)
  else goDetail(id)
}

// ─── 批量操作 ───
const onBulkArchive = async (archive) => {
  const action = archive ? '归档' : '取消归档'
  const ok = await confirmDialog({
    title: `批量${action}`,
    message: `确认对选中的 ${selected.value.length} 个工单执行「${action}」？`,
    confirmText: `确认${action}`,
  })
  if (!ok) return
  let okCount = 0, failCount = 0
  for (const id of selected.value) {
    try {
      await archiveTicket(id, archive)
      okCount++
    } catch (e) { failCount++ }
  }
  toast.success(`${action}完成：成功 ${okCount}${failCount ? `，失败 ${failCount}` : ''}`)
  exitBulk()
  await reload()
}

const onBulkDelete = async () => {
  const ok = await confirmDialog({
    title: '批量删除',
    message: `确认删除选中的 ${selected.value.length} 个工单？删除后不可在列表中查看。`,
    confirmText: '确认删除',
    danger: true,
  })
  if (!ok) return
  try {
    const r = await batchDeleteTickets([...selected.value])
    toast.success(`批量删除完成：${r.message || ''}`)
    exitBulk()
    await reload()
  } catch (e) {}
}

// ─── 模拟提交 ───
const openCreate = () => {
  createForm.question = ''
  createForm.user_id = `sess_${Date.now()}`
  createForm.rag_result = ''
  createVisible.value = true
}

const doCreate = async () => {
  if (!createForm.question.trim()) { toast.warning('请填写问题'); return }
  if (!createForm.user_id.trim()) { toast.warning('请填写用户 ID'); return }
  creating.value = true
  try {
    const r = await submitTicket(createForm)
    createVisible.value = false
    toast.success(`工单已创建：${r.ticket_id}`)
    router.push(`/tickets/${r.ticket_id}`)
  } catch (e) {} finally {
    creating.value = false
  }
}

// ─── 导出 ───
const openExport = () => {
  exportFormat.value = 'csv'
  exportVisible.value = true
}

const doExport = async () => {
  exporting.value = true
  try {
    const params = buildQuery()
    delete params.page
    delete params.page_size
    const r = await exportTickets({ ...params, format: exportFormat.value })
    const blob = r.data
    const ext = exportFormat.value
    const ts = new Date().toISOString().replace(/[:.]/g, '-').slice(0, 19)
    const filename = `tickets-${ts}.${ext}`
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    a.click()
    URL.revokeObjectURL(url)
    toast.success(`已导出 ${filename}（${(blob.size / 1024).toFixed(1)} KB）`)
    exportVisible.value = false
  } catch (e) {
    toast.error('导出失败：' + (e?.response?.data?.error_message || e.message))
  } finally {
    exporting.value = false
  }
}

onMounted(reload)
</script>

<style scoped>
.ticket-list { display: flex; flex-direction: column; gap: 18px; }

/* ─── 页头 ─── */
.list-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}
.head-text h1 {
  margin: 0 0 3px;
  font-size: 23px;
  font-weight: 700;
  letter-spacing: -0.5px;
  color: var(--text-1);
}
.head-text p { margin: 0; font-size: 13px; color: var(--text-3); }
.head-actions { display: flex; gap: 10px; flex-shrink: 0; }
@media (max-width: 767px) {
  .head-text h1 { font-size: 19px; }
  .head-text p { font-size: 12px; }
}

/* ─── 统计 ─── */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
}
@media (max-width: 960px) { .stats-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 480px) { .stats-grid { gap: 10px; } }

/* ─── 工具栏 ─── */
.toolbar { padding: 12px 16px; }
.toolbar-main {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}
.search-input { width: 230px; }
@media (max-width: 767px) {
  .toolbar-main { flex-direction: column; align-items: stretch; }
  .toolbar-right { width: 100%; }
  .search-input { flex: 1; width: auto; }
}

.filter-on {
  border-color: var(--primary) !important;
  color: var(--primary) !important;
  background: var(--primary-soft) !important;
}

/* 高级筛选折叠 */
.adv-wrap {
  display: grid;
  grid-template-rows: 0fr;
  transition: grid-template-rows var(--d-base) var(--ease-out);
}
.adv-wrap.open { grid-template-rows: 1fr; }
.adv-inner {
  overflow: hidden;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(170px, 1fr));
  gap: 12px 16px;
  align-items: end;
}
.adv-wrap.open .adv-inner { padding-top: 14px; }
.adv-field label {
  display: block;
  font-size: 12px;
  color: var(--text-3);
  margin-bottom: 5px;
  font-weight: 500;
}
.adv-field .ui-select, .adv-field .ui-field { width: 100%; }
.adv-actions { display: flex; justify-content: flex-end; }

/* ─── 列表卡片 ─── */
.list-card { overflow: hidden; }
.list-loading { padding: 8px 20px; }

/* 桌面表格 */
.t-head, .t-row {
  display: grid;
  grid-template-columns: 40px 158px minmax(0, 1fr) 110px 116px 150px 30px;
  align-items: center;
  gap: 12px;
  padding: 0 20px;
}
.t-head {
  height: 42px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-3);
  border-bottom: 1px solid var(--border-soft);
  background: var(--bg-sink);
}
.t-row {
  min-height: 62px;
  padding-top: 10px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--border-soft);
  cursor: pointer;
  transition: background var(--d-fast) var(--ease-out);
}
.t-row:last-child { border-bottom: none; }
.t-row:hover { background: var(--bg-hover); }
.t-row.selected { background: var(--primary-soft); }
.t-cell { min-width: 0; display: flex; align-items: center; }
.c-check { justify-content: center; }
.c-arrow { justify-content: center; color: var(--text-3); opacity: 0; transition: all var(--d-fast) var(--ease-out); }
.t-row:hover .c-arrow { opacity: 1; transform: translateX(2px); }

.tid {
  font-size: 12px;
  color: var(--primary);
  font-weight: 600;
  letter-spacing: -0.2px;
}
.q-text { font-size: 13.5px; color: var(--text-1); line-height: 1.55; }
.c-status { gap: 6px; }

.archived-mark {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  color: var(--text-3);
  font-size: 11px;
}

/* 移动端卡片 */
.m-card {
  display: flex;
  gap: 12px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border-soft);
  cursor: pointer;
  transition: background var(--d-fast) var(--ease-out);
}
.m-card:last-child { border-bottom: none; }
.m-card:active { background: var(--bg-hover); }
.m-card.selected { background: var(--primary-soft); }
.m-check { display: flex; align-items: center; }
.m-body { flex: 1; min-width: 0; }
.m-top {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 7px;
}
.m-time { margin-left: auto; }
.m-q {
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-1);
  margin-bottom: 8px;
}
.m-bottom {
  display: flex;
  align-items: center;
  gap: 8px;
}
.m-chev { margin-left: auto; color: var(--text-3); }

/* ─── 分页 ─── */
.list-pager {
  display: flex;
  justify-content: flex-end;
  padding: 14px 20px;
  border-top: 1px solid var(--border-soft);
}
@media (max-width: 767px) { .list-pager { justify-content: center; } }

/* ─── 浮动批量条 ─── */
.bulk-bar {
  position: fixed;
  left: 50%;
  bottom: 26px;
  transform: translateX(-50%);
  z-index: 60;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  border-radius: var(--r-full);
  box-shadow: var(--shadow-lg);
}
@media (max-width: 767px) {
  .bulk-bar { bottom: calc(var(--bottomnav-h) + env(safe-area-inset-bottom, 0px) + 14px); }
}
.bulk-count { font-size: 13px; color: var(--text-2); padding: 0 6px; white-space: nowrap; }
.bulk-count strong { color: var(--primary); }
.bulk-divider { width: 1px; height: 18px; background: var(--border); margin: 0 3px; }
.bulk-danger { color: var(--danger) !important; }
.bulk-danger:hover { background: var(--danger-soft) !important; }
.bulk-close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border: none;
  border-radius: 50%;
  background: var(--bg-hover);
  color: var(--text-2);
  cursor: pointer;
  transition: all var(--d-fast) var(--ease-out);
}
.bulk-close:hover { background: var(--bg-active); color: var(--text-1); }

/* ─── 表单 ─── */
.form-stack { display: flex; flex-direction: column; gap: 16px; }
.form-item { display: flex; flex-direction: column; gap: 7px; }
.form-label { font-size: 13px; font-weight: 500; color: var(--text-2); }
.form-label em { color: var(--danger); font-style: normal; }

.export-note {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin: 0;
  padding: 11px 13px;
  border-radius: var(--r-md);
  background: var(--primary-soft);
  color: var(--text-2);
  font-size: 12.5px;
  line-height: 1.6;
}
.export-note svg { flex-shrink: 0; margin-top: 2px; color: var(--primary); }
</style>
