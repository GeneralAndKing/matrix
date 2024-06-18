<template>
  <q-layout view="hHh Lpr lff">
    <q-header elevated>
      <q-bar class="q-electron-drag bg-primary">
        <q-icon name="laptop_chromebook"/>
        <div>Matrix</div>

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
          <q-item v-for="item in menus" :key="item.name" :to="item.path" exact clickable v-ripple>
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
        <q-item clickable v-ripple>
          <q-item-section avatar>
            <q-avatar size="sm">
              <img alt="avatar" src="https://cdn.quasar.dev/img/boy-avatar.png">
            </q-avatar>
          </q-item-section>
          <q-item-section>
            <q-item-label lines="1">zhongyue</q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <q-scroll-area :horizontal-thumb-style="{ opacity: '1' }"
                     style="height: calc(100vh - 36px); width: 100%;">
        <q-page class="row items-center justify-evenly">
          <router-view/>
        </q-page>
      </q-scroll-area>
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { isElectron } from 'src/utils/action'
import { ref } from 'vue'

const drawer = ref(true)
const miniState = ref(false)

const menus = [
  { name: '创作集合', icon: 'video_library', path: '/' }
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
</script>
