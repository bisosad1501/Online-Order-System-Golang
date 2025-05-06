"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { formatDate } from "@/lib/utils"
import { MapPin, Package, Truck, CheckCircle, AlertCircle } from "lucide-react"
import { useToast } from "@/hooks/use-toast"

type OrderTrackingProps = {
  orderId: string
  shipment: {
    id: string
    order_id: string
    carrier: string
    tracking_number: string
    status: "PENDING" | "IN_TRANSIT" | "DELIVERED"
    shipping_address: string
    customer_id: string
    estimated_delivery: string
    created_at: string
    updated_at: string
  } | null
}

type TrackingEvent = {
  id: string
  status: string
  location: string
  timestamp: string
  description: string
}

export default function OrderTracking({ orderId, shipment }: OrderTrackingProps) {
  const { toast } = useToast()
  const [trackingEvents, setTrackingEvents] = useState<TrackingEvent[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (shipment) {
      // In a real app, this would fetch tracking events from the shipping service
      // For now, we'll generate mock data
      setTrackingEvents(generateMockTrackingEvents(shipment))
    }
  }, [shipment])

  const handleTrackPackage = () => {
    setLoading(true)

    // Simulate API call
    setTimeout(() => {
      setLoading(false)
      toast({
        title: "Tracking information updated",
        description: "The latest tracking information has been loaded",
      })
    }, 1500)
  }

  if (!shipment) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Tracking Information</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col items-center justify-center py-6 text-center">
            <Package className="mb-2 h-12 w-12 text-gray-400" />
            <h3 className="mb-1 text-lg font-medium">No tracking information available</h3>
            <p className="text-sm text-muted-foreground">
              Tracking information will be available once your order has been shipped.
            </p>
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Thông tin vận chuyển</CardTitle>
        <Badge
          variant="outline"
          className={
            shipment.status === "DELIVERED"
              ? "bg-green-50 text-green-700"
              : shipment.status === "IN_TRANSIT"
              ? "bg-blue-50 text-blue-700"
              : "bg-yellow-50 text-yellow-700"
          }
        >
          {shipment.status === "DELIVERED"
            ? "Đã giao hàng"
            : shipment.status === "IN_TRANSIT"
            ? "Đang giao hàng"
            : "Đã lên lịch"}
        </Badge>
      </CardHeader>
      <CardContent>
        <div className="mb-4 space-y-2">
          <div className="flex justify-between">
            <span className="text-sm text-muted-foreground">Đơn vị vận chuyển</span>
            <span className="font-medium">{shipment.carrier}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-sm text-muted-foreground">Mã vận đơn</span>
            <span className="font-medium">{shipment.tracking_number}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-sm text-muted-foreground">Ngày giao hàng dự kiến</span>
            <span className="font-medium">
              {shipment.status === "DELIVERED"
                ? "Đã giao hàng"
                : shipment.estimated_delivery
                  ? new Date(shipment.estimated_delivery).toLocaleDateString()
                  : new Date(new Date(shipment.created_at).getTime() + 5 * 24 * 60 * 60 * 1000).toLocaleDateString()}
            </span>
          </div>
        </div>

        <Separator className="my-4" />

        <div className="relative mt-6 pl-6">
          <div className="absolute left-0 top-0 bottom-0 w-px bg-gray-200 dark:bg-gray-700"></div>

          {trackingEvents.map((event, index) => (
            <div key={event.id} className="mb-6 relative">
              <div
                className={`absolute left-0 top-0 -ml-3 h-6 w-6 rounded-full flex items-center justify-center ${
                  index === 0
                    ? "bg-green-100 text-green-600 dark:bg-green-900 dark:text-green-300"
                    : "bg-gray-100 text-gray-500 dark:bg-gray-800 dark:text-gray-400"
                }`}
              >
                {index === 0 && event.status === "Delivered" ? (
                  <CheckCircle className="h-4 w-4" />
                ) : event.status === "In Transit" ? (
                  <Truck className="h-4 w-4" />
                ) : event.status === "Scheduled" ? (
                  <Package className="h-4 w-4" />
                ) : (
                  <MapPin className="h-4 w-4" />
                )}
              </div>

              <div className="pl-6">
                <h4 className="font-medium">{event.status}</h4>
                <p className="text-sm text-muted-foreground">{event.description}</p>
                <div className="mt-1 flex items-center text-xs text-muted-foreground">
                  <span>{formatDate(event.timestamp)}</span>
                  <span className="mx-2">•</span>
                  <span>{event.location}</span>
                </div>
              </div>
            </div>
          ))}
        </div>

        <Button onClick={handleTrackPackage} disabled={loading} className="mt-2 w-full">
          {loading ? "Updating..." : "Update Tracking Information"}
        </Button>
      </CardContent>
    </Card>
  )
}

