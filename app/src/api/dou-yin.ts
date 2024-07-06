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
  getHotspot: async (
    keyword: string = ''
  ): Promise<DouYinHotspotResponse[]> => {
    return await api.get('/utils/douyin/hotspot', {
      params: { keyword }
    })
  },
  getChallenge: async (
    keyword: string = ''
  ): Promise<DouYinChallengeSugResponse[]> => {
    return await api.get('/utils/douyin/challenge', {
      params: { keyword }
    })
  },
  getActivity: async (): Promise<DouYinActivityResponse[]> => {
    return await api.get('/utils/douyin/activity')
  },
  getFlashmob: async (
    keyword: string = ''
  ): Promise<DouYinFlashmobResponse[]> => {
    return await api.get('/utils/douyin/flashmob', {
      params: { keyword }
    })
  }
}
