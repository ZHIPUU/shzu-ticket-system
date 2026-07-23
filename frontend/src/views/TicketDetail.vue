<template>
  <div class="ticket-detail" v-loading="loading">
    <!-- 面包屑 -->
    <div class="breadcrumb">
      <span class="crumb" @click="$router.push('/tickets')">工单列表</span>
      <span class="sep">/</span>
      <span class="crumb current">{{ ticket?.ticket_id || '加载中...' }}</span>
    </div>

    <template v-if="ticket">
      <!-- 头部信息卡 -->
      <el-card class="header-card" shadow="never">
        <div class="header-content">
          <div class="header-main">
            <div class="ticket-id-large">
              <span class="status-dot" :data-status="ticket.status"></span>
              <span class="id-text">{{ ticket.ticket_id }}</span>
              <span class="status-pill" :data-status="ticket.status">{{ statusLabel(ticket.status) }}</span>
            </div>
            <div class="meta-row">
              <span class="meta-item">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/></svg>
                {{ formatTime(ticket.created_at) }}
              </span>
              <span class="meta-item">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
                {{ sourceLabel(ticket.source) }}
              </span>
              <span class="meta-item">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                {{ ticket.user_id }}
              </span>
            </div>
          </div>
          <div class="header-actions" v-if="ticket.status !== 'closed'">
            <el-button type="primary" size="default" @click="answerVisible = true">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:16px;height:16px;margin-right:6px"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
              答复
            </el-button>
            <el-button @click="onClose" plain size="default">关闭</el-button>
          </div>
        </div>
      </el-card>

      <!-- 用户问题 -->
      <el-card class="section-card" shadow="never">
        <template #header>
          <div class="section-header">
            <span class="section-icon">💬</span>
            <span class="section-title">用户问题</span>
          </div>
        </template>
        <div class="question-box">{{ ticket.question }}</div>
        <div v-if="ticket.rag_result" class="rag-result">
          <div class="rag-label">知识库检索结果</div>
          <div class="rag-content">{{ ticket.rag_result }}</div>
        </div>
        <div v-else class="rag-empty">知识库完全无结果</div>
      </el-card>

      <!-- 人工答复 -->
      <el-card class="section-card" shadow="never">
        <template #header>
          <div class="section-header">
            <span class="section-icon">✨</span>
            <span class="section-title">人工答复</span>
            <span v-if="ticket.answered_at" class="section-meta">
              {{ ticket.answered_by }} · {{ formatTime(ticket.answered_at) }}
            </span>
          </div>
        </template>
        <div v-if="ticket.answer" class="answer-box">
          {{ ticket.answer }}
        </div>
        <el-empty v-else description="尚未答复" :image-size="80" />
      </el-card>

      <!-- 关闭信息 -->
      <el-card v-if="ticket.status === 'closed' && ticket.close_reason" class="section-card closed-card" shadow="never">
        <template #header>
          <div class="section-header">
            <span class="section-icon">🔒</span>
            <span class="section-title">工单已关闭</span>
          </div>
        </template>
        <div class="close-reason">{{ ticket.close_reason }}</div>
      </el-card>
    </template>

    <!-- 答复弹窗 -->
    <el-dialog v-model="answerVisible" title="人工答复" :width="dialogWidth" class="answer-dialog" :destroy-on-close="true">
      <el-form :model="answerForm" label-width="84px" label-position="right">
        <el-form-item label="答复内容" required>
          <el-input v-model="answerForm.answer" type="textarea" :rows="8" maxlength="5000" show-word-limit placeholder="请输入标准答案" />
        </el-form-item>
        <el-form-item label="答复人" required>
          <el-input v-model="answerForm.operator" maxlength="64" placeholder="姓名或工号" />
        </el-form-item>
        <el-form-item label="同步知识库">
          <el-switch v-model="answerForm.sync_to_kb" />
          <span style="margin-left: 12px; color: var(--text-tertiary); font-size: 12px">（预留接口）</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="answerVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="doAnswer">提交答复</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getTicket, answerTicket, closeTicket } from '../api/ticket'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const loading = ref(true)
const ticket = ref(null)
const answerVisible = ref(false)
const submitting = ref(false)

const isMobile = computed(() => window.innerWidth < 768)
const dialogWidth = computed(() => isMobile.value ? '95vw' : '640px')

const answerForm = reactive({
  answer: '',
  operator: '',
  sync_to_kb: false,
})

