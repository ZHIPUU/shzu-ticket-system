<template>
  <div class="ticket-detail" v-loading="loading">
    <!-- 面包屑 -->
    <div class="breadcrumb">
      <span class="crumb" @click="$router.push('/tickets')">工单列表</span>
      <ChevronRight :size="14" :stroke-width="2" class="sep-icon" />
      <span class="crumb current">{{ ticket?.ticket_id || '加载中...' }}</span>
    </div>

    <template v-if="ticket">
      <!-- 头部摘要卡 -->
      <SectionCard :padded="false">
        <div class="header-pad">
          <div class="header-row">
            <div class="header-left">
              <div class="ticket-id-row">
                <code class="ticket-id">{{ ticket.ticket_id }}</code>
                <StatusBadge :status="ticket.status" />
                <ArchivedBadge v-if="ticket.archived" :archived-at="ticket.archived_at" />
                <CategoryBadge :category="ticket.category" />
              </div>
              <div class="header-meta">
                <span class="meta-item">
                  <Clock :size="14" :stroke-width="2" />
                  <span>创建于</span>
                  <DateTimeText :value="ticket.created_at" mode="datetime" />
                </span>
                <span v-if="ticket.answered_at" class="meta-item">
                  <MessageSquareCheck :size="14" :stroke-width="2" />
                  <span>答复于</span>
                  <DateTimeText :value="ticket.answered_at" mode="datetime" />
                </span>
                <span v-if="ticket.status === 'closed' && ticket.archived_at" class="meta-item">
                  <Archive :size="14" :stroke-width="2" />
                  <span>归档于</span>
                  <DateTimeText :value="ticket.archived_at" mode="datetime" />
                </span>
              </div>
            </div>
            <div class="header-right">
              <el-button v-if="ticket.archived" size="default" @click="onUnarchive" :loading="busy">
                <ArchiveRestore :size="16" :stroke-width="2" />
                <span>取消归档</span>
              </el-button>
              <el-button v-else size="default" @click="onArchive" :loading="busy">
                <Archive :size="16" :stroke-width="2" />
                <span>归档</span>
              </el-button>
              <el-dropdown trigger="click" @command="onAction">
                <el-button :loading="busy">
                  <span>更多</span>
                  <ChevronDown :size="14" :stroke-width="2" />
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item v-if="ticket.status === 'closed'" command="reopen">
                      <RotateCcw :size="14" :stroke-width="2" />
                      <span>重开工单</span>
                    </el-dropdown-item>
                    <el-dropdown-item command="delete" divided>
                      <span class="danger-text">
                        <Trash2 :size="14" :stroke-width="2" />
                        删除工单
                      </span>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </div>
      </SectionCard>

      <!-- 两列布局：左侧详情 / 右侧操作面板 -->
      <div class="detail-grid">
        <!-- 左 -->
        <div class="detail-main">
          <!-- 用户问题 -->
          <SectionCard title="用户问题" icon="HelpCircle" :padded="false">
            <div class="content-pad">
              <div class="question-box">{{ ticket.question }}</div>
            </div>
          </SectionCard>

          <!-- 答复（内联编辑器） -->
          <SectionCard title="人工答复" icon="MessageSquareText" :padded="false">
            <template #header>
              <span v-if="ticket.answered_at" class="header-tag">
                <MessageSquareCheck :size="12" :stroke-width="2.4" />
                <span>已答复</span>
              </span>
              <span v-else class="header-tag warn">
                <MessageSquare :size="12" :stroke-width="2.4" />
                <span>待答复</span>
              </span>
            </template>
            <div class="content-pad">
              <div v-if="!editingAnswer" class="answer-view">
                <div v-if="ticket.answer" class="answer-text">{{ ticket.answer }}</div>
                <EmptyState
                  v-else
                  icon="MessageSquare"
                  title="尚未答复"
                  description="点击下方按钮开始答复"
                />
                <div class="answer-actions">
                  <el-button v-if="ticket.status !== 'closed'" type="primary" @click="enterEditAnswer">
                    <Pencil v-if="ticket.answer" :size="14" :stroke-width="2" />
                    <Plus v-else :size="14" :stroke-width="2" />
                    <span style="margin-left: 4px">{{ ticket.answer ? '编辑答复' : '撰写答复' }}</span>
                  </el-button>
                  <el-button v-else disabled>
                    工单已关闭
                  </el-button>
                </div>
              </div>
              <div v-else class="answer-edit">
                <el-input
                  v-model="answerDraft"
                  type="textarea"
                  :rows="8"
                  maxlength="5000"
                  show-word-limit
                  placeholder="请输入标准答复内容（可重复编辑覆盖）"
                  ref="answerInputRef"
                />
                <div class="edit-actions">
                  <el-button @click="cancelEditAnswer">取消</el-button>
                  <el-button type="primary" :loading="answerSubmitting" @click="saveAnswer">
                    <Save :size="14" :stroke-width="2" />
                    <span style="margin-left: 4px">保存答复</span>
                  </el-button>
                </div>
              </div>
            </div>
          </SectionCard>

          <!-- 关闭信息 -->
          <SectionCard v-if="ticket.status === 'closed' && ticket.close_reason" title="关闭原因" icon="Lock" :padded="false">
            <div class="content-pad">
              <div class="close-reason">{{ ticket.close_reason }}</div>
            </div>
          </SectionCard>
        </div>

        <!-- 右 -->
        <div class="detail-side">
          <!-- 元信息 -->
          <SectionCard title="工单信息" icon="Info" :padded="false">
            <div class="content-pad">
              <el-descriptions :column="1" border size="small" class="info-desc">
                <el-descriptions-item label="提交人">
                  <span class="mono">{{ ticket.user_id }}</span>
                </el-descriptions-item>
                <el-descriptions-item label="来源">
                  <span>{{ SOURCE_LABEL[ticket.source] || ticket.source || '—' }}</span>
                </el-descriptions-item>
                <el-descriptions-item label="分类">
                  <CategoryBadge :category="ticket.category" />
                </el-descriptions-item>
                <el-descriptions-item label="状态">
                  <StatusBadge :status="ticket.status" />
                </el-descriptions-item>
                <el-descriptions-item label="归档">
                  <span v-if="ticket.archived" class="archived-text">已归档</span>
                  <span v-else class="muted">未归档</span>
                </el-descriptions-item>
                <el-descriptions-item label="创建时间">
                  <DateTimeText :value="ticket.created_at" mode="datetime" />
                </el-descriptions-item>
                <el-descriptions-item v-if="ticket.answered_at" label="答复时间">
                  <DateTimeText :value="ticket.answered_at" mode="datetime" />
                </el-descriptions-item>
                <el-descriptions-item v-if="ticket.archived_at" label="归档时间">
                  <DateTimeText :value="ticket.archived_at" mode="datetime" />
                </el-descriptions-item>
              </el-descriptions>
            </div>
          </SectionCard>

          <!-- 操作 -->
          <SectionCard title="操作" icon="Settings2" :padded="false">
            <div class="content-pad">
              <div class="action-list">
                <el-button
                  v-if="ticket.status !== 'closed'"
                  class="action-btn"
                  type="primary"
                  :loading="busy"
                  @click="enterEditAnswer"
                >
                  <MessageSquareText :size="16" :stroke-width="2" />
                  <span>{{ ticket.answer ? '编辑答复' : '答复工单' }}</span>
                </el-button>
                <el-button
                  v-if="ticket.status !== 'closed'"
                  class="action-btn"
                  :loading="busy"
                  @click="onClose"
                >
                  <Lock :size="16" :stroke-width="2" />
                  <span>关闭工单</span>
                </el-button>
                <el-button
                  v-if="ticket.status === 'closed'"
                  class="action-btn"
                  :loading="busy"
                  @click="onReopen"
                >
                  <RotateCcw :size="16" :stroke-width="2" />
                  <span>重开工单</span>
                </el-button>
                <el-button
                  v-if="!ticket.archived"
                  class="action-btn"
                  :loading="busy"
                  @click="onArchive"
                >
                  <Archive :size="16" :stroke-width="2" />
                  <span>归档工单</span>
                </el-button>
                <el-button
                  v-else
                  class="action-btn"
                  :loading="busy"
                  @click="onUnarchive"
                >
                  <ArchiveRestore :size="16" :stroke-width="2" />
                  <span>取消归档</span>
                </el-button>
                <el-button
                  class="action-btn danger"
                  :loading="busy"
                  @click="onDelete"
                >
                  <Trash2 :size="16" :stroke-width="2" />
                  <span>删除工单</span>
                </el-button>
              </div>
            </div>
          </SectionCard>

          <!-- 关联信息（RAG 检索结果） -->
          <SectionCard v-if="ticket.rag_result" title="知识库检索" icon="BookOpen" :padded="false">
            <div class="content-pad">
              <pre class="rag-content">{{ ticket.rag_result }}</pre>
            </div>
          </SectionCard>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ChevronRight,
  ChevronDown,
  Clock,
  MessageSquare,
  MessageSquareCheck,
  MessageSquareText,
  HelpCircle,
  Archive,
  ArchiveRestore,
  RotateCcw,
  Trash2,
  Pencil,
  Plus,
  Save,
  Lock,
} from '@lucide/vue'
import {
  getTicket,
  answerTicket,
  closeTicket,
  reopenTicket,
  archiveTicket,
  deleteTicket,
} from '../api/ticket'
import SectionCard from '../components/common/SectionCard.vue'
import StatusBadge from '../components/common/StatusBadge.vue'
import CategoryBadge from '../components/common/CategoryBadge.vue'
import ArchivedBadge from '../components/common/ArchivedBadge.vue'
import EmptyState from '../components/common/EmptyState.vue'
import DateTimeText from '../components/common/DateTimeText.vue'

