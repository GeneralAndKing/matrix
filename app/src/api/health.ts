import { healthApi } from 'boot/request'
import { AxiosResponse } from 'axios'

export const HealthApi = {
  ping: async (): Promise<AxiosResponse<string>> => {
    return await healthApi.get<string>('/ping')
  }
}
