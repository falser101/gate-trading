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

export const getUser = () => get<User>('/user')
export const bindApiKey = (data: ApiKeyRequest) => post('/user/api-key', data)
export const updateApiKey = (data: ApiKeyRequest) => put('/user/api-key', data)