const route = useRoute()
const router = useRouter()
const loading = ref(true)
const ticket = ref(null)
const busy = ref(false)
const editingAnswer = ref(false)
const answerDraft = ref('')
const answerSubmitting = ref(false)
const answerInputRef = ref(null)

const SOURCE_LABEL = {
  hiagent_chat: 'HiAgent 对话',
  wechat_service: '微信服务号',
  wechat_subscribe: '微信订阅号',
  feishu: '飞书',
  yiban: '易班',
}

const fetch = async () => {
  loading.value = true
  try {
    ticket.value = await getTicket(route.params.id)
  } catch (e) {
    router.push('/tickets')
  } finally {
    loading.value = false
  }
}

const enterEditAnswer = async () => {
  if (ticket.value?.status === 'closed') {
    ElMessage.warning('工单已关闭，无法答复')
    return
  }
  answerDraft.value = ticket.value?.answer || ''
  editingAnswer.value = true
  await nextTick()
  if (answerInputRef.value?.focus) {
    answerInputRef.value.focus()
  }
}

const cancelEditAnswer = () => {
  editingAnswer.value = false
  answerDraft.value = ''
}

const saveAnswer = async () => {
  const text = answerDraft.value.trim()
  if (!text) {
    ElMessage.warning('请填写答复内容')
    return
  }
  answerSubmitting.value = true
  try {
    await answerTicket(route.params.id, { answer: text })
    ElMessage.success('答复已保存')
    editingAnswer.value = false
    await fetch()
  } catch (e) {} finally {
    answerSubmitting.value = false
  }
}

