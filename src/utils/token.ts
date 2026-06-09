/** Token 存取工具。支持 localStorage（记住我）和 sessionStorage（默认），两者互斥。 */
const TOKEN_KEY = 'admin_token'

/** getToken 从 localStorage 或 sessionStorage 中获取 Token */
export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY) || sessionStorage.getItem(TOKEN_KEY)
}

/** setToken 存储 Token。remember=true 存入 localStorage，否则存入 sessionStorage */
export function setToken(token: string, remember = false): void {
  const storage = remember ? localStorage : sessionStorage
  storage.setItem(TOKEN_KEY, token)
  if (remember) sessionStorage.removeItem(TOKEN_KEY)
  else localStorage.removeItem(TOKEN_KEY)
}

/** removeToken 同时清除 localStorage 和 sessionStorage 中的 Token */
export function removeToken(): void {
  localStorage.removeItem(TOKEN_KEY)
  sessionStorage.removeItem(TOKEN_KEY)
}
