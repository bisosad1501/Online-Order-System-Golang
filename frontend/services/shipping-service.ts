import api from "./api"

export const shippingService = {
  // Get shipment by order ID
  async getShipmentByOrderId(orderId: string) {
    const response = await api.get(`/shipments/order/${orderId}`)
    return response.data
  },
}
