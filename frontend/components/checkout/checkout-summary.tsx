import Image from "next/image"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { Badge } from "@/components/ui/badge"
import { ShoppingCart, Package, CreditCard, Truck } from "lucide-react"

type CartItem = {
  id: string
  product_id: string
  quantity: number
  price: number
  product_name?: string
  product_image?: string
}

export default function CheckoutSummary({
  items,
}: {
  items: CartItem[]
}) {
  const subtotal = items.reduce((total, item) => total + item.price * item.quantity, 0)

  // In a real app, you might calculate these dynamically
  const shipping = 0 // Free shipping
  const tax = subtotal * 0.1 // 10% tax
  const total = subtotal + shipping + tax

  return (
    <Card>
      <CardHeader>
        <CardTitle>Tóm tắt đơn hàng</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <div className="space-y-3">
            {items.map((item) => (
              <div key={item.id} className="flex items-start gap-3">
                <div className="relative h-16 w-16 overflow-hidden rounded-md bg-gray-100">
                  <Image
                    src={item.product_image || `/images/placeholder.svg`}
                    alt={item.product_name || `Product ${item.product_id.substring(0, 8)}`}
                    fill
                    className="object-cover"
                    onError={(e) => {
                      // Fallback to placeholder if image fails to load
                      const target = e.target as HTMLImageElement;
                      target.src = '/images/placeholder.svg';
                    }}
                  />
                </div>
                <div className="flex-1">
                  <p className="text-sm font-medium">
                    {item.product_name || `Sản phẩm ${item.product_id.substring(0, 8)}`}
                  </p>
                  <p className="text-xs text-muted-foreground">Số lượng: {item.quantity}</p>
                  <p className="text-xs text-muted-foreground">Mã SP: {item.product_id.substring(0, 8)}...</p>
                </div>
                <div className="text-right">
                  <p className="text-sm font-medium">${item.price.toFixed(2)}</p>
                  <p className="text-xs text-muted-foreground">${(item.price * item.quantity).toFixed(2)}</p>
                </div>
              </div>
            ))}
          </div>

          <Separator />

          <div className="space-y-1.5">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Tạm tính</span>
              <span>${subtotal.toFixed(2)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Phí vận chuyển</span>
              <span>{shipping === 0 ? "Miễn phí" : `$${shipping.toFixed(2)}`}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Thuế (10%)</span>
              <span>${tax.toFixed(2)}</span>
            </div>
            <Separator className="my-2" />
            <div className="flex justify-between font-medium">
              <span>Tổng cộng</span>
              <span>${total.toFixed(2)}</span>
            </div>
          </div>

          <div className="mt-6 rounded-lg border border-blue-100 bg-blue-50 p-4">
            <h3 className="mb-2 text-sm font-medium text-blue-800">Quy trình đặt hàng</h3>
            <div className="space-y-2">
              <div className="flex items-center gap-2">
                <ShoppingCart className="h-4 w-4 text-blue-600" />
                <span className="text-xs text-blue-700">Bước 1-2: Đặt hàng và xác nhận</span>
              </div>
              <div className="flex items-center gap-2">
                <Package className="h-4 w-4 text-blue-600" />
                <span className="text-xs text-blue-700">Bước 3-5: Kiểm tra tồn kho</span>
              </div>
              <div className="flex items-center gap-2">
                <CreditCard className="h-4 w-4 text-blue-600" />
                <span className="text-xs text-blue-700">Bước 6-8: Xử lý thanh toán</span>
              </div>
              <div className="flex items-center gap-2">
                <Truck className="h-4 w-4 text-blue-600" />
                <span className="text-xs text-blue-700">Bước 9-14: Giao hàng và xác nhận</span>
              </div>
            </div>
          </div>

          <div className="mt-4">
            <Badge variant="outline" className="w-full justify-center bg-green-50 py-1 text-green-700">
              Đảm bảo hoàn tiền trong 30 ngày
            </Badge>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
