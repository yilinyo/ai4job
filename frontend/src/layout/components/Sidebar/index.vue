<template>
  <div :class="{'has-logo': showLogo}">
    <sidebar-logo
      v-if="showLogo"
      :collapse="isCollapse"
    />
    <el-scrollbar wrap-class="scrollbar-wrapper">
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :background-color="variables.menuBg"
        :text-color="variables.menuText"
        :active-text-color="menuActiveTextColor"
        :unique-opened="false"
        :collapse-transition="false"
        mode="vertical"
      >
        <sidebar-item
          v-for="route in routes"
          :key="route.path"
          :item="route"
          :base-path="route.path"
          :is-collapse="isCollapse"
        />
      </el-menu>
    </el-scrollbar>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { AppModule } from '@/store/modules/app'
import { PermissionModule } from '@/store/modules/permission'
import { SettingsModule } from '@/store/modules/settings'
import SidebarItem from './SidebarItem.vue'
import SidebarLogo from './SidebarLogo.vue'
import variables from '@/styles/_variables.scss'

@Component({
  name: 'SideBar',
  components: {
    SidebarItem,
    SidebarLogo
  }
})
export default class extends Vue {
  get sidebar() {
    return AppModule.sidebar
  }

  get routes() {
    return PermissionModule.routes
  }

  get showLogo() {
    return SettingsModule.showSidebarLogo
  }

  get menuActiveTextColor() {
    if (SettingsModule.sidebarTextTheme) {
      return SettingsModule.theme
    } else {
      return variables.menuActiveText
    }
  }

  get variables() {
    return variables
  }

  get activeMenu() {
    const route = this.$route
    const { meta, path } = route
    // if set path, the sidebar will highlight the path you set
    if (meta.activeMenu) {
      return meta.activeMenu
    }
    return path
  }

  get isCollapse() {
    return !this.sidebar.opened
  }
}
</script>

<style lang="scss">
.sidebar-container {
  background: linear-gradient(180deg, #f7fbff 0%, #eef5ff 100%);
  border-right: 1px solid rgba(211, 223, 243, 0.95);
  box-shadow: 0 10px 30px rgba(31, 45, 61, 0.08);

  // reset element-ui css
  .horizontal-collapse-transition {
    transition: 0s width ease-in-out, 0s padding-left ease-in-out, 0s padding-right ease-in-out;
  }

  .scrollbar-wrapper {
    overflow-x: hidden !important;
  }

  .el-scrollbar__view {
    height: 100%
  }

  .el-menu {
    padding: 12px 10px 16px;
    background: transparent !important;
  }

  .el-menu-item,
  .el-submenu__title {
    height: 46px;
    line-height: 46px;
    margin-bottom: 6px;
    color: #42526e !important;
    border-radius: 14px;
    transition: all 0.2s ease;
  }

  .el-menu-item:hover,
  .el-submenu__title:hover {
    background: #e8f0ff !important;
    color: #1d63ed !important;
  }

  .el-menu-item.is-active {
    background: rgba(255, 255, 255, 0.96) !important;
    color: #1d63ed !important;
    box-shadow: 0 10px 24px rgba(67, 114, 255, 0.14);
  }

  .el-menu-item.is-active .svg-icon,
  .el-submenu.is-active > .el-submenu__title .svg-icon,
  .el-menu-item:hover .svg-icon,
  .el-submenu__title:hover .svg-icon {
    color: #1d63ed;
  }

  .svg-icon {
    color: #7f8da3;
    transition: color 0.2s ease;
  }

  .el-scrollbar__bar {
    &.is-vertical {
      right: 0px;
    }

    &.is-horizontal {
      display: none;
    }
  }
}
</style>

<style lang="scss" scoped>
.el-scrollbar {
  height: 100%
}

.has-logo {
  .el-scrollbar {
    height: calc(100% - 50px);
  }
}

.el-menu {
  border: none;
  height: 100%;
  background: transparent;
  width: 100% !important;
}
</style>
