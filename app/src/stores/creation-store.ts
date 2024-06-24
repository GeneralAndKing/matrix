import { defineStore } from 'pinia'
import { reactive } from 'vue'
import { CreationApi, CreationInformation } from 'src/api/creation'

interface Data {
  information: CreationInformation | null
}

export const useInformationStore = defineStore('information', () => {
  const styleData = reactive({
    tab: 'douYin',
    innerTab: 'douYin',
    splitterModel: 20
  })

  const data = reactive<Data>({
    information: null
  })

  const handleInformation = (id: number) => {
    CreationApi.getInformation(id)
      .then(res => {
        data.information = res
      })
  }

  return {
    styleData,
    handleInformation,
    data
  }
}, {
  persist: true
})
