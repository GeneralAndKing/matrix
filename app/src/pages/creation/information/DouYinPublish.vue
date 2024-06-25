<script setup lang="ts">
import { date, QDate, QSelect, QTime } from 'quasar'
import { useInformationStore } from 'stores/creation-store'
import { DouYinUser, DouYinUserApi } from 'src/api/user'
import { DouYinAccountRelation } from 'src/api/creation'
import { onMounted, reactive } from 'vue'
import {
  DouYinActivityResponse,
  DouYinApi,
  DouYinChallengeSugResponse,
  DouYinFlashmobResponse,
  DouYinHotspotResponse
} from 'src/api/dou-yin'
import { formatChineseNumber, formatNumber } from 'src/utils'

interface Data {
  userList: DouYinUser[],
  filterUserList: DouYinUser[],
  confirmPublishList: number[],
  publishList: number[],
  activityList: DouYinActivityResponse[],
  challengeList: DouYinChallengeSugResponse[],
  challengeLoading: boolean,
  hotspotList: DouYinHotspotResponse[],
  hotspotLoading: boolean,
  flashmobList: DouYinFlashmobResponse[],
  flashmobLoading: boolean,
  batchConfig: DouYinAccountRelation & {
    videoCoverFile?: File,
    challengeList?: string[],
    releaseType: 0 | 1,
    releaseDateString: string,
  }
}

const store = useInformationStore()
const data = reactive<Data>({
  publishList: [],
  userList: [],
  filterUserList: [],
  confirmPublishList: [],
  activityList: [],
  challengeList: [],
  challengeLoading: false,
  hotspotList: [],
  hotspotLoading: false,
  flashmobList: [],
  flashmobLoading: false,
  batchConfig: {
    id: '0',
    title: store.data.information?.creation.title ?? '',
    description: store.data.information?.creation.description ?? '',
    associatedHotspot: '',
    syncToToutiao: false,
    allowedToSave: true,
    whoCanWatch: 1,
    releaseTime: 0,
    releaseType: 0,
    releaseDateString: date.formatDate(date.addToDate(new Date(), { hour: 2 }), 'YYYY-MM-DD HH:mm')
  }
})

const handleInitDouYinData = () => {
  DouYinUserApi.getAll().then(res => {
    data.userList = res
  })
  DouYinApi.getActivity().then(res => {
    data.activityList = res.filter(item => item.challenges.length > 0)
  })
  data.hotspotLoading = true
  DouYinApi.getHotspot()
    .then(res => {
      data.hotspotList = res
    })
    .then(() => {
      data.hotspotLoading = false
    })
  DouYinApi.getChallenge().then(res => {
    data.challengeList = res
  })
  data.flashmobLoading = true
  DouYinApi.getFlashmob()
    .then(res => {
      data.flashmobList = res
    })
    .finally(() => {
      data.flashmobLoading = false
    })
}

onMounted(() => {
  handleInitDouYinData()
})

const handleFilterDouYinUser: QSelect['onFilter'] = (val, update) => {
  if (val === '') {
    update(() => {
      data.filterUserList = data.userList
    })
    return
  }

  update(() => {
    const needle = val.toLowerCase()
    data.filterUserList = data.userList.filter(v => (v.name.toLowerCase().indexOf(needle) > -1))
  })
}

const handleFilterFlashmobList: QSelect['onFilter'] = async (val, update) => {
  try {
    data.flashmobLoading = true
    const res = await DouYinApi.getFlashmob(val)
    update(() => {
      data.flashmobList = res
    })
  } finally {
    data.flashmobLoading = false
  }
}

const handleFilterHotspotList: QSelect['onFilter'] = async (val, update) => {
  try {
    data.hotspotLoading = true
    const res = await DouYinApi.getHotspot(val)
    update(() => {
      data.hotspotList = res
    })
  } finally {
    data.hotspotLoading = false
  }
}

const handleFilterChallengeList: QSelect['onFilter'] = async (val, update) => {
  try {
    data.challengeLoading = true
    const res = await DouYinApi.getChallenge(val)
    update(() => {
      data.challengeList = res
    })
  } finally {
    data.challengeLoading = false
  }
}

const handleAddHotspot = (activity: DouYinActivityResponse, id: number | undefined) => {
  console.log(activity)
  if (id) {
    return
  }
  const value = ` #${activity.challenges[0]} `
  if (data.batchConfig.description.includes(value)) {
    data.batchConfig.description.replaceAll(value, '')
  } else {
    data.batchConfig.description += value
  }
}

