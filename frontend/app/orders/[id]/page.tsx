"use client"

import { useEffect, useState } from "react"
import { useParams } from "next/navigation"
import Image from "next/image"
import Link from "next/link"
import { orderService } from "@/services/order-service"
import { paymentService } from "@/services/payment-service"
import { shippingService } from "@/services/shipping-service"
import OrderStatusStepper from "@/components/checkout/order-status-stepper"
import OrderTracking from "@/components/order/order-tracking"
import PaymentRetryForm from "@/components/order/payment-retry-form"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import { formatDate } from "@/lib/utils"
import { useToast } from "@/hooks/use-toast"
import { Download, Truck, Star } from "lucide-react"

// Dựa trên schema trong tài liệu phân tích thiết kế
type OrderDetail = {
  id: string
  customer_id: string
  status: "CREATED" | "INVENTORY_CHECKED" | "PAYMENT_PROCESSED" | "CONFIRMED" | "SHIPPING_SCHEDULED" | "SHIPPED" | "DELIVERED" | "CANCELLED" | "FAILED"
  total_amount: number
  shipping_address: string
  items: {
    id: string
    order_id: string
    product_id: string
    quantity: number
    price: number
    product_name?: string
  }[]
  created_at: string
  updated_at: string
  inventory_locked: boolean
  payment_processed: boolean
  shipping_scheduled: boolean
  failure_reason: string | null
  estimated_delivery_date?: string
}

type Payment = {
  id: string
  order_id: string
  amount: number
  status: "SUCCEEDED" | "FAILED"
  payment_method: string
  transaction_id: string
  created_at: string
}

// Dựa trên schema trong tài liệu phân tích thiết kế
type Shipment = {
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
  last_update?: {
    status: string
    description: string
    timestamp: string
  }
}

