<script setup lang="ts">

import { useInformationStore } from 'stores/creation-store'
import { onMounted, reactive } from 'vue'
import { DouYinUser, DouYinUserApi } from 'src/api/user'
import { QSelect } from 'quasar'

interface Data {
  name: string,
  douYinUserList: DouYinUser[],
  douYinFilterUserList: DouYinUser[],
  douYinPublishList: number[]
}

const store = useInformationStore()
const data = reactive<Data>({
  name: '',
  douYinPublishList: [],
  douYinUserList: [],
  douYinFilterUserList: []
})

onMounted(() => {
  DouYinUserApi.getAll().then(res => {
    data.douYinUserList = res
  })
})

const handleFilterDouYinUser: QSelect['onFilter'] = (val, update) => {
  if (val === '') {
    update(() => {
      data.douYinFilterUserList = data.douYinUserList

      // here you have access to "ref" which
      // is the Vue reference of the QSelect
    })
    return
  }

  update(() => {
    const needle = val.toLowerCase()
    data.douYinFilterUserList = data.douYinUserList.filter(v => (v.name.toLowerCase().indexOf(needle) > -1))
  })
}

</script>

<template>
  <q-tabs
    v-model="store.styleData.tab"
    dense
    class="text-grey q-mt-sm"
    active-color="primary"
    indicator-color="primary"
    align="left"
  >
    <q-tab name="mails">
      <div class="row q-gutter-x-sm">
        <q-img alt="douyin" src="~assets/logo/douyin.svg" width="22px" height="22px"/>
        <div>抖音</div>
      </div>
    </q-tab>
  </q-tabs>
  <q-separator/>
  <q-tab-panels v-model="store.styleData.tab" animated style="flex: 1">
    <q-tab-panel name="mails" class="q-pa-none q-pt-sm">
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
            <q-tab style="justify-content: start" name="global" label="全局配置"/>
            <q-tab
              style="justify-content: start"
              v-for="user in data.douYinUserList"
              :key="`config-dou-yin-tab-${user.ID}`"
              :name="`${user.ID}-${user.name}`"
              :label="user.name"
              content-class="account-list-tab"
            >
              <q-badge v-if="data.douYinPublishList.includes(user.ID)" class="q-ml-sm" label="已配置" />
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
                    v-model="data.douYinPublishList"
                    multiple
                    outlined
                    placeholder="搜索或选择需要发布的账号信息"
                    emit-value
                    map-options
                    @filter="handleFilterDouYinUser"
                    dense
                    use-input
                    input-debounce="0"
                    :options="data.douYinFilterUserList"
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
                  <q-input class="q-mt-sm" outlined dense :model-value="data.name"
                           placeholder="好的作品标题可以获得更多浏览"/>
                  <div>作品简介</div>
                  <q-input outlined dense :model-value="data.name" type="textarea"/>
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
                        <div class="bg-grey-3 q-pa-sm" style="width: 140px">
                          <div class="ellipsis text-body2 text-bold">崩坏：星穹铁道2.3版本创作者激励活动</div>
                          <div class="text-tip text-caption">热度:1.6k</div>
                        </div>
                        <div class="bg-grey-3 q-pa-sm" style="width: 140px">
                          <div class="ellipsis text-body2 text-bold">崩坏：星穹铁道2.3版本创作者激励活动</div>
                          <div class="text-tip text-caption">热度:1.6k</div>
                        </div>
                        <div class="bg-grey-3 q-pa-sm" style="width: 140px">
                          <div class="ellipsis text-body2 text-bold">崩坏：星穹铁道2.3版本创作者激励活动</div>
                          <div class="text-tip text-caption">热度:1.6k</div>
                        </div>
                      </div>
                    </q-scroll-area>
                  </div>
                  <div>设置视频封面</div>
                  <q-file
                    :model-value="[]"
                    placeholder="选择一个视频"
                    outlined
                    dense
                    label="选择 png/jpeg/jpg 格式图片"
                  />
                  <div>设置标签</div>
                  <q-input :model-value="data.name" dense outlined placeholder="地理位置" />
                  <div class="row items-center q-mt-md">
                    <q-toggle checked-icon="check" unchecked-icon="clear" dense :model-value="true">
                      <span>同步到其他平台</span>
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
                    <q-toggle checked-icon="check" unchecked-icon="clear" dense :model-value="true">
                      <span>允许他人保存视频到相册</span>
                    </q-toggle>
                  </div>
                  <div class="q-mt-md">可见范围</div>
                  <div class="q-gutter-x-md">
                    <q-radio dense :model-value="true" :val="true" label="公开" />
                    <q-radio dense :model-value="false" val="true" label="好友可见" />
                    <q-radio dense :model-value="false" val="true" label="仅自己可见" />
                  </div>
                  <div class="q-mt-md">发布时间</div>
                  <div class="q-gutter-x-md">
                    <q-radio dense :model-value="true" :val="true" label="立即发布" />
                    <q-radio dense :model-value="false" val="true" label="定时发布">
                      <q-icon class="q-ml-xs text-tip" name="help">
                        <q-tooltip anchor="top middle" self="center middle">
                          可选时间支持设置到2小时后及14天内
                        </q-tooltip>
                      </q-icon>
                    </q-radio>
                    <q-btn color="primary" flat label="选择时间">
                      <q-menu self="top start">
                        <div class="q-gutter-md row items-start">
                          <q-date :model-value="data.name" mask="YYYY-MM-DD HH:mm" />
                          <q-time :model-value="data.name" mask="YYYY-MM-DD HH:mm" />
                        </div>
                      </q-menu>
                    </q-btn>
                  </div>
                </div>
              </q-scroll-area>
            </q-tab-panel>
            <q-tab-panel
              v-for="user in data.douYinUserList"
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
    </q-tab-panel>
  </q-tab-panels>
</template>

<style scoped lang="scss">
:deep(.account-list-tab) {
  flex-direction: row !important;
  flex-wrap: nowrap;
}
</style>
