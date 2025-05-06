"use client"

import { Suspense } from "react"
import Image from "next/image"
import Link from "next/link"
import ProductList from "@/components/product/product-list"
import CategoryFilter from "@/components/product/category-filter"
import { Skeleton } from "@/components/ui/skeleton"
import { Button } from "@/components/ui/button"
import { ArrowRight, Truck, Shield, Clock, Check } from "lucide-react"
import HeroCarousel from "@/components/home/hero-carousel"
import BrandSlider from "@/components/home/brand-slider"
import ProductCard3D from "@/components/product/product-card-3d"
import PromoCountdown from "@/components/home/promo-countdown"
import FeatureShowcase from "@/components/home/feature-showcase"

export default function Home() {
  return (
    <div className="animate-fade-in">
      {/* Hero Section with Carousel */}
      <HeroCarousel />

      {/* Features Section */}
      <section className="relative z-10 bg-background py-12">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
            <div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-lg hover:shadow-xl transition-shadow duration-300 flex items-center group">
              <div className="rounded-full bg-blue-100 dark:bg-blue-900/50 p-4 mr-4 group-hover:bg-blue-200 dark:group-hover:bg-blue-800/50 transition-colors">
                <Truck className="h-6 w-6 text-blue-600 dark:text-blue-400" />
              </div>
              <div>
                <h3 className="font-bold text-lg">Giao hàng miễn phí</h3>
                <p className="text-gray-600 dark:text-gray-400">Cho đơn hàng từ 500.000đ</p>
              </div>
            </div>
            <div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-lg hover:shadow-xl transition-shadow duration-300 flex items-center group">
              <div className="rounded-full bg-purple-100 dark:bg-purple-900/50 p-4 mr-4 group-hover:bg-purple-200 dark:group-hover:bg-purple-800/50 transition-colors">
                <Shield className="h-6 w-6 text-purple-600 dark:text-purple-400" />
              </div>
              <div>
                <h3 className="font-bold text-lg">Bảo hành 12 tháng</h3>
                <p className="text-gray-600 dark:text-gray-400">Đổi trả dễ dàng</p>
              </div>
            </div>
            <div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-lg hover:shadow-xl transition-shadow duration-300 flex items-center group">
              <div className="rounded-full bg-green-100 dark:bg-green-900/50 p-4 mr-4 group-hover:bg-green-200 dark:group-hover:bg-green-800/50 transition-colors">
                <Check className="h-6 w-6 text-green-600 dark:text-green-400" />
              </div>
              <div>
                <h3 className="font-bold text-lg">Sản phẩm chính hãng</h3>
                <p className="text-gray-600 dark:text-gray-400">100% hàng chính hãng</p>
              </div>
            </div>
            <div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-lg hover:shadow-xl transition-shadow duration-300 flex items-center group">
              <div className="rounded-full bg-amber-100 dark:bg-amber-900/50 p-4 mr-4 group-hover:bg-amber-200 dark:group-hover:bg-amber-800/50 transition-colors">
                <Clock className="h-6 w-6 text-amber-600 dark:text-amber-400" />
              </div>
              <div>
                <h3 className="font-bold text-lg">Hỗ trợ 24/7</h3>
                <p className="text-gray-600 dark:text-gray-400">Luôn sẵn sàng phục vụ</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Brand Slider */}
      <section className="bg-gray-50 dark:bg-gray-900 py-12">
        <div className="container mx-auto px-4">
          <h2 className="text-2xl font-bold text-center mb-8">Thương hiệu nổi bật</h2>
          <BrandSlider />
        </div>
      </section>

      {/* Featured Products Section */}
      <section className="container mx-auto px-4 py-16">
        <div className="flex items-center justify-between mb-10">
          <div>
            <h2 className="text-3xl font-bold mb-2">Sản phẩm nổi bật</h2>
            <div className="w-20 h-1 bg-blue-600 rounded-full"></div>
          </div>
          <Link
            href="/products"
            className="group flex items-center text-blue-600 font-medium hover:text-blue-800 transition-colors"
          >
            Xem tất cả
            <ArrowRight className="ml-2 h-4 w-4 transform group-hover:translate-x-1 transition-transform" />
          </Link>
        </div>
        <Suspense fallback={<FeaturedProductsSkeleton />}>
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
            <ProductCard3D
              product={{
                id: "1",
                name: "Điện thoại thông minh X1",
                price: 12990000,
                image: "/images/placeholder.svg",
                category: "electronics",
                tags: ["new", "bestseller"],
              }}
            />
            <ProductCard3D
              product={{
                id: "2",
                name: "Laptop Ultrabook Pro",
                price: 24990000,
                image: "/images/placeholder.svg",
                category: "electronics",
                tags: ["premium"],
              }}
            />
            <ProductCard3D
              product={{
                id: "3",
                name: "Tai nghe không dây",
                price: 2990000,
                image: "/images/placeholder.svg",
                category: "electronics",
                tags: ["new"],
              }}
            />
            <ProductCard3D
              product={{
                id: "4",
                name: "Đồng hồ thông minh",
                price: 4990000,
                image: "/images/placeholder.svg",
                category: "electronics",
                tags: ["bestseller"],
              }}
            />
          </div>
        </Suspense>
      </section>

      {/* Promo Banner */}
      <section className="py-16 bg-gradient-to-r from-blue-600 to-purple-600 relative overflow-hidden">
        <div className="absolute inset-0 opacity-10">
          <svg className="w-full h-full" viewBox="0 0 100 100" preserveAspectRatio="none">
            <defs>
              <pattern id="grid" width="8" height="8" patternUnits="userSpaceOnUse">
                <path d="M 8 0 L 0 0 0 8" fill="none" stroke="white" strokeWidth="0.5" />
              </pattern>
            </defs>
            <rect width="100" height="100" fill="url(#grid)" />
          </svg>
        </div>

        <div className="container mx-auto px-4 relative z-10">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8 items-center">
            <div className="text-white">
              <h2 className="text-4xl font-bold mb-4">Ưu đãi mùa hè cực lớn</h2>
              <p className="text-xl mb-6 text-white/90">
                Giảm giá lên đến 50% cho tất cả sản phẩm điện tử. Nhanh tay kẻo lỡ!
              </p>
              <div className="mb-8">
                <PromoCountdown targetDate="2023-12-31T23:59:59" />
              </div>
              <Button size="lg" className="bg-white text-blue-600 hover:bg-blue-50">
                Mua ngay
              </Button>
            </div>
            <div className="relative">
              <div className="relative h-[400px] w-full">
                <Image
                  src="/images/promo-electronics.png"
                  alt="Electronics Promotion"
                  fill
                  className="object-contain"
                />
              </div>
              <div className="absolute top-0 right-0 bg-yellow-500 text-black font-bold py-2 px-4 rounded-full transform rotate-12">
                -50%
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Categories Banner */}
      <section className="bg-background py-16">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold mb-10 text-center">Danh mục sản phẩm</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <Link href="/products?category=electronics" className="group">
              <div className="relative h-80 rounded-2xl overflow-hidden shadow-lg hover:shadow-xl transition-shadow duration-300">
                <Image
                  src="/images/category-electronics.jpg"
                  alt="Electronics"
                  fill
                  className="object-cover group-hover:scale-110 transition-transform duration-700"
                />
                <div className="absolute inset-0 bg-gradient-to-t from-black/80 to-transparent"></div>
                <div className="absolute bottom-0 left-0 p-6">
                  <h3 className="text-2xl font-bold text-white mb-2">Điện tử</h3>
                  <p className="text-white/80 mb-4">Khám phá các thiết bị công nghệ mới nhất</p>
                  <span className="inline-flex items-center text-blue-400 font-medium group-hover:text-blue-300 transition-colors">
                    Khám phá ngay
                    <ArrowRight className="ml-2 h-4 w-4 transform group-hover:translate-x-1 transition-transform" />
                  </span>
                </div>
              </div>
            </Link>
            <Link href="/products?category=clothing" className="group">
              <div className="relative h-80 rounded-2xl overflow-hidden shadow-lg hover:shadow-xl transition-shadow duration-300">
                <Image
                  src="/images/category-fashion.jpg"
                  alt="Fashion"
                  fill
                  className="object-cover group-hover:scale-110 transition-transform duration-700"
                />
                <div className="absolute inset-0 bg-gradient-to-t from-black/80 to-transparent"></div>
                <div className="absolute bottom-0 left-0 p-6">
                  <h3 className="text-2xl font-bold text-white mb-2">Thời trang</h3>
                  <p className="text-white/80 mb-4">Phong cách mới cho mọi mùa</p>
                  <span className="inline-flex items-center text-blue-400 font-medium group-hover:text-blue-300 transition-colors">
                    Khám phá ngay
                    <ArrowRight className="ml-2 h-4 w-4 transform group-hover:translate-x-1 transition-transform" />
                  </span>
                </div>
              </div>
            </Link>
            <Link href="/products?category=home" className="group">
              <div className="relative h-80 rounded-2xl overflow-hidden shadow-lg hover:shadow-xl transition-shadow duration-300">
                <Image
                  src="/images/category-home.jpg"
                  alt="Home & Living"
                  fill
                  className="object-cover group-hover:scale-110 transition-transform duration-700"
                />
                <div className="absolute inset-0 bg-gradient-to-t from-black/80 to-transparent"></div>
                <div className="absolute bottom-0 left-0 p-6">
                  <h3 className="text-2xl font-bold text-white mb-2">Nhà cửa & Đời sống</h3>
                  <p className="text-white/80 mb-4">Nâng cấp không gian sống của bạn</p>
                  <span className="inline-flex items-center text-blue-400 font-medium group-hover:text-blue-300 transition-colors">
                    Khám phá ngay
                    <ArrowRight className="ml-2 h-4 w-4 transform group-hover:translate-x-1 transition-transform" />
                  </span>
                </div>
              </div>
            </Link>
          </div>
        </div>
      </section>

      {/* Feature Showcase */}
      <FeatureShowcase />

      {/* New Arrivals Section */}
      <section className="container mx-auto px-4 py-16">
        <div className="flex items-center justify-between mb-10">
          <div>
            <h2 className="text-3xl font-bold mb-2">Sản phẩm mới nhất</h2>
            <div className="w-20 h-1 bg-blue-600 rounded-full"></div>
          </div>
          <div className="flex items-center gap-4">
            <CategoryFilter />
          </div>
        </div>
        <Suspense fallback={<ProductListSkeleton />}>
          <ProductList limit={8} />
        </Suspense>
        <div className="mt-10 text-center">
          <Button
            asChild
            size="lg"
            className="bg-blue-600 hover:bg-blue-700 px-8 py-6 rounded-lg shadow-lg hover:shadow-blue-500/30 transition-all duration-300"
          >
            <Link href="/products">Xem tất cả sản phẩm</Link>
          </Button>
        </div>
      </section>

      {/* Testimonials Section */}
      <section className="bg-gradient-to-r from-blue-900 to-purple-900 py-16 text-white">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold mb-2 text-center">Khách hàng nói gì về chúng tôi</h2>
          <p className="text-center text-blue-200 mb-12">Trải nghiệm mua sắm từ khách hàng thực tế</p>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className="bg-white/10 backdrop-blur-lg rounded-xl p-6 shadow-lg hover:shadow-xl transition-shadow duration-300 hover:bg-white/20"
              >
                <div className="flex items-center mb-4">
                  <Image
                    src="/images/placeholder.svg"
                    alt={`Customer ${i}`}
                    width={48}
                    height={48}
                    className="rounded-full mr-4"
                  />
                  <div>
                    <h3 className="font-bold">Khách hàng {i}</h3>
                    <div className="flex">
                      {[...Array(5)].map((_, j) => (
                        <svg key={j} className="w-4 h-4 text-yellow-400 fill-current" viewBox="0 0 20 20">
                          <path d="M10 15l-5.878 3.09 1.123-6.545L.489 6.91l6.572-.955L10 0l2.939 5.955 6.572.955-4.756 4.635 1.123 6.545z" />
                        </svg>
                      ))}
                    </div>
                  </div>
                </div>
                <p className="italic text-white/90">
                  "Sản phẩm chất lượng cao, giao hàng nhanh và dịch vụ khách hàng tuyệt vời. Tôi rất hài lòng với trải
                  nghiệm mua sắm tại đây và chắc chắn sẽ quay lại."
                </p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Newsletter Section */}
      <section className="container mx-auto px-4 py-16">
        <div className="bg-gray-50 dark:bg-gray-800 rounded-2xl p-8 md:p-12 shadow-lg">
          <div className="max-w-3xl mx-auto text-center">
            <h2 className="text-3xl font-bold mb-4">Đăng ký nhận thông tin</h2>
            <p className="text-gray-600 dark:text-gray-300 mb-8">
              Nhận thông tin về sản phẩm mới, khuyến mãi đặc biệt và các sự kiện độc quyền.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 max-w-lg mx-auto">
              <input
                type="email"
                placeholder="Nhập email của bạn"
                className="flex-1 px-4 py-3 rounded-lg border border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
              />
              <Button className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg shadow-lg hover:shadow-blue-500/30 transition-all duration-300">
                Đăng ký
              </Button>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}

function FeaturedProductsSkeleton() {
  return (
    <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
      {[...Array(4)].map((_, i) => (
        <div key={i} className="rounded-xl border p-4 shadow-md">
          <Skeleton className="h-60 w-full rounded-lg" />
          <Skeleton className="mt-4 h-6 w-3/4" />
          <Skeleton className="mt-2 h-4 w-1/2" />
          <div className="mt-4 flex gap-2">
            <Skeleton className="h-10 flex-1 rounded-md" />
            <Skeleton className="h-10 w-10 rounded-md" />
          </div>
        </div>
      ))}
    </div>
  )
}

function ProductListSkeleton() {
  return (
    <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
      {[...Array(8)].map((_, i) => (
        <div key={i} className="rounded-xl border p-4 shadow-md">
          <Skeleton className="h-60 w-full rounded-lg" />
          <Skeleton className="mt-4 h-6 w-3/4" />
          <Skeleton className="mt-2 h-4 w-1/2" />
          <div className="mt-4 flex gap-2">
            <Skeleton className="h-10 flex-1 rounded-md" />
            <Skeleton className="h-10 w-10 rounded-md" />
          </div>
        </div>
      ))}
    </div>
  )
}
