<template>
  <el-card shadow="never" class="todo-panel">
    <div slot="header" class="todo-panel__header">
      <span>待办事项</span>
      <el-button type="primary" size="mini" icon="el-icon-plus" @click="handleCreate">
        新建待办
      </el-button>
    </div>

    <el-empty
      v-if="!todos.length"
      description="暂无待办事项"
      :image-size="100"
    />
    <div v-else class="todo-list">
      <div
        v-for="todo in todos"
        :key="todo.id"
        :class="['todo-item', {'is-completed': isCompleted(todo), 'is-clickable': todo.route}]"
        @click="handleOpen(todo)"
      >
        <div class="todo-item__main">
          <div class="todo-item__title-row">
            <span class="todo-item__title">{{ todo.title }}</span>
            <el-tag :type="todo.source === 'system' ? 'warning' : 'success'" size="mini">
              {{ todo.source === 'system' ? '系统' : '手动' }}
            </el-tag>
          </div>
          <div v-if="todo.description" class="todo-item__description">
            {{ todo.description }}
          </div>
          <div class="todo-item__meta">
            <span v-if="todo.dueAt">
              <i class="el-icon-time" /> {{ formatDateTime(todo.dueAt) }}
            </span>
            <span v-if="todo.targetType && todo.targetId">
              <i class="el-icon-connection" /> {{ targetTypeLabel(todo.targetType) }} #{{ todo.targetId }}
            </span>
            <span v-if="todo.completedAt">
              <i class="el-icon-circle-check" /> 已完成 {{ formatDateTime(todo.completedAt) }}
            </span>
          </div>
        </div>
        <div class="todo-item__actions" @click.stop>
          <el-button
            size="mini"
            type="success"
            :disabled="isCompleted(todo) || todo.source === 'system'"
            @click="handleComplete(todo)"
          >
            完成
          </el-button>
          <el-button
            size="mini"
            :disabled="isCompleted(todo) || todo.source === 'system'"
            @click="handleEdit(todo)"
          >
            编辑
          </el-button>
          <el-button
            size="mini"
            type="danger"
            :disabled="isCompleted(todo) || todo.source === 'system'"
            @click="handleDelete(todo)"
          >
            删除
          </el-button>
        </div>
      </div>
    </div>

    <el-dialog
      :title="dialogTitle"
      :visible.sync="dialogVisible"
      width="520px"
      @close="resetForm"
    >
      <el-form
        ref="todoForm"
        :model="form"
        :rules="rules"
        label-width="90px"
      >
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" maxlength="200" show-word-limit />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            maxlength="1000"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="截止时间">
          <el-date-picker
            v-model="form.dueAt"
            type="datetime"
            value-format="yyyy-MM-ddTHH:mm:ss.SSSZ"
            placeholder="选择截止时间"
            style="width: 100%;"
          />
        </el-form-item>
        <el-form-item label="目标类型">
          <el-select v-model="form.targetType" clearable placeholder="选择目标类型" style="width: 100%;">
            <el-option label="投递" value="application" />
            <el-option label="面试" value="interview" />
            <el-option label="题目" value="question" />
            <el-option label="简历" value="resume" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标 ID">
          <el-input v-model="form.targetId" placeholder="请输入目标 ID" clearable />
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button :disabled="submitting" @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">保存</el-button>
      </div>
    </el-dialog>
  </el-card>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import { ElForm } from 'element-ui/types/form'
import {
  CreateDashboardTodoRequest,
  DashboardTargetType,
  DashboardTodo,
  UpdateDashboardTodoRequest
} from '@/models'
import { parseTime } from '@/utils'

const formatDashboardDateTime = (value?: string) => {
  if (!value) return '-'
  const date = new Date(value)
  if (!Number.isNaN(date.getTime())) {
    const pad = (num: number) => String(num).padStart(2, '0')
    return `${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`
  }
  return parseTime(value, '{m}-{d} {h}:{i}') || '-'
}

interface TodoForm {
  title: string
  description: string
  dueAt: string
  targetType: DashboardTargetType | ''
  targetId: string
}

@Component({
  name: 'TodoPanel'
})
export default class extends Vue {
  @Prop({ required: true, default: () => [] }) private todos!: DashboardTodo[]
  @Prop({ default: false }) private submitting!: boolean
  @Prop({ default: 0 }) private submitSuccessKey!: number

  private dialogVisible = false
  private editingTodoId = ''
  private form: TodoForm = this.getEmptyForm()
  private rules = {
    title: [
      { required: true, message: '请输入待办标题', trigger: 'blur' },
      { max: 200, message: '标题不能超过 200 个字符', trigger: 'blur' }
    ]
  }

