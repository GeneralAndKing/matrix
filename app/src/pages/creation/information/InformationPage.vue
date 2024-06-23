<script setup lang="ts">
import { useRoute } from 'vue-router'
import { onMounted, reactive } from 'vue'
import InformationAccount from 'pages/creation/information/InformationAccount.vue'

interface Data {
  id: number,
  informationSplitterModel: number,
  tab: 'DouYin'
}

const route = useRoute()
const data = reactive<Data>({
  id: 0,
  informationSplitterModel: 30,
  tab: 'DouYin'
})

onMounted(() => {
  data.id = parseInt(route.params.id as string)
  // document.body.style.zoom = '1'
})

</script>

<template>
  <q-card flat bordered style="height: calc(100vh - 68px)" class="q-pa-md q-pr-sm flex column">
    <q-splitter
      class="full-width flex row no-wrap "
      v-model="data.informationSplitterModel"
      style="max-height: 400px"
      :limits="[30, 40]"
    >
      <template v-slot:before>
        <div class="q-mr-md">
          <div class="text-h6">Our Changing Planet
            <q-badge label="视频" color="teal-6"/>
          </div>
          <div>
            Lorem ipsum dolor sit, amet consectetur adipisicing
            elit. Quis praesentium cumque magnam odio iure quidem, quod illum numquam possimus obcaecati commodi
            minima assumenda consectetur culpa fuga nulla ullam. In, libero.
          </div>
          <div class="q-mt-md">
            <q-chip size="12px" icon="event">
              创建时间 2021-02-12 11:21:21
            </q-chip>
          </div>
        </div>
      </template>
      <template v-slot:separator>
        <q-avatar color="primary" text-color="white" size="20px">
          12
        </q-avatar>
      </template>
      <template v-slot:after>
        <q-scroll-area lass="fit" style="height: 150px; width: 100%">
          <div class="q-mx-md resource">
            <q-img style="border: 1px solid #000" width="200px" height="200px"
                   src="C:\Users\zyue\Pictures\QQ截图20240128011140.png"/>
            <q-img style="border: 1px solid #000" width="200px" height="200px"
                   src="C:\Users\zyue\Pictures\webwxgetmsgimg.jpg"/>
          </div>
        </q-scroll-area>
      </template>
    </q-splitter>
    <q-separator class="q-mt-md"/>
    <InformationAccount/>
  </q-card>
</template>

<style scoped lang="scss">
.resource {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 0px;
}
</style>
