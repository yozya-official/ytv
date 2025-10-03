import { toast } from '@yuelioi/toast'
import dayjs from 'dayjs'

import 'dayjs/locale/zh-cn'

dayjs.locale('zh-cn')

export const formatDate = (date: string): string => {
  return dayjs(date).format('YYYY年M月D日')
}

/**
 * 通用错误处理函数
 */
export function handleApiError(error: unknown, action = '操作') {
  let message = ''

  if (typeof error === 'object' && error !== null) {
    const err = error as {
      response?: { data?: { error?: string; message?: string } }
      message?: string
    }

    if (err.response?.data?.error) {
      message = err.response.data.error
    } else if (err.response?.data?.message) {
      message = err.response.data.message
    } else if (err.message) {
      message = err.message
    } else {
      message = JSON.stringify(error)
    }
  } else if (typeof error === 'string') {
    message = error
  } else {
    message = String(error)
  }

  toast.error(`${action}失败: ${message}`)
}
