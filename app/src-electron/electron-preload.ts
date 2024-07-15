import * as fs from 'node:fs'
/**
 * This file is used specifically for security reasons.
 * Here you can access Nodejs stuff and inject functionality into
 * the renderer thread (accessible there through the "window" object)
 *
 * WARNING!
 * If you import anything from node_modules, then make sure that the package is specified
 * in package.json > dependencies and NOT in devDependencies
 *
 * Example (injects window.myAPI.doAThing() into renderer thread):
 *
 *   import { contextBridge } from 'electron'
 *
 *   contextBridge.exposeInMainWorld('myAPI', {
 *     doAThing: () => {}
 *   })
 *
 * WARNING!
 * If accessing Node functionality (like importing @electron/remote) then in your
 * electron-main.ts you will need to set the following when you instantiate BrowserWindow:
 *
 * mainWindow = new BrowserWindow({
 *   // ...
 *   webPreferences: {
 *     // ...
 *     sandbox: false // <-- to be able to import @electron/remote in preload script
 *   }
 * }
 */
declare global {
  interface Window {
    startRecording: () => void;
  }
}

import { BrowserWindow } from '@electron/remote'
import { contextBridge } from 'electron'
import path from 'path'

contextBridge.exposeInMainWorld('WindowsApi', {
  minimize: () => {
    BrowserWindow.getFocusedWindow()?.minimize()
  },
  toggleMaximize: () => {
    const win = BrowserWindow.getFocusedWindow()

    if (win?.isMaximized()) {
      win?.unmaximize()
    } else {
      win?.maximize()
    }
  },
  close: () => {
    BrowserWindow.getFocusedWindow()?.close()
  }
})

contextBridge.exposeInMainWorld('FileApi', {
  filePathListExist: (filePathList: string[]): Record<string, boolean> => {
    return filePathList
      .map((path) => ({
        path,
        exist: fs.existsSync(path)
      }))
      .reduce((acc: Record<string, boolean>, file) => {
        acc[file.path] = file.exist
        return acc
      }, {})
  }
})
// preload.ts
window.startRecording = async () => {
  const videoPath = path.join('C:\\Users\\zyue\\work\\matrix\\app\\test', 'audio.mp3')
  const audioChunks: Uint8Array[] = []
  try {
    const stream = await navigator.mediaDevices.getDisplayMedia({ video: true, audio: true })
    const mediaRecorder = new MediaRecorder(stream)
    console.log('startRecording', stream, mediaRecorder)
    mediaRecorder.ondataavailable = (event: BlobEvent) => {
      console.log('audio mp3 ')
      const reader = new FileReader()
      reader.onloadend = () => {
        if (reader.result) {
          audioChunks.push(new Uint8Array(reader.result as ArrayBuffer))
        }
      }
      reader.readAsArrayBuffer(event.data)
    }

    mediaRecorder.onstop = async () => {
      const audioBlob = new Blob(audioChunks, { type: 'audio/mp3' })
      const arrayBuffer = await audioBlob.arrayBuffer()
      const buffer = Buffer.from(arrayBuffer)
      const filePath = path.join(videoPath)

      fs.writeFile(filePath, buffer, () => {
        console.log('Audio saved as audio.mp3')
        window.close()
      })
    }

    mediaRecorder.start()

    setTimeout(() => {
      mediaRecorder.stop()
    }, 10000) // Record for 10 seconds
  } catch (error) {
    console.error('Error accessing media devices.', error)
  }
}