const handleReleaseDateLimit: QDate['options'] = (currentDate: string) => {
  const now = new Date()
  const start = date.addToDate(now, { hour: 2 })
  const end = date.addToDate(now, { day: 14 })
  return date.isBetweenDates(currentDate, start, end)
}

const handleReleaseTimeLimit: QTime['options'] = (hr, min) => {
  const now = new Date()
  const start = date.addToDate(now, { hour: 2 })
  const end = date.addToDate(now, { day: 14 })
  const dateString = data.batchConfig.releaseDateString.split(' ')[0]
  const select = new Date(`${dateString} ${hr}:${min ?? 0}`)
  return date.isBetweenDates(select, start, end)
}
</script>

<template>
  <q-splitter
    class="full-height"
    v-model="store.styleData.splitterModel"
  >
    <template v-slot:before>
      <q-tabs
        v-model="store.styleData.innerTab"
        vertical
        dense
        active-color="primary"
        indicator-color="primary"
        class="text-grey"
        align="left"
      >
        <q-tab style="justify-content: start" name="global" label="批量设置"/>
        <q-tab
          style="justify-content: start"
          v-for="user in data.userList"
          :key="`config-dou-yin-tab-${user.ID}`"
          :name="`${user.ID}-${user.name}`"
          :label="user.name"
          content-class="account-list-tab"
        >
          <q-badge v-if="data.confirmPublishList.includes(user.ID)" class="q-ml-sm" label="已配置"/>
        </q-tab>
      </q-tabs>
    </template>
    <template v-slot:after>
      <q-tab-panels
        v-model="store.styleData.innerTab"
        animated
        transition-prev="fade"
        transition-next="fade"
        class="full-height"
      >
        <q-tab-panel name="global" class="q-pr-none q-py-none">
          <q-scroll-area class="full-height full-width">
            <div class="flex column q-gutter-y-sm" style="max-width: 550px">
              <div>发布账号</div>
              <q-select
                v-model="data.publishList"
                multiple
                outlined
                placeholder="搜索或选择需要设置的账号信息"
                emit-value
                map-options
                @filter="handleFilterDouYinUser"
                dense
                use-input
                input-debounce="0"
                :options="data.filterUserList"
                option-value="ID"
                option-label="name"
                autocomplete=""
              >
                <template v-slot:no-option>
                  <q-item>
                    <q-item-section class="text-grey">
                      没有找到符合要求的账号哦~
                    </q-item-section>
                  </q-item>
                </template>
                <template v-slot:option="scope">
                  <q-item v-bind="scope.itemProps">
                    <q-item-section avatar>
                      <q-avatar>
                        <q-img :src="scope.opt.avatar"/>
                      </q-avatar>
                    </q-item-section>
                    <q-item-section>
                      <q-item-label>{{ scope.opt.name }}</q-item-label>
                      <q-item-label caption>{{ scope.opt.douyinId }}</q-item-label>
                    </q-item-section>
                  </q-item>
                </template>
                <template v-slot:selected-item="scope">
                  <q-chip
                    removable
                    dense
                    @remove="scope.removeAtIndex(scope.index)"
                    :tabindex="scope.tabindex"
                    class="q-ma-none"
                  >
                    <q-avatar>
                      <q-img :src="scope.opt.avatar"/>
                    </q-avatar>
                    {{ scope.opt.name }}
                  </q-chip>
                </template>
              </q-select>
              <div>作品名称</div>
              <q-input class="q-mt-sm" outlined dense v-model="data.batchConfig.title"
                       placeholder="好的作品标题可以获得更多浏览"/>
              <div>作品简介</div>
              <q-input outlined dense v-model="data.batchConfig.description" type="textarea"/>
              <div>作品标签</div>
              <q-select
                outlined
                dense
                emit-value
                map-options
                clearable
                usein
                use-input
                multiple
                @filter="handleFilterChallengeList"
                :placeholder="(data.batchConfig.challengeList?.length ?? 0) > 0 ? '' : '搜索或选择标签'"
                option-value="name"
                max-values="5"
                option-label="name"
                :popup-content-style="{height: '300px'}"
                v-model="data.batchConfig.challengeList"
                :options="data.challengeList"
                autocomplete="true"
                :loading="data.challengeLoading"
              >
                <template v-slot:option="scope">
                  <q-item v-bind="scope.itemProps">
                    <div class="flex row full-width justify-between items-center">
                      <div>#{{ scope.opt.name }}</div>
                      <div class="text-tip">{{ formatChineseNumber(scope.opt.viewCount) }}跟拍</div>
                    </div>
                  </q-item>
                </template>
              </q-select>
              <div class="row items-center">
                作品活动奖励
                <q-icon name="help" class="cursor-pointer q-ml-xs text-tip">
                  <q-tooltip anchor="top middle" self="center middle">
                    添加活动将有机会获得流量奖励
                  </q-tooltip>
                </q-icon>
                <q-badge class="bg-red-1 q-ml-sm text-red" label="NEW"/>
                <q-space/>
                <div class="text-tip text-caption">了解更多官方活动 ></div>
              </div>
              <div class="q-mt-sm" style="width: 100%; height: 60px">
                <q-scroll-area class="full-height full-width no-wrap">
                  <div class="row q-gutter-x-md no-wrap">
                    <div
                      v-for="item in data.activityList.slice(0, 3)"
                      :key="item.name"
                      :class="`cursor-pointer bg-grey-3 q-pa-sm activity ${data.batchConfig.description.includes(item.challenges[0])}`"
                      style="width: 153px"
                      @click="() => handleAddHotspot(item, undefined)"
                    >
                      <div class="ellipsis text-body2 text-bold">{{ item.name }}</div>
                      <div class="hot-value text-tip text-caption">热度:{{
                          formatNumber(parseInt(item.hotScore))
                        }}
                      </div>
                      <div class="action text-tip flex justify-between text-caption">
                        <div>添加</div>
                        <div>活动详情</div>
                      </div>
                    </div>
                    <div class="cursor-pointer hotspot bg-grey-3 q-pa-sm flex flex-center text-tip" style="flex: 1">
                      +{{ data.hotspotList.length - 3 }}
                    </div>
                  </div>
                </q-scroll-area>
              </div>
              <div>设置视频封面</div>
              <q-file
                v-model="data.batchConfig.videoCoverFile"
                outlined
                dense
                label="选择 png/jpeg/jpg 格式图片"
              />
              <div>设置标签</div>
              <q-input v-model="data.batchConfig.location" dense outlined placeholder="地理位置"/>
              <div class="row items-center">
                添加挑战贴纸
                <q-icon name="help" class="cursor-pointer q-ml-xs text-tip">
                  <q-tooltip anchor="top middle" self="center middle">
                    优质作品参与挑战，跟拍引爆流量
                  </q-tooltip>
                </q-icon>
              </div>
              <q-select
                outlined
                dense
                emit-value
                map-options
                clearable
                usein
                use-input
                @filter="handleFilterFlashmobList"
                placeholder="选择一个"
                option-value="name"
                option-label="name"
                :popup-content-style="{height: '300px'}"
                v-model="data.batchConfig.collectionName"
                :options="data.flashmobList"
                autocomplete="true"
                :loading="data.flashmobLoading"
              >
                <template v-slot:option="scope">
                  <q-item v-bind="scope.itemProps">
                    <q-item-section avatar>
                      <q-avatar>
                        <q-img :src="scope.opt.cover"/>
                      </q-avatar>
                    </q-item-section>
                    <div class="flex row full-width justify-between items-center">
                      <div>{{ scope.opt.name }}</div>
                      <div class="text-tip">{{ formatChineseNumber(scope.opt.count) }}跟拍</div>
                    </div>
                  </q-item>
                </template>
              </q-select>
              <div class="row items-center">
                申请关联热点
                <q-icon name="help" class="cursor-pointer q-ml-xs text-tip">
                  <q-tooltip anchor="top middle" self="center middle">
                    你可以申请和一个热点做关联，如果视频确实和热点非常相关，将会进入抖音热点榜，若不相关则不会生效。
                  </q-tooltip>
                </q-icon>
              </div>
              <q-select
                outlined
                dense
                emit-value
                map-options
                clearable
                usein
                use-input
                @filter="handleFilterHotspotList"
                placeholder="选择一个"
                option-value="word"
                option-label="word"
                v-model="data.batchConfig.associatedHotspot"
                :options="data.hotspotList"
                :popup-content-style="{height: '300px'}"
                autocomplete="true"
                :loading="data.hotspotLoading"
              >
                <template v-slot:option="scope">
                  <q-item v-bind="scope.itemProps">
                    <q-item-section avatar>
                      <q-avatar>
                        <q-img :src="scope.opt.cover"/>
                      </q-avatar>
                    </q-item-section>
                    <div class="flex row full-width justify-between items-center">
                      <div>{{ scope.opt.word }}</div>
                      <div class="text-tip">{{ formatChineseNumber(scope.opt.hotValue) }}在看</div>
                    </div>
                  </q-item>
                </template>
              </q-select>
              <div class="row items-center">
                申请关联热点
                <q-icon name="help" class="cursor-pointer q-ml-xs text-tip">
                  <q-tooltip anchor="top middle" self="center middle">
                    你可以申请和一个热点做关联，如果视频确实和热点非常相关，将会进入抖音热点榜，若不相关则不会生效。
                  </q-tooltip>
                </q-icon>
              </div>
              <div class="row items-center q-mt-md">
                <q-toggle checked-icon="check" unchecked-icon="clear" dense :model-value="true">
                  <span>允许同步到其他平台</span>
                  <q-icon class="text-tip q-ml-xs cursor-pointer" name="help">
                    <q-tooltip anchor="top middle" self="center middle">
                      <div>- 已自动同步了抖音的视频描述与封面1</div>
                      <div>- 如你在抖音中删除内容，西瓜中相应内容仍会保留</div>
                      <div>- 在抖音允许被下载的内容，在西瓜也可以被下载到「本地相册」</div>
                    </q-tooltip>
                  </q-icon>
                </q-toggle>
              </div>
              <div class="row items-center">
                <q-toggle checked-icon="check" unchecked-icon="clear" dense v-model="data.batchConfig.allowedToSave">
                  <span>允许他人保存视频到相册</span>
                </q-toggle>
              </div>
              <div class="q-mt-md">可见范围</div>
              <div class="q-gutter-x-md">
                <q-radio dense v-model="data.batchConfig.whoCanWatch" :val="1" label="公开"/>
                <q-radio dense v-model="data.batchConfig.whoCanWatch" :val="2" label="好友可见"/>
                <q-radio dense v-model="data.batchConfig.whoCanWatch" :val="3" label="仅自己可见"/>
              </div>
              <div class="q-mt-md">发布时间</div>
              <div class="q-gutter-x-md q-mb-md flex row">
                <q-radio dense v-model="data.batchConfig.releaseType" :val="0" label="立即发布"/>
                <q-radio dense v-model="data.batchConfig.releaseType" :val="1" label="定时发布">
                  <q-icon class="q-ml-xs text-tip" name="help">
                    <q-tooltip anchor="top middle" self="center middle">
                      可选时间支持设置到2小时后及14天内
                    </q-tooltip>
                  </q-icon>
                </q-radio>
                <div v-if="data.batchConfig.releaseType === 1" style="max-width: 300px" class="inline-block">
                  <q-input outlined dense v-model="data.batchConfig.releaseDateString">
                    <template v-slot:prepend>
                      <q-icon name="event" class="cursor-pointer">
                        <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                          <q-date v-model="data.batchConfig.releaseDateString" :options="handleReleaseDateLimit" mask="YYYY-MM-DD HH:mm" bordered>
                            <div class="row items-center justify-end">
                              <q-btn v-close-popup label="关闭" color="primary" flat />
                            </div>
                          </q-date>
                        </q-popup-proxy>
                      </q-icon>
                    </template>

                    <template v-slot:append>
                      <q-icon name="access_time" class="cursor-pointer">
                        <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                          <q-time v-model="data.batchConfig.releaseDateString" :options="handleReleaseTimeLimit" mask="YYYY-MM-DD HH:mm" format24h bordered>
                            <div class="row items-center justify-end">
                              <q-btn v-close-popup label="关闭" color="primary" flat />
                            </div>
                          </q-time>
                        </q-popup-proxy>
                      </q-icon>
                    </template>
                  </q-input>
                </div>
              </div>
            </div>
          </q-scroll-area>
        </q-tab-panel>
        <q-tab-panel
          v-for="user in data.userList"
          :key="`config-dou-yin-content-${user.ID}`"
          :name="`${user.ID}-${user.name}`"
          class="q-pr-none q-py-none"
        >
          <q-scroll-area class="full-height full-width">
            <div>test</div>
          </q-scroll-area>
        </q-tab-panel>
      </q-tab-panels>
    </template>

  </q-splitter>
</template>

<style scoped lang="scss">
.activity {
  transition: 0.25s all;

  .action {
    display: none;
  }
}

.activity:hover, .activity-active {
  background: rgb(255, 232, 233) !important;

  .hot-value {
    display: none;
  }

  .action {
    display: flex;
  }
}
</style>
