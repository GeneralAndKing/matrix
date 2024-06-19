// eslint-disable-next-line @typescript-eslint/no-unused-vars
import axios from 'axios'

declare module 'axios' {
  export interface AxiosResponse<T> extends Promise<T> {}
}
