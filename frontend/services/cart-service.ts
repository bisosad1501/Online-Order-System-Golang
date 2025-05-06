import api from "./api"

export const cartService = {
  // Get current user's cart
  async getCart() {
    try {
      const cartId = localStorage.getItem("cartId")
      if (cartId) {
        const response = await api.get(`/carts/${cartId}`)
        return response.data
      } else {
        // Nếu không có cartId trong localStorage, thử lấy từ user ID
        const userId = localStorage.getItem("userId")
        if (userId) {
          return this.getCartByUserId(userId)
        }
        return null
      }
    } catch (error: any) {
      // Nếu lỗi 404 (không tìm thấy giỏ hàng), trả về null
      if (error.response && error.response.status === 404) {
        return null
      }
      throw error
    }
  },
  // Get cart by user ID
  async getCartByUserId(userId: string) {
    try {
      const response = await api.get(`/carts/user/${userId}`)
      return response.data
    } catch (error: any) {
      // Nếu lỗi 404 (không tìm thấy giỏ hàng), trả về null thay vì ném lỗi
      if (error.response && error.response.status === 404) {
        return null
      }
      throw error
    }
  },

  // Create a new cart
  async createCart(customerId: string) {
    const response = await api.post("/carts", { customer_id: customerId })
    return response.data
  },

  // Add item to cart
  async addCartItem(
    cartId: string,
    item: {
      product_id: string
      quantity: number
      price: number
    },
  ) {
    const response = await api.post(`/carts/${cartId}/items`, item)
    return response.data
  },

  // Update cart item quantity
  async updateCartItem(cartId: string, itemId: string, quantity: number) {
    const response = await api.put(`/carts/${cartId}/items/${itemId}`, {
      quantity,
    })
    return response.data
  },

  // Remove item from cart
  async removeCartItem(cartId: string, itemId: string) {
    const response = await api.delete(`/carts/${cartId}/items/${itemId}`)
    return response.data
  },

  // Delete cart
  async deleteCart(cartId: string) {
    const response = await api.delete(`/carts/${cartId}`)
    return response.data
  },
}
