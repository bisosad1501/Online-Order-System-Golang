"use client"

import { useState, useEffect } from "react"
import Link from "next/link"
import Image from "next/image"
import { motion } from "framer-motion"
import { useInView } from "react-intersection-observer"
import { productService } from "@/services/product-service"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardFooter } from "@/components/ui/card"
import AddToCartButton from "@/components/cart/add-to-cart-button"
import { Eye } from "lucide-react"

type Product = {
  id: string
  name: string
  description: string
  category_id: string
  price: number
  tags: string[]
  created_at: string
  updated_at: string
}

export default function ProductList({
  page = 1,
  limit = 12,
  category,
  search,
}: {
  page?: number
  limit?: number
  category?: string
  search?: string
}) {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const { ref, inView } = useInView({
    triggerOnce: true,
    threshold: 0.1,
  })

  useEffect(() => {
    async function fetchProducts() {
      try {
        setLoading(true)
        setError(null)
        const data = await productService.getProducts({ page, limit, category, search })
        setProducts(data)
      } catch (error) {
        console.error("Failed to fetch products:", error)
        setError(error instanceof Error ? error.message : "Đã xảy ra lỗi khi tải dữ liệu sản phẩm")
        setProducts([])
      } finally {
        setLoading(false)
      }
    }

    fetchProducts()
  }, [page, limit, category, search])

  const container = {
    hidden: { opacity: 0 },
    show: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1,
      },
    },
  }

  const item = {
    hidden: { opacity: 0, y: 20 },
    show: { opacity: 1, y: 0, transition: { duration: 0.5 } },
  }

  if (loading) {
    return (
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
        {[...Array(limit)].map((_, i) => (
          <div key={i} className="rounded-xl border p-4 shadow-md animate-pulse">
            <div className="aspect-square w-full rounded-lg bg-gray-200 dark:bg-gray-700"></div>
            <div className="mt-4 h-6 w-3/4 rounded-md bg-gray-200 dark:bg-gray-700"></div>
            <div className="mt-2 h-4 w-1/2 rounded-md bg-gray-200 dark:bg-gray-700"></div>
            <div className="mt-4 flex gap-2">
              <div className="h-10 flex-1 rounded-md bg-gray-200 dark:bg-gray-700"></div>
              <div className="h-10 w-10 rounded-md bg-gray-200 dark:bg-gray-700"></div>
            </div>
          </div>
        ))}
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-red-200 bg-red-50 p-12 text-center dark:border-red-800 dark:bg-red-900/20">
        <h2 className="mb-2 text-xl font-semibold text-red-600 dark:text-red-400">Đã xảy ra lỗi</h2>
        <p className="mb-6 text-red-500 dark:text-red-400">
          {error}
        </p>
        <Button
          variant="outline"
          onClick={() => window.location.reload()}
          className="border-red-200 text-red-600 hover:bg-red-50 hover:text-red-700 dark:border-red-800 dark:text-red-400 dark:hover:bg-red-900/30"
        >
          Thử lại
        </Button>
      </div>
    )
  }

  if (!products.length) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed p-12 text-center">
        <h2 className="mb-2 text-xl font-semibold">Không tìm thấy sản phẩm</h2>
        <p className="mb-6 text-gray-500 dark:text-gray-400">
          Vui lòng điều chỉnh tìm kiếm hoặc bộ lọc để tìm thấy sản phẩm phù hợp.
        </p>
      </div>
    )
  }

  return (
    <motion.div
      ref={ref}
      variants={container}
      initial="hidden"
      animate={inView ? "show" : "hidden"}
      className="grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4"
    >
      {products.map((product) => (
        <motion.div key={product.id} variants={item} className="product-card-3d">
          <Card className="group overflow-hidden rounded-xl border-0 shadow-lg hover:shadow-xl transition-all duration-300 bg-white dark:bg-gray-800">
            <div className="relative aspect-square overflow-hidden">
              <Image
                src={product.image_url || `/default-product-blue.svg`}
                alt={product.name}
                fill
                priority
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
              <p className="line-clamp-2 text-sm text-muted-foreground mt-1">{product.description}</p>
              <div className="mt-2 flex items-center justify-between">
                <p className="text-xl font-bold text-blue-600 dark:text-blue-400">${product.price.toFixed(2)}</p>
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
      ))}
    </motion.div>
  )
}