const sourceLabel = (s) => ({
  hiagent_chat: 'HiAgent 对话',
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
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
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

const doAnswer = async () => {
  if (!answerForm.answer.trim()) { ElMessage.warning('请填写答复内容'); return }
  if (!answerForm.operator.trim()) { ElMessage.warning('请填写答复人'); return }
  submitting.value = true
  try {
    await answerTicket(route.params.id, answerForm)
    ElMessage.success('答复成功')
    answerVisible.value = false
    answerForm.answer = ''
    answerForm.operator = ''
    await fetch()
  } catch (e) {} finally {
    submitting.value = false
  }
}

const onClose = () => {
  ElMessageBox.prompt('请输入关闭原因（可选）', '关闭工单', {
    confirmButtonText: '确认关闭',
    cancelButtonText: '取消',
    inputPlaceholder: '如：重复工单、无效问题等',
  })
    .then(async ({ value }) => {
      await closeTicket(route.params.id, value || '')
      ElMessage.success('工单已关闭')
      await fetch()
    })
    .catch(() => {})
}

onMounted(fetch)
</script>

<style scoped>
.ticket-detail {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
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
.sep { color: var(--text-tertiary); }

/* 头部 */
.header-card {
  background: var(--bg-surface);
  border: 1px solid var(--border-soft);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}
.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}
@media (max-width: 768px) {
  .header-content { flex-direction: column; align-items: flex-start; gap: 12px; }
}

.ticket-id-large {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}
.id-text {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  letter-spacing: -0.3px;
}
@media (max-width: 768px) { .id-text { font-size: 16px; } }

.status-dot {
  display: inline-block;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}
.status-dot[data-status="pending"] { background: #fdcb6e; box-shadow: 0 0 0 3px rgba(253, 203, 110, 0.2); }
.status-dot[data-status="processing"] { background: #74b9ff; box-shadow: 0 0 0 3px rgba(116, 185, 255, 0.2); }
.status-dot[data-status="answered"] { background: #00b894; box-shadow: 0 0 0 3px rgba(0, 184, 148, 0.2); }
.status-dot[data-status="closed"] { background: #b2bec3; box-shadow: 0 0 0 3px rgba(178, 190, 195, 0.2); }

.status-pill {
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  background: var(--bg-hover);
  color: var(--text-secondary);
}
.status-pill[data-status="pending"] { background: rgba(253, 203, 110, 0.15); color: #d68f00; }
.status-pill[data-status="processing"] { background: rgba(116, 185, 255, 0.15); color: #0984e3; }
.status-pill[data-status="answered"] { background: rgba(0, 184, 148, 0.15); color: #00b894; }
.status-pill[data-status="closed"] { background: rgba(178, 190, 195, 0.15); color: #636e72; }

.meta-row {
  display: flex;
  gap: 20px;
  margin-top: 12px;
  flex-wrap: wrap;
}
@media (max-width: 768px) { .meta-row { gap: 12px; } }
.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--text-secondary);
}
@media (max-width: 768px) { .meta-item { font-size: 12px; } }
.meta-item svg { width: 14px; height: 14px; opacity: 0.6; flex-shrink: 0; }

.header-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}
@media (max-width: 768px) { .header-actions { width: 100%; } }
@media (max-width: 768px) { .header-actions .el-button { flex: 1; } }

/* Section */
.section-card {
  background: var(--bg-surface);
  border: 1px solid var(--border-soft);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}
.section-card.closed-card { background: var(--bg-base); }
.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
.section-icon { font-size: 16px; }
.section-title {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 15px;
}
.section-meta {
  margin-left: auto;
  font-size: 12px;
  color: var(--text-tertiary);
  font-weight: normal;
}
@media (max-width: 768px) { .section-meta { margin-left: 0; width: 100%; } }

.question-box {
  background: var(--gradient-soft);
  border-left: 3px solid var(--color-primary);
  padding: 16px 20px;
  border-radius: 8px;
  white-space: pre-wrap;
  line-height: 1.7;
  color: var(--text-primary);
  font-size: 15px;
}
@media (max-width: 768px) { .question-box { padding: 12px 14px; font-size: 14px; } }

.rag-result {
  margin-top: 16px;
  padding: 12px 16px;
  background: var(--bg-base);
  border-radius: 8px;
  border: 1px solid var(--border-soft);
}
.rag-label {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-bottom: 6px;
  font-weight: 500;
}
.rag-content {
  font-size: 13px;
  color: var(--text-secondary);
  white-space: pre-wrap;
  font-family: 'SF Mono', Menlo, Consolas, monospace;
}
.rag-empty {
  margin-top: 16px;
  padding: 12px 16px;
  background: var(--bg-base);
  border-radius: 8px;
  border: 1px dashed var(--border-color);
  color: var(--text-tertiary);
  font-size: 13px;
  text-align: center;
}

.answer-box {
  background: linear-gradient(135deg, #f0fbf8 0%, #e6f9f5 100%);
  border-left: 3px solid var(--color-primary);
  padding: 20px 24px;
  border-radius: 8px;
  white-space: pre-wrap;
  line-height: 1.8;
  color: var(--text-primary);
  font-size: 15px;
}
html.dark .answer-box {
  background: rgba(0, 184, 148, 0.08);
  border-left-color: var(--color-primary);
}
@media (max-width: 768px) { .answer-box { padding: 14px 16px; font-size: 14px; } }

.close-reason {
  padding: 12px 16px;
  background: var(--bg-base);
  border-radius: 8px;
  color: var(--text-secondary);
  font-size: 14px;
}
</style>
