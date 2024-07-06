<script setup lang="ts">
import { DouYinLive, DouYinLiveApi } from 'src/api/live'
import { onMounted, reactive } from 'vue'

interface Data {
  liveList: DouYinLive[]
}

const data = reactive<Data>({
  liveList: []
})

const handleData = () => {
  DouYinLiveApi.getAll().then((res) => {
    data.liveList = res
    console.log(res)
  })
}

onMounted(() => {
  handleData()
})
</script>

<template>
  <div>
    Live page
  </div>
</template>

<style scoped lang="scss">

</style>
