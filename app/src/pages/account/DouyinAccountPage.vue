<script setup lang="ts">

import { onMounted, reactive } from 'vue'
import { DouYinUserApi, DouYinUser } from 'src/api/user'
import { date, QSelectProps, QTableProps, useQuasar } from 'quasar'
import { LabelApi } from 'src/api/label'

interface ActionState {
  refreshLoading?: boolean
  deleteLoading?: boolean
}

interface Data {
  userList: (DouYinUser & ActionState)[],
  labelList: string[],
  filterLabelList: string[],
  labelLoading: boolean,
  loading: boolean,
  search: string
}

const $q = useQuasar()
const data = reactive<Data>({
  userList: [],
  labelList: [],
  filterLabelList: [],
  labelLoading: false,
  loading: false,
  search: ''
})

const columns: QTableProps['columns'] = [
  { name: 'avatar', label: '', align: 'left', field: 'avatar' },
  { name: 'name', label: '用户名称', align: 'left', field: 'name', sortable: true },
  { name: 'douyinId', label: '抖音 ID', align: 'left', field: 'douyinId', sortable: true },
  { name: 'description', label: '描述', align: 'left', field: 'description' },
  { name: 'labels', label: '标签', align: 'left', field: 'labels' },
  { name: 'status', label: '状态', align: 'center', field: 'ID', sortable: true },
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

const handleLabels = () => {
  data.labelLoading = true
  LabelApi.getAll().then(res => {
    data.labelList = (res ?? []).map(item => item.Name)
    data.filterLabelList = (res ?? []).map(item => item.Name)
  }).finally(() => {
    data.labelLoading = false
  })
}

onMounted(() => {
  handleData()
  handleLabels()
})

const handleAdd = async () => {
  DouYinUserApi.add()
    .then(handleData)
}

const handleRefresh = (row: DouYinUser, rowIndex: number) => {
  data.userList[rowIndex].refreshLoading = true
  DouYinUserApi.refresh(row.ID)
    .finally(() => {
      data.userList[rowIndex].refreshLoading = false
    })
}

const handleManager = async (row: DouYinUser) => {
  await DouYinUserApi.manage(row.ID)
}

const handleDelete = (row: DouYinUser) => {
  $q.dialog({
    title: '确认操作',
    message: '真的要删除当前账号吗？删除后不可恢复哦！',
    cancel: { color: 'primary', outline: true },
    ok: { color: 'negative', outline: true }
  }).onOk(() => {
    DouYinUserApi.delete(row.ID)
      .then(handleData)
  })
}

const handleAddLabel: QSelectProps['onNewValue'] = (val: string, done) => {
  if (val.length > 0) {
    if (!data.labelList.includes(val)) {
      data.labelList.push(val)
    }
    done(val, 'toggle')
  }
}

const handleFilterLabel: QSelectProps['onFilter'] = (val: string, doneFn) => {
  doneFn(() => {
    if (val === '') {
      data.filterLabelList = data.labelList
    } else {
      const needle = val.toLowerCase()
      data.filterLabelList = data.labelList.filter(
        v => v.toLowerCase().indexOf(needle) > -1
      )
    }
  })
}

const handleUpdateLabel = (row: DouYinUser, label: string[], initialValue: string[]) => {
  DouYinUserApi.update(row.ID, label)
    .catch(() => { row.labels = initialValue })
}

</script>

<template>
  <div class="q-mt-md column">
    <div class="row justify-between items-center">
      <q-input debounce="300" v-model="data.search" placeholder="搜索关键字" dense>
        <template v-slot:append>
          <q-icon name="search"/>
        </template>
      </q-input>
      <div>
        <q-btn style="max-width: 200px" size="sm" outline color="primary" label="新增账号" @click="handleAdd"/>
      </div>
    </div>
    <q-table
      :rows="data.userList"
      :columns="columns"
      row-key="name"
      :filter="data.search"
      :loading="data.loading"
      flat
    >
      <template v-slot:header-cell-labels="props">
        <q-th :props="props">
          {{ props.col.label }}
          <q-icon name="info" color="primary" size="1rem">
            <q-tooltip anchor="top middle" self="center middle">
              点击单元格可以进行编辑哦！
            </q-tooltip>
          </q-icon>
        </q-th>
      </template>
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
          <q-td key="labels" class="cursor" auto-width :props="props">
            <div class="cursor-pointer">
              <span v-if="props.row.labels && props.row.labels.length > 0">
                <q-chip v-for="label in props.row.labels" :key="label" color="primary"
                        outline square dense :label="label"/>
              </span>
              <span v-else>
                -
              </span>
              <q-popup-edit @save="(value, initialValue) => handleUpdateLabel(props.row, value, initialValue)" v-model="props.row.labels" buttons label-set="保存" label-cancel="取消" v-slot="scope">
                <q-select
                  v-model="scope.value"
                  dense
                  use-input
                  use-chips
                  multiple
                  input-debounce="0"
                  @new-value="handleAddLabel"
                  :options="data.filterLabelList"
                  @filter="handleFilterLabel"
                  style="width: 250px"
                />
              </q-popup-edit>
            </div>
          </q-td>
          <q-td key="status" :props="props" auto-width>
            <q-badge rounded :color="props.row.expired ? 'red' : 'green'">
              <q-tooltip anchor="top middle" self="center middle">
                {{props.row.expired ? '登录已过期': '已登录'}} 最后更新 {{ date.formatDate(new Date(props.row.UpdatedAt), 'YYYY-MM-DD HH:mm:ss') }}
              </q-tooltip>
            </q-badge>
          </q-td>
          <q-td key="action" :props="props" auto-width>
            <q-btn :loading="props.row.refreshLoading" flat round color="secondary" icon="refresh" size="sm" dense
                   @click="() => handleRefresh(props.row, props.rowIndex)">
              <q-tooltip anchor="top middle" self="center middle">刷新</q-tooltip>
            </q-btn>
            <q-btn flat round color="accent" icon="manage_accounts" size="sm" dense
                   @click="() => handleManager(props.row)">
              <q-tooltip anchor="top middle" self="center middle">账号管理</q-tooltip>
            </q-btn>
            <q-btn :loading="props.row.deleteLoading" flat round color="red" icon="delete" size="sm" dense
                   @click="() => handleDelete(props.row)">
              <q-tooltip anchor="top middle" self="center middle">删除</q-tooltip>
            </q-btn>
          </q-td>
        </q-tr>
      </template>
      <template v-slot:no-data="{ icon }">
        <div class="full-width row flex-center q-gutter-sm">
          <q-icon size="2em" :name="icon"/>
          <template v-if="data.search.trim() === ''">
            <span>没有查询到任何账号数据哦~快添加一个试试吧！</span>
          </template>
          <template v-else>
            <span>没有符合条件的账号信息哦~修改条件重新查询下吧！</span>
          </template>
        </div>
      </template>
    </q-table>
  </div>
</template>

<style scoped>

</style>
