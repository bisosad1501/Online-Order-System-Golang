"use client"

import { CheckCircle, Clock, Package, ShoppingCart, Truck, CreditCard, AlertCircle } from "lucide-react"

type OrderStatus =
  | "CREATED"
  | "INVENTORY_CHECKED"
  | "PAYMENT_PROCESSED"
  | "CONFIRMED"
  | "SHIPPING_SCHEDULED"
  | "SHIPPED"
  | "DELIVERED"
  | "CANCELLED"
  | "FAILED"

interface OrderStatusStepperProps {
  status: OrderStatus
  failureReason?: string
}

export default function OrderStatusStepper({ status, failureReason }: OrderStatusStepperProps) {
  // Các bước dựa trên quy trình 14 bước trong tài liệu phân tích thiết kế
  const steps = [
    {
      id: "cart",
      label: "Đặt hàng",
      icon: ShoppingCart,
      completed: true,
      description: "Bước 1: Khách hàng đặt hàng và thêm vào giỏ hàng"
    },
    {
      id: "order_confirmation",
      label: "Xác nhận đơn hàng",
      icon: CheckCircle,
      completed: ["CREATED", "INVENTORY_CHECKED", "PAYMENT_PROCESSED", "CONFIRMED", "SHIPPING_SCHEDULED", "SHIPPED", "DELIVERED"].includes(status),
      description: "Bước 2: Khách hàng xác nhận đơn hàng và nhập thông tin giao hàng"
    },
    {
      id: "inventory",
      label: "Kiểm tra tồn kho",
      icon: Package,
      completed: ["INVENTORY_CHECKED", "PAYMENT_PROCESSED", "CONFIRMED", "SHIPPING_SCHEDULED", "SHIPPED", "DELIVERED"].includes(status),
      description: "Bước 3-5: Hệ thống kiểm tra và khóa tồn kho"
    },
    {
      id: "payment",
      label: "Thanh toán",
      icon: CreditCard,
      completed: ["PAYMENT_PROCESSED", "CONFIRMED", "SHIPPING_SCHEDULED", "SHIPPED", "DELIVERED"].includes(status),
      description: "Bước 6-7: Xử lý thanh toán đơn hàng"
    },
    {
      id: "order_confirmed",
      label: "Xác nhận đơn hàng",
      icon: CheckCircle,
      completed: ["CONFIRMED", "SHIPPING_SCHEDULED", "SHIPPED", "DELIVERED"].includes(status),
      description: "Bước 8: Đơn hàng được xác nhận sau khi thanh toán thành công"
    },
    {
      id: "shipping_scheduled",
      label: "Lên lịch giao hàng",
      icon: Clock,
      completed: ["SHIPPING_SCHEDULED", "SHIPPED", "DELIVERED"].includes(status),
      description: "Bước 9-10: Hệ thống lên lịch và tạo đơn giao hàng"
    },
    {
      id: "shipping_info",
      label: "Đang giao hàng",
      icon: Truck,
      completed: ["SHIPPED", "DELIVERED"].includes(status),
      description: "Bước 11-12: Đơn hàng đang được vận chuyển"
    },
    {
      id: "delivered",
      label: "Đã giao hàng",
      icon: CheckCircle,
      completed: ["DELIVERED"].includes(status),
      description: "Bước 13-14: Xác nhận giao hàng thành công và gửi thông báo"
    },
  ]

  if (status === "FAILED") {
    return (
      <div className="rounded-lg border border-red-200 bg-red-50 p-4">
        <div className="flex items-center gap-3">
          <AlertCircle className="h-6 w-6 text-red-500" />
          <h3 className="text-lg font-semibold text-red-700">Đơn hàng không thành công</h3>
        </div>
        <p className="mt-2 text-red-600">
          {failureReason === "inventory_unavailable"
            ? "Bước 4: Một số sản phẩm trong đơn hàng của bạn hiện không có sẵn trong kho. Hệ thống đề xuất các sản phẩm tương tự."
            : failureReason === "payment_failed"
            ? "Bước 7: Thanh toán không thành công. Vui lòng kiểm tra thông tin thanh toán và thử lại."
            : "Đã xảy ra lỗi khi xử lý đơn hàng của bạn."}
        </p>
      </div>
    )
  }

  if (status === "CANCELLED") {
    return (
      <div className="rounded-lg border border-gray-200 bg-gray-50 p-4">
        <div className="flex items-center gap-3">
          <AlertCircle className="h-6 w-6 text-gray-500" />
          <h3 className="text-lg font-semibold text-gray-700">Đơn hàng đã bị hủy</h3>
        </div>
        <p className="mt-2 text-gray-600">
          Đơn hàng này đã bị hủy. Bạn có thể đặt lại đơn hàng mới nếu muốn.
        </p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="relative">
        <div className="absolute left-5 top-0 h-full w-0.5 bg-gray-200"></div>
        <div className="space-y-8">
          {steps.map((step, index) => (
            <div key={step.id} className="relative">
              <div className={`absolute left-5 top-5 h-0.5 w-5 ${step.completed ? "bg-blue-500" : "bg-gray-200"}`}></div>
              <div className="flex items-start">
                <div className={`relative z-10 flex h-10 w-10 items-center justify-center rounded-full ${
                  step.completed ? "bg-blue-500 text-white" : "bg-gray-200 text-gray-500"
                }`}>
                  <step.icon className="h-5 w-5" />
                </div>
                <div className="ml-4">
                  <h3 className={`text-lg font-semibold ${step.completed ? "text-gray-900" : "text-gray-500"}`}>
                    {step.label}
                  </h3>
                  <p className="mt-1 text-sm text-gray-500">{step.description}</p>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
