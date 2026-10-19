<template>
  <div class="api-key-management">
    <div class="toolbar"><span>API Key</span><el-button type="primary" size="small" icon="el-icon-plus" @click="open()">创建 Key</el-button></div>
    <el-table v-loading="loading" :data="keys" size="small">
      <el-table-column prop="name" label="名称" min-width="120" />
      <el-table-column prop="displayPrefix" label="Key" min-width="145" />
      <el-table-column label="状态" width="82"><template slot-scope="{row}"><el-tag :type="status(row) === '有效' ? 'success' : 'info'" size="mini">{{ status(row) }}</el-tag></template></el-table-column>
      <el-table-column label="创建时间" min-width="155"><template slot-scope="{row}">{{ time(row.createdAt) }}</template></el-table-column>
      <el-table-column label="过期时间" min-width="155"><template slot-scope="{row}">{{ row.expiresAt ? time(row.expiresAt) : '永不过期' }}</template></el-table-column>
      <el-table-column label="操作" width="150"><template slot-scope="{row}"><el-button type="text" size="small" @click="open(row)">轮换</el-button><el-button type="text" class="danger" size="small" @click="revoke(row)">撤销</el-button></template></el-table-column>
    </el-table>
    <el-dialog :title="selected ? '轮换 API Key' : '创建 API Key'" :visible.sync="dialog" width="420px" append-to-body><el-form label-width="86px"><el-form-item label="名称"><el-input v-model="form.name" maxlength="100" /></el-form-item><el-form-item label="有效期"><el-radio-group v-model="form.expiresInDays"><el-radio :label="30">30 天</el-radio><el-radio :label="60">60 天</el-radio><el-radio :label="90">90 天</el-radio><el-radio :label="null">永不过期</el-radio></el-radio-group></el-form-item></el-form><div slot="footer"><el-button @click="dialog=false">取消</el-button><el-button type="primary" :loading="saving" @click="save">确定</el-button></div></el-dialog>
    <el-dialog title="请立即保存 API Key" :visible.sync="plainVisible" width="520px" :close-on-click-modal="false" append-to-body><el-alert title="此完整 Key 只显示一次，关闭后无法再次查看。" type="warning" :closable="false" /><el-input class="plain-key" :value="plain" readonly><el-button slot="append" icon="el-icon-document-copy" @click="copy" /></el-input><div slot="footer"><el-button type="primary" @click="plainVisible=false">我已保存</el-button></div></el-dialog>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { APIKeyRecord, CreateAPIKeyRequest } from '@/models'
import { APIKeyModule } from '@/store/modules/api-key'
@Component({ name: 'ApiKeyManagement' })
export default class extends Vue {
  private keys: APIKeyRecord[] = []; private loading = false; private saving = false; private dialog = false; private plainVisible = false; private plain = ''; private selected: APIKeyRecord | null = null; private form: CreateAPIKeyRequest = { name: '', expiresInDays: 30 }
  created() { this.load() }
  private async load() { this.loading = true; try { await APIKeyModule.Load(); this.keys = APIKeyModule.keys } finally { this.loading = false } }
  private open(key?: APIKeyRecord) { this.selected = key || null; this.form = { name: key ? key.name : '', expiresInDays: 30 }; this.dialog = true }
  private async save() { if (!this.form.name.trim()) { this.$message.error('请输入 Key 名称'); return } this.saving = true; try { const data = this.selected ? await APIKeyModule.Rotate({ id: this.selected.id, data: this.form }) : await APIKeyModule.Create(this.form); this.dialog = false; this.plain = data.key; this.plainVisible = true; this.keys = APIKeyModule.keys } finally { this.saving = false } }
  private revoke(key: APIKeyRecord) { this.$confirm(`撤销“${key.name}”后将立即无法用于 MCP。`, '确认撤销', { type: 'warning' }).then(async() => { await APIKeyModule.Revoke(key.id); this.keys = APIKeyModule.keys }).catch(() => undefined) }
  private copy() { navigator.clipboard.writeText(this.plain); this.$message.success('已复制') }
  private status(key: APIKeyRecord) { return key.status === 'revoked' ? '已撤销' : key.expiresAt && new Date(key.expiresAt) <= new Date() ? '已过期' : '有效' }
  private time(value: string) { return value.replace('T', ' ').slice(0, 16) }
}
</script>
<style lang="scss" scoped>
.toolbar { display:flex; justify-content:space-between; align-items:center; margin-bottom:16px; font-weight:600; }.danger { color:#f56c6c; }.plain-key { margin-top:16px; }
</style>
