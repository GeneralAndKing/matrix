import { BaseModal } from 'src/api/type'
import { api } from 'boot/request'

export enum DouYinLiveMonitor {
  PauseLiveMonitorStatus = 1,
  RunningLiveMonitorStatus,
  NotExistLiveMonitorStatus
}

export interface DouYinLive extends BaseModal {
  liveId: string
  name: string
  douYinId: string
  avatar: string
  labels: string[]
  monitor: DouYinLiveMonitor
}

export const DouYinLiveApi = {
  getAll: async (): Promise<DouYinLive[]> => {
    return await api.get('/live/douyin')
  }
}
