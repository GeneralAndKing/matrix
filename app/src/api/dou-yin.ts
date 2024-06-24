import { api } from 'boot/request'

export interface DouYinHotspotResponse {
  eventTime: string
  hotValue: number
  videoCount: number
  word: string
  cover: string
}

export interface DouYinChallengeSugResponse {
  name: string
  viewCount: number
}

export interface DouYinActivityResponse {
  cover: string
  hotScore: string
  name: string
  challenges: string[]
  startTime: string
  endTime: string
}

export interface DouYinFlashmobResponse {
  name: string
  count: number
  cover: string
}

export const DouYinApi = {
  getHotspot: async (): Promise<DouYinHotspotResponse[]> => {
    return await api.get<DouYinHotspotResponse[]>('/douyin/hotspot')
  },
  getChallenge: async (): Promise<DouYinChallengeSugResponse[]> => {
    return await api.get<DouYinChallengeSugResponse[]>('/douyin/challenge')
  },
  getActivity: async (): Promise<DouYinActivityResponse[]> => {
    return await api.get<DouYinActivityResponse[]>('/douyin/activity')
  }
}
