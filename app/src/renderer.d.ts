export interface WindowsApi {
  minimize: () => void
  toggleMaximize: () => void
  close: () => void
}

declare global {
  interface Window {
    WindowsApi: WindowsApi
  }
}
