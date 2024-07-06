<script setup lang="ts">
import { QEditor, QFile, useDialogPluginComponent, useQuasar } from 'quasar'
import {
  Creation,
  CreationApi,
  CreationType,
  CreationTypeMap
} from 'src/api/creation'
import { computed, onMounted, reactive, ref } from 'vue'
import _ from 'lodash'

interface Data {
  creation: Creation
  tags: string[]
}

const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } =
  useDialogPluginComponent()

const props = defineProps<{
  item: Creation | null
}>()

const defaultCreation: Creation = {
  ID: 0,
  UpdatedAt: '',
  CreatedAt: '',
  DeletedAt: null,
  // 需要设置的字段
  type: 0,
  title: '',
  description: '',
  paths: []
}
const $q = useQuasar()
const descriptionEditorRef = ref<QEditor | null>(null)
const fileRef = ref<QFile | null>(null)
const data = reactive<Data>({
  creation: defaultCreation,
  tags: [
    '#小狗快快',
    '#我不想再来',
    '#你会怎么看这世界👓',
    '#如果这个雨一直下💧',
    '#让我重生一次'
  ]
})

defineEmits([...useDialogPluginComponent.emits])

const dataHook = computed(() => {
  if (props.item) {
    return {
      title: '编辑',
      creation: props.item,
      saveAction: async (item: Creation) => {
        await CreationApi.add(item)
        onDialogOK()
      }
    }
  }
  return {
    title: '新增',
    creation: defaultCreation,
    saveAction: async (item: Creation) => {
      await CreationApi.add(item)
      onDialogOK()
    }
  }
})

onMounted(() => {
  data.creation = _.cloneDeep(dataHook.value.creation)
})

const handleSave = async () => {
  const tempDiv = document.createElement('div')
  tempDiv.innerHTML = data.creation.description
  const description = tempDiv.innerText
    .replaceAll('close', '')
    .replace(/\s+/g, ' ')
  await CreationApi.add({
    ...data.creation,
    description
  })
  onDialogOK()
}

const handleDescriptionAddTag = (tag: string) => {
  const edit = descriptionEditorRef.value
  if (!edit) {
    $q.notify({
      message: '操作失败，请重新打开窗口尝试'
    })
    return
  }
  edit.runCmd(
    'insertHTML',
    `&nbsp;<div class="tag-info row inline items-center" contenteditable="false">&nbsp;<span>${tag}</span>&nbsp;<i class="q-icon material-icons cursor-pointer" onclick="this.parentNode.parentNode.removeChild(this.parentNode)">close</i></div>&nbsp;`
  )
  edit.focus()
}

