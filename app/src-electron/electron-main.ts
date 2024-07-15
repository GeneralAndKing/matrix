import os from 'os'
import path from 'path'
import { fileURLToPath } from 'url'
import { enable, initialize } from '@electron/remote/main/index.js'
import { app, BrowserWindow, HandlerDetails, WindowOpenHandlerResponse } from 'electron'
import ffmpeg from 'fluent-ffmpeg'
import ffmpegPath from '@ffmpeg-installer/ffmpeg'
import { PassThrough } from 'stream'
// app.commandLine.appendSwitch('remote-debugging-port', '9222')
// app.commandLine.appendSwitch('remote-allow-origins', 'http://localhost:9222')
// https://quasar.dev/quasar-cli-vite/developing-electron-apps/frameless-electron-window#setting-frameless-window
initialize()
// needed in case process is undefined under Linux
const platform = process.platform || os.platform()

const currentDir = fileURLToPath(new URL('.', import.meta.url))

let mainWindow: BrowserWindow | undefined

let kernel = false
for (const arg of process.argv) {
  if (arg === '--kernel') {
    kernel = true
  }
}
console.log(`Current application ${kernel ? '' : 'not '}run Kernel mode.`)
const videoPath = path.join('C:\\Users\\zyue\\work\\matrix\\app\\test', 'output.mp4')
const frameStream = new PassThrough()
let isRecording: boolean = false
let captureTimer: ReturnType<typeof setInterval> | null = null
ffmpeg.setFfmpegPath(ffmpegPath.path)
const captureInterval = 10
const captureFrame = async () => {
  if (!isRecording || !mainWindow) return
  const image = await mainWindow.capturePage()
  if (!isRecording) return
  const buffer = image.toPNG()
  frameStream.write(buffer)
}
const startRecording = () => {
  isRecording = true
  createVideo()
  captureTimer = setInterval(captureFrame, captureInterval)
  // mainWindow?.hide()
  setTimeout(() => {
    mainWindow?.close()
  }, 30000)
}

const stopRecording = () => {
  if (!isRecording) return
  isRecording = false
  if (captureTimer) {
    clearInterval(captureTimer)
    captureTimer = null
  }
  frameStream.end()
}

const createVideo = () => {
  ffmpeg()
    .input(frameStream)
    .inputFormat('image2pipe')
    .inputOptions('-framerate 20')
    .outputOptions('-pix_fmt yuv420p')
    .on('start', (commandLine) => {
      console.log('start: Spawned Ffmpeg with command: ' + commandLine)
    })
    .on('error', (err, stdout, stderr) => {
      console.error('Error: ' + err.message)
      console.error('ffmpeg stderr: ' + stderr)
    })
    .on('end', () => {
      console.log('End: Finished processing')
    })
    .save(videoPath)
}

const createWindow = async () => {
  mainWindow = new BrowserWindow({
    icon: path.resolve(currentDir, 'icons/icon.png'), // tray icon
    minHeight: 800,
    minWidth: 1400,
    width: 1000,
    height: 600,
    useContentSize: true,
    frame: false,
    webPreferences: {
      sandbox: false,
      contextIsolation: true,
      webSecurity: false,
      // offscreen: true,
      // More info: https://v2.quasar.dev/quasar-cli-vite/developing-electron-apps/electron-preload-script
      preload: path.resolve(
        currentDir,
        path.join(
          process.env.QUASAR_ELECTRON_PRELOAD_FOLDER,
          (kernel ? 'stealth.min' : 'electron-preload') +
            process.env.QUASAR_ELECTRON_PRELOAD_EXTENSION
        )
      )
    }
  })

  enable(mainWindow.webContents)

  if (kernel) {
    mainWindow.maximize()
    void mainWindow.loadURL('about:blank')
  } else if (process.env.DEV) {
    // void mainWindow.loadURL(process.env.APP_URL)
    void mainWindow.loadURL('https://live.douyin.com/4467754')
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
  mainWindow.webContents.setWindowOpenHandler(
    (details: HandlerDetails): WindowOpenHandlerResponse => {
      if (details.url) {
        return {
          action: 'allow',
          overrideBrowserWindowOptions: {
            title: details.frameName,
            autoHideMenuBar: true,
            parent: mainWindow,
            resizable: false,
            maximizable: true,
            webPreferences: {
              preload: path.resolve(
                currentDir,
                path.join(
                  process.env.QUASAR_ELECTRON_PRELOAD_FOLDER,
                  (kernel ? 'stealth.min' : 'electron-preload') +
                    process.env.QUASAR_ELECTRON_PRELOAD_EXTENSION
                )
              ),
              sandbox: false,
              contextIsolation: true,
              webSecurity: false
            }
          }
        }
      }
      return { action: 'deny' }
    }
  )
  mainWindow.webContents.on('did-finish-load', () => {
    mainWindow?.webContents.setFrameRate(20)
    mainWindow?.webContents.executeJavaScript('startRecording()')
  })
  mainWindow.on('close', async (event) => {
    event.preventDefault() // 阻止窗口关闭
    stopRecording()
    mainWindow?.destroy() // 确保在保存完成后关闭窗口
  })
  startRecording()
}

app.whenReady().then(async () => {
  await createWindow()
})

app.on('window-all-closed', () => {
  if (platform !== 'darwin') {
    app.quit()
  }
})

app.on('activate', async () => {
  if (mainWindow === undefined) {
    await createWindow()
  }
})
