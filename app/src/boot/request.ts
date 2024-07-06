import axios from 'axios'
import { LoadingBar, Notify } from 'quasar'

const BASE_URL = 'http://localhost:8080'

const api = axios.create({
  baseURL: BASE_URL
})

const healthApi = axios.create({
  baseURL: BASE_URL
})

LoadingBar.setDefaults({
  hijackFilter (url) {
    return url.endsWith('/health')
  }
})

api.interceptors.request.use(
  (config) => {
    if (config.method !== 'GET') {
      LoadingBar.start()
    }
    return config
  },
  (error) => {
    LoadingBar.stop()
    Notify.create({
      type: 'negative',
      message: `请求失败：${error.message}`
    })
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  async (response) => {
    LoadingBar.stop()
    const { status, data } = response
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
  },
  (error) => {
    const { response } = error
    LoadingBar.stop()
    console.error(error)
    if (response.status >= 400 && response.status < 500) {
      Notify.create({
        type: 'warning',
        message: '获取数据失败：请检查参数是否合法'
      })
      return Promise.reject(error.response)
    }
    Notify.create({
      type: 'negative',
      message: '请求失败：服务器内部错误'
    })
    return Promise.reject(error.response)
  }
)

export { axios, api, healthApi }
