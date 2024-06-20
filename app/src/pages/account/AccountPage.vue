<script setup lang="ts">

import { ref } from 'vue'

const tab = ref('mails')

</script>

<template>
  <div class="flex column">
    <q-tabs
      v-model="tab"
      dense
      align="left"
      class="text-grey"
      active-color="primary"
      indicator-color="primary"
      :breakpoint="0"
    >
      <q-route-tab name="mails" to="/account" exact>
        <div class="row q-gutter-x-sm">
          <q-img alt="douyin" src="~assets/logo/douyin.svg" width="22px" height="22px"/>
          <div>抖音</div>
        </div>
      </q-route-tab>
    </q-tabs>
    <q-separator/>
    <div>
      <router-view/>
    </div>
  </div>
</template>

<style scoped>

</style>
