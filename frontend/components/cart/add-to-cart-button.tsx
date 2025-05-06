"use client"

import type React from "react"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { useCart } from "@/context/cart-context"
import { useAuth } from "@/context/auth-context"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { useToast } from "@/hooks/use-toast"
import { Minus, Plus, ShoppingCart } from "lucide-react"

type Product = {
  id: string
  name: string
  price: number
  [key: string]: any
}

export default function AddToCartButton({
  product,
  variant = "default",
}: {
  product: Product
  variant?: "default" | "icon"
}) {
  const router = useRouter()
  const { cart, fetchCart, addToCart } = useCart()
  const { user, isAuthenticated } = useAuth()
  const { toast } = useToast()
  const [quantity, setQuantity] = useState(1)
  const [isLoading, setIsLoading] = useState(false)
  const [showLoginDialog, setShowLoginDialog] = useState(false)

  const handleQuantityChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = Number.parseInt(e.target.value)
    if (!isNaN(value) && value > 0) {
      setQuantity(value)
    }
  }

  const increaseQuantity = () => {
    setQuantity((prev) => prev + 1)
  }

  const decreaseQuantity = () => {
    setQuantity((prev) => (prev > 1 ? prev - 1 : 1))
  }

  const handleAddToCart = async () => {
    if (!isAuthenticated) {
      setShowLoginDialog(true)
      return
    }

    setIsLoading(true)
    try {
      // Nếu giỏ hàng chưa được tải và người dùng đã đăng nhập, tải giỏ hàng
      if ((!cart || !cart.items) && user) {
        try {
          await fetchCart(user.id);
          // Không cần đợi 500ms nữa vì fetchCart đã có xử lý đồng bộ
        } catch (fetchError) {
          console.error("Failed to fetch cart:", fetchError);
          toast({
            title: "Lỗi",
            description: "Không thể tải giỏ hàng. Đang tạo giỏ hàng mới...",
            variant: "destructive",
          });
          // Tiếp tục thực hiện, vì addToCart sẽ tạo giỏ hàng mới nếu cần
        }
      }

      // Kiểm tra lại cart sau khi đã tải
      if (!cart && user) {
        // Nếu vẫn không có giỏ hàng, thử tạo mới
        await fetchCart(user.id);
      }

      // Nếu vẫn không có giỏ hàng sau khi thử tạo mới, báo lỗi
      if (!cart) {
        throw new Error("Không thể tạo giỏ hàng. Vui lòng tải lại trang.");
      }

      await addToCart({
        product_id: product.id,
        quantity,
        price: product.price,
      })

      // Hiển thị thông báo thành công với nút xem giỏ hàng
      toast({
        title: "Đã thêm vào giỏ hàng",
        description: (
          <div className="flex flex-col space-y-2">
            <div>{quantity} x {product.name} đã được thêm vào giỏ hàng của bạn</div>
            <Button
              size="sm"
              variant="outline"
              className="mt-2 w-full"
              onClick={() => router.push('/cart')}
            >
              Xem giỏ hàng
            </Button>
          </div>
        ),
        duration: 5000,
      })

      // Reset số lượng về 1 sau khi thêm thành công
      setQuantity(1)
    } catch (error) {
      // Hiển thị thông báo lỗi chi tiết hơn
      let errorMessage = "Không thể thêm sản phẩm vào giỏ hàng";

      if (error instanceof Error) {
        errorMessage = error.message || errorMessage;
      }

      toast({
        title: "Lỗi",
        description: errorMessage,
        variant: "destructive",
      })
    } finally {
      setIsLoading(false)
    }
  }

  if (variant === "icon") {
    return (
      <>
        <Button size="icon" onClick={handleAddToCart} disabled={isLoading} className="bg-blue-600 hover:bg-blue-700">
          <ShoppingCart className="h-4 w-4" />
          <span className="sr-only">Thêm vào giỏ hàng</span>
        </Button>

        <Dialog open={showLoginDialog} onOpenChange={setShowLoginDialog}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Yêu cầu đăng nhập</DialogTitle>
              <DialogDescription>Bạn cần đăng nhập để thêm sản phẩm vào giỏ hàng.</DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <Button variant="outline" onClick={() => setShowLoginDialog(false)}>
                Hủy
              </Button>
              <Button onClick={() => router.push("/login")} className="bg-blue-600 hover:bg-blue-700">
                Đăng nhập
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </>
    )
  }

  return (
    <>
      <div className="flex w-full flex-col gap-4">
        <div className="flex items-center gap-2">
          <Button variant="outline" size="icon" onClick={decreaseQuantity} disabled={quantity <= 1}>
            <Minus className="h-4 w-4" />
            <span className="sr-only">Giảm số lượng</span>
          </Button>

          <Input type="number" min="1" value={quantity} onChange={handleQuantityChange} className="w-20 text-center" />

          <Button variant="outline" size="icon" onClick={increaseQuantity}>
            <Plus className="h-4 w-4" />
            <span className="sr-only">Tăng số lượng</span>
          </Button>
        </div>

        <Button className="w-full bg-blue-600 hover:bg-blue-700" onClick={handleAddToCart} disabled={isLoading}>
          {isLoading ? "Đang thêm..." : "Thêm vào giỏ hàng"}
        </Button>
      </div>

      <Dialog open={showLoginDialog} onOpenChange={setShowLoginDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Yêu cầu đăng nhập</DialogTitle>
            <DialogDescription>Bạn cần đăng nhập để thêm sản phẩm vào giỏ hàng.</DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowLoginDialog(false)}>
              Hủy
            </Button>
            <Button onClick={() => router.push("/login")} className="bg-blue-600 hover:bg-blue-700">
              Đăng nhập
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
}
