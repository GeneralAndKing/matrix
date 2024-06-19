<script setup lang="ts">

import { onMounted, reactive } from 'vue'
import { DouYinUserApi, DouYinUser } from 'src/api/user'
import { QTableProps } from 'quasar'

interface ActionState {
  refreshLoading?: boolean
  deleteLoading?: boolean
}

interface Data {
  userList: (DouYinUser & ActionState)[]
  loading: boolean
}

const data = reactive<Data>({
  userList: [],
  loading: false
})

const columns: QTableProps['columns'] = [
  { name: 'avatar', label: '', align: 'left', field: 'avatar' },
  { name: 'name', label: '用户名称', align: 'left', field: 'name', sortable: true },
  { name: 'douyinId', label: '抖音 ID', align: 'left', field: 'douyinId', sortable: true },
  { name: 'description', label: '描述', align: 'left', field: 'description', sortable: true },
  { name: 'action', label: '操作', align: 'center', field: 'ID' }
]

const handleData = () => {
  data.loading = true
  DouYinUserApi.getAll().then(res => {
    data.userList = res ?? []
  }).finally(() => {
    data.loading = false
  })
}

onMounted(() => {
  handleData()
})

const handleAdd = async () => {
  await DouYinUserApi.add()
}

const handleRefresh = (row: DouYinUser, rowIndex: number) => {
  data.userList[rowIndex].refreshLoading = true
  DouYinUserApi.refresh(row.ID)
    .finally(() => {
      data.userList[rowIndex].refreshLoading = false
    })
}

const handleDelete = (row: DouYinUser, rowIndex: number) => {
  console.log(row, rowIndex)
}
</script>

<template>
  <div class="q-mt-md column">
    <div>
      <q-btn style="max-width: 200px" size="sm" outline color="primary" label="新增账号" @click="handleAdd"/>
    </div>
    <q-table
      :rows="data.userList"
      :columns="columns"
      row-key="name"
      :loading="data.loading"
      flat
    >
      <template v-slot:body="props">
        <q-tr :props="props">
          <q-td key="avatar" :props="props" auto-width>
            <q-avatar class="shadow-3">
              <q-img :src="props.row.avatar"/>
            </q-avatar>
          </q-td>
          <q-td key="name" :props="props" auto-width>
            {{ props.row.name }}
          </q-td>
          <q-td key="douyinId" :props="props" auto-width>
            {{ props.row.douyinId }}
          </q-td>
          <q-td key="description" :props="props">
            {{ props.row.description }}
          </q-td>
          <q-td key="action" :props="props" auto-width>
            <q-btn :loading="props.row.refreshLoading" flat round color="secondary" icon="refresh" size="sm" dense @click="() => handleRefresh(props.row, props.rowIndex)">
              <q-tooltip anchor="top middle" self="center middle">刷新</q-tooltip>
            </q-btn>
            <q-btn :loading="props.row.deleteLoading" flat round color="red" icon="delete" size="sm" dense @click="() => handleDelete(props.row, props.rowIndex)">
              <q-tooltip anchor="top middle" self="center middle">删除</q-tooltip>
            </q-btn>
          </q-td>
        </q-tr>
      </template>
      <template v-slot:no-data="{ icon }">
        <div class="full-width row flex-center q-gutter-sm">
          <q-icon size="2em" :name="icon"/>
          <span>没有查询到任何数据哦~</span>
          <span class="cursor-pointer text-primary" @click="handleAdd">快点我添加吧！</span>
        </div>
      </template>
    </q-table>
  </div>
</template>

<style scoped>

</style>
