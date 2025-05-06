"use client"

import { createContext, useContext, useState, useCallback, type ReactNode } from "react"
import { notificationService } from "@/services/notification-service"

type Notification = {
  id: string
  customer_id: string
  type: string
  subject: string
  content: string
  recipient: string
  status: "UNREAD" | "READ"
  created_at: string
  updated_at: string
}

type NotificationContextType = {
  notifications: Notification[]
  loading: boolean
  fetchNotifications: (userId: string) => Promise<void>
  markAsRead: (notificationId: string) => Promise<void>
  unreadCount: number
}

const NotificationContext = createContext<NotificationContextType | undefined>(undefined)

export function NotificationProvider({ children }: { children: ReactNode }) {
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [loading, setLoading] = useState(false)

  const fetchNotifications = useCallback(async (userId: string) => {
    setLoading(true)
    try {
      const data = await notificationService.getNotificationsByUserId(userId)
      setNotifications(data)
    } catch (error) {
      console.error("Failed to fetch notifications:", error)
    } finally {
      setLoading(false)
    }
  }, [])

  const markAsRead = async (notificationId: string) => {
    setLoading(true)
    try {
      await notificationService.markAsRead(notificationId)

      // Update local state
      setNotifications(
        notifications.map((notification) =>
          notification.id === notificationId ? { ...notification, status: "READ" } : notification,
        ),
      )
    } catch (error) {
      console.error("Failed to mark notification as read:", error)
    } finally {
      setLoading(false)
    }
  }

  const unreadCount = notifications.filter((notification) => notification.status === "UNREAD").length

  return (
    <NotificationContext.Provider
      value={{
        notifications,
        loading,
        fetchNotifications,
        markAsRead,
        unreadCount,
      }}
    >
      {children}
    </NotificationContext.Provider>
  )
}

export function useNotifications() {
  const context = useContext(NotificationContext)
  if (context === undefined) {
    throw new Error("useNotifications must be used within a NotificationProvider")
  }
  return context
}
