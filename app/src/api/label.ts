import { api } from 'boot/request'
import { BaseModal } from 'src/api/type'

export interface Label extends BaseModal {
  Name: string
}

export const LabelApi = {
  getAll: async (): Promise<Label[]> => {
    return await api.get<Label[]>('/label')
  }
}
