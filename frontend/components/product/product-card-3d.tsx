"use client"

import type React from "react"

import { useState, useRef } from "react"
import Image from "next/image"
import Link from "next/link"
import { motion } from "framer-motion"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardFooter } from "@/components/ui/card"
import AddToCartButton from "@/components/cart/add-to-cart-button"
import { Eye } from "lucide-react"

type ProductCardProps = {
  product: {
    id: string
    name: string
    price: number
    image: string
    category: string
    tags?: string[]
  }
}

export default function ProductCard3D({ product }: ProductCardProps) {
  const [rotateX, setRotateX] = useState(0)
  const [rotateY, setRotateY] = useState(0)
  const [scale, setScale] = useState(1)
  const cardRef = useRef<HTMLDivElement>(null)

  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    if (!cardRef.current) return

    const card = cardRef.current
    const { left, top, width, height } = card.getBoundingClientRect()
    const x = e.clientX - left
    const y = e.clientY - top

    // Calculate rotation based on mouse position
    const rotateY = ((x - width / 2) / width) * 10 // -10 to 10 degrees
    const rotateX = ((height / 2 - y) / height) * 10 // -10 to 10 degrees

    setRotateX(rotateX)
    setRotateY(rotateY)
    setScale(1.05)
  }

  const handleMouseLeave = () => {
    setRotateX(0)
    setRotateY(0)
    setScale(1)
  }

  return (
    <motion.div
      ref={cardRef}
      className="product-card-3d"
      onMouseMove={handleMouseMove}
      onMouseLeave={handleMouseLeave}
      style={{ transformStyle: "preserve-3d" }}
    >
      <motion.div
        className="product-card-inner"
        animate={{
          rotateX,
          rotateY,
          scale,
        }}
        transition={{ type: "spring", stiffness: 300, damping: 20 }}
      >
        <Card className="overflow-hidden rounded-xl border-0 shadow-lg hover:shadow-xl transition-all duration-300 bg-white dark:bg-gray-800">
          <div className="relative aspect-square overflow-hidden">
            <Image
              src={product.image || "/placeholder.svg"}
              alt={product.name}
              fill
              className="object-cover transition-transform duration-700 group-hover:scale-110"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>

            <div className="absolute inset-0 flex items-center justify-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
              <Button asChild size="icon" variant="secondary" className="rounded-full h-12 w-12 shadow-lg">
                <Link href={`/products/${product.id}`}>
                  <Eye className="h-5 w-5" />
                </Link>
              </Button>
              <AddToCartButton product={product} variant="icon" />
            </div>

            {product.tags && product.tags.includes("bestseller") && (
              <div className="absolute top-4 left-4 bg-red-500 text-white text-xs font-bold px-3 py-1 rounded-full shadow-lg">
                Bán chạy
              </div>
            )}

            {product.tags && product.tags.includes("new") && (
              <div className="absolute top-4 left-4 bg-green-500 text-white text-xs font-bold px-3 py-1 rounded-full shadow-lg">
                Mới
              </div>
            )}

            {product.tags && product.tags.includes("premium") && (
              <div className="absolute top-4 left-4 bg-purple-500 text-white text-xs font-bold px-3 py-1 rounded-full shadow-lg">
                Cao cấp
              </div>
            )}

            <div className="product-card-shine"></div>
          </div>

          <CardContent className="p-5">
            <div className="mb-2 flex flex-wrap gap-1">
              {product.tags &&
                product.tags.slice(0, 2).map((tag) => (
                  <Badge
                    key={tag}
                    variant="secondary"
                    className="text-xs bg-blue-100 text-blue-800 hover:bg-blue-200 dark:bg-blue-900 dark:text-blue-300"
                  >
                    {tag}
                  </Badge>
                ))}
            </div>
            <Link href={`/products/${product.id}`} className="block">
              <h3 className="line-clamp-1 text-lg font-semibold group-hover:text-blue-600 transition-colors duration-200">
                {product.name}
              </h3>
            </Link>
            <div className="mt-2 flex items-center justify-between">
              <p className="text-xl font-bold text-blue-600 dark:text-blue-400">
                {new Intl.NumberFormat("vi-VN", { style: "currency", currency: "VND" }).format(product.price)}
              </p>
              <div className="flex">
                {[...Array(5)].map((_, i) => (
                  <svg key={i} className="w-4 h-4 text-yellow-400 fill-current" viewBox="0 0 20 20">
                    <path d="M10 15l-5.878 3.09 1.123-6.545L.489 6.91l6.572-.955L10 0l2.939 5.955 6.572.955-4.756 4.635 1.123 6.545z" />
                  </svg>
                ))}
              </div>
            </div>
          </CardContent>

          <CardFooter className="p-5 pt-0">
            <Button
              asChild
              variant="outline"
              className="w-full group-hover:bg-blue-50 transition-colors duration-200 dark:hover:bg-blue-900"
            >
              <Link href={`/products/${product.id}`} className="flex items-center justify-center">
                <Eye className="mr-2 h-4 w-4" /> Xem chi tiết
              </Link>
            </Button>
          </CardFooter>
        </Card>
      </motion.div>
    </motion.div>
  )
}
