"use client"

import { useEffect, useState } from "react"
import { useAuth } from "@/context/auth-context"
import { useNotifications } from "@/context/notification-context"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Badge } from "@/components/ui/badge"
import { formatDate } from "@/lib/utils"
import { Bell, CheckCircle, Package, Truck } from "lucide-react"

export default function NotificationsPage() {
  const { user } = useAuth()
  const { notifications, loading, markAsRead, fetchNotifications } = useNotifications()
  const [activeTab, setActiveTab] = useState("all")

  useEffect(() => {
    if (user) {
      fetchNotifications(user.id)
    }
  }, [user, fetchNotifications])

  const filteredNotifications =
    activeTab === "all"
      ? notifications
      : notifications.filter((notification) => notification.type.includes(activeTab.toUpperCase()))

  const getNotificationIcon = (type: string) => {
    switch (type) {
      case "ORDER_CONFIRMED":
        return <CheckCircle className="h-5 w-5 text-green-500" />
      case "SHIPPING_UPDATED":
        return <Truck className="h-5 w-5 text-blue-500" />
      case "SHIPPING_COMPLETED":
        return <Package className="h-5 w-5 text-purple-500" />
      case "ORDER_CREATED":
        return <CheckCircle className="h-5 w-5 text-blue-500" />
      case "PAYMENT_PROCESSED":
        return <CheckCircle className="h-5 w-5 text-green-500" />
      case "INVENTORY_CHECKED":
        return <CheckCircle className="h-5 w-5 text-orange-500" />
      default:
        return <Bell className="h-5 w-5 text-gray-500" />
    }
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="mb-6 text-3xl font-bold">Notifications</h1>
        <div className="flex justify-center">
          <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-t-2 border-blue-600"></div>
        </div>
      </div>
    )
  }

  if (!notifications.length) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="mb-6 text-3xl font-bold">Notifications</h1>
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed p-12 text-center">
          <Bell className="mb-4 h-16 w-16 text-gray-400" />
          <h2 className="mb-2 text-xl font-semibold">No notifications</h2>
          <p className="mb-6 text-gray-500">You don&apos;t have any notifications yet.</p>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="mb-6 text-3xl font-bold">Notifications</h1>

      <Tabs defaultValue="all" onValueChange={setActiveTab}>
        <TabsList className="mb-6">
          <TabsTrigger value="all">All</TabsTrigger>
          <TabsTrigger value="order">Orders</TabsTrigger>
          <TabsTrigger value="shipping">Shipping</TabsTrigger>
        </TabsList>

        <TabsContent value={activeTab} className="space-y-4">
          {filteredNotifications.map((notification) => (
            <Card
              key={notification.id}
              className={notification.status === "UNREAD" ? "border-blue-200 bg-blue-50" : ""}
            >
              <CardHeader className="flex flex-row items-start justify-between pb-2">
                <div className="flex items-start gap-4">
                  <div className="mt-1">{getNotificationIcon(notification.type)}</div>
                  <div>
                    <CardTitle className="text-lg">{notification.subject}</CardTitle>
                    <CardDescription>{formatDate(notification.created_at)}</CardDescription>
                  </div>
                </div>
                {notification.status === "UNREAD" && <Badge className="bg-blue-100 text-blue-800">New</Badge>}
              </CardHeader>
              <CardContent>
                <p className="mb-4">{notification.content}</p>
                {notification.status === "UNREAD" && (
                  <Button variant="outline" size="sm" onClick={() => markAsRead(notification.id)}>
                    Mark as read
                  </Button>
                )}
              </CardContent>
            </Card>
          ))}
        </TabsContent>
      </Tabs>
    </div>
  )
}