export default function OrderDetailPage() {
  const params = useParams()
  const { toast } = useToast()
  const [order, setOrder] = useState<OrderDetail | null>(null)
  const [payment, setPayment] = useState<Payment | null>(null)
  const [shipment, setShipment] = useState<Shipment | null>(null)
  const [loading, setLoading] = useState(true)
  const [showPaymentRetry, setShowPaymentRetry] = useState(false)

  const orderId = params.id as string

  useEffect(() => {
    async function fetchOrderDetails() {
      try {
        setLoading(true)

        // Fetch order details
        const orderData = await orderService.getOrderById(orderId)
        setOrder(orderData)

        // Fetch payment details
        try {
          const paymentData = await paymentService.getPaymentByOrderId(orderId)
          setPayment(paymentData)
        } catch (error) {
          console.log("No payment data found")
        }

        // Fetch shipment details
        try {
          const shipmentData = await shippingService.getShipmentByOrderId(orderId)
          setShipment(shipmentData)
        } catch (error) {
          console.log("No shipment data found")
        }
      } catch (error) {
        toast({
          title: "Error",
          description: "Failed to load order details",
          variant: "destructive",
        })
      } finally {
        setLoading(false)
      }
    }

    fetchOrderDetails()
  }, [orderId, toast])

  const handleCancelOrder = async () => {
    try {
      await orderService.cancelOrder(orderId)
      toast({
        title: "Order cancelled",
        description: "Your order has been cancelled successfully",
      })
      // Refresh order data
      const orderData = await orderService.getOrderById(orderId)
      setOrder(orderData)
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to cancel order",
        variant: "destructive",
      })
    }
  }

  const handleDownloadInvoice = () => {
    toast({
      title: "Invoice download",
      description: "Your invoice is being generated and will download shortly",
    })
    // In a real app, this would trigger a download
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center">
          <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-t-2 border-blue-600"></div>
        </div>
      </div>
    )
  }

  if (!order) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="rounded-lg border border-red-200 bg-red-50 p-6 text-center">
          <h2 className="text-xl font-semibold text-red-800">Order not found</h2>
          <p className="mt-2 text-red-600">
            The order you are looking for does not exist or you don&apos;t have permission to view it.
          </p>
        </div>
      </div>
    )
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case "CREATED":
        return "bg-blue-100 text-blue-800"
      case "INVENTORY_CHECKED":
        return "bg-indigo-100 text-indigo-800"
      case "PAYMENT_PROCESSED":
        return "bg-violet-100 text-violet-800"
      case "CONFIRMED":
        return "bg-purple-100 text-purple-800"
      case "SHIPPING_SCHEDULED":
        return "bg-cyan-100 text-cyan-800"
      case "SHIPPED":
        return "bg-teal-100 text-teal-800"
      case "DELIVERED":
        return "bg-green-100 text-green-800"
      case "CANCELLED":
        return "bg-gray-100 text-gray-800"
      case "FAILED":
        return "bg-red-100 text-red-800"
      default:
        return "bg-gray-100 text-gray-800"
    }
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-6 flex flex-col justify-between gap-4 md:flex-row md:items-center">
        <div>
          <h1 className="text-3xl font-bold">Order #{order.id.substring(0, 8)}</h1>
          <p className="text-muted-foreground">Placed on {formatDate(order.created_at)}</p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button variant="outline" onClick={handleDownloadInvoice}>
            <Download className="mr-2 h-4 w-4" />
            Download Invoice
          </Button>
          {order.status === "DELIVERED" && (
            <Button asChild>
              <Link href={`/orders/${order.id}/review`}>
                <Star className="mr-2 h-4 w-4" />
                Review Products
              </Link>
            </Button>
          )}
          {order.status === "CREATED" && (
            <Button variant="destructive" onClick={handleCancelOrder}>
              Cancel Order
            </Button>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <div className="lg:col-span-2">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Order Status</CardTitle>
                <Badge className={getStatusColor(order.status)}>{order.status}</Badge>
              </div>
              {order.failure_reason && (
                <CardDescription className="text-red-500">Failure reason: {order.failure_reason}</CardDescription>
              )}
            </CardHeader>
            <CardContent>
              <OrderStatusStepper
                status={order.status}
                failureReason={order.failure_reason}
              />

              {order.status === "FAILED" && order.failure_reason?.includes("payment") && (
                <div className="mt-6">
                  {showPaymentRetry ? (
                    <PaymentRetryForm
                      orderId={order.id}
                      onSuccess={() => {
                        setShowPaymentRetry(false)
                        // Refresh order data
                        orderService.getOrderById(orderId).then(setOrder)
                      }}
                      onCancel={() => setShowPaymentRetry(false)}
                    />
                  ) : (
                    <Button onClick={() => setShowPaymentRetry(true)}>Retry Payment</Button>
                  )}
                </div>
              )}
            </CardContent>
          </Card>

          <Card className="mt-6">
            <CardHeader>
              <CardTitle>Order Items</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {order.items.map((item) => (
                  <div key={item.id} className="flex items-start gap-4">
                    <div className="relative h-16 w-16 overflow-hidden rounded-md bg-gray-100">
                      <Image
                        src={`/default-product.svg`}
                        alt={item.product_name || `Product ${item.product_id.substring(0, 8)}`}
                        fill
                        className="object-cover"
                      />
                    </div>
                    <div className="flex-1">
                      <p className="font-medium">{item.product_name || `Product ${item.product_id.substring(0, 8)}`}</p>
                      <p className="text-sm text-muted-foreground">Quantity: {item.quantity}</p>
                    </div>
                    <div className="text-right">
                      <p className="font-medium">${item.price.toFixed(2)}</p>
                      <p className="text-sm text-muted-foreground">${(item.price * item.quantity).toFixed(2)}</p>
                    </div>
                  </div>
                ))}
              </div>

              <Separator className="my-4" />

              <div className="space-y-1.5">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Subtotal</span>
                  <span>${order.total_amount.toFixed(2)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Shipping</span>
                  <span>Free</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Tax</span>
                  <span>Included</span>
                </div>
                <Separator className="my-2" />
                <div className="flex justify-between font-medium">
                  <span>Total</span>
                  <span>${order.total_amount.toFixed(2)}</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <div className="space-y-6">
          {/* Shipping Information */}
          <Card>
            <CardHeader>
              <CardTitle>Shipping Information</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div>
                  <h3 className="text-sm font-medium text-muted-foreground">Shipping Address</h3>
                  <p className="mt-1">{order.shipping_address}</p>
                </div>

                {shipment && (
                  <>
                    <div>
                      <h3 className="text-sm font-medium text-muted-foreground">Đơn vị vận chuyển</h3>
                      <p className="mt-1">{shipment.carrier}</p>
                    </div>

                    <div>
                      <h3 className="text-sm font-medium text-muted-foreground">Mã vận đơn</h3>
                      <p className="mt-1">{shipment.tracking_number}</p>
                    </div>

                    {shipment.estimated_delivery && (
                      <div>
                        <h3 className="text-sm font-medium text-muted-foreground">Ngày giao hàng dự kiến</h3>
                        <p className="mt-1">{formatDate(shipment.estimated_delivery)}</p>
                      </div>
                    )}

                    <div>
                      <h3 className="text-sm font-medium text-muted-foreground">Trạng thái giao hàng</h3>
                      <div className="mt-1 flex items-center gap-2">
                        <Badge variant="outline" className={
                          shipment.status === "DELIVERED" ? "bg-green-50" :
                          shipment.status === "IN_TRANSIT" ? "bg-blue-50" : "bg-yellow-50"
                        }>
                          {shipment.status === "PENDING" ? "Đã lên lịch" :
                           shipment.status === "IN_TRANSIT" ? "Đang giao hàng" : "Đã giao hàng"}
                        </Badge>
                        {shipment.status === "IN_TRANSIT" && <Truck className="h-4 w-4 text-blue-500" />}
                      </div>
                    </div>

                    {shipment.last_update && (
                      <div>
                        <h3 className="text-sm font-medium text-muted-foreground">Cập nhật mới nhất</h3>
                        <p className="mt-1">{shipment.last_update.description}</p>
                        <p className="text-xs text-muted-foreground">{formatDate(shipment.last_update.timestamp)}</p>
                      </div>
                    )}
                  </>
                )}
              </div>
            </CardContent>
          </Card>

          {/* Order Tracking */}
          {(shipment || order.status === "CONFIRMED" || order.status === "SHIPPING_SCHEDULED" || order.status === "SHIPPED" || order.status === "DELIVERED") && (
            <Card className="mt-6">
              <OrderTracking orderId={order.id} shipment={shipment} />
            </Card>
          )}

          <Card>
            <CardHeader>
              <CardTitle>Payment Information</CardTitle>
            </CardHeader>
            <CardContent>
              {payment ? (
                <div className="space-y-4">
                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground">Payment Method</h3>
                    <p className="mt-1 capitalize">{payment.payment_method.replace("_", " ")}</p>
                  </div>

                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground">Transaction ID</h3>
                    <p className="mt-1">{payment.transaction_id}</p>
                  </div>

                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground">Payment Status</h3>
                    <Badge variant="outline" className={payment.status === "SUCCEEDED" ? "bg-green-50" : "bg-red-50"}>
                      {payment.status}
                    </Badge>
                  </div>

                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground">Payment Date</h3>
                    <p className="mt-1">{formatDate(payment.created_at)}</p>
                  </div>
                </div>
              ) : (
                <p className="text-muted-foreground">
                  {order.status === "FAILED" && order.failure_reason?.includes("payment")
                    ? "Payment failed. Please retry payment."
                    : "Payment information not available yet."}
                </p>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
