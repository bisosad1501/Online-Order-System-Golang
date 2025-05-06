"use client"

import type React from "react"

import { useState, useEffect } from "react"
import Image from "next/image"
import Link from "next/link"
import { useCart } from "@/context/cart-context"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { useToast } from "@/hooks/use-toast"
import { Minus, Plus, Trash2 } from "lucide-react"
import { productService } from "@/services/product-service"

type CartItem = {
  id: string
  cart_id: string
  product_id: string
  quantity: number
  price: number
  product_name?: string
}

type Product = {
  id: string
  name: string
  description: string
  price: number
  [key: string]: any
}

export default function CartItem({ item }: { item: CartItem }) {
  const { updateCartItem, removeFromCart } = useCart()
  const { toast } = useToast()
  const [quantity, setQuantity] = useState(item.quantity)
  const [isUpdating, setIsUpdating] = useState(false)
  const [isRemoving, setIsRemoving] = useState(false)
  const [product, setProduct] = useState<Product | null>(null)
  const [isLoadingProduct, setIsLoadingProduct] = useState(true)

  useEffect(() => {
    async function fetchProductDetails() {
      setIsLoadingProduct(true)
      try {
        const productData = await productService.getProductById(item.product_id)
        setProduct(productData)
      } catch (error) {
        console.error("Failed to fetch product details:", error)
      } finally {
        setIsLoadingProduct(false)
      }
    }

    fetchProductDetails()
  }, [item.product_id])

  const handleQuantityChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = Number.parseInt(e.target.value)
    if (!isNaN(value) && value > 0) {
      setQuantity(value)
      updateQuantity(value)
    }
  }

  const increaseQuantity = () => {
    const newQuantity = quantity + 1
    setQuantity(newQuantity)
    updateQuantity(newQuantity)
  }

  const decreaseQuantity = () => {
    if (quantity > 1) {
      const newQuantity = quantity - 1
      setQuantity(newQuantity)
      updateQuantity(newQuantity)
    }
  }

  const updateQuantity = async (newQuantity: number) => {
    setIsUpdating(true)
    try {
      await updateCartItem(item.id, newQuantity)
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to update quantity",
        variant: "destructive",
      })
      // Reset to original quantity on error
      setQuantity(item.quantity)
    } finally {
      setIsUpdating(false)
    }
  }

  const handleRemove = async () => {
    setIsRemoving(true)
    try {
      await removeFromCart(item.id)
      toast({
        title: "Item removed",
        description: "Item has been removed from your cart",
      })
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to remove item from cart",
        variant: "destructive",
      })
    } finally {
      setIsRemoving(false)
    }
  }

  // Hiển thị skeleton loader khi đang tải thông tin sản phẩm
  if (isLoadingProduct) {
    return (
      <div className="flex items-start gap-4 p-4">
        <div className="relative h-20 w-20 flex-shrink-0 overflow-hidden rounded-md bg-gray-200 animate-pulse"></div>
        <div className="flex flex-1 flex-col gap-2">
          <div className="h-5 w-3/4 bg-gray-200 rounded animate-pulse"></div>
          <div className="h-4 w-1/4 bg-gray-200 rounded animate-pulse"></div>
          <div className="mt-2 flex items-center gap-2">
            <div className="h-8 w-8 bg-gray-200 rounded animate-pulse"></div>
            <div className="h-8 w-16 bg-gray-200 rounded animate-pulse"></div>
            <div className="h-8 w-8 bg-gray-200 rounded animate-pulse"></div>
          </div>
        </div>
        <div className="flex flex-col items-end gap-2">
          <div className="h-5 w-16 bg-gray-200 rounded animate-pulse"></div>
          <div className="h-8 w-8 bg-gray-200 rounded animate-pulse"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex items-start gap-4 p-4">
      <div className="relative h-20 w-20 flex-shrink-0 overflow-hidden rounded-md bg-gray-100">
        <Image
          src={product?.image_url || `/default-product.svg`}
          alt={product?.name || `Product ${item.product_id.substring(0, 8)}`}
          fill
          className="object-cover"
        />
      </div>

      <div className="flex flex-1 flex-col">
        <Link href={`/products/${item.product_id}`} className="font-medium hover:underline">
          {product?.name || `Product ${item.product_id.substring(0, 8)}`}
        </Link>

        <p className="text-sm text-muted-foreground">${item.price.toFixed(2)}</p>

        <div className="mt-2 flex items-center gap-2">
          <Button
            variant="outline"
            size="icon"
            className="h-8 w-8"
            onClick={decreaseQuantity}
            disabled={quantity <= 1 || isUpdating}
          >
            <Minus className="h-3 w-3" />
            <span className="sr-only">Decrease quantity</span>
          </Button>

          <Input
            type="number"
            min="1"
            value={quantity}
            onChange={handleQuantityChange}
            className="h-8 w-16 text-center"
            disabled={isUpdating}
          />

          <Button variant="outline" size="icon" className="h-8 w-8" onClick={increaseQuantity} disabled={isUpdating}>
            <Plus className="h-3 w-3" />
            <span className="sr-only">Increase quantity</span>
          </Button>
        </div>
      </div>

      <div className="flex flex-col items-end gap-2">
        <p className="font-medium">${(item.price * quantity).toFixed(2)}</p>

        <Button
          variant="ghost"
          size="icon"
          className="h-8 w-8 text-red-500 hover:bg-red-50 hover:text-red-600"
          onClick={handleRemove}
          disabled={isRemoving}
        >
          <Trash2 className="h-4 w-4" />
          <span className="sr-only">Remove item</span>
        </Button>
      </div>
    </div>
  )
}
