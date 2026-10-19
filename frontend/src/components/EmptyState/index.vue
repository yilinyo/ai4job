<template>
  <div class="empty-state">
    <i class="el-icon-document" />
    <p>{{ description }}</p>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'

@Component({
  name: 'EmptyState'
})
export default class extends Vue {
  @Prop({ default: '暂无数据' }) private description!: string
}
</script>

<style lang="scss" scoped>
.empty-state {
  min-height: 180px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #909399;
  text-align: center;

  i {
    margin-bottom: 12px;
    font-size: 48px;
    color: #c0c4cc;
  }

  p {
    margin: 0;
    font-size: 14px;
    line-height: 22px;
  }
}
</style>
