import { get, post, put } from './request'

export interface User {
  id: number
  email: string
  api_key_set: boolean
  created_at: string
}

export interface ApiKeyRequest {
  api_key: string
  api_secret: string
}

// 获取当前用户信息
export const getUser = () => {
  return get<User>('/user')
}

// 绑定 API Key
export const bindApiKey = (data: ApiKeyRequest) => {
  return post('/user/api-key', data)
}

// 更新 API Key
export const updateApiKey = (data: ApiKeyRequest) => {
  return put('/user/api-key', data)
}
