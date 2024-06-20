import axios from 'axios'
import { LoadingBar, Notify } from 'quasar'

const BASE_URL = 'http://localhost:8080'

const api = axios.create({
  baseURL: BASE_URL,
  timeout: 10000
})

const healthApi = axios.create({
  baseURL: BASE_URL,
  timeout: 3000
})

LoadingBar.setDefaults({
  hijackFilter (url) {
    return url.endsWith('/health')
  }
})

api.interceptors.request.use(config => {
  if (config.method !== 'GET') {
    LoadingBar.start()
  }
  return config
}, error => {
  LoadingBar.stop()
  Notify.create({
    type: 'negative',
    message: `请求失败：${error.message}`
  })
  return Promise.reject(error)
})

api.interceptors.response.use(async (response) => {
  LoadingBar.stop()
  const { status, data } = response
  if (status >= 400 && status < 500) {
    console.error(status, data)
    throw new Error(`获取数据失败：${status}`)
  }
  if (status >= 500) {
    console.error(status, data)
    Notify.create({
      type: 'negative',
      message: '请求失败：服务器内部错误'
    })
    throw new Error(`服务器错误：${data}`)
  }
  if (status === 201) {
    Notify.create({
      type: 'positive',
      message: '创建成功'
    })
    return undefined
  }
  if (status === 204) {
    Notify.create({
      type: 'positive',
      message: '操作成功'
    })
    return undefined
  }
  return data
}, function (error) {
  LoadingBar.stop()
  Notify.create({
    type: 'negative',
    message: `请求失败：${error.message}`
  })
  return Promise.reject(error)
})

export { axios, api, healthApi }
