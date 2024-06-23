import { BaseModal, EnumModal } from 'src/api/type'
import { api } from 'boot/request'

export interface CreationRequest {
  type: number
  title: string
  description: string
  paths: string[]
}

export interface Creation extends BaseModal, CreationRequest {
}

export const CreationType: EnumModal[] = [
  { value: 1, label: '视频', color: 'teal-6' },
  { value: 2, label: '图文', color: 'cyan-6' }
]

export const CreationTypeMap: Record<number, EnumModal> =
  CreationType.reduce((acc: Record<number, EnumModal>, data) => {
    acc[data.value] = data
    return acc
  }, {})

export const CreationApi = {
  getAll: async (): Promise<Creation[]> => {
    return api.get('/creation')
  },
  add: async (data: CreationRequest): Promise<void> => {
    return api.post('/creation', data)
  }
}
