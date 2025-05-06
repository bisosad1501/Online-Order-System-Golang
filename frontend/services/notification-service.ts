import api from "./api"

export const notificationService = {
  // Get notifications by user ID
  async getNotificationsByUserId(userId: string) {
    const response = await api.get(`/notifications/user/${userId}`)
    return response.data
  },

  // Mark notification as read
  async markAsRead(notificationId: string) {
    const response = await api.put(`/notifications/${notificationId}/status`, {
      status: "READ",
    })
    return response.data
  },
}
