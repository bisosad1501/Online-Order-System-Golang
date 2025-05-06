"use client"

import { createContext, useContext, useState, useCallback, useEffect, type ReactNode } from "react"
import { cartService } from "@/services/cart-service"
import { useAuth } from "@/context/auth-context"

type CartItem = {
  id: string
  cart_id: string
  product_id: string
  quantity: number
  price: number
}

type Cart = {
  id: string
  customer_id: string
  items: CartItem[]
  created_at: string
  updated_at: string
}

type NewCartItem = {
  product_id: string
  quantity: number
  price: number
}

type CartContextType = {
  cart: Cart | null
  loading: boolean
  fetchCart: (userId: string) => Promise<void>
  addToCart: (item: NewCartItem) => Promise<void>
  updateCartItem: (itemId: string, quantity: number) => Promise<void>
  removeFromCart: (itemId: string) => Promise<void>
  clearCart: () => Promise<void>
}

const CartContext = createContext<CartContextType | undefined>(undefined)

export function CartProvider({ children }: { children: ReactNode }) {
  const [cart, setCart] = useState<Cart | null>(null)
  const [loading, setLoading] = useState(false)
  const { user, isAuthenticated } = useAuth()

  const fetchCart = useCallback(async (userId: string) => {
    setLoading(true)
    try {
      console.log("Fetching cart for user:", userId)
      const cartData = await cartService.getCartByUserId(userId)

      // Nếu không tìm thấy giỏ hàng (cartData là null), tạo giỏ hàng mới
      if (!cartData) {
        console.log("No cart found, creating new cart for user:", userId)
        const newCart = await cartService.createCart(userId)
        console.log("New cart created:", newCart)

        // Đảm bảo cart mới có thuộc tính items
        if (newCart && !newCart.items) {
          newCart.items = []
        }
        setCart(newCart)
        // Lưu cartId vào localStorage
        if (newCart && newCart.id) {
          localStorage.setItem("cartId", newCart.id)
        }
        return
      }

      console.log("Cart data received:", cartData)
      // Đảm bảo cart có thuộc tính items
      if (cartData && !cartData.items) {
        cartData.items = []
      }
      setCart(cartData)
      // Lưu cartId vào localStorage
      if (cartData && cartData.id) {
        localStorage.setItem("cartId", cartData.id)
      }
    } catch (error) {
      console.error("Failed to fetch or create cart:", error)
      // Nếu có lỗi, đặt cart thành null
      setCart(null)
    } finally {
      setLoading(false)
    }
  }, [])

  const addToCart = async (item: NewCartItem) => {
    if (!cart) {
      console.error("Cannot add to cart: Cart is null")
      throw new Error("Giỏ hàng không tồn tại. Vui lòng tải lại trang.")
    }

    // Đảm bảo cart có thuộc tính items
    if (!cart.items) {
      cart.items = []
    }

    setLoading(true)
    try {
      console.log("Adding item to cart:", item, "Cart ID:", cart.id)
      // Gọi API để thêm sản phẩm vào giỏ hàng
      // API sẽ trả về giỏ hàng đã cập nhật
      const updatedCart = await cartService.addCartItem(cart.id, item)
      console.log("Updated cart received:", updatedCart)

      // Đảm bảo updatedCart có thuộc tính items
      if (updatedCart && !updatedCart.items) {
        updatedCart.items = []
      }

      // Cập nhật state với giỏ hàng đã cập nhật từ server
      if (updatedCart) {
        setCart(updatedCart)
        // Lưu cartId vào localStorage
        if (updatedCart.id) {
          localStorage.setItem("cartId", updatedCart.id)
        }
      }
    } catch (error) {
      console.error("Failed to add item to cart:", error)
      throw error
    } finally {
      setLoading(false)
    }
  }

  const updateCartItem = async (itemId: string, quantity: number) => {
    if (!cart) return

    setLoading(true)
    try {
      // Gọi API để cập nhật số lượng sản phẩm trong giỏ hàng
      // API sẽ trả về giỏ hàng đã cập nhật
      const updatedCart = await cartService.updateCartItem(cart.id, itemId, quantity)

      // Cập nhật state với giỏ hàng đã cập nhật từ server
      if (updatedCart) {
        setCart(updatedCart)
      }
    } catch (error) {
      console.error("Failed to update cart item:", error)
      throw error
    } finally {
      setLoading(false)
    }
  }

  const removeFromCart = async (itemId: string) => {
    if (!cart) return

    setLoading(true)
    try {
      // Gọi API để xóa sản phẩm khỏi giỏ hàng
      // API sẽ trả về giỏ hàng đã cập nhật
      const updatedCart = await cartService.removeCartItem(cart.id, itemId)

      // Cập nhật state với giỏ hàng đã cập nhật từ server
      if (updatedCart) {
        setCart(updatedCart)
      }
    } catch (error) {
      console.error("Failed to remove item from cart:", error)
      throw error
    } finally {
      setLoading(false)
    }
  }

  const clearCart = async () => {
    if (!cart) return

    setLoading(true)
    try {
      // Gọi API để xóa giỏ hàng
      await cartService.deleteCart(cart.id)

      // Tạo giỏ hàng mới sau khi xóa
      const newCart = await cartService.createCart(cart.customer_id)

      // Cập nhật state với giỏ hàng mới
      if (newCart) {
        setCart(newCart)
      }
    } catch (error) {
      console.error("Failed to clear cart:", error)
      // Fallback: Cập nhật state local nếu API gọi thất bại
      setCart({ ...cart, items: [] })
    } finally {
      setLoading(false)
    }
  }

  // Tự động tải giỏ hàng khi người dùng đăng nhập hoặc khi component được mount
  useEffect(() => {
    const loadCart = async () => {
      // Nếu người dùng đã đăng nhập, lấy giỏ hàng theo user ID
      if (isAuthenticated && user) {
        fetchCart(user.id)
      } else {
        // Nếu chưa đăng nhập, thử lấy giỏ hàng từ localStorage
        const cartId = localStorage.getItem("cartId")
        if (cartId) {
          try {
            const cartData = await cartService.getCart()
            if (cartData) {
              if (!cartData.items) {
                cartData.items = []
              }
              setCart(cartData)
            } else {
              setCart(null)
            }
          } catch (error) {
            console.error("Failed to load cart from localStorage:", error)
            setCart(null)
          }
        } else {
          setCart(null)
        }
      }
    }

    loadCart()
  }, [isAuthenticated, user, fetchCart])

  // Lắng nghe sự kiện cart-loaded từ AuthContext
  useEffect(() => {
    const handleCartLoaded = (event: CustomEvent) => {
      setCart(event.detail)
    }

    const handleUserLogout = () => {
      setCart(null)
    }

    window.addEventListener('cart-loaded', handleCartLoaded as EventListener)
    window.addEventListener('user-logout', handleUserLogout)

    return () => {
      window.removeEventListener('cart-loaded', handleCartLoaded as EventListener)
      window.removeEventListener('user-logout', handleUserLogout)
    }
  }, [])

  return (
    <CartContext.Provider
      value={{
        cart,
        loading,
        fetchCart,
        addToCart,
        updateCartItem,
        removeFromCart,
        clearCart,
      }}
    >
      {children}
    </CartContext.Provider>
  )
}

export function useCart() {
  const context = useContext(CartContext)
  if (context === undefined) {
    throw new Error("useCart must be used within a CartProvider")
  }
  return context
}
