import { defineStore } from 'pinia'
import { onMounted, ref } from 'vue'
import { businessSocket, messageSocket } from 'boot/socket'

export const useHealthStore = defineStore('health', () => {
  const isOk = ref(false)

  const handleHealth = () => {
    messageSocket.onopen = () => {
      isOk.value = true
    }
    messageSocket.onclose = () => {
      isOk.value = false
    }
    messageSocket.onerror = () => {
      console.log('error')
    }
    messageSocket.onmessage = () => {
      console.log('message')
    }
    businessSocket.onopen = () => {
      isOk.value = true
    }
    businessSocket.onclose = () => {
      isOk.value = false
    }
  }
  onMounted(() => {
    handleHealth()
  })
  return {
    isOk
  }
})