  get dialogTitle() {
    return this.editingTodoId ? '编辑待办' : '新建待办'
  }

  private getEmptyForm(): TodoForm {
    return {
      title: '',
      description: '',
      dueAt: '',
      targetType: '',
      targetId: ''
    }
  }

  private isCompleted(todo: DashboardTodo) {
    return Boolean(todo.completedAt)
  }

  private formatDateTime(value?: string) {
    return formatDashboardDateTime(value)
  }

  @Watch('submitSuccessKey')
  private handleSubmitSuccess() {
    if (this.dialogVisible) {
      this.dialogVisible = false
    }
  }

  private targetTypeLabel(type: DashboardTargetType) {
    const labels: Record<DashboardTargetType, string> = {
      application: '投递',
      interview: '面试',
      question: '题目',
      resume: '简历'
    }
    return labels[type]
  }

  private handleCreate() {
    this.editingTodoId = ''
    this.form = this.getEmptyForm()
    this.dialogVisible = true
  }

  private handleEdit(todo: DashboardTodo) {
    if (todo.source === 'system' || this.isCompleted(todo)) return
    this.editingTodoId = todo.id
    this.form = {
      title: todo.title,
      description: todo.description || '',
      dueAt: todo.dueAt || '',
      targetType: todo.targetType || '',
      targetId: todo.targetId || ''
    }
    this.dialogVisible = true
  }

  private handleComplete(todo: DashboardTodo) {
    if (todo.source === 'system' || this.isCompleted(todo)) return
    this.$emit('complete', todo.id)
  }

  private async handleDelete(todo: DashboardTodo) {
    if (todo.source === 'system' || this.isCompleted(todo)) return
    try {
      await this.$confirm('确定要删除这个待办事项吗？删除后无法恢复！', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      this.$emit('delete', todo.id)
    } catch (error) {
      if (error !== 'cancel') {
        this.$message.error('删除操作被中断')
      }
    }
  }

  private handleOpen(todo: DashboardTodo) {
    if (!todo.route) return
    this.$emit('open', todo.route)
  }

  private handleSubmit() {
    if (this.submitting) return
    (this.$refs.todoForm as ElForm).validate((valid: boolean) => {
      if (!valid) return
      const payload = this.buildPayload()
      if (this.editingTodoId) {
        this.$emit('update', {
          id: this.editingTodoId,
          data: payload
        })
      } else {
        this.$emit('create', payload)
      }
    })
  }

  private buildPayload(): CreateDashboardTodoRequest | UpdateDashboardTodoRequest {
    const targetType = this.form.targetType || undefined
    const targetId = this.form.targetId.trim() || undefined
    const payload: CreateDashboardTodoRequest | UpdateDashboardTodoRequest = {
      title: this.form.title.trim(),
      description: this.form.description.trim()
    }

    if (this.form.dueAt) {
      payload.dueAt = this.form.dueAt
    } else if (this.editingTodoId) {
      (payload as UpdateDashboardTodoRequest).clearDueAt = true
    }

    if (targetType && targetId) {
      payload.targetType = targetType
      payload.targetId = targetId
    } else if (this.editingTodoId) {
      (payload as UpdateDashboardTodoRequest).clearTarget = true
    }

    return payload
  }

  private resetForm() {
    this.editingTodoId = ''
    this.form = this.getEmptyForm()
    const formRef = this.$refs.todoForm as ElForm | undefined
    if (formRef) {
      formRef.clearValidate()
    }
  }
}
</script>

<style lang="scss" scoped>
.todo-panel {
  margin-bottom: 20px;
}

.todo-panel__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
  color: #303133;
}

.todo-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.todo-item {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 14px;
  border: 1px solid #ebeef5;
  border-radius: 10px;
  background: #fff;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;

  &.is-clickable {
    cursor: pointer;

    &:hover {
      border-color: #409eff;
      box-shadow: 0 4px 12px rgba(64, 158, 255, 0.12);
    }
  }

  &.is-completed {
    opacity: 0.68;
  }
}

.todo-item__main {
  min-width: 0;
  flex: 1;
}

.todo-item__title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.todo-item__title {
  color: #303133;
  font-weight: 600;
  line-height: 1.4;
}

.todo-item__description {
  color: #606266;
  line-height: 1.5;
  margin-bottom: 8px;
}

.todo-item__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  color: #909399;
  font-size: 12px;
}

.todo-item__actions {
  display: flex;
  align-items: flex-start;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 8px;

  .el-button + .el-button {
    margin-left: 0;
  }
}

@media (max-width: 768px) {
  .todo-item {
    flex-direction: column;
  }

  .todo-item__actions {
    justify-content: flex-start;
  }
}
</style>
