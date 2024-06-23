<template>
  <q-layout view="hHh Lpr lff">
    <q-header elevated>
      <q-bar class="q-electron-drag bg-primary">
        <q-icon name="laptop_chromebook"/>
        <div>Matrix</div>
        <q-badge rounded :color="healthStore.isOk ? 'green' : 'red'" style="transition: all 1s;"/>
        <q-space/>

        <q-btn dense flat icon="minimize" @click="handleMinimize"/>
        <q-btn dense flat icon="crop_square" @click="handleToggleMaximize"/>
        <q-btn dense flat icon="close" @click="handleCloseApp"/>
      </q-bar>
    </q-header>

    <q-drawer
      v-model="drawer"
      show-if-above
      :mini="miniState"
      @mouseover="miniState = false"
      @mouseout="miniState = true"
      mini-to-overlay
      :width="200"
      :breakpoint="500"
      bordered
      :class="`overflow-hidden column ${$q.dark.isActive ? 'bg-grey-9' : 'bg-grey-3'}`"
    >
      <q-scroll-area class="fit" :horizontal-thumb-style="{ opacity: '0' }" style="flex: 1">
        <q-list padding>
          <q-item v-for="item in menus" :key="item.name" :exact="item.exact" :to="item.path"  clickable v-ripple>
            <q-item-section avatar>
              <q-icon :name="item.icon"/>
            </q-item-section>
            <q-item-section> {{ item.name }}</q-item-section>
          </q-item>
        </q-list>
      </q-scroll-area>
      <q-list>
        <q-item clickable v-ripple>
          <q-item-section avatar>
            <q-icon name="settings"/>
          </q-item-section>
          <q-item-section>设置</q-item-section>
        </q-item>
        <q-item clickable v-ripple to="/account">
          <q-item-section avatar>
            <q-icon name="account_circle"/>
          </q-item-section>
          <q-item-section>
            <q-item-label lines="1">账号管理</q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <q-scroll-area style="height: 100vh; width: calc(100% - 5px);">
        <q-page class="q-pa-md">
          <router-view/>
        </q-page>
      </q-scroll-area>
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { isElectron } from 'src/utils/action'
import { onMounted, ref } from 'vue'
import { useHealthStore } from 'stores/application-store'
import { useQuasar } from 'quasar'

const drawer = ref(false)
const miniState = ref(true)
const healthStore = useHealthStore()
const $q = useQuasar()

const menus = [
  { name: '创作集合', icon: 'video_library', path: '/creation', exact: false }
]

const handleMinimize = () => {
  if (isElectron()) {
    window.WindowsApi.minimize()
  }
}

const handleToggleMaximize = () => {
  if (isElectron()) {
    window.WindowsApi.toggleMaximize()
  }
}

const handleCloseApp = () => {
  if (isElectron()) {
    window.WindowsApi.close()
  }
}

onMounted(() => {
  // $q.localStorage.set('test', 'test')
  console.log($q.localStorage.getItem('test'))
})
</script>
