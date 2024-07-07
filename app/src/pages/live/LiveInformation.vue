<script setup lang="ts">

import { useRoute } from 'vue-router'

const route = useRoute()

</script>

<template>
  <q-tabs
    dense
    class="text-grey"
    active-color="primary"
    indicator-color="primary"
    align="left">
    <q-route-tab :to="{name: 'LiveDashboard', param: {id: route.params.id}}" exact label="仪表盘"/>
    <q-route-tab :to="{name: 'LiveContent', param: {id: route.params.id}}" exact label="直播现况"/>
  </q-tabs>
  <q-separator/>
  <div class="q-mt-sm full-height" style="flex: 1">
    <router-view/>
  </div>
</template>

<style scoped lang="scss">

</style>
