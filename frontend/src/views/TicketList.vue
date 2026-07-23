<template>
  <div class="ticket-list">
    <header class="page-header">
      <div>
        <h1 class="page-title">工单列表</h1>
        <p class="page-desc">管理来自各渠道的用户提问工单</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" class="action-btn" @click="openCreate">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:16px;height:16px;margin-right:6px">
            <path d="M12 5v14M5 12h14" />
          </svg>
          <span class="btn-label">模拟提交</span>
        </el-button>
      </div>
    </header>

    <!-- 统计卡片 -->
    <section class="stats-row">
      <div class="stat-card" v-for="s in stats" :key="s.key" :style="{ '--accent': s.color }">
        <div class="stat-label">{{ s.label }}</div>
        <div class="stat-value">{{ s.value }}</div>
        <div class="stat-meta">{{ s.meta }}</div>
      </div>
    </section>

    <!-- 筛选 -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-header" @click="filterExpanded = !filterExpanded">
        <span>筛选条件</span>
        <svg :class="{ rotated: filterExpanded }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:16px;height:16px">
          <path d="M6 9l6 6 6-6" />
        </svg>
      </div>
      <div class="filter-body" :class="{ expanded: filterExpanded }">
        <el-form inline :model="filters" @submit.prevent>
          <el-form-item label="状态">
            <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 130px" @change="reload">
              <el-option label="待处理" value="pending" />
              <el-option label="处理中" value="processing" />
              <el-option label="已回答" value="answered" />
              <el-option label="已关闭" value="closed" />
            </el-select>
          </el-form-item>
          <el-form-item label="来源">
            <el-select v-model="filters.source" placeholder="全部渠道" clearable style="width: 150px" @change="reload">
              <el-option label="HiAgent" value="hiagent_chat" />
              <el-option label="微信服务号" value="wechat_service" />
              <el-option label="微信订阅号" value="wechat_subscribe" />
              <el-option label="飞书" value="feishu" />
              <el-option label="易班" value="yiban" />
            </el-select>
          </el-form-item>
          <el-form-item label="时间">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              value-format="YYYY-MM-DD"
              start-placeholder="开始"
              end-placeholder="结束"
              range-separator="至"
              @change="onDateChange"
            />
          </el-form-item>
          <el-form-item>
            <el-button @click="reset" plain>重置</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>

    <!-- 表格卡片 -->
    <el-card class="table-card" shadow="never">
      <div class="table-wrapper">
        <el-table
          :data="rows"
          v-loading="loading"
          stripe
          class="ticket-table"
          :empty-text="loading ? '加载中...' : '暂无工单'"
        >
          <el-table-column prop="ticket_id" label="工单号" width="160">
            <template #default="{ row }">
              <span class="ticket-id" @click="$router.push(`/tickets/${row.ticket_id}`)">
                {{ row.ticket_id }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="question" label="用户问题" show-overflow-tooltip min-width="240">
            <template #default="{ row }">
              <div class="question-cell">{{ row.question }}</div>
            </template>
          </el-table-column>
          <el-table-column prop="source" label="来源" width="100">
            <template #default="{ row }">
              <span class="source-tag" :data-source="row.source">{{ shortSource(row.source) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="80">
            <template #default="{ row }">
              <span class="status-dot" :data-status="row.status"></span>
              <span class="status-text">{{ statusLabel(row.status) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="提交时间" width="155">
            <template #default="{ row }">
              <span class="time-text">{{ formatTime(row.created_at) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="80" fixed="right" align="center">
            <template #default="{ row }">
              <el-button size="small" link type="primary" @click="$router.push(`/tickets/${row.ticket_id}`)">
                查看
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <el-pagination
        v-model:current-page="filters.page"
        v-model:page-size="filters.page_size"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        :layout="paginationLayout"
        class="pagination"
        @current-change="reload"
        @size-change="reload"
      />
    </el-card>

    <!-- 模拟提交对话框 -->
    <el-dialog v-model="createVisible" title="模拟智能体提交工单" :width="dialogWidth" class="create-dialog" :destroy-on-close="true">
      <el-form :model="createForm" label-width="80px" label-position="right">
        <el-form-item label="问题" required>
          <el-input v-model="createForm.question" type="textarea" :rows="3" maxlength="500" show-word-limit placeholder="例如：石河子大学2026年计算机学院宿舍分配在哪？" />
        </el-form-item>
        <el-form-item label="用户 ID" required>
          <el-input v-model="createForm.user_id" maxlength="128" placeholder="HiAgent 平台 session_id" />
        </el-form-item>
        <el-form-item label="来源">
          <el-select v-model="createForm.source" style="width: 100%">
            <el-option label="HiAgent 对话" value="hiagent_chat" />
            <el-option label="微信服务号" value="wechat_service" />
            <el-option label="飞书" value="feishu" />
            <el-option label="易班" value="yiban" />
          </el-select>
        </el-form-item>
        <el-form-item label="RAG 结果">
          <el-input v-model="createForm.rag_result" type="textarea" :rows="2" placeholder="留空表示完全无结果" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" @click="doCreate">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { listTickets, submitTicket } from '../api/ticket'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const createVisible = ref(false)
const dateRange = ref([])
const filterExpanded = ref(true)

const isMobile = computed(() => window.innerWidth < 768)
const dialogWidth = computed(() => isMobile.value ? '95vw' : '540px')
const paginationLayout = computed(() => isMobile.value ? 'total, prev, pager, next' : 'total, sizes, prev, pager, next, jumper')

const filters = reactive({
  status: '',
  source: '',
  startDate: '',
  endDate: '',
  page: 1,
  page_size: 20,
})

const createForm = reactive({
  question: '',
  user_id: `sess_${Date.now()}`,
  source: 'hiagent_chat',
  rag_result: '',
})

const shortSource = (s) => ({
  hiagent_chat: 'HiAgent',
  wechat_service: '微信',
  wechat_subscribe: '订阅号',
  feishu: '飞书',
  yiban: '易班',
})[s] || s

const sourceLabel = (s) => ({
  hiagent_chat: 'HiAgent',
  wechat_service: '微信服务号',
  wechat_subscribe: '微信订阅号',
  feishu: '飞书',
  yiban: '易班',
})[s] || s

const statusLabel = (s) => ({
  pending: '待处理',
  processing: '处理中',
  answered: '已回答',
  closed: '已关闭',
})[s] || s

const formatTime = (iso) => {
  if (!iso) return ''
  const d = new Date(iso)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

// 统计
const stats = computed(() => {
  const pending = rows.value.filter((r) => r.status === 'pending').length
  const answered = rows.value.filter((r) => r.status === 'answered').length
  const closed = rows.value.filter((r) => r.status === 'closed').length
  return [
    { key: 'total', label: '当前页工单数', value: rows.value.length, meta: '当前筛选下', color: '#00b894' },
    { key: 'pending', label: '待处理', value: pending, meta: '需尽快响应', color: '#fdcb6e' },
    { key: 'answered', label: '已回答', value: answered, meta: '等待复问', color: '#00cec9' },
    { key: 'total_all', label: '总工单数', value: total.value, meta: '全部数据', color: '#6c5ce7' },
  ]
})

const reload = async () => {
  loading.value = true
  try {
    const data = await listTickets({
      status: filters.status || undefined,
      source: filters.source || undefined,
      start_date: filters.startDate || undefined,
      end_date: filters.endDate || undefined,
      page: filters.page,
      page_size: filters.page_size,
    })
    rows.value = data.items
    total.value = data.total
  } catch (e) {} finally {
    loading.value = false
  }
}

const onDateChange = (val) => {
  if (val && val.length === 2) {
    filters.startDate = val[0]
    filters.endDate = val[1]
  } else {
    filters.startDate = ''
    filters.endDate = ''
  }
  reload()
}

const reset = () => {
  filters.status = ''
  filters.source = ''
  filters.startDate = ''
  filters.endDate = ''
  dateRange.value = []
  filters.page = 1
  reload()
}

const openCreate = () => {
  createForm.question = ''
  createForm.user_id = `sess_${Date.now()}`
  createForm.source = 'hiagent_chat'
  createForm.rag_result = ''
  createVisible.value = true
}

const doCreate = async () => {
  if (!createForm.question.trim()) { ElMessage.warning('请填写问题'); return }
  if (!createForm.user_id.trim()) { ElMessage.warning('请填写用户 ID'); return }
  try {
    const r = await submitTicket(createForm)
    ElMessageBox.alert(r.message, `工单已创建：${r.ticket_id}`, {
      confirmButtonText: '查看详情',
    })
      .then(() => location.assign(`/tickets/${r.ticket_id}`))
      .catch(() => location.reload())
    createVisible.value = false
  } catch (e) {}
}

onMounted(() => {
  reload()
  // 根据窗口宽度决定筛选折叠初始状态
  if (window.innerWidth < 768) filterExpanded.value = false
})
</script>

<style scoped>
.ticket-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 4px;
}
.page-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 4px;
  letter-spacing: -0.5px;
}
@media (max-width: 768px) { .page-title { font-size: 18px; } }
.page-desc {
  color: var(--text-secondary);
  margin: 0;
  font-size: 13px;
}
@media (max-width: 768px) { .page-desc { font-size: 12px; } }

.btn-label { display: inline; }
@media (max-width: 768px) { .btn-label { display: none; } }

.action-btn {
  background: var(--gradient-header);
  border: none;
  box-shadow: 0 2px 8px rgba(0, 184, 148, 0.25);
  transition: all var(--transition-base);
}
.action-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 184, 148, 0.35);
  background: var(--gradient-header) !important;
}

/* 统计 */
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
}
@media (max-width: 900px) { .stats-row { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 480px) { .stats-row { grid-template-columns: repeat(2, 1fr); gap: 10px; } }

.stat-card { position: relative; }
.stat-card .stat-label { color: var(--text-secondary); font-size: 13px; margin-bottom: 6px; }
@media (max-width: 480px) { .stat-card .stat-label { font-size: 11px; } }
.stat-card .stat-value {
  color: var(--text-primary);
  font-size: 28px;
  font-weight: 600;
  letter-spacing: -0.5px;
  line-height: 1.2;
}
@media (max-width: 480px) { .stat-card .stat-value { font-size: 22px; } }
.stat-card .stat-meta { color: var(--text-tertiary); font-size: 12px; margin-top: 6px; }
@media (max-width: 480px) { .stat-card .stat-meta { font-size: 10px; } }
.stat-card::after {
  content: '';
  position: absolute;
  left: 0; top: 0; bottom: 0;
  width: 3px;
  background: var(--accent, var(--color-primary));
  opacity: 0;
  transition: opacity var(--transition-base);
}
.stat-card:hover::after { opacity: 1; }

/* 筛选折叠 */
.filter-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
  user-select: none;
}
.filter-header svg { transition: transform 0.2s; }
.filter-header svg.rotated { transform: rotate(180deg); }

.filter-body { max-height: 0; overflow: hidden; transition: max-height 0.3s ease; }
.filter-body.expanded { max-height: 300px; }

@media (min-width: 769px) {
  .filter-header { display: none; }
  .filter-body { max-height: 300px; }
}

/* 卡片 */
:deep(.filter-card),
:deep(.table-card) {
  border: 1px solid var(--border-soft);
  box-shadow: var(--shadow-sm);
  background: var(--bg-surface);
  border-radius: var(--radius-lg);
}

/* 表格横向滚动 */
.table-wrapper {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

/* 表格 */
:deep(.ticket-table .el-table__row) {
  transition: background var(--transition-fast);
}
:deep(.ticket-table .el-table__row:hover > td) {
  background: var(--color-primary-soft) !important;
}

.ticket-id {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--color-primary);
  cursor: pointer;
  font-weight: 500;
  transition: color var(--transition-fast);
  word-break: keep-all;
}
.ticket-id:hover { color: var(--color-primary-dark); text-decoration: underline; }

.question-cell {
  color: var(--text-primary);
  font-size: 14px;
  line-height: 1.5;
}

.source-tag {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 500;
  background: var(--bg-hover);
  color: var(--text-secondary);
}
.source-tag[data-source="hiagent_chat"] { background: rgba(0, 184, 148, 0.1); color: #00b894; }
.source-tag[data-source="feishu"] { background: rgba(108, 92, 231, 0.1); color: #6c5ce7; }
.source-tag[data-source="wechat_service"] { background: rgba(0, 184, 148, 0.1); color: #00b894; }
.source-tag[data-source="wechat_subscribe"] { background: rgba(253, 203, 110, 0.15); color: #e17055; }
.source-tag[data-source="yiban"] { background: rgba(116, 185, 255, 0.15); color: #0984e3; }

.status-dot {
  display: inline-block;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  margin-right: 5px;
  vertical-align: middle;
}
.status-dot[data-status="pending"] { background: #fdcb6e; box-shadow: 0 0 0 2px rgba(253, 203, 110, 0.2); }
.status-dot[data-status="processing"] { background: #74b9ff; box-shadow: 0 0 0 2px rgba(116, 185, 255, 0.2); }
.status-dot[data-status="answered"] { background: #00b894; box-shadow: 0 0 0 2px rgba(0, 184, 148, 0.2); }
.status-dot[data-status="closed"] { background: #b2bec3; box-shadow: 0 0 0 2px rgba(178, 190, 195, 0.2); }
.status-text { color: var(--text-primary); font-size: 12px; }

.time-text {
  color: var(--text-secondary);
  font-size: 12px;
  font-variant-numeric: tabular-nums;
}

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
  display: flex;
}
@media (max-width: 768px) {
  .pagination {
    justify-content: center;
    flex-wrap: wrap;
  }
}
</style>
