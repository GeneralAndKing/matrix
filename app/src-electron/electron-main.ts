import { app, BrowserWindow } from 'electron'
import path from 'path'
import os from 'os'
import { fileURLToPath } from 'url'
import { enable, initialize } from '@electron/remote/main/index.js'

// https://quasar.dev/quasar-cli-vite/developing-electron-apps/frameless-electron-window#setting-frameless-window
initialize()
// needed in case process is undefined under Linux
const platform = process.platform || os.platform()

const currentDir = fileURLToPath(new URL('.', import.meta.url))

let mainWindow: BrowserWindow | undefined

const createWindow = () => {
  mainWindow = new BrowserWindow({
    icon: path.resolve(currentDir, 'icons/icon.png'), // tray icon
    minHeight: 600,
    minWidth: 1000,
    width: 1000,
    height: 600,
    useContentSize: true,
    frame: false,
    webPreferences: {
      sandbox: false,
      contextIsolation: true,
      webSecurity: false,
      // More info: https://v2.quasar.dev/quasar-cli-vite/developing-electron-apps/electron-preload-script
      preload: path.resolve(
        currentDir,
        path.join(process.env.QUASAR_ELECTRON_PRELOAD_FOLDER, 'electron-preload' + process.env.QUASAR_ELECTRON_PRELOAD_EXTENSION)
      )
    }
  })

  enable(mainWindow.webContents)
  if (process.env.DEV) {
    void mainWindow.loadURL(process.env.APP_URL)
  } else {
    void mainWindow.loadFile('index.html')
  }

  if (process.env.DEBUGGING) {
    // if on DEV or Production with debug enabled
    mainWindow.webContents.openDevTools()
  } else {
    // we're on production; no access to devtools pls
    mainWindow.webContents.on('devtools-opened', () => {
      // mainWindow?.webContents.closeDevTools()
    })
  }

  mainWindow.on('closed', () => {
    mainWindow = undefined
  })
  mainWindow.webContents.on('did-finish-load', () => {
    // mainWindow?.webContents.executeJavaScript("\n" +
    //   "const codeToInject = 'Object.defineProperty(navigator,\"language\", {" +
    //   "  get: function () { return \"zh-Hans-CN\"; }, " +
    //   "  set: function (a) {}" +
    //   " });';" +
    //   "const script = document.createElement('script');" +
    //   "script.appendChild(document.createTextNode(codeToInject));" +
    //   "(document.head || document.documentElement).appendChild(script);" +
    //   "script.parentNode?.removeChild(script);" +
    //   "console.log('hello');" +
    //   "alert('test')"
    // )

    // mainWindow?.webContents.debugger.sendCommand('Page.addScriptToEvaluateOnNewDocument', {
    //   'source': 'console.log("hello!")'
    // })
    // mainWindow?.webContents.debugger.sendCommand('Browser.getVersion')
    //   .then(result => {
    //     console.log('Browser version:', result);
    //   })
    //   .catch(error => {
    //     console.error('Failed to get browser version:', error);
    //   });
  })

  // session.defaultSession.webRequest.onBeforeSendHeaders((details, callback) => {
  //   details.requestHeaders['User-Agent'] = 'Mozilla/8.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
  //   callback({ requestHeaders: details.requestHeaders })
  // })
  // mainWindow?.webContents.debugger.on('message', (event, method, params) => {
  //   console.log(event, method, params)
  // })
  // mainWindow?.webContents.debugger.on('detach', (event, reason) => {
  //   console.log('Debugger detached due to: ', event, reason);
  // })
}

app.whenReady().then(async () => {
  createWindow()
})

app.on('window-all-closed', () => {
  if (platform !== 'darwin') {
    app.quit()
  }
})

app.on('activate', () => {
  if (mainWindow === undefined) {
    createWindow()
  }
})
