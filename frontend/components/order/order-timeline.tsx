import { CheckCircle, Clock, Package, Truck, XCircle } from "lucide-react"
import { formatDate } from "@/lib/utils"

export default function OrderTimeline({
  status,
  inventoryLocked,
  paymentProcessed,
  shippingScheduled,
  createdAt,
  updatedAt,
}: {
  status: string
  inventoryLocked: boolean
  paymentProcessed: boolean
  shippingScheduled: boolean
  createdAt: string
  updatedAt: string
}) {
  const steps = [
    {
      title: "Đặt hàng",
      description: `Đơn hàng được đặt vào ${formatDate(createdAt)}`,
      icon: <Clock className="h-5 w-5" />,
      completed: true,
    },
    {
      title: "Kiểm tra tồn kho",
      description: inventoryLocked
        ? "Sản phẩm có sẵn và đã được đặt trước cho đơn hàng của bạn"
        : "Đang kiểm tra tình trạng sản phẩm",
      icon: <Package className="h-5 w-5" />,
      completed: inventoryLocked,
    },
    {
      title: "Thanh toán",
      description: paymentProcessed
        ? "Thanh toán thành công"
        : status === "FAILED" && !paymentProcessed
          ? "Thanh toán thất bại"
          : "Đang xử lý thanh toán",
      icon: paymentProcessed ? (
        <CheckCircle className="h-5 w-5" />
      ) : status === "FAILED" && !paymentProcessed ? (
        <XCircle className="h-5 w-5" />
      ) : (
        <Clock className="h-5 w-5" />
      ),
      completed: paymentProcessed,
      failed: status === "FAILED" && !paymentProcessed,
    },
    {
      title: "Vận chuyển",
      description: shippingScheduled ? "Đơn hàng của bạn đã được gửi đi" : "Đang chờ vận chuyển",
      icon: <Truck className="h-5 w-5" />,
      completed: shippingScheduled,
    },
    {
      title: "Giao hàng",
      description: status === "DELIVERED" ? `Đơn hàng đã được giao vào ${formatDate(updatedAt)}` : "Đang chờ giao hàng",
      icon: status === "DELIVERED" ? <CheckCircle className="h-5 w-5" /> : <Clock className="h-5 w-5" />,
      completed: status === "DELIVERED",
    },
  ]

  return (
    <div className="space-y-4">
      {steps.map((step, index) => (
        <div key={index} className="flex items-start gap-4">
          <div
            className={`flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full ${
              step.failed
                ? "bg-red-100 text-red-600"
                : step.completed
                  ? "bg-green-100 text-green-600"
                  : "bg-gray-100 text-gray-500"
            }`}
          >
            {step.icon}
          </div>
          <div className="flex-1">
            <h3 className="font-medium">{step.title}</h3>
            <p className="text-sm text-muted-foreground">{step.description}</p>
          </div>
        </div>
      ))}
    </div>
  )
}