// Helper function to generate mock tracking events
function generateMockTrackingEvents(shipment: OrderTrackingProps["shipment"]): TrackingEvent[] {
  if (!shipment) return []

  const events: TrackingEvent[] = []
  const createdDate = new Date(shipment.created_at)

  // Add events based on shipment status
  if (shipment.status === "DELIVERED") {
    // Delivered - Bước 13
    events.push({
      id: "1",
      status: "Đã giao hàng",
      location: "Địa chỉ nhận hàng",
      timestamp: shipment.updated_at,
      description: "Đơn hàng đã được giao thành công đến người nhận",
    })

    // Out for delivery - Bước 12
    events.push({
      id: "2",
      status: "Đang giao hàng",
      location: "Trung tâm phân phối địa phương",
      timestamp: new Date(new Date(shipment.updated_at).getTime() - 8 * 60 * 60 * 1000).toISOString(),
      description: "Đơn hàng đang được giao đến địa chỉ người nhận",
    })

    // In transit - Bước 11
    events.push({
      id: "3",
      status: "Đang vận chuyển",
      location: "Trung tâm phân phối",
      timestamp: new Date(new Date(shipment.updated_at).getTime() - 24 * 60 * 60 * 1000).toISOString(),
      description: "Đơn hàng đang được vận chuyển đến trung tâm phân phối địa phương",
    })

    // Shipping created - Bước 10
    events.push({
      id: "4",
      status: "Đã tạo vận đơn",
      location: "Kho hàng",
      timestamp: new Date(new Date(shipment.updated_at).getTime() - 2 * 24 * 60 * 60 * 1000).toISOString(),
      description: "Đơn hàng đã được tạo vận đơn và chuẩn bị giao cho đơn vị vận chuyển",
    })

    // Order confirmed - Bước 9
    events.push({
      id: "5",
      status: "Đã xác nhận",
      location: "Hệ thống",
      timestamp: shipment.created_at,
      description: "Đơn hàng đã được xác nhận và lên lịch giao hàng",
    })
  } else if (shipment.status === "IN_TRANSIT") {
    // In transit - Bước 11-12
    events.push({
      id: "1",
      status: "Đang vận chuyển",
      location: "Trung tâm phân phối",
      timestamp: shipment.updated_at,
      description: "Đơn hàng đang được vận chuyển đến địa chỉ người nhận",
    })

    // Shipping created - Bước 10
    events.push({
      id: "2",
      status: "Đã tạo vận đơn",
      location: "Kho hàng",
      timestamp: new Date(new Date(shipment.updated_at).getTime() - 1 * 24 * 60 * 60 * 1000).toISOString(),
      description: "Đơn hàng đã được tạo vận đơn và chuẩn bị giao cho đơn vị vận chuyển",
    })

    // Order confirmed - Bước 9
    events.push({
      id: "3",
      status: "Đã xác nhận",
      location: "Hệ thống",
      timestamp: shipment.created_at,
      description: "Đơn hàng đã được xác nhận và lên lịch giao hàng",
    })
  } else {
    // Pending - Bước 9
    events.push({
      id: "1",
      status: "Đã lên lịch giao hàng",
      location: "Hệ thống",
      timestamp: shipment.created_at,
      description: "Đơn hàng đã được lên lịch giao hàng và đang chờ xử lý",
    })
  }

  return events
}
