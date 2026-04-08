import { post } from './request'

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  expires_at: string
}

// 用户登录
export const login = (data: LoginRequest) => {
  return post<LoginResponse>('/auth/login', data)
}

// 用户注册
export const register = (data: RegisterRequest) => {
  return post<LoginResponse>('/auth/register', data)
}
