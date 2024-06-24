<script setup lang="ts">

import { onMounted, reactive } from 'vue'
import { Creation, CreationApi, CreationTypeMap } from 'src/api/creation'
import { QTableProps, useQuasar } from 'quasar'
import CreationEditDialog from 'pages/creation/list/CreationEditDialog.vue'
import { useRouter } from 'vue-router'

const $q = useQuasar()
const router = useRouter()
const columns: QTableProps['columns'] = [
  { name: 'type', label: '状态', align: 'left', field: 'type', sortable: true },
  { name: 'title', label: '标题', align: 'left', field: 'title' },
  { name: 'description', label: '描述', align: 'left', field: 'description' },
  { name: 'count', label: '资源', align: 'left', field: 'ID' },
  { name: 'action', label: '操作', align: 'left', field: 'ID' }
]

interface Data {
  creationList: Creation[],
  loading: boolean,
  search: string
}

const data = reactive<Data>({
  creationList: [],
  loading: false,
  search: ''
})

const handleData = () => {
  CreationApi.getAll()
    .then(res => { data.creationList = res })
}

onMounted(() => {
  handleData()
})

const handleAdd = () => {
  $q.dialog({
    component: CreationEditDialog,
    componentProps: {
      item: null
    }
  }).onOk(() => {
    handleData()
  })
}

const handleToInformation = (id: number) => {
  void router.push(`creation/information/${id}`)
}

</script>

<template>
  <div>
    <q-table
      :rows="data.creationList"
      :columns="columns"
      row-key="name"
      :loading="data.loading"
      :filter="data.search"
      flat
      bordered
    >
      <template v-slot:top-left>
        <div class="text-h6">创作集合</div>
        <q-btn class="q-ml-sm" icon="add_box" dense rounded flat color="primary" @click="handleAdd"/>
      </template>
      <template v-slot:top-right>
        <q-input borderless dense debounce="300" v-model="data.search" placeholder="输入关键字进行搜索">
          <template v-slot:append>
            <q-icon name="search" />
          </template>
        </q-input>
      </template>
      <template v-slot:body="props">
        <q-tr class="cursor-pointer" :props="props" @click="() => handleToInformation(props.row.ID)">
          <q-td key="type" :props="props" auto-width>
            <q-badge :color="CreationTypeMap[props.row.type]?.color">
              {{ CreationTypeMap[props.row.type]?.label ?? '未知' }}
            </q-badge>
          </q-td>
          <q-td key="title" :props="props" auto-width>
            {{ props.row.title }}
          </q-td>
          <q-td key="description" :props="props">
            {{ props.row.description }}
          </q-td>
          <q-td key="count" :props="props" auto-width>
            <span class="text-green">1</span>
            <span class="text-gray" style="opacity: 0.2">-</span>
            <span class="text-red">2</span>
          </q-td>
          <q-td key="action" :props="props" auto-width>
            -
          </q-td>
        </q-tr>
      </template>
    </q-table>
  </div>
</template>
