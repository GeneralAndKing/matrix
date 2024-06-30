import { defineStore } from 'pinia'
import { reactive } from 'vue'
import { CreationApi, CreationInformation, CreationPublishRequest } from 'src/api/creation'
import { useQuasar } from 'quasar'

interface Data extends CreationPublishRequest{
  information: CreationInformation | null
}

export const useInformationStore = defineStore('information', () => {
  const $q = useQuasar()
  const styleData = reactive({
    tab: 'douYin',
    douYinTab: 'global',
    splitterModel: 20
  })

  const data = reactive<Data>({
    information: null,
    douyin: []
  })

  const handleInformation = async (id: string) => {
    const res = await CreationApi.getInformation(id)
    data.information = res
    // data.douyin = res.douyin ?? []
    return res
  }

  const handlePublish = async (id: string) => {
    if (data.douyin.length === 0) {
      $q.notify({ type: 'warning', message: '至少需要配置一个账号才可以发布哦~' })
      return
    }
    await CreationApi.publish(`${data.information?.creation.ID}`, {
      douyin: data.douyin
    })
    await handleInformation(id)
  }

  return {
    styleData,
    handleInformation,
    data,
    handlePublish
  }
}, {
  persist: true
})
