<script setup lang="ts">

</script>

<template>
  <div class="full-height full-width flex row">
    <q-card flat bordered class="column q-pa-md" style="flex: 1">
      <div class="text-subtitle1 text-bold text-primary flex items-center">
        <q-icon name="dashboard" class="q-mr-sm"/>
        直播信息
      </div>
      <div class="row q-mt-sm full-width no-wrap">
        <div class="row justify-between" style="width: 70%;">
            <div
              style="background: rgba(255,135,84, .1); border: 1px solid rgba(255,135,84, .5); width: 23%;height: 90px"
              class="q-py-md q-px-md"
            >
              <div class="row items-center text-h5 overflow-hidden no-wrap">
                <q-icon name="live_tv" color="primary" class="q-mr-sm"/>
                <span>5</span>
              </div>
              <div class="text-tip q-mt-sm ellipsis">直播场次</div>
            </div>
            <div
              style="background: rgba(132,84,255, .1); border: 1px solid rgba(132,84,255,.5); width: 23%;height: 90px"
              class="q-py-md q-px-md"
            >
              <div class="row items-center text-h5 overflow-hidden no-wrap">
                <q-icon name="currency_yen" color="primary" class="q-mr-sm"/>
                <span>165206</span>
              </div>
              <div class="text-tip q-mt-sm ellipsis">礼物总收益</div>
            </div>
            <div
              style="background: rgba(195,53,253, .1); border: 1px solid rgba(195,53,253,.5); width: 23%;height: 90px"
              class="q-py-md q-px-md"
            >
              <div class="row items-center text-h5 overflow-hidden no-wrap">
                <q-icon name="group" color="primary" class="q-mr-sm"/>
                <span>5</span>
              </div>
              <div class="text-tip q-mt-sm ellipsis">粉丝团总人数</div>
            </div>
            <div
              style="background: rgba(255,48,54, .1); border: 1px solid rgba(255,48,54, .5); width: 23%;height: 90px"
              class="q-py-md q-px-md"
            >
              <div class="row items-center text-h5 overflow-hidden no-wrap">
                <q-icon name="favorite" color="primary" class="q-mr-sm"/>
                <span>5</span>
              </div>
              <div class="text-tip q-mt-sm ellipsis">会员总人数</div>
            </div>
          </div>
        <q-separator class="q-mx-md" vertical />
        <div style="flex: 1">
          <div class="text-subtitle1 text-bold text-accent">真爱粉</div>
          <q-list class="q-mt-md">
            <q-item v-for="item in 5" :key="item" class="q-px-none">
              <q-item-section top avatar>
                <q-avatar color="primary" text-color="white" icon="bluetooth" />
              </q-item-section>
              <q-item-section class="q-mt-xs">
                <q-item-label class="flex">
                  <div>肯定假粉丝</div>
                  <q-badge rounded class="q-ml-xs" label="15"/>
                </q-item-label>
                <q-item-label class="text-caption text-tip ellipsis" caption lines="2">送礼 3 次 2000 音浪 200 弹幕</q-item-label>
              </q-item-section>
              <q-item-section side top>
                <q-item-label caption>12</q-item-label>
                <q-icon name="star" color="yellow" />
              </q-item-section>
            </q-item>
          </q-list>
        </div>
      </div>
    </q-card>
    <q-card flat bordered class="column q-pa-md q-ml-lg" style="width: 20rem">
      <div class="text-subtitle1 text-bold text-primary flex items-center">
        <q-icon name="dashboard" class="q-mr-sm"/>
        直播记录
      </div>
      <div class="q-mt-sm text-subtitle1 text-tip text-caption">12 月 总计 20 场</div>
      <div class="q-pa-sm q-my-sm row secondary-container">
        <div class="flex justify-start items-end column text-right q-ml-sm" style="width: 80px;">
          <div class="text-subtitle1 text-bold">09:00</div>
          <div class="text-caption text-tip">2024-12-25</div>
          <div class="text-tip ellipsis">第一次直播</div>
        </div>
        <q-separator color="secondary" spaced size="5px" vertical />
        <div style="flex: 1">
          <div>直播间标题</div>
          <div class="text-tip text-caption ellipsis"><q-icon name="person" /> 200 人观看 100 人送礼 </div>
          <div class="text-tip text-caption ellipsis"><q-icon name="schedule" /> 时长约 80 分钟</div>
          <div class="text-tip text-caption ellipsis"><q-icon name="favorite" /> 送礼20次 3000 音浪</div>
        </div>
      </div>
    </q-card>
  </div>
</template>

<style scoped lang="scss">

</style>
