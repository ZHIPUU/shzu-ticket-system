<template>
  <div class="detail">
    <!-- 加载骨架 -->
    <div v-if="loading" class="detail-skeleton">
      <UiSkeleton type="card" class="anim-fade-up" />
      <div class="sk-grid">
        <UiSkeleton type="lines" :count="5" />
        <UiSkeleton type="lines" :count="4" />
      </div>
    </div>

    <template v-else-if="ticket">
      <!-- ═══ 返回 + 头部 ═══ -->
      <button class="back-link anim-fade-up" @click="$router.push('/tickets')">
        <ArrowLeft :size="15" :stroke-width="2.2" />
        <span>返回列表</span>
      </button>

      <div class="head-card card anim-stagger" style="--i: 1">
        <div class="head-main">
          <div class="head-id-row">
            <button class="tid-chip mono" :title="'点击复制 ' + ticket.ticket_id" @click="copyId">
              {{ ticket.ticket_id }}
              <component :is="copied ? Check : Copy" :size="12" :stroke-width="2.2" class="copy-icon" :class="{ ok: copied }" />
            </button>
            <StatusPill :status="ticket.status" />
            <span v-if="ticket.archived" class="archived-pill" :title="`归档于 ${ticket.archived_at || ''}`">
              <Archive :size="11" :stroke-width="2.4" />已归档
            </span>
            <CategoryTag :category="ticket.category" />
          </div>
          <div class="head-meta">
            <span class="meta-item">
              <Clock :size="13" :stroke-width="2" />
              创建于 <TimeText :value="ticket.created_at" />
            </span>
            <span v-if="ticket.answered_at" class="meta-item">
              <MessageSquareCheck :size="13" :stroke-width="2" />
              答复于 <TimeText :value="ticket.answered_at" />
            </span>
            <span v-if="ticket.archived_at" class="meta-item">
              <Archive :size="13" :stroke-width="2" />
              归档于 <TimeText :value="ticket.archived_at" />
            </span>
          </div>
        </div>

        <div class="head-actions">
          <UiButton
            v-if="ticket.status !== 'closed'"
            variant="primary"
            :loading="busy"
            @click="enterEditAnswer"
          >
            <template #icon><MessageSquareText :size="15" :stroke-width="2" /></template>
            {{ ticket.answer ? '编辑答复' : '答复工单' }}
          </UiButton>
          <UiButton v-if="!ticket.archived" variant="ghost" :icon-only="isMobile" :loading="busy" @click="onArchive">
            <template #icon><Archive :size="15" :stroke-width="2" /></template>
            <span v-if="!isMobile">归档</span>
          </UiButton>
          <UiButton v-else variant="ghost" :icon-only="isMobile" :loading="busy" @click="onUnarchive">
            <template #icon><ArchiveRestore :size="15" :stroke-width="2" /></template>
            <span v-if="!isMobile">取消归档</span>
          </UiButton>
          <UiDropdown :items="moreActions" placement="end" @select="onMoreAction">
            <template #trigger>
              <UiButton variant="ghost" icon-only aria-label="更多操作">
                <template #icon><Ellipsis :size="16" :stroke-width="2" /></template>
              </UiButton>
            </template>
          </UiDropdown>
        </div>
      </div>

      <!-- ═══ 双列布局 ═══ -->
      <div class="detail-grid">
        <!-- 左列 -->
        <div class="col-main">
          <!-- 用户问题 -->
          <section class="card block anim-stagger" style="--i: 2">
            <header class="block-head">
              <span class="block-icon"><CircleQuestionMark :size="16" :stroke-width="2" /></span>
              <h3>用户问题</h3>
              <span class="from-chip mono" :title="ticket.user_id">
                <User :size="12" :stroke-width="2" />{{ ticket.user_id }}
              </span>
            </header>
            <div class="block-body">
              <blockquote class="question-quote">{{ ticket.question }}</blockquote>
            </div>
          </section>

          <!-- 人工答复 -->
          <section class="card block anim-stagger" style="--i: 3">
            <header class="block-head">
              <span class="block-icon"><MessageSquareText :size="16" :stroke-width="2" /></span>
              <h3>人工答复</h3>
              <span class="answer-state" :class="ticket.answer ? 'done' : 'todo'">
                <component :is="ticket.answer ? MessageSquareCheck : MessageSquare" :size="11" :stroke-width="2.4" />
                {{ ticket.answer ? '已答复' : '待答复' }}
              </span>
            </header>
            <div class="block-body">
              <!-- 查看态 -->
              <div v-if="!editingAnswer" class="answer-view">
                <div v-if="ticket.answer" class="answer-text">{{ ticket.answer }}</div>
                <UiEmpty
                  v-else
                  icon="MessageSquare"
                  title="尚未答复"
                  description="用户还在等待人工回复，尽快处理吧"
                />
                <div class="answer-actions">
                  <UiButton v-if="ticket.status !== 'closed'" variant="primary" @click="enterEditAnswer">
                    <template #icon>
                      <Pencil v-if="ticket.answer" :size="14" :stroke-width="2" />
                      <Plus v-else :size="14" :stroke-width="2.2" />
                    </template>
                    {{ ticket.answer ? '编辑答复' : '撰写答复' }}
                  </UiButton>
                  <span v-else class="closed-tip">
                    <Lock :size="13" :stroke-width="2" />工单已关闭，如需答复请先重开
                  </span>
                </div>
              </div>

              <!-- 编辑态 -->
              <div v-else class="answer-edit">
                <UiInput
                  ref="answerInputRef"
                  v-model="answerDraft"
                  type="textarea"
                  :rows="8"
                  :maxlength="5000"
                  show-word-limit
                  placeholder="请输入标准答复内容（可重复编辑覆盖）"
                />
                <div class="edit-actions">
                  <UiButton variant="ghost" @click="cancelEditAnswer">取消</UiButton>
                  <UiButton variant="primary" :loading="answerSubmitting" @click="saveAnswer">
                    <template #icon><Save :size="14" :stroke-width="2" /></template>
                    保存答复
                  </UiButton>
                </div>
              </div>
            </div>
          </section>

          <!-- 关闭原因 -->
          <section v-if="ticket.status === 'closed' && ticket.close_reason" class="card block anim-stagger" style="--i: 4">
            <header class="block-head">
              <span class="block-icon warn"><Lock :size="16" :stroke-width="2" /></span>
              <h3>关闭原因</h3>
            </header>
            <div class="block-body">
              <p class="close-reason">{{ ticket.close_reason }}</p>
            </div>
          </section>
        </div>

        <!-- 右列 -->
        <div class="col-side">
          <!-- 工单信息 -->
          <section class="card block anim-stagger" style="--i: 3">
            <header class="block-head">
              <span class="block-icon"><Info :size="16" :stroke-width="2" /></span>
              <h3>工单信息</h3>
            </header>
            <div class="block-body info-list">
              <div class="info-row">
                <span class="info-label">提交人</span>
                <span class="info-value mono break">{{ ticket.user_id }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">来源</span>
                <span class="info-value">{{ SOURCE_LABEL[ticket.source] || ticket.source || '—' }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">分类</span>
                <span class="info-value"><CategoryTag :category="ticket.category" /></span>
              </div>
              <div class="info-row">
                <span class="info-label">状态</span>
                <span class="info-value"><StatusPill :status="ticket.status" /></span>
              </div>
              <div class="info-row">
                <span class="info-label">归档</span>
                <span class="info-value" :class="ticket.archived ? '' : 'muted'">{{ ticket.archived ? '已归档' : '未归档' }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">创建时间</span>
                <span class="info-value"><TimeText :value="ticket.created_at" /></span>
              </div>
              <div v-if="ticket.answered_at" class="info-row">
                <span class="info-label">答复时间</span>
                <span class="info-value"><TimeText :value="ticket.answered_at" /></span>
              </div>
              <div v-if="ticket.archived_at" class="info-row">
                <span class="info-label">归档时间</span>
                <span class="info-value"><TimeText :value="ticket.archived_at" /></span>
              </div>
            </div>
          </section>

          <!-- 快捷操作 -->
          <section class="card block anim-stagger" style="--i: 4">
            <header class="block-head">
              <span class="block-icon"><ListChecks :size="16" :stroke-width="2" /></span>
              <h3>操作</h3>
            </header>
            <div class="block-body action-list">
              <UiButton v-if="ticket.status !== 'closed'" variant="ghost" block :loading="busy" @click="onClose">
                <template #icon><Lock :size="15" :stroke-width="2" /></template>
                关闭工单
              </UiButton>
              <UiButton v-else variant="ghost" block :loading="busy" @click="onReopen">
                <template #icon><RotateCcw :size="15" :stroke-width="2" /></template>
                重开工单
              </UiButton>
              <UiButton v-if="!ticket.archived" variant="ghost" block :loading="busy" @click="onArchive">
                <template #icon><Archive :size="15" :stroke-width="2" /></template>
                归档工单
              </UiButton>
              <UiButton v-else variant="ghost" block :loading="busy" @click="onUnarchive">
                <template #icon><ArchiveRestore :size="15" :stroke-width="2" /></template>
                取消归档
              </UiButton>
              <UiButton variant="danger-soft" block :loading="busy" @click="onDelete">
                <template #icon><Trash2 :size="15" :stroke-width="2" /></template>
                删除工单
              </UiButton>
            </div>
          </section>

          <!-- 知识库检索 -->
          <section v-if="ticket.rag_result" class="card block anim-stagger" style="--i: 5">
            <header class="block-head clickable" @click="ragOpen = !ragOpen">
              <span class="block-icon"><BookOpen :size="16" :stroke-width="2" /></span>
              <h3>知识库检索</h3>
              <ChevronDown :size="15" :stroke-width="2" class="rag-chev" :class="{ open: ragOpen }" />
            </header>
            <div class="rag-wrap" :class="{ open: ragOpen }">
              <div class="rag-inner">
                <pre class="rag-content mono">{{ ticket.rag_result }}</pre>
              </div>
            </div>
          </section>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowLeft, Copy, Check, Clock, MessageSquare, MessageSquareCheck, MessageSquareText,
  CircleQuestionMark, Archive, ArchiveRestore, RotateCcw, Trash2, Pencil, Plus, Save,
  Lock, Info, BookOpen, ChevronDown, Ellipsis, User, ListChecks,
} from '@lucide/vue'
import {
  getTicket, answerTicket, closeTicket, reopenTicket, archiveTicket, deleteTicket,
} from '../api/ticket'
import UiButton from '../ui/UiButton.vue'
import UiInput from '../ui/UiInput.vue'
import UiDropdown from '../ui/UiDropdown.vue'
import UiSkeleton from '../ui/UiSkeleton.vue'
import UiEmpty from '../ui/UiEmpty.vue'
import StatusPill from '../ui/StatusPill.vue'
import CategoryTag from '../ui/CategoryTag.vue'
import TimeText from '../ui/TimeText.vue'
import { toast } from '../ui/toast'
import { confirmDialog, promptDialog } from '../ui/confirm'
import { useIsMobile } from '../composables/useMediaQuery'

const route = useRoute()
const router = useRouter()
const isMobile = useIsMobile()

const loading = ref(true)
const ticket = ref(null)
const busy = ref(false)
const editingAnswer = ref(false)
const answerDraft = ref('')
const answerSubmitting = ref(false)
const answerInputRef = ref(null)
const ragOpen = ref(true)
const copied = ref(false)

const SOURCE_LABEL = {
  hiagent_chat: 'HiAgent 对话',
  wechat_service: '微信服务号',
  wechat_subscribe: '微信订阅号',
  feishu: '飞书',
  yiban: '易班',
}

const moreActions = computed(() => {
  const items = []
  if (ticket.value?.status === 'closed') {
    items.push({ label: '重开工单', value: 'reopen', icon: RotateCcw })
  } else {
    items.push({ label: '关闭工单', value: 'close', icon: Lock })
  }
  items.push({ label: '删除工单', value: 'delete', icon: Trash2, danger: true, divided: true })
  return items
})

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

const copyId = async () => {
  try {
    await navigator.clipboard.writeText(ticket.value.ticket_id)
    copied.value = true
    toast.success('工单号已复制')
    setTimeout(() => (copied.value = false), 1500)
  } catch {
    toast.warning('复制失败，请手动选择复制')
  }
}

// ─── 答复 ───
const enterEditAnswer = async () => {
  if (ticket.value?.status === 'closed') {
    toast.warning('工单已关闭，无法答复')
    return
  }
  answerDraft.value = ticket.value?.answer || ''
  editingAnswer.value = true
  await nextTick()
  answerInputRef.value?.focus()
}

const cancelEditAnswer = () => {
  editingAnswer.value = false
  answerDraft.value = ''
}

const saveAnswer = async () => {
  const text = answerDraft.value.trim()
  if (!text) { toast.warning('请填写答复内容'); return }
  answerSubmitting.value = true
  try {
    await answerTicket(route.params.id, { answer: text })
    toast.success('答复已保存')
    editingAnswer.value = false
    await fetch()
  } catch (e) {} finally {
    answerSubmitting.value = false
  }
}

// ─── 状态操作 ───
const onClose = async () => {
  const reason = await promptDialog({
    title: '关闭工单',
    message: '可填写关闭原因（选填），如：重复工单、无效问题等',
    placeholder: '关闭原因…',
    confirmText: '确认关闭',
  })
  if (reason === null) return
  busy.value = true
  try {
    await closeTicket(route.params.id, reason || '')
    toast.success('工单已关闭')
    await fetch()
  } catch (e) {} finally { busy.value = false }
}

const onReopen = async () => {
  busy.value = true
  try {
    await reopenTicket(route.params.id)
    toast.success('工单已重开')
    await fetch()
  } catch (e) {} finally { busy.value = false }
}

const onArchive = async () => {
  busy.value = true
  try {
    await archiveTicket(route.params.id, true)
    toast.success('工单已归档')
    await fetch()
  } catch (e) {} finally { busy.value = false }
}

const onUnarchive = async () => {
  busy.value = true
  try {
    await archiveTicket(route.params.id, false)
    toast.success('已取消归档')
    await fetch()
  } catch (e) {} finally { busy.value = false }
}

const onDelete = async () => {
  const ok = await confirmDialog({
    title: '删除工单',
    message: '确认删除该工单？删除后不可在列表中查看。',
    confirmText: '确认删除',
    danger: true,
  })
  if (!ok) return
  busy.value = true
  try {
    await deleteTicket(route.params.id, false)
    toast.success('工单已删除')
    router.push('/tickets')
  } catch (e) {} finally { busy.value = false }
}

const onMoreAction = (cmd) => {
  if (cmd === 'close') onClose()
  else if (cmd === 'reopen') onReopen()
  else if (cmd === 'delete') onDelete()
}

onMounted(fetch)
</script>

<style scoped>
.detail { display: flex; flex-direction: column; gap: 16px; }

.detail-skeleton { display: flex; flex-direction: column; gap: 16px; }
.sk-grid {
  display: grid;
  grid-template-columns: 6fr 4fr;
  gap: 16px;
  padding: 8px;
}
@media (max-width: 1023px) { .sk-grid { grid-template-columns: 1fr; } }

/* ─── 返回 ─── */
.back-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  align-self: flex-start;
  border: none;
  background: transparent;
  color: var(--text-3);
  font-size: 13px;
  cursor: pointer;
  padding: 4px 8px 4px 4px;
  border-radius: var(--r-sm);
  transition: all var(--d-fast) var(--ease-out);
}
.back-link:hover { color: var(--primary); background: var(--primary-soft); }

/* ─── 头部卡 ─── */
.head-card {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  padding: 18px 22px;
  flex-wrap: wrap;
}
.head-main { min-width: 0; flex: 1; }
.head-id-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
.tid-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  border: 1px solid var(--border);
  background: var(--bg-sink);
  color: var(--text-1);
  font-size: 13px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: var(--r-sm);
  cursor: pointer;
  transition: all var(--d-fast) var(--ease-out);
  letter-spacing: -0.2px;
}
.tid-chip:hover { border-color: var(--primary); color: var(--primary); }
.copy-icon { color: var(--text-3); transition: color var(--d-fast); }
.copy-icon.ok { color: var(--success); }

