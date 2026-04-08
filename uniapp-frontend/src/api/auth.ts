import { post } from './request'

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  expires_at: string
}

export const login = (data: LoginRequest) => {
  return post<LoginResponse>('/auth/login', data)
}

export const register = (data: LoginRequest) => {
  return post<LoginResponse>('/auth/register', data)
}
