export interface WindowsApi {
  minimize: () => void
  toggleMaximize: () => void
  close: () => void
}

export interface FileApi {
  filePathListExist: (filePathList: string[]) => Record<string, boolean>
}

declare global {
  interface Window {
    WindowsApi: WindowsApi,
    FileApi: FileApi
  }
}
