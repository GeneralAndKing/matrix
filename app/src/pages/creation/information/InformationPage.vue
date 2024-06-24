<script setup lang="ts">
import { useRoute } from 'vue-router'
import { onMounted } from 'vue'
import { useInformationStore } from 'stores/creation-store'
import DouYinPublish from 'pages/creation/information/DouYinPublish.vue'

const route = useRoute()
const store = useInformationStore()

onMounted(() => {
  const id = parseInt(route.params.id as string)
  store.handleInformation(id)
})

</script>

<template>
  <q-card flat bordered style="height: calc(100vh - 68px)" class="q-px-md q-pr-sm flex column">
    <q-tabs
      v-model="store.styleData.tab"
      dense
      class="text-grey q-mt-sm"
      active-color="primary"
      indicator-color="primary"
      align="left"
    >
      <q-tab name="mails">
        <div class="row q-gutter-x-sm">
          <q-img alt="douyin" src="~assets/logo/douyin.svg" width="22px" height="22px"/>
          <div>抖音</div>
        </div>
      </q-tab>
    </q-tabs>
    <q-separator/>
    <q-tab-panels v-model="store.styleData.tab" animated style="flex: 1">
      <q-tab-panel name="mails" class="q-pa-none q-pt-sm">
        <DouYinPublish />
      </q-tab-panel>
    </q-tab-panels>
  </q-card>
</template>

<style scoped lang="scss">
:deep(.account-list-tab) {
  flex-direction: row !important;
  flex-wrap: nowrap;
}
</style>
