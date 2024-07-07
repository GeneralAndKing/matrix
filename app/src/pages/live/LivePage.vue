<script setup lang="ts">
import { DouYinLive, DouYinLiveApi } from 'src/api/live'
import { onMounted, reactive } from 'vue'
import { QTableProps, useQuasar } from 'quasar'
import { useRouter } from 'vue-router'

interface Data {
  liveList: DouYinLive[]
  search: string,
  loading: {
    data: boolean
    add: boolean
  }
}

const router = useRouter()
const $q = useQuasar()
const data = reactive<Data>({
  liveList: [],
  search: '',
  loading: {
    data: false,
    add: false
  }
})
const columns: QTableProps['columns'] = [
  { name: 'avatar', label: '头像', align: 'left', field: 'name', sortable: true },
  { name: 'liveId', label: '直播间号', align: 'left', field: 'liveId' },
  { name: 'labels', label: '标签', align: 'left', field: 'labels' },
  { name: 'monitor', label: '状态', align: 'left', field: 'monitor', sortable: true }
]

const handleData = () => {
  data.loading.data = true
  DouYinLiveApi.getAll().then((res) => {
    data.liveList = res
    console.log(res)
  }).finally(() => {
    data.loading.data = false
  })
}

onMounted(() => {
  handleData()
})

const handleAdd = () => {
  $q.dialog({
    title: '添加直播间号',
    message: '直接输入直播间号或者链接地址，例如 xxxx 或者 https://www.douyin.com/follow/live/xxxx 或者 https://live.douyin.com/xxxx',
    prompt: {
      model: '',
      type: 'text'
    },
    cancel: true
  }).onOk((param: string) => {
    let liveId = param
    if (liveId.startsWith('https://www.douyin.com/follow/live') ||
      liveId.startsWith('https://live.douyin.com')) {
      if (liveId.includes('?')) {
        liveId = liveId.substring(liveId.lastIndexOf('/') + 1, liveId.lastIndexOf('?'))
      } else {
        liveId = liveId.substring(liveId.lastIndexOf('/') + 1)
      }
    }
    data.loading.add = true
    DouYinLiveApi.add(liveId)
      .then(handleData)
      .finally(() => {
        data.loading.add = false
      })
  })
}

const handleDelete = (id: number) => {
  $q.dialog({
    title: '确认删除吗？',
    message: '当前操作会删除指定项的所有数据且不可逆，请谨慎操作！',
    cancel: { color: 'red', flat: true },
    focus: 'cancel'
  }).onOk(async () => {
    await DouYinLiveApi.delete(id)
    handleData()
  })
}

const handleToInformation = (id: number) => {
  router.push({ name: 'LiveDashboard', params: { id } })
}
</script>

<template>
  <div>
    <div class="flex row items-center">
      <div class="text-h6 text-bold">抖音直播列表</div>
      <q-space/>
      <q-input class="q-mr-md" model-value="" outlined dense placeholder="输入关键字" />
      <q-btn :loading="data.loading.add" outline color="primary" label="添加" @click="handleAdd"/>
    </div>
    <q-table class="q-mt-md" :rows="data.liveList" :columns="columns" row-key="liveId" :loading="data.loading.data" flat bordered>
      <template v-slot:body="props">
        <q-tr :props="props" class="cursor-pointer" @click="() => handleToInformation(props.row.ID)">
          <q-td key="avatar" :props="props" auto-width>
            <q-item>
              <q-item-section top thumbnail class="q-ml-none">
                <q-avatar class="shadow-2">
                  <q-img img-class="full-height full-width" fit="cover" :src="props.row.avatar"/>
                </q-avatar>
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ props.row.name }}</q-item-label>
                <q-item-label caption>
                  <q-badge>ID: {{props.row.douYinId}}</q-badge>
                </q-item-label>
              </q-item-section>
            </q-item>
          </q-td>
          <q-td key="liveId" :props="props">
            {{props.row.liveId}}
          </q-td>
          <q-td key="labels" :props="props">
            <template v-if="!props.row.labels || props.row.labels?.length === 0">
              -
            </template>
            <template v-else>
            </template>
          </q-td>
          <q-td key="monitor" :props="props">
            <span v-if="props.row.monitor === 1"><q-badge rounded class="q-mr-xs" color="primary" />未开播</span>
            <span v-else-if="props.row.monitor === 2"><q-badge rounded class="q-mr-xs" color="secondary" />正在直播</span>
            <span v-else><q-badge rounded class="q-mr-xs" color="red" />错误状态</span>
          </q-td>
          <q-menu
            touch-position
            context-menu
          >
            <q-list dense style="min-width: 100px">
              <q-item clickable v-close-popup>
                <q-item-section class="text-red" @click="() => handleDelete(props.row.ID)">删除</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-tr>
      </template>
    </q-table>
  </div>
</template>

<style scoped lang="scss">

</style>
