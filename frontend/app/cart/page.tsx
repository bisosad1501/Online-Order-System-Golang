"use client"

import { useEffect } from "react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { useCart } from "@/context/cart-context"
import { useAuth } from "@/context/auth-context"
import CartItem from "@/components/cart/cart-item"
import CartSummary from "@/components/cart/cart-summary"
import { Button } from "@/components/ui/button"
import { ShoppingBag } from "lucide-react"

export default function CartPage() {
  const { cart, loading, fetchCart } = useCart()
  const { user, isAuthenticated, isLoading: authLoading } = useAuth()
  const router = useRouter()

  useEffect(() => {
    // Chỉ fetch cart khi đã xác thực thành công và có thông tin user
    if (isAuthenticated && user) {
      fetchCart(user.id)
    }
    // Nếu đã hoàn thành quá trình xác thực nhưng không có user, chuyển hướng đến trang đăng nhập
    else if (!authLoading && !isAuthenticated) {
      router.push('/login?redirect=/cart')
    }
  }, [isAuthenticated, user, fetchCart, authLoading, router])

  // Hiển thị loading khi đang tải thông tin giỏ hàng hoặc đang xác thực
  if (loading || authLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="mb-6 text-3xl font-bold">Your Cart</h1>
        <div className="flex flex-col items-center justify-center">
          <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-t-2 border-blue-600 mb-4"></div>
          <p className="text-gray-500">{authLoading ? "Verifying your account..." : "Loading your cart..."}</p>
        </div>
      </div>
    )
  }

  // Hiển thị thông báo nếu chưa đăng nhập
  if (!isAuthenticated && !authLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="mb-6 text-3xl font-bold">Your Cart</h1>
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed p-12 text-center">
          <ShoppingBag className="mb-4 h-16 w-16 text-gray-400" />
          <h2 className="mb-2 text-xl font-semibold">Please sign in to view your cart</h2>
          <p className="mb-6 text-gray-500">You need to be signed in to view and manage your shopping cart.</p>
          <Button asChild>
            <Link href="/login?redirect=/cart">Sign In</Link>
          </Button>
        </div>
      </div>
    )
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="mb-6 text-3xl font-bold">Your Cart</h1>
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed p-12 text-center">
          <ShoppingBag className="mb-4 h-16 w-16 text-gray-400" />
          <h2 className="mb-2 text-xl font-semibold">Your cart is empty</h2>
          <p className="mb-6 text-gray-500">Looks like you haven&apos;t added any products to your cart yet.</p>
          <Button asChild>
            <Link href="/products">Browse Products</Link>
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="mb-6 text-3xl font-bold">Your Cart</h1>

      <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
        <div className="lg:col-span-2">
          <div className="rounded-lg border">
            <div className="border-b p-4">
              <h2 className="text-lg font-semibold">Cart Items ({cart.items.length})</h2>
            </div>
            <div className="divide-y">
              {cart.items.map((item) => (
                <CartItem key={item.id} item={item} />
              ))}
            </div>
          </div>
        </div>

        <div>
          <CartSummary items={cart.items} onCheckout={() => router.push("/checkout")} />
        </div>
      </div>
    </div>
  )
}