const onClose = () => {
  ElMessageBox.prompt('请输入关闭原因（可选）', '关闭工单', {
    confirmButtonText: '确认关闭',
    cancelButtonText: '取消',
    inputPlaceholder: '如：重复工单、无效问题等',
  })
    .then(async ({ value }) => {
      busy.value = true
      try {
        await closeTicket(route.params.id, value || '')
        ElMessage.success('工单已关闭')
        await fetch()
      } catch (e) {} finally { busy.value = false }
    })
    .catch(() => {})
}

const onReopen = async () => {
  busy.value = true
  try {
    await reopenTicket(route.params.id)
    ElMessage.success('工单已重开')
    await fetch()
  } catch (e) {} finally { busy.value = false }
}

const onArchive = async () => {
  busy.value = true
  try {
    await archiveTicket(route.params.id, true)
    ElMessage.success('工单已归档')
    await fetch()
  } catch (e) {} finally { busy.value = false }
}

const onUnarchive = async () => {
  busy.value = true
  try {
    await archiveTicket(route.params.id, false)
    ElMessage.success('已取消归档')
    await fetch()
  } catch (e) {} finally { busy.value = false }
}

const onDelete = () => {
  ElMessageBox.confirm(
    '确认删除该工单？软删后可在筛选中包含已删除查看。',
    '删除工单',
    { type: 'warning', confirmButtonText: '确认删除', cancelButtonText: '取消' }
  )
    .then(async () => {
      busy.value = true
      try {
        await deleteTicket(route.params.id, false)
        ElMessage.success('工单已删除')
        router.push('/tickets')
      } catch (e) {} finally { busy.value = false }
    })
    .catch(() => {})
}

