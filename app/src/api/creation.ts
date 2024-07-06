import { BaseModal, EnumModal } from 'src/api/type'
import { api } from 'boot/request'

export interface CreationRequest {
  type: number
  title: string
  description: string
  paths: string[]
}

export interface Creation extends BaseModal, CreationRequest {}

export interface DouYinAccountRelation {
  id: string
  title: string
  description: string
  videoCoverPath?: string
  location?: string
  flashmob?: string
  collectionName?: string
  collectionNum?: number
  associatedHotspot: string
  syncToToutiao: boolean
  allowedToSave: boolean
  whoCanWatch: 1 | 2 | 3
  // 时间错
  releaseTime: number
}

export interface CreationInformation {
  creation: Creation
  douyin: DouYinAccountRelation[]
}

export interface CreationPublishRequest {
  douyin: DouYinAccountRelation[]
}

export const CreationType: EnumModal[] = [
  { value: 1, label: '视频', color: 'teal-6' },
  { value: 2, label: '图文', color: 'cyan-6' }
]

export const CreationTypeMap: Record<number, EnumModal> = CreationType.reduce(
  (acc: Record<number, EnumModal>, data) => {
    acc[data.value] = data
    return acc
  },
  {}
)

export const CreationApi = {
  getAll: async (): Promise<Creation[]> => {
    return api.get('/creation')
  },
  add: async (data: CreationRequest): Promise<void> => {
    return api.post('/creation', {
      Type: data.type,
      Title: data.title,
      Description: data.description,
      Paths: data.paths
    })
  },
  getInformation: async (id: string): Promise<CreationInformation> => {
    return api.get<CreationInformation>(`/creation/${id}`)
  },
  publish: async (
    creationId: string,
    data: CreationPublishRequest
  ): Promise<void> => {
    return api.post(
      '/creation/publish',
      { data },
      { params: { id: creationId } }
    )
  }
}
