"use client"

import { useState, useEffect, useRef, Suspense } from "react"
import Image from "next/image"
import Link from "next/link"
import { useParams } from "next/navigation"
import { motion, useScroll, useTransform, AnimatePresence } from "framer-motion"
import { productService } from "@/services/product-service"
import AddToCartButton from "@/components/cart/add-to-cart-button"
import ProductRecommendations from "@/components/product/product-recommendations"
import ProductReviews from "@/components/product/product-reviews"
import ProductImageGallery from "@/components/product/product-image-gallery"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Button } from "@/components/ui/button"
import { Heart, Share2, Award, Package, RefreshCw, Truck, ChevronDown, Check, Star } from "lucide-react"
import { useToast } from "@/hooks/use-toast"

// Fallback 3D Model component that doesn't require Three.js
function ProductViewer3D({ product }) {
  return (
    <div className="relative h-[500px] w-full rounded-xl overflow-hidden border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 flex items-center justify-center">
      <div className="text-center p-8">
        <div className="animate-float mb-6">
          <Image
            src={`/placeholder.svg?height=200&width=200&text=${encodeURIComponent("3D View")}`}
            alt="3D View"
            width={200}
            height={200}
            className="mx-auto"
          />
        </div>
        <h3 className="text-xl font-bold mb-2">Xem sản phẩm 3D</h3>
        <p className="text-gray-500 dark:text-gray-400 mb-4">
          Tính năng xem 3D hiện không khả dụng. Vui lòng xem hình ảnh sản phẩm.
        </p>
        <Button variant="outline" size="sm">
          Quay lại hình ảnh
        </Button>
      </div>
    </div>
  )
}

