import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { DouYinAccountRelation } from 'src/api/creation'
import { date } from 'quasar'
import { useRoute } from 'vue-router'

export interface DouYinPublishForm extends DouYinAccountRelation {
  videoCoverFile?: File,
  challengeList?: string[],
  releaseType: 0 | 1,
  releaseDateString: string,
}
export const DOU_YIN_GLOBAL_TAB = 'global'

export const useDouYinPublishStore = defineStore('douYinPublish', () => {
  const batchConfigMap = ref<Record<string, DouYinPublishForm>>({})
  const route = useRoute()
  const getDefaultConfig = (): DouYinPublishForm => ({
    id: '0',
    title: '',
    description: '',
    associatedHotspot: '',
    syncToToutiao: false,
    allowedToSave: true,
    whoCanWatch: 1,
    releaseTime: 0,
    releaseType: 0,
    releaseDateString: date.formatDate(date.addToDate(new Date(), { hour: 2 }), 'YYYY-MM-DD HH:mm')
  })
  const toConfigForm = (config: DouYinAccountRelation): DouYinPublishForm => {
    return {
      ...config,
      releaseType: config.releaseTime === 0 ? 0 : 1,
      // TODO: 需要获取格式化的结果
      releaseDateString: date.formatDate(date.addToDate(new Date(), { hour: 2 }), 'YYYY-MM-DD HH:mm')
    }
  }
  const toRequestParam = (userId: number | string, form: DouYinPublishForm): DouYinAccountRelation => {
    return {
      id: `${userId}`,
      title: form.title,
      description: form.description,
      location: form.location,
      flashmob: form.flashmob,
      collectionName: form.collectionName,
      collectionNum: form.collectionNum,
      associatedHotspot: form.associatedHotspot,
      syncToToutiao: form.syncToToutiao,
      allowedToSave: form.allowedToSave,
      whoCanWatch: form.whoCanWatch,
      videoCoverPath: form.videoCoverFile?.path,
      releaseTime: form.releaseType === 0 ? 0 : date.extractDate(form.releaseDateString, 'YYYY-MM-DD HH:mm').getTime()
    }
  }
  const currentBatchConfig = computed({
    get () {
      return batchConfigMap.value[route.params.id as string] ?? getDefaultConfig()
    },
    set (newValue) {
      batchConfigMap.value[route.params.id as string] = newValue
    }
  })

  return {
    batchConfigMap,
    currentBatchConfig,
    getDefaultConfig,
    toConfigForm,
    toRequestParam
  }
}, {
  persist: true
})