.archived-pill {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 9px;
  border-radius: var(--r-full);
  font-size: 11.5px;
  font-weight: 500;
  background: var(--info-soft);
  color: var(--info);
}

.head-meta {
  display: flex;
  gap: 18px;
  margin-top: 11px;
  flex-wrap: wrap;
}
.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12.5px;
  color: var(--text-3);
}

.head-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
  flex-wrap: wrap;
}

/* ─── 双列 ─── */
.detail-grid {
  display: grid;
  grid-template-columns: minmax(0, 6fr) minmax(0, 4fr);
  gap: 16px;
  align-items: start;
}
@media (max-width: 1023px) { .detail-grid { grid-template-columns: 1fr; } }
.col-main, .col-side {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

/* ─── 通用块 ─── */
.block { overflow: hidden; }
.block-head {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 14px 20px;
  border-bottom: 1px solid var(--border-soft);
}
.block-head.clickable { cursor: pointer; user-select: none; }
.block-head h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-1);
  letter-spacing: -0.2px;
}
.block-icon {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  background: var(--primary-soft);
  color: var(--primary);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.block-icon.warn { background: var(--warning-soft); color: var(--warning); }
.block-body { padding: 18px 20px; }

.from-chip {
  margin-left: auto;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 11.5px;
  color: var(--text-3);
  background: var(--bg-sink);
  border: 1px solid var(--border-soft);
  padding: 3px 9px;
  border-radius: var(--r-full);
  max-width: 46%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* ─── 问题 ─── */
.question-quote {
  margin: 0;
  padding: 16px 18px;
  border-left: 3px solid var(--primary);
  border-radius: 0 var(--r-md) var(--r-md) 0;
  background: linear-gradient(90deg, var(--primary-soft), transparent 75%);
  color: var(--text-1);
  font-size: 15px;
  line-height: 1.75;
  white-space: pre-wrap;
  word-break: break-word;
}
@media (max-width: 767px) { .question-quote { padding: 13px 14px; font-size: 14px; } }

/* ─── 答复 ─── */
.answer-state {
  margin-left: auto;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 10px;
  border-radius: var(--r-full);
  font-size: 11.5px;
  font-weight: 600;
}
.answer-state.done { background: var(--success-soft); color: var(--success); }
.answer-state.todo { background: var(--warning-soft); color: var(--warning); }

.answer-view { display: flex; flex-direction: column; gap: 14px; }
.answer-text {
  padding: 16px 18px;
  border-radius: var(--r-md);
  background: var(--bg-sink);
  border: 1px solid var(--border-soft);
  color: var(--text-1);
  font-size: 14.5px;
  line-height: 1.8;
  white-space: pre-wrap;
  word-break: break-word;
}
.answer-actions { display: flex; justify-content: flex-end; align-items: center; }
.closed-tip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12.5px;
  color: var(--text-3);
}
.answer-edit { display: flex; flex-direction: column; gap: 12px; }
.edit-actions { display: flex; justify-content: flex-end; gap: 8px; }

.close-reason {
  margin: 0;
  padding: 13px 16px;
  border-radius: var(--r-md);
  background: var(--warning-soft);
  color: var(--text-2);
  font-size: 13.5px;
  line-height: 1.7;
  white-space: pre-wrap;
}

/* ─── 信息列表 ─── */
.info-list { display: flex; flex-direction: column; }
.info-row {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 9px 0;
  border-bottom: 1px solid var(--border-soft);
}
.info-row:last-child { border-bottom: none; }
.info-label {
  width: 64px;
  flex-shrink: 0;
  font-size: 12.5px;
  color: var(--text-3);
}
.info-value { font-size: 13px; color: var(--text-1); min-width: 0; }
.info-value.break { word-break: break-all; font-size: 12.5px; }

/* ─── 操作 ─── */
.action-list { display: flex; flex-direction: column; gap: 9px; }

/* ─── RAG ─── */
.rag-chev {
  margin-left: auto;
  color: var(--text-3);
  transition: transform var(--d-base) var(--ease-out);
}
.rag-chev.open { transform: rotate(180deg); }
.rag-wrap {
  display: grid;
  grid-template-rows: 0fr;
  transition: grid-template-rows var(--d-base) var(--ease-out);
}
.rag-wrap.open { grid-template-rows: 1fr; }
.rag-inner { overflow: hidden; }
.rag-content {
  margin: 16px 20px 18px;
  padding: 13px 15px;
  border-radius: var(--r-md);
  background: var(--bg-sink);
  border: 1px solid var(--border-soft);
  font-size: 12px;
  color: var(--text-2);
  line-height: 1.65;
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 300px;
  overflow: auto;
}
</style>
