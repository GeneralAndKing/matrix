import { defineStore } from 'pinia'
import { HealthApi } from 'src/api/health'
import { onMounted, ref } from 'vue'
import { useInterval } from 'quasar'

export const useHealthStore = defineStore('counter', () => {
  const isOk = ref(false)

  const {
    registerInterval
  } = useInterval()

  const handleHealth = () => {
    HealthApi.ping()
      .then((response) => {
        isOk.value = response.data === 'pong'
      })
      .catch(() => { isOk.value = false })
  }
  onMounted(() => {
    handleHealth()
    registerInterval(handleHealth, 2000)
  })
  return {
    isOk
  }
})