onMounted(fetch)
</script>

<style scoped>
.ticket-detail {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--text-tertiary);
  margin-bottom: 4px;
}
.crumb {
  cursor: pointer;
  transition: color var(--transition-fast);
}
.crumb:hover { color: var(--color-primary); }
.crumb.current { color: var(--text-primary); cursor: default; }
.crumb.current:hover { color: var(--text-primary); }
.sep-icon { color: var(--text-tertiary); flex-shrink: 0; }

.header-pad { padding: 16px 20px; }

.header-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  flex-wrap: wrap;
}
.header-left { flex: 1; min-width: 0; }
.header-right {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
  flex-wrap: wrap;
}

.ticket-id-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
.ticket-id {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  background: var(--bg-base);
  padding: 3px 8px;
  border-radius: var(--radius-sm);
  letter-spacing: -0.2px;
}

.header-meta {
  display: flex;
  gap: 16px;
  margin-top: 10px;
  flex-wrap: wrap;
}
.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12.5px;
  color: var(--text-secondary);
}
.meta-item :first-child { color: var(--text-tertiary); }

.header-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--color-success);
  background: var(--color-success-soft);
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}
.header-tag.warn {
  color: var(--color-warning);
  background: var(--color-warning-soft);
}

.content-pad { padding: 20px; }

.detail-grid {
  display: grid;
  grid-template-columns: minmax(0, 6fr) minmax(0, 4fr);
  gap: 16px;
  align-items: start;
}
.detail-main, .detail-side {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}
@media (max-width: 1024px) {
  .detail-grid { grid-template-columns: 1fr; }
}

.question-box {
  background: var(--gradient-soft);
  border-left: 3px solid var(--color-primary);
  padding: 16px 20px;
  border-radius: var(--radius-md);
  white-space: pre-wrap;
  line-height: 1.7;
  color: var(--text-primary);
  font-size: 15px;
  word-break: break-word;
}
@media (max-width: 768px) { .question-box { padding: 12px 14px; font-size: 14px; } }

.answer-view {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.answer-text {
  background: var(--color-primary-soft);
  border-left: 3px solid var(--color-primary);
  padding: 16px 20px;
  border-radius: var(--radius-md);
  white-space: pre-wrap;
  line-height: 1.8;
  color: var(--text-primary);
  font-size: 15px;
  word-break: break-word;
}
@media (max-width: 768px) { .answer-text { padding: 14px 16px; font-size: 14px; } }

.answer-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.answer-edit {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.edit-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.close-reason {
  padding: 12px 16px;
  background: var(--bg-base);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.7;
  white-space: pre-wrap;
}

.mono {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12.5px;
  color: var(--text-primary);
  word-break: break-all;
}
.muted { color: var(--text-tertiary); }
.archived-text { color: var(--text-secondary); }
.danger-text {
  color: var(--color-danger);
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.info-desc :deep(.el-descriptions__label) {
  color: var(--text-secondary);
  width: 80px;
  font-size: 12.5px;
}

.action-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.action-btn {
  justify-content: flex-start;
  width: 100%;
  height: 38px;
}
.action-btn.danger {
  color: var(--color-danger);
  border-color: var(--border-color);
  background: var(--bg-surface);
}
.action-btn.danger:hover {
  background: var(--color-danger-soft);
  border-color: var(--color-danger);
  color: var(--color-danger);
}
.action-btn :deep(svg) { flex-shrink: 0; }

.rag-content {
  background: var(--bg-base);
  border: 1px solid var(--border-soft);
  padding: 12px 14px;
  border-radius: var(--radius-md);
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12.5px;
  color: var(--text-secondary);
  white-space: pre-wrap;
  word-break: break-word;
  margin: 0;
  max-height: 320px;
  overflow: auto;
}
</style>