function extractHashtagStrings (input: string): string[] {
  const matches = input.match(/(?<=(&nbsp;|\s))#(\S+?)(?=(&nbsp;|\s))/g)
  if (matches) {
    return matches.map((match) => match.trim())
  }
  return []
}

const handleUpdateDescription = (value: string) => {
  console.log(value)
  const tagList = extractHashtagStrings(value)
  let description = value
  if (tagList.length === 0) {
    data.creation.description = description
    descriptionEditorRef.value?.focus()
    return
  }
  for (const tag of tagList) {
    description = description.replaceAll(` ${tag} `, '')
    description = description.replaceAll(`${tag}`, '')
    description = description.replaceAll(`&nbsp;${tag}&nbsp;`, '')
    description = description.replaceAll(` ${tag}&nbsp;`, '')
    description = description.replace(/\s+/g, ' ')
    description += `&nbsp;<div class="tag-info row inline items-center" contenteditable="false">&nbsp;<span>${tag}</span>&nbsp;<i class="q-icon material-icons cursor-pointer" onclick="this.parentNode.parentNode.removeChild(this.parentNode)">close</i></div>&nbsp;`
  }
  const edit = descriptionEditorRef.value
  edit?.runCmd('innerHTML', description)
  edit?.focus()
  const editorContent = edit?.$el.querySelector('.q-editor__content')
  if (editorContent) {
    const range = document.createRange()
    const sel = window.getSelection()
    range.setStart(editorContent, 0)
    range.collapse(true)

    sel?.removeAllRanges()
    sel?.addRange(range)
  }
  data.creation.description = description
}

const handlePaste = (event: ClipboardEvent) => {
  const target = event.target as HTMLElement
  if (target?.nodeName === 'INPUT') return
  let text: string | null = ''
  event.preventDefault()
  event.stopPropagation()
  if (event.clipboardData) {
    text = event.clipboardData.getData('text/plain')
    descriptionEditorRef.value?.runCmd('insertText', text)
  }
}

const handleAddFile = async (files: FileList) => {
  for (const index in files) {
    const path = files[index].path
    if (!data.creation.paths.includes(path)) {
      data.creation.paths.push(path)
    }
  }
}

const handleRemoveFile = (item: string) => {
  const index = data.creation.paths.indexOf(item)
  data.creation.paths.splice(index, 1)
}

const filePathExist = computed(() => {
  const booleans = window.FileApi.filePathListExist(
    data.creation.paths.map((item) => `${item}`)
  )
  console.log(booleans)
  return booleans
})
</script>

<template>
  <q-dialog class="creation-edit" ref="dialogRef" @hide="onDialogHide">
    <q-card class="q-dialog-plugin">
      <q-card-section>
        <div class="text-h6">{{ dataHook.title }}</div>
      </q-card-section>
      <q-card-section>
        <q-input outlined dense label="标题" v-model="data.creation.title"/>
        <form
          autocomplete="off"
        >
          <q-editor
            class="q-mt-md"
            ref="descriptionEditorRef"
            placeholder="请输入描述信息"
            :model-value="data.creation.description"
            @update:model-value="handleUpdateDescription"
            :toolbar="[]"
            content-class="editor-container"
            :content-style="{border: '1px #ddd solid', minHeight: '5rem'}"
            flat
            height="5px"
            @paste="handlePaste"
          >
          </q-editor>
        </form>
        <div class="q-mt-sm" style="width: 100%; height: 30px">
          <q-scroll-area class="full-height full-width no-wrap">
            <div class="row no-wrap">
              <q-chip size="10px" dense square v-for="tag in data.tags" clickable :key="`tag-${tag}`"
                      color="grey-8" text-color="white" @click="() => handleDescriptionAddTag(tag)">
                {{ tag }}
              </q-chip>
            </div>
          </q-scroll-area>
        </div>
        <div class="q-gutter-sm">
          <q-radio v-for="type in CreationType" :key="`select-type-${type.value}`" v-model="data.creation.type"
                   :val="type.value" :label="type.label" :color="type.color"/>
        </div>
        <q-file
          class="hidden"
          ref="fileRef"
          label="选择文件"
          use-chips
          multiple
          :model-value="[]"
          @update:model-value="handleAddFile"
          accept="*"
        />
        <div>
          <div class="q-mt-md vertical-middle">{{CreationTypeMap[data.creation.type]?.label}}列表
            <span class="text-caption text-tip">(点击文件名称进行删除)</span>
            <q-icon color="primary" name="add" class="cursor-pointer q-ml-sm" @click="() => fileRef?.pickFiles()"/>
          </div>
          <q-list separator class="q-mt-sm" dense style="min-height: 50px; border: 1px solid rgba(0,0,0,0.2)">
            <q-item clickable v-ripple
                    :class="`flex items-center ${filePathExist[item] ? 'text-green-9' : 'text-red-9'}`"
                    v-for="item in data.creation.paths"
                    :style="{'line-break': 'anywhere'}" :key="`file-${item}`"
                    @click="() => handleRemoveFile(item)">
              {{ item }}
            </q-item>
            <div class="cursor-pointer full-width full-height row flex justify-center items-center text-grey-5"
                 style="min-height: 50px"
                 v-if="data.creation.paths.length === 0" @click="() => fileRef?.pickFiles()">
              暂无数据, 点击添加
            </div>
          </q-list>
        </div>
      </q-card-section>
      <q-card-actions align="right">
        <q-btn outline color="red" label="取消" @click="onDialogCancel"/>
        <q-btn outline color="primary" label="保存" @click="handleSave"/>
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<style scoped lang="scss">
.creation-edit {
  &:deep(.q-editor__toolbar-group) {
    width: 100%;
  }
}

</style>
