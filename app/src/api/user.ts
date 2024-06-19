import { api } from 'boot/request'
import { BaseModal } from 'src/api/type'

export interface DouYinUser extends BaseModal{
  name: string
  description?: string
  douyinId: string
  avatar: string
  labels?: string[]
  expired: boolean
}

export const DouYinUserApi = {
  getAll: async (): Promise<DouYinUser[]> => {
    return await api.get<DouYinUser[]>('/user/douyin')
  },
  add: async (): Promise<void> => {
    return await api.post('/user/douyin', [])
  },
  refresh: async (id: number): Promise<void> => {
    return await api.post(`/user/douyin/${id}/refresh`)
  }
}
