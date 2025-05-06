import api from "./api"

export const paymentService = {
  // Process payment
  async processPayment(paymentData: {
    order_id: string
    amount: number
    payment_method: string
    card_details?: {
      card_number: string
      expiry_month: string
      expiry_year: string
      cvv: string
    }
  }) {
    const response = await api.post("/payments", paymentData)
    return response.data
  },

  // Get payment by order ID
  async getPaymentByOrderId(orderId: string) {
    const response = await api.get(`/payments/order/${orderId}`)
    return response.data
  },
}
