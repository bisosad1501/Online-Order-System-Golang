import Link from "next/link"
import Image from "next/image"
import { productService } from "@/services/product-service"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardFooter } from "@/components/ui/card"
import AddToCartButton from "@/components/cart/add-to-cart-button"
import { Eye } from "lucide-react"

export default async function FeaturedProducts() {
  // Thêm xử lý lỗi khi lấy sản phẩm
  let products = []
  try {
    products = await productService.getProducts({ limit: 4 })
  } catch (error) {
    console.error("Failed to fetch featured products:", error)
  }

  return (
    <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
      {products.map((product) => (
        <Card
          key={product.id}
          className="group overflow-hidden rounded-xl border-0 shadow-lg hover:shadow-xl transition-all duration-300"
        >
          <div className="relative aspect-square overflow-hidden">
            <Image
              src={product.image_url || `/default-product.svg`}
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
          </div>

          <CardContent className="p-5">
            <div className="mb-2 flex flex-wrap gap-1">
              {product.tags &&
                product.tags.slice(0, 2).map((tag) => (
                  <Badge key={tag} variant="secondary" className="text-xs bg-blue-100 text-blue-800 hover:bg-blue-200">
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
              <p className="text-xl font-bold text-blue-600">${product.price.toFixed(2)}</p>
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
            <Button asChild variant="outline" className="w-full group-hover:bg-blue-50 transition-colors duration-200">
              <Link href={`/products/${product.id}`} className="flex items-center justify-center">
                <Eye className="mr-2 h-4 w-4" /> Xem chi tiết
              </Link>
            </Button>
          </CardFooter>
        </Card>
      ))}
    </div>
  )
}