export default function ProductDetailPage() {
  const params = useParams()
  const { toast } = useToast()
  const [product, setProduct] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)
  const [activeTab, setActiveTab] = useState("description")
  const [quantity, setQuantity] = useState(1)
  const [selectedColor, setSelectedColor] = useState("black")
  const [selectedSize, setSelectedSize] = useState("m")
  const [showFullDescription, setShowFullDescription] = useState(false)
  const [showModel3D, setShowModel3D] = useState(false)

  const ref = useRef(null)
  const { scrollYProgress } = useScroll({
    target: ref,
    offset: ["start start", "end start"],
  })

  const opacity = useTransform(scrollYProgress, [0, 0.5], [1, 0])
  const scale = useTransform(scrollYProgress, [0, 0.5], [1, 0.8])

  useEffect(() => {
    async function fetchProduct() {
      try {
        setLoading(true)
        setError(false)

        // Attempt to fetch the product
        const productData = await productService.getProductById(params.id as string)

        if (productData) {
          setProduct(productData)
        } else {
          // If no product data is returned, create a fallback product
          setProduct(generateFallbackProduct(params.id as string))
          console.warn("Product not found, using fallback data")
        }
      } catch (error) {
        console.error("Error fetching product:", error)
        setError(true)
        // Create a fallback product in case of error
        setProduct(generateFallbackProduct(params.id as string))
      } finally {
        setLoading(false)
      }
    }

    fetchProduct()
  }, [params.id])

  // Function to generate fallback product data
  function generateFallbackProduct(id: string) {
    return {
      id,
      name: `Sản phẩm ${id}`,
      description:
        "Thông tin chi tiết về sản phẩm này hiện không có sẵn. Vui lòng liên hệ với chúng tôi để biết thêm chi tiết.",
      category_id: "general",
      price: 999000,
      tags: ["general"],
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }
  }

  const handleAddToWishlist = () => {
    toast({
      title: "Đã thêm vào danh sách yêu thích",
      description: "Sản phẩm đã được thêm vào danh sách yêu thích của bạn",
    })
  }

  const handleShare = () => {
    navigator.clipboard.writeText(window.location.href)
    toast({
      title: "Đã sao chép liên kết",
      description: "Liên kết sản phẩm đã được sao chép vào clipboard",
    })
  }

  if (loading) {
    return <ProductDetailSkeleton />
  }

  // Mảng hình ảnh mẫu cho sản phẩm
  const productImages = [
    product.image_url || `/default-product.svg`,
    `/default-product-blue.svg`,
    `/default-product.svg`,
    `/default-product-blue.svg`,
  ]

  // Thông số kỹ thuật mẫu
  const specifications = [
    { name: "Thương hiệu", value: "Brand X" },
    { name: "Xuất xứ", value: "Việt Nam" },
    { name: "Bảo hành", value: "12 tháng" },
    { name: "Chất liệu", value: "Cao cấp" },
    { name: "Kích thước", value: "M, L, XL" },
    { name: "Màu sắc", value: "Đen, Trắng, Xanh" },
  ]

  // Màu sắc sản phẩm
  const colors = [
    { id: "black", name: "Đen", value: "bg-black" },
    { id: "white", name: "Trắng", value: "bg-white border border-gray-300" },
    { id: "blue", name: "Xanh", value: "bg-blue-600" },
    { id: "red", name: "Đỏ", value: "bg-red-600" },
  ]

  // Kích thước sản phẩm
  const sizes = [
    { id: "s", name: "S" },
    { id: "m", name: "M" },
    { id: "l", name: "L" },
    { id: "xl", name: "XL" },
  ]

  return (
    <div className="bg-gradient-to-b from-white to-gray-50 dark:from-gray-900 dark:to-gray-800">
      {/* Breadcrumb */}
      <div className="bg-white dark:bg-gray-900 border-b dark:border-gray-800">
        <div className="container mx-auto px-4 py-3">
          <nav className="flex text-sm">
            <Link
              href="/"
              className="text-gray-500 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
            >
              Trang chủ
            </Link>
            <span className="mx-2 text-gray-400 dark:text-gray-600">/</span>
            <Link
              href="/products"
              className="text-gray-500 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
            >
              Sản phẩm
            </Link>
            <span className="mx-2 text-gray-400 dark:text-gray-600">/</span>
            <span className="text-gray-900 dark:text-white font-medium">{product.name}</span>
          </nav>
        </div>
      </div>

      {/* Error Message */}
      {error && (
        <div className="container mx-auto px-4 py-4">
          <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4 mb-6">
            <p className="text-yellow-800 dark:text-yellow-200">
              Đã xảy ra lỗi khi tải thông tin sản phẩm. Một số thông tin có thể không chính xác.
            </p>
          </div>
        </div>
      )}

      {/* Product Detail */}
      <div className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-10">
          {/* Product Images */}
          <motion.div
            ref={ref}
            style={{ opacity, scale }}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
            className="sticky top-24"
          >
            <div className="flex justify-end mb-4">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowModel3D(!showModel3D)}
                className="flex items-center gap-2"
              >
                {showModel3D ? "Xem hình ảnh" : "Xem mô hình 3D"}
                <Badge
                  variant="secondary"
                  className="ml-2 bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300"
                >
                  Mới
                </Badge>
              </Button>
            </div>

            <AnimatePresence mode="wait">
              {showModel3D ? (
                <motion.div
                  key="3d-model"
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  exit={{ opacity: 0 }}
                  transition={{ duration: 0.3 }}
                >
                  <ProductViewer3D product={product} />
                </motion.div>
              ) : (
                <motion.div
                  key="image-gallery"
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  exit={{ opacity: 0 }}
                  transition={{ duration: 0.3 }}
                >
                  <ProductImageGallery images={productImages} productName={product.name} />
                </motion.div>
              )}
            </AnimatePresence>
          </motion.div>

          {/* Product Info */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
            className="flex flex-col"
          >
            {/* Product Header */}
            <div className="mb-6">
              <div className="flex flex-wrap gap-2 mb-3">
                {product.tags &&
                  product.tags.map((tag: string) => (
                    <Badge
                      key={tag}
                      variant="secondary"
                      className="bg-blue-100 text-blue-800 hover:bg-blue-200 dark:bg-blue-900 dark:text-blue-300"
                    >
                      {tag}
                    </Badge>
                  ))}
              </div>
              <h1 className="text-3xl md:text-4xl font-bold mb-2">{product.name}</h1>

              <div className="flex items-center gap-4 mb-4">
                <div className="flex">
                  {[...Array(5)].map((_, i) => (
                    <Star
                      key={i}
                      className={`h-5 w-5 ${i < 4 ? "text-yellow-400 fill-yellow-400" : "text-gray-300 dark:text-gray-600"}`}
                    />
                  ))}
                </div>
                <Link href="#reviews" className="text-sm text-blue-600 dark:text-blue-400 font-medium hover:underline">
                  124 đánh giá
                </Link>
                <span className="text-sm text-green-600 dark:text-green-400 font-medium">Đã bán 1.5k+</span>
              </div>

              <div className="flex items-baseline gap-4">
                <span className="text-3xl font-bold text-blue-600 dark:text-blue-400">
                  {new Intl.NumberFormat("vi-VN", { style: "currency", currency: "VND" }).format(product.price)}
                </span>
                <span className="text-lg text-gray-500 dark:text-gray-400 line-through">
                  {new Intl.NumberFormat("vi-VN", { style: "currency", currency: "VND" }).format(product.price * 1.2)}
                </span>
                <span className="bg-red-100 dark:bg-red-900/50 text-red-600 dark:text-red-400 text-sm font-medium px-2 py-1 rounded">
                  Giảm 20%
                </span>
              </div>
            </div>

            {/* Product Description */}
            <div className="mb-6">
              <h3 className="font-semibold mb-2">Mô tả sản phẩm</h3>
              <div className={`text-gray-700 dark:text-gray-300 ${showFullDescription ? "" : "line-clamp-3"}`}>
                <p>{product.description}</p>
                <p className="mt-2">
                  Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam auctor, nisl eget ultricies tincidunt,
                  nisl nisl aliquam nisl, eget aliquam nisl nisl eget nisl.
                </p>
              </div>
              <Button
                variant="link"
                onClick={() => setShowFullDescription(!showFullDescription)}
                className="p-0 h-auto text-blue-600 dark:text-blue-400 mt-1 flex items-center"
              >
                {showFullDescription ? "Thu gọn" : "Xem thêm"}
                <ChevronDown
                  className={`ml-1 h-4 w-4 transition-transform ${showFullDescription ? "rotate-180" : ""}`}
                />
              </Button>
            </div>

            {/* Color Selection */}
            <div className="mb-6">
              <h3 className="font-semibold mb-2">Màu sắc</h3>
              <div className="flex gap-3">
                {colors.map((color) => (
                  <button
                    key={color.id}
                    type="button"
                    className={`w-10 h-10 rounded-full ${color.value} flex items-center justify-center transition-all ${
                      selectedColor === color.id
                        ? "ring-2 ring-offset-2 ring-blue-600 dark:ring-blue-400"
                        : "hover:ring-2 hover:ring-offset-2 hover:ring-gray-300 dark:hover:ring-gray-600"
                    }`}
                    onClick={() => setSelectedColor(color.id)}
                    aria-label={`Màu ${color.name}`}
                  >
                    {selectedColor === color.id && (
                      <Check className={`h-5 w-5 ${color.id === "white" ? "text-black" : "text-white"}`} />
                    )}
                  </button>
                ))}
              </div>
              <p className="text-sm text-gray-500 dark:text-gray-400 mt-2">
                Đã chọn:{" "}
                <span className="font-medium text-gray-900 dark:text-white capitalize">
                  {colors.find((c) => c.id === selectedColor)?.name}
                </span>
              </p>
            </div>

            {/* Size Selection */}
            <div className="mb-6">
              <h3 className="font-semibold mb-2">Kích thước</h3>
              <div className="flex gap-3">
                {sizes.map((size) => (
                  <button
                    key={size.id}
                    type="button"
                    className={`w-10 h-10 rounded-md flex items-center justify-center transition-all ${
                      selectedSize === size.id
                        ? "bg-blue-600 dark:bg-blue-700 text-white"
                        : "bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-white hover:bg-gray-200 dark:hover:bg-gray-700"
                    }`}
                    onClick={() => setSelectedSize(size.id)}
                  >
                    {size.name}
                  </button>
                ))}
              </div>
              <p className="text-sm text-gray-500 dark:text-gray-400 mt-2">
                Đã chọn: <span className="font-medium text-gray-900 dark:text-white uppercase">{selectedSize}</span>
              </p>
            </div>

            {/* Product Actions */}
            <div className="mb-8">
              <AddToCartButton product={product} />

              <div className="flex gap-3 mt-4">
                <Button
                  variant="outline"
                  size="icon"
                  className="rounded-full hover:bg-red-50 dark:hover:bg-red-900/20 hover:text-red-600 dark:hover:text-red-400 transition-colors"
                  onClick={handleAddToWishlist}
                >
                  <Heart className="h-5 w-5 text-red-500" />
                  <span className="sr-only">Thêm vào yêu thích</span>
                </Button>
                <Button
                  variant="outline"
                  size="icon"
                  className="rounded-full hover:bg-blue-50 dark:hover:bg-blue-900/20 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                  onClick={handleShare}
                >
                  <Share2 className="h-5 w-5 text-blue-500" />
                  <span className="sr-only">Chia sẻ</span>
                </Button>
              </div>
            </div>

            {/* Product Features */}
            <div className="bg-blue-50 dark:bg-blue-900/20 rounded-xl p-4 mb-8">
              <h3 className="font-semibold mb-3">Đặc điểm nổi bật</h3>
              <ul className="space-y-2">
                <li className="flex items-center gap-2">
                  <Award className="h-5 w-5 text-blue-600 dark:text-blue-400 flex-shrink-0" />
                  <span>Sản phẩm chính hãng 100%</span>
                </li>
                <li className="flex items-center gap-2">
                  <Package className="h-5 w-5 text-blue-600 dark:text-blue-400 flex-shrink-0" />
                  <span>Đóng gói cẩn thận, an toàn</span>
                </li>
                <li className="flex items-center gap-2">
                  <RefreshCw className="h-5 w-5 text-blue-600 dark:text-blue-400 flex-shrink-0" />
                  <span>Đổi trả trong vòng 30 ngày</span>
                </li>
                <li className="flex items-center gap-2">
                  <Truck className="h-5 w-5 text-blue-600 dark:text-blue-400 flex-shrink-0" />
                  <span>Giao hàng nhanh toàn quốc</span>
                </li>
              </ul>
            </div>

            {/* Delivery Info */}
            <div className="border dark:border-gray-700 rounded-xl p-4">
              <h3 className="font-semibold mb-3">Thông tin giao hàng</h3>
              <div className="flex items-center gap-2 mb-2">
                <Truck className="h-5 w-5 text-green-600 dark:text-green-400" />
                <span className="text-green-600 dark:text-green-400 font-medium">Giao hàng miễn phí</span>
              </div>
              <p className="text-sm text-gray-600 dark:text-gray-300">Giao hàng miễn phí cho đơn hàng từ 500.000đ</p>
            </div>
          </motion.div>
        </div>

        {/* Product Details Tabs */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.4 }}
          className="mt-16"
        >
          <Tabs defaultValue="description" value={activeTab} onValueChange={setActiveTab} className="w-full">
            <TabsList className="grid w-full grid-cols-3 lg:w-auto lg:inline-flex">
              <TabsTrigger value="description">Mô tả sản phẩm</TabsTrigger>
              <TabsTrigger value="specifications">Thông số kỹ thuật</TabsTrigger>
              <TabsTrigger value="reviews" id="reviews">
                Đánh giá (124)
              </TabsTrigger>
            </TabsList>
            <div className="mt-6 bg-white dark:bg-gray-800 rounded-xl border dark:border-gray-700 p-6">
              <TabsContent value="description" className="space-y-4">
                <h3 className="text-xl font-bold">Chi tiết sản phẩm</h3>
                <p className="text-gray-700 dark:text-gray-300 leading-relaxed">{product.description}</p>
                <p className="text-gray-700 dark:text-gray-300 leading-relaxed">
                  Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam auctor, nisl eget ultricies tincidunt,
                  nisl nisl aliquam nisl, eget aliquam nisl nisl eget nisl. Nullam auctor, nisl eget ultricies
                  tincidunt, nisl nisl aliquam nisl, eget aliquam nisl nisl eget nisl.
                </p>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mt-8">
                  <Image
                    src={`/default-product.svg`}
                    alt="Product detail 1"
                    width={600}
                    height={400}
                    className="rounded-lg"
                  />
                  <Image
                    src={`/default-product-blue.svg`}
                    alt="Product detail 2"
                    width={600}
                    height={400}
                    className="rounded-lg"
                  />
                </div>
              </TabsContent>

              <TabsContent value="specifications">
                <h3 className="text-xl font-bold mb-4">Thông số kỹ thuật</h3>
                <div className="overflow-hidden rounded-lg border dark:border-gray-700">
                  <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                    <tbody className="divide-y divide-gray-200 dark:divide-gray-700 bg-white dark:bg-gray-800">
                      {specifications.map((spec, index) => (
                        <tr
                          key={index}
                          className={index % 2 === 0 ? "bg-gray-50 dark:bg-gray-900/50" : "bg-white dark:bg-gray-800"}
                        >
                          <td className="px-6 py-4 text-sm font-medium text-gray-900 dark:text-white">{spec.name}</td>
                          <td className="px-6 py-4 text-sm text-gray-500 dark:text-gray-300">{spec.value}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </TabsContent>

              <TabsContent value="reviews">
                <ProductReviews productId={params.id as string} />
              </TabsContent>
            </div>
          </Tabs>
        </motion.div>

        {/* Recommendations */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.6 }}
          className="mt-16"
        >
          <h2 className="text-2xl font-bold mb-6">Sản phẩm tương tự</h2>
          <Suspense fallback={<RecommendationsSkeleton />}>
            <ProductRecommendations productId={params.id as string} />
          </Suspense>
        </motion.div>
      </div>
    </div>
  )
}

