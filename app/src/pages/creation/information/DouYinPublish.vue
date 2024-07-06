<script setup lang="ts">
import { date, QDate, QSelect, QTime, useQuasar } from 'quasar'
import { useInformationStore } from 'stores/creation-store'
import { DouYinUser, DouYinUserApi } from 'src/api/user'
import { onMounted, reactive, toRaw, watch } from 'vue'
import {
  DouYinActivityResponse,
  DouYinApi,
  DouYinChallengeSugResponse,
  DouYinFlashmobResponse,
  DouYinHotspotResponse
} from 'src/api/dou-yin'
import { formatChineseNumber, formatNumber } from 'src/utils'
import {
  DOU_YIN_GLOBAL_TAB,
  DouYinPublishForm,
  useDouYinPublishStore
} from 'stores/publish-store'
import _ from 'lodash'
import { useRoute } from 'vue-router'

interface Data {
  userList: DouYinUser[]
  filterUserList: DouYinUser[]
  publishList: number[]
  activityList: DouYinActivityResponse[]
  challengeList: DouYinChallengeSugResponse[]
  challengeLoading: boolean
  hotspotList: DouYinHotspotResponse[]
  hotspotLoading: boolean
  flashmobList: DouYinFlashmobResponse[]
  flashmobLoading: boolean
  currentConfig: DouYinPublishForm
}

const $q = useQuasar()
const route = useRoute()
const store = useInformationStore()
const publishStore = useDouYinPublishStore()
const data = reactive<Data>({
  publishList: [],
  userList: [],
  filterUserList: [],
  activityList: [],
  challengeList: [],
  challengeLoading: false,
  hotspotList: [],
  hotspotLoading: false,
  flashmobList: [],
  flashmobLoading: false,
  currentConfig: {
    id: '0',
    title: '',
    description: '',
    associatedHotspot: '',
    syncToToutiao: false,
    allowedToSave: true,
    whoCanWatch: 1,
    releaseTime: 0,
    releaseType: 0,
    releaseDateString: date.formatDate(
      date.addToDate(new Date(), { hour: 2 }),
      'YYYY-MM-DD HH:mm'
    )
  }
})

const handleInitDouYinData = () => {
  DouYinUserApi.getAll().then((res) => {
    data.userList = res
  })
  DouYinApi.getActivity().then((res) => {
    data.activityList = res.filter((item) => item.challenges.length > 0)
  })
  data.hotspotLoading = true
  DouYinApi.getHotspot()
    .then((res) => {
      data.hotspotList = res
    })
    .then(() => {
      data.hotspotLoading = false
    })
  DouYinApi.getChallenge().then((res) => {
    data.challengeList = res
  })
  data.flashmobLoading = true
  DouYinApi.getFlashmob()
    .then((res) => {
      data.flashmobList = res
    })
    .finally(() => {
      data.flashmobLoading = false
    })
}

onMounted(async () => {
  handleInitDouYinData()
  data.currentConfig = _.cloneDeep({
    ...publishStore.currentBatchConfig,
    title: store.data.information?.creation.title ?? '',
    description: store.data.information?.creation.description ?? ''
  })
})

const getConfigById = (id: string) => {
  const existConfig = store.data.information?.douyin.find(
    (item) => `${item.id}` === `${id}`
  )
  let config = publishStore.getDefaultConfig()
  if (existConfig) {
    config = publishStore.toConfigForm(existConfig)
  }
  config.title = store.data.information?.creation.title ?? ''
  config.description = store.data.information?.creation.description ?? ''
  console.debug(config)
  return config
}

watch(
  () => store.styleData.douYinTab,
  (value, oldValue) => {
    if (oldValue === DOU_YIN_GLOBAL_TAB) {
      // 原来为全局配置的情况
      data.currentConfig = _.cloneDeep(getConfigById(value))
    } else if (value === DOU_YIN_GLOBAL_TAB) {
      // 现在切换为全局配置的情况
      data.currentConfig = _.cloneDeep(publishStore.currentBatchConfig)
    }
  }
)

const handleInnerTabChange = (value: string) => {
  // 全局配置切换到特定配置
  if (store.styleData.douYinTab === DOU_YIN_GLOBAL_TAB) {
    if (_.isEqual(publishStore.currentBatchConfig, data.currentConfig)) {
      store.styleData.douYinTab = value
    } else {
      $q.dialog({
        title: '确认切换吗？',
        message: '当前配置尚未保存，切换后会丢失当前配置哦',
        cancel: true,
        persistent: true
      }).onOk(() => {
        store.styleData.douYinTab = value
      })
    }
    return
  }
  // 特定配置切换到全局配置或者特定配置，先找到原来是否有配置
  const existConfig = store.data.douyin.find(
    (item) => `${item.id}` === `${store.styleData.douYinTab}`
  )
  if (
    existConfig &&
    _.isEqual(
      toRaw(existConfig),
      publishStore.toRequestParam(store.styleData.douYinTab, data.currentConfig)
    )
  ) {
    // 如果与原来的相等
    store.styleData.douYinTab = value
  } else {
    $q.dialog({
      title: '确认切换吗？',
      message: '当前配置尚未保存，切换后会丢失当前配置哦',
      cancel: true
    }).onOk(() => {
      store.styleData.douYinTab = value
    })
  }
}

