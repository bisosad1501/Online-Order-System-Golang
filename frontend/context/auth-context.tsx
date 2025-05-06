"use client"

import { createContext, useContext, useState, useEffect, type ReactNode } from "react"
import { authService } from "@/services/auth-service"
import { useToast } from "@/hooks/use-toast"
import { cartService } from "@/services/cart-service"

type User = {
  id: string
  email: string
  first_name: string
  last_name: string
  phone: string
  address: string
  created_at: string
  updated_at: string
}

type LoginCredentials = {
  id?: string
  email: string
  address?: string
}

type AuthContextType = {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (credentials: LoginCredentials) => Promise<void>
  logout: () => void
  updateUserData: (data: Partial<User>) => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const { toast } = useToast()

  useEffect(() => {
    const checkAuth = async () => {
      const token = localStorage.getItem("token")
      if (token) {
        try {
          const userData = await authService.getCurrentUser()
          setUser(userData)
        } catch (error) {
          console.error("Failed to get current user:", error)
          localStorage.removeItem("token")
        }
      }
      setIsLoading(false)
    }

    checkAuth()
  }, [])

  const login = async (credentials: LoginCredentials) => {
    try {
      const { token, user } = await authService.login(credentials)
      localStorage.setItem("token", token)
      localStorage.setItem("userId", user.id)
      setUser(user)

      // Tải giỏ hàng sau khi đăng nhập
      try {
        const cart = await cartService.getCart()
        if (cart && cart.id) {
          localStorage.setItem("cartId", cart.id)
          // Dispatch một event để thông báo cho CartContext biết giỏ hàng đã được tải
          window.dispatchEvent(new CustomEvent('cart-loaded', { detail: cart }))
        }
      } catch (cartError) {
        console.error("Failed to load cart after login:", cartError)
      }
    } catch (error) {
      console.error("Login failed:", error)
      throw error
    }
  }

  const logout = () => {
    localStorage.removeItem("token")
    localStorage.removeItem("userId")
    localStorage.removeItem("cartId")
    setUser(null)
    // Dispatch một event để thông báo cho CartContext biết người dùng đã đăng xuất
    window.dispatchEvent(new CustomEvent('user-logout'))
    toast({
      title: "Đã đăng xuất",
      description: "Bạn đã đăng xuất thành công",
    })
  }

  const updateUserData = (data: Partial<User>) => {
    if (user) {
      setUser({ ...user, ...data })
    }
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!user,
        isLoading,
        login,
        logout,
        updateUserData,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider")
  }
  return context
}