function ProductDetailSkeleton() {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-10">
        <div>
          <Skeleton className="aspect-square w-full rounded-xl" />
          <div className="grid grid-cols-4 gap-2 mt-2">
            {[...Array(4)].map((_, i) => (
              <Skeleton key={i} className="aspect-square w-full rounded-lg" />
            ))}
          </div>
        </div>

        <div>
          <Skeleton className="h-8 w-3/4 rounded-md mb-4" />
          <Skeleton className="h-6 w-1/2 rounded-md mb-6" />
          <Skeleton className="h-10 w-1/3 rounded-md mb-6" />
          <Skeleton className="h-20 w-full rounded-md mb-6" />
          <Skeleton className="h-12 w-full rounded-md mb-4" />
          <div className="flex gap-2">
            <Skeleton className="h-10 w-10 rounded-full" />
            <Skeleton className="h-10 w-10 rounded-full" />
          </div>
          <Skeleton className="h-32 w-full rounded-xl mt-6" />
        </div>
      </div>

      <div className="mt-16">
        <Skeleton className="h-10 w-96 rounded-md mb-6" />
        <Skeleton className="h-60 w-full rounded-xl" />
      </div>

      <div className="mt-16">
        <Skeleton className="h-8 w-64 rounded-md mb-6" />
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
          {[...Array(4)].map((_, i) => (
            <Skeleton key={i} className="h-80 rounded-xl" />
          ))}
        </div>
      </div>
    </div>
  )
}

function RecommendationsSkeleton() {
  return (
    <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
      {[...Array(4)].map((_, i) => (
        <div key={i} className="rounded-lg border p-4">
          <Skeleton className="h-40 w-full rounded-md" />
          <Skeleton className="mt-4 h-6 w-3/4" />
          <Skeleton className="mt-2 h-4 w-1/2" />
        </div>
      ))}
    </div>
  )
}