const handleFilterDouYinUser: QSelect['onFilter'] = (val, update) => {
  if (val === '') {
    update(() => {
      data.filterUserList = data.userList
    })
    return
  }

  update(() => {
    const needle = val.toLowerCase()
    data.filterUserList = data.userList.filter(
      (v) => v.name.toLowerCase().indexOf(needle) > -1
    )
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

const handleAddHotspot = (
  activity: DouYinActivityResponse,
  id: number | undefined
) => {
  console.log(activity)
  if (id) {
    return
  }
  const value = ` #${activity.challenges[0]} `
  if (data.currentConfig.description.includes(value)) {
    data.currentConfig.description.replaceAll(value, '')
  } else {
    data.currentConfig.description += value
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
  const dateString = data.currentConfig.releaseDateString.split(' ')[0]
  const select = new Date(`${dateString} ${hr}:${min ?? 0}`)
  return date.isBetweenDates(select, start, end)
}

const handleApplyGlobalConfig = () => {
  if (data.publishList.length === 0) {
    $q.notify({
      type: 'warning',
      message: '至少需要选择一个账号哦~'
    })
    return
  }
  // 需要覆盖的 id 列表
  const overrideIdList = store.data.douyin
    .map((item) => parseInt(item.id))
    .filter((id) => data.publishList.includes(id))
  const needNewIdList = _.difference(data.publishList, overrideIdList)
  if (overrideIdList.length > 0) {
    const nameList = data.userList
      .filter((item) => overrideIdList.includes(item.ID))
      .map((item) => item.name)
    $q.dialog({
      title: '存在已经拥有配置的账号，确认继续吗？',
      message: `账号 ${nameList.join(',')} 都已经配置过了，再次进行全局配置会覆盖到每个项的配置。`,
      cancel: true
    }).onOk(() => {
      // 覆盖数据
      for (const overrideId of overrideIdList) {
        const index = store.data.douyin.findIndex(
          (item) => parseInt(item.id) === overrideId
        )
        if (index >= 0) {
          store.data.douyin[index] = publishStore.toRequestParam(
            overrideId,
            data.currentConfig
          )
        }
      }
      needNewIdList.forEach((userId) =>
        store.data.douyin.push(
          publishStore.toRequestParam(userId, data.currentConfig)
        )
      )
      publishStore.currentBatchConfig = _.cloneDeep(data.currentConfig)
      $q.notify({ type: 'success', message: '操作成功' })
    })
  } else {
    needNewIdList.forEach((userId) =>
      store.data.douyin.push(
        publishStore.toRequestParam(userId, data.currentConfig)
      )
    )
    publishStore.currentBatchConfig = _.cloneDeep(data.currentConfig)
    $q.notify({ type: 'success', message: '操作成功' })
  }
}

const handleApplyCurrentConfig = () => {
  const index = store.data.douyin.findIndex(
    (item) => `${item.id}` === `${store.styleData.douYinTab}`
  )
  const config = publishStore.toRequestParam(
    store.styleData.douYinTab,
    data.currentConfig
  )
  if (index >= 0) {
    store.data.douyin[index] = config
  } else {
    store.data.douyin.push(config)
  }
  $q.notify({ type: 'success', message: '操作成功' })
}

const handleResetConfig = () => {
  data.currentConfig = _.cloneDeep({
    ...publishStore.currentBatchConfig,
    title: store.data.information?.creation.title ?? '',
    description: store.data.information?.creation.description ?? ''
  })
  $q.notify({ type: 'success', message: '操作成功' })
}

const handleRemoveConfig = () => {
  const index = store.data.douyin.findIndex(
    (item) => `${item.id}` === `${store.styleData.douYinTab}`
  )
  if (index >= 0) {
    _.pullAt(store.data.douyin, [index])
    $q.notify({ type: 'success', message: '操作成功' })
  }
}
</script>

<template>
  <q-splitter
    class="full-height"
    v-model="store.styleData.splitterModel"
  >
    <template v-slot:before>
      <q-tabs
        :model-value="store.styleData.douYinTab"
        @update:model-value="handleInnerTabChange"
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
          :name="user.ID"
          :label="user.name"
          content-class="account-list-tab"
        >
          <q-badge v-if="store.data.douyin.map(item => item.id).includes(`${user.ID}`)" class="q-ml-sm" label="已配置"/>
        </q-tab>
      </q-tabs>
    </template>
    <template v-slot:after>
      <q-scroll-area class="full-width q-pl-md" style="height: calc(100% - 50px)">
        <div class="flex column q-gutter-y-sm" style="max-width: 550px">
          <div v-if="store.styleData.douYinTab === DOU_YIN_GLOBAL_TAB">发布账号</div>
          <q-select
            v-if="store.styleData.douYinTab === DOU_YIN_GLOBAL_TAB"
            v-model="data.publishList"
            multiple
            outlined
            :placeholder="data.publishList.length > 0 ? '' : '搜索或选择需要设置的账号信息'"
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
          <q-input class="q-mt-sm" outlined dense v-model="data.currentConfig.title"
                   placeholder="好的作品标题可以获得更多浏览"/>
          <div>作品简介</div>
          <q-input outlined dense v-model="data.currentConfig.description" type="textarea"/>
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
            :placeholder="(data.currentConfig.challengeList?.length ?? 0) > 0 ? '' : '搜索或选择标签'"
            option-value="name"
            max-values="5"
            option-label="name"
            :popup-content-style="{height: '300px'}"
            v-model="data.currentConfig.challengeList"
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
                  :class="`cursor-pointer bg-grey-3 q-pa-sm activity ${data.currentConfig.description.includes(item.challenges[0])}`"
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
            v-model="data.currentConfig.videoCoverFile"
            outlined
            dense
            label="选择 png/jpeg/jpg 格式图片"
          />
          <div>设置标签</div>
          <q-input v-model="data.currentConfig.location" dense outlined placeholder="地理位置"/>
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
            :placeholder="(data.currentConfig.flashmob?.length ?? 0) > 0 ? '' : '输入进行搜索'"
            option-value="name"
            option-label="name"
            :popup-content-style="{height: '300px'}"
            v-model="data.currentConfig.flashmob"
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
            :placeholder="(data.currentConfig.associatedHotspot?.length ?? 0) > 0 ? '' : '输入进行搜索'"
            option-value="word"
            option-label="word"
            v-model="data.currentConfig.associatedHotspot"
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
            <q-toggle checked-icon="check" unchecked-icon="clear" dense v-model="data.currentConfig.allowedToSave">
              <span>允许他人保存视频到相册</span>
            </q-toggle>
          </div>
          <div class="q-mt-md">可见范围</div>
          <div class="q-gutter-x-md">
            <q-radio dense v-model="data.currentConfig.whoCanWatch" :val="1" label="公开"/>
            <q-radio dense v-model="data.currentConfig.whoCanWatch" :val="2" label="好友可见"/>
            <q-radio dense v-model="data.currentConfig.whoCanWatch" :val="3" label="仅自己可见"/>
          </div>
          <div class="q-mt-md">发布时间</div>
          <div class="q-gutter-x-md q-mb-md flex row">
            <q-radio dense v-model="data.currentConfig.releaseType" :val="0" label="立即发布"/>
            <q-radio dense v-model="data.currentConfig.releaseType" :val="1" label="定时发布">
              <q-icon class="q-ml-xs text-tip" name="help">
                <q-tooltip anchor="top middle" self="center middle">
                  可选时间支持设置到2小时后及14天内
                </q-tooltip>
              </q-icon>
            </q-radio>
            <div v-if="data.currentConfig.releaseType === 1" style="max-width: 300px" class="inline-block">
              <q-input outlined dense v-model="data.currentConfig.releaseDateString">
                <template v-slot:prepend>
                  <q-icon name="event" class="cursor-pointer">
                    <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                      <q-date v-model="data.currentConfig.releaseDateString" :options="handleReleaseDateLimit"
                              mask="YYYY-MM-DD HH:mm" bordered>
                        <div class="row items-center justify-end">
                          <q-btn v-close-popup label="关闭" color="primary" flat/>
                        </div>
                      </q-date>
                    </q-popup-proxy>
                  </q-icon>
                </template>

                <template v-slot:append>
                  <q-icon name="access_time" class="cursor-pointer">
                    <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                      <q-time v-model="data.currentConfig.releaseDateString" :options="handleReleaseTimeLimit"
                              mask="YYYY-MM-DD HH:mm" format24h bordered>
                        <div class="row items-center justify-end">
                          <q-btn v-close-popup label="关闭" color="primary" flat/>
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
      <div class="q-ml-md q-pt-sm full-width justify-end flex">
        <q-btn outline color="secondary" v-if="store.styleData.douYinTab === DOU_YIN_GLOBAL_TAB" @click="handleResetConfig" label="重置" />
        <q-btn outline color="red" v-else @click="handleRemoveConfig" label="取消当前账号配置" />
        <q-btn class="q-ml-md" @click="handleApplyGlobalConfig" v-if="store.styleData.douYinTab === DOU_YIN_GLOBAL_TAB" outline color="primary" label="应用配置到指定账号" />
        <q-btn class="q-ml-md" @click="handleApplyCurrentConfig" v-else outline color="primary" label="保存当前账号配置" />
        <q-btn class="q-ml-md" @click="() => store.handlePublish(route.params.id as string)" color="primary" label="发布" />
      </div>
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
