import api from "./api"

export const orderService = {
  // Create a new order
  async createOrder(orderData: {
    customer_id: string
    shipping_address: string
    items: {
      product_id: string
      quantity: number
      price: number
    }[]
  }) {
    console.log("Creating order with data:", JSON.stringify(orderData, null, 2));
    try {
      // Sử dụng endpoint đúng theo cấu hình Traefik
      const response = await api.post("/orders", orderData);
      console.log("Order created successfully:", response.data);
      return response.data;
    } catch (error: any) {
      console.error("Error creating order:", error);
      console.error("Error response:", error.response?.data);
      console.error("Error status:", error.response?.status);
      console.error("Error headers:", error.response?.headers);
      throw error;
    }
  },

  // Get all orders for the current user
  async getOrders() {
    const response = await api.get("/orders")
    return response.data
  },

  // Get order by ID
  async getOrderById(orderId: string) {
    const response = await api.get(`/orders/${orderId}`)
    return response.data
  },

  // Retry payment for an order
  async retryPayment(
    orderId: string,
    paymentData: {
      payment_method: string
      card_number?: string
      expiry_month?: string
      expiry_year?: string
      cvv?: string
    },
  ) {
    const response = await api.post(`/orders/${orderId}/retry-payment`, paymentData)
    return response.data
  },

  // Cancel an order
  async cancelOrder(orderId: string) {
    const response = await api.post(`/orders/${orderId}/cancel`)
    return response.data
  },

  // Process payment
  async processPayment(paymentData: {
    order_id: string
    amount: number
    payment_method: string
    card_number?: string
    expiry_month?: string
    expiry_year?: string
    cvv?: string
  }) {
    console.log("Processing payment with data:", JSON.stringify(paymentData, null, 2));
    try {
      const response = await api.post("/payments", paymentData);
      console.log("Payment processed successfully:", response.data);
      return response.data;
    } catch (error: any) {
      console.error("Error processing payment:", error);
      console.error("Error response:", error.response?.data);
      console.error("Error status:", error.response?.status);
      console.error("Error headers:", error.response?.headers);
      throw error;
    }
  },
}
