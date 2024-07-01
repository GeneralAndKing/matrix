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
import { BrowserWindow } from '@electron/remote'
import { contextBridge } from 'electron'
import * as fs from 'node:fs'

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
    return filePathList.map(path => ({
      path,
      exist: fs.existsSync(path)
    })).reduce((acc: Record<string, boolean>, file) => {
      acc[file.path] = file.exist
      return acc
    }, {})
  }
})
let get = function (target, key) {
  console.log(target,key)
  let id = null
  switch (key) {
    case 'height':
      id = 3123
      break
    case 'width':
      id = 3123
      break
    case 'colorDepth':
      id = 3213
      break
    case 'pixelDepth':
      id = 3123
      break
  }
  if (id != null) {
    console.log(id)
  }
  let res = target[key]
  if (typeof res === 'function') {
    return res.bind(target)
  } else {
    return res
  }
}

Object.defineProperty(window, 'screen', {
  value: new Proxy(window.screen, { get })
})
console.log(window)
// const generateFingerprint = () => {
//   console.log('test')
//   Object.defineProperty(navigator, 'userAgent', {
//     value: 'Mozilla/12.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
//     writable: false
//   })
//   // 修改 navigator.platform
//   console.log('test2')
// }
//
// generateFingerprint()
//
// console.log('BrowserWindow.getFocusedWindow()', BrowserWindow.getFocusedWindow())
// document.addEventListener('DOMContentLoaded', () => {
//   console.log('test3')
//   const codeToInject = 'var inject = function () {\
//   console.log("inject");\
//   const toBlob = HTMLCanvasElement.prototype.toBlob;\
//   const toDataURL = HTMLCanvasElement.prototype.toDataURL;\
//   const getImageData = CanvasRenderingContext2D.prototype.getImageData;\
//   var noisify = function (canvas, context) {\
//     if (context) {\
//       const shift = {\
//         \'r\': Math.floor(Math.random() * 10) - 5,\
//         \'g\': Math.floor(Math.random() * 10) - 5,\
//         \'b\': Math.floor(Math.random() * 10) - 5,\
//         \'a\': Math.floor(Math.random() * 10) - 5\
//       };\
//       const width = canvas.width;\
//       const height = canvas.height;\
//       if (width && height) {\
//         const imageData = getImageData.apply(context, [0, 0, width, height]);\
//         for (let i = 0; i < height; i++) {\
//           for (let j = 0; j < width; j++) {\
//             const n = ((i * (width * 4)) + (j * 4));\
//             imageData.data[n + 0] = imageData.data[n + 0] + shift.r;\
//             imageData.data[n + 1] = imageData.data[n + 1] + shift.g;\
//             imageData.data[n + 2] = imageData.data[n + 2] + shift.b;\
//             imageData.data[n + 3] = imageData.data[n + 3] + shift.a;\
//           }\
//         }\
//         window.top.postMessage("canvas-fingerprint-defender-alert", \'*\');\
//         context.putImageData(imageData, 0, 0);\
//       }\
//     }\
//   };\
//   Object.defineProperty(HTMLCanvasElement.prototype, "toBlob", {\
//     "value": function () {\
//       noisify(this, this.getContext("2d"));\
//       return toBlob.apply(this, arguments);\
//     }\
//   });\
//   Object.defineProperty(HTMLCanvasElement.prototype, "toDataURL", {\
//     "value": function () {\
//       noisify(this, this.getContext("2d"));\
//       return toDataURL.apply(this, arguments);\
//     }\
//   });\
//   Object.defineProperty(CanvasRenderingContext2D.prototype, "getImageData", {\
//     "value": function () {\
//       noisify(this.canvas, this);\
//       return getImageData.apply(this, arguments);\
//     }\
//   });\
//   document.documentElement.dataset.cbscriptallow = true;\
// };\
// inject();'
//   const script = document.createElement('script')
//   script.appendChild(document.createTextNode(codeToInject));
//   (document.head || document.documentElement).appendChild(script)
//   console.log('test3')
//   const script1 = document.createElement('script');
//   script1.textContent = "console.log('Script added to head!');";
//   document.head.appendChild(script);
// })
// console.log('document.head', document.head)
