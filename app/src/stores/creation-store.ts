import { defineStore } from 'pinia'
import { reactive } from 'vue'

export const useInformationStore = defineStore('information', () => {
  const styleData = reactive({
    tab: 'douYin',
    innerTab: 'douYin',
    splitterModel: 20
  })
  return {
    styleData
  }
}, {
  persist: true
})
