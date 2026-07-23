<template>
  <div class="ticket-list">
    <PageHeader title="工单列表" description="管理来自各渠道的用户提问工单">
      <template #actions>
        <el-button type="primary" @click="openCreate">
          <Plus :size="16" :stroke-width="2.4" />
          <span>模拟提交</span>
        </el-button>
      </template>
    </PageHeader>

    <!-- 统计卡片 -->
    <div class="stats-row">
      <StatCard
        icon="Inbox"
        label="当前页工单数"
        :value="rows.length"
        meta="当前筛选下"
        accent="#2563EB"
      />
      <StatCard
        icon="Clock"
        label="待处理"
        :value="counts.pending"
        meta="需尽快响应"
        accent="#F59E0B"
      />
      <StatCard
        icon="MessageSquareCheck"
        label="已答复"
        :value="counts.answered"
        meta="等待复问"
        accent="#10B981"
      />
      <StatCard
        icon="Database"
        label="总工单数"
        :value="total"
        meta="全部数据"
        accent="#8B5CF6"
      />
    </div>

    <!-- 筛选 -->
    <SectionCard :padded="false">
      <div class="filter-pad">
        <div class="filter-header" @click="filterExpanded = !filterExpanded">
          <div class="filter-title">
            <SlidersHorizontal :size="14" :stroke-width="2" />
            <span>筛选条件</span>
          </div>
          <ChevronDown :size="16" :stroke-width="2" :class="{ rotated: filterExpanded }" class="chev" />
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
              <el-button @click="reset" plain>
                <RotateCcw :size="14" :stroke-width="2" />
                <span style="margin-left: 4px">重置</span>
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </SectionCard>

    <!-- 表格卡片 -->
    <SectionCard :padded="false">
      <div class="table-wrapper">
        <el-table
          :data="rows"
          v-loading="loading"
          stripe
          class="ticket-table"
          :empty-text="loading ? '加载中...' : '暂无工单'"
        >
          <el-table-column prop="ticket_id" label="工单号" width="180">
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
          <el-table-column prop="status" label="状态" width="90">
            <template #default="{ row }">
              <StatusBadge :status="row.status" />
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="155">
            <template #default="{ row }">
              <DateTimeText :value="row.created_at" mode="datetime" />
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

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="filters.page"
          v-model:page-size="filters.page_size"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          :layout="paginationLayout"
          @current-change="reload"
          @size-change="reload"
        />
      </div>
    </SectionCard>

    <!-- 模拟提交对话框 -->
    <el-dialog v-model="createVisible" title="模拟智能体提交工单" :width="dialogWidth" class="create-dialog" :destroy-on-close="true">
      <el-form :model="createForm" label-width="80px" label-position="right">
        <el-form-item label="问题" required>
          <el-input v-model="createForm.question" type="textarea" :rows="3" maxlength="500" show-word-limit placeholder="例如：石河子大学2026年计算机学院宿舍分配在哪？" />
        </el-form-item>
        <el-form-item label="用户 ID" required>
          <el-input v-model="createForm.user_id" maxlength="128" placeholder="HiAgent 平台 session_id" />
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
import { Plus, SlidersHorizontal, ChevronDown, RotateCcw } from '@lucide/vue'
import PageHeader from '../components/common/PageHeader.vue'
import SectionCard from '../components/common/SectionCard.vue'
import StatCard from '../components/common/StatCard.vue'
import StatusBadge from '../components/common/StatusBadge.vue'
import DateTimeText from '../components/common/DateTimeText.vue'

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const createVisible = ref(false)
const dateRange = ref([])
const filterExpanded = ref(true)

const isMobile = computed(() => window.innerWidth < 768)
const dialogWidth = computed(() => isMobile.value ? '95vw' : '540px')
const paginationLayout = computed(() =>
  isMobile.value ? 'total, prev, pager, next' : 'total, sizes, prev, pager, next, jumper'
)

const filters = reactive({
  status: '',
  startDate: '',
  endDate: '',
  page: 1,
  page_size: 20,
})

const createForm = reactive({
  question: '',
  user_id: `sess_${Date.now()}`,
  rag_result: '',
})

const counts = computed(() => ({
  pending: rows.value.filter((r) => r.status === 'pending').length,
  answered: rows.value.filter((r) => r.status === 'answered').length,
  closed: rows.value.filter((r) => r.status === 'closed').length,
}))

const reload = async () => {
  loading.value = true
  try {
    const data = await listTickets({
      status: filters.status || undefined,
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
  filters.startDate = ''
  filters.endDate = ''
  dateRange.value = []
  filters.page = 1
  reload()
}

const openCreate = () => {
  createForm.question = ''
  createForm.user_id = `sess_${Date.now()}`
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
  if (window.innerWidth < 768) filterExpanded.value = false
})
</script>

<style scoped>
.ticket-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.btn-label { display: inline; }
@media (max-width: 768px) { .btn-label { display: none; } }

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}
@media (max-width: 900px) { .stats-row { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 480px) { .stats-row { grid-template-columns: repeat(2, 1fr); gap: 8px; } }

.filter-pad { padding: 0; }
.filter-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
  user-select: none;
}
.filter-title { display: inline-flex; align-items: center; gap: 6px; }
.chev { transition: transform 0.2s; color: var(--text-tertiary); }
.chev.rotated { transform: rotate(180deg); }

.filter-body {
  max-height: 0;
  overflow: hidden;
  transition: max-height 0.3s ease;
  border-top: 1px solid transparent;
}
.filter-body.expanded {
  max-height: 300px;
  padding: 0 20px 16px;
  border-top-color: var(--border-soft);
}

@media (min-width: 769px) {
  .filter-header { display: none; }
  .filter-body { max-height: 300px; padding: 0 20px 16px; border-top: 1px solid var(--border-soft); }
}

.table-wrapper {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

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

.pagination-wrap {
  padding: 16px 20px;
  display: flex;
  justify-content: flex-end;
}
@media (max-width: 768px) {
  .pagination-wrap { justify-content: center; }
}
</style>
