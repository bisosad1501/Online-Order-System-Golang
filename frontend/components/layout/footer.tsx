import Link from "next/link"
import { Package } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Facebook, Twitter, Instagram, Youtube, MapPin, Phone, Mail } from "lucide-react"

export default function Footer() {
  return (
    <footer className="bg-gray-900 text-white">
      <div className="container mx-auto px-4 py-12">
        {/* Newsletter Section */}
        <div className="mb-12 rounded-xl bg-gradient-to-r from-blue-600 to-purple-600 p-8 shadow-xl">
          <div className="flex flex-col md:flex-row items-center justify-between gap-6">
            <div className="max-w-md">
              <h3 className="text-2xl font-bold mb-2">Đăng ký nhận thông tin</h3>
              <p className="text-white/80">Nhận thông báo về sản phẩm mới và khuyến mãi đặc biệt.</p>
            </div>
            <div className="flex w-full max-w-md flex-col sm:flex-row gap-2">
              <Input
                type="email"
                placeholder="Email của bạn"
                className="bg-white/20 border-white/30 text-white placeholder:text-white/60 focus:border-white"
              />
              <Button className="bg-white text-blue-600 hover:bg-blue-50 whitespace-nowrap">Đăng ký ngay</Button>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 gap-8 md:grid-cols-4 lg:grid-cols-5">
          <div className="lg:col-span-2">
            <Link href="/" className="mb-6 flex items-center space-x-2">
              <Package className="h-8 w-8 text-blue-400" />
              <span className="font-bold text-2xl">Cửa hàng trực tuyến</span>
            </Link>
            <p className="mb-6 text-gray-400 max-w-md">
              Cửa hàng trực tuyến đáp ứng mọi nhu cầu của bạn. Sản phẩm chất lượng, giao hàng nhanh chóng và dịch vụ
              khách hàng tuyệt vời.
            </p>
            <div className="flex space-x-4">
              <Link href="#" className="rounded-full bg-gray-800 p-2 hover:bg-blue-600 transition-colors">
                <Facebook className="h-5 w-5" />
                <span className="sr-only">Facebook</span>
              </Link>
              <Link href="#" className="rounded-full bg-gray-800 p-2 hover:bg-blue-400 transition-colors">
                <Twitter className="h-5 w-5" />
                <span className="sr-only">Twitter</span>
              </Link>
              <Link href="#" className="rounded-full bg-gray-800 p-2 hover:bg-pink-600 transition-colors">
                <Instagram className="h-5 w-5" />
                <span className="sr-only">Instagram</span>
              </Link>
              <Link href="#" className="rounded-full bg-gray-800 p-2 hover:bg-red-600 transition-colors">
                <Youtube className="h-5 w-5" />
                <span className="sr-only">Youtube</span>
              </Link>
            </div>
          </div>

          <div>
            <h3 className="mb-6 text-lg font-semibold">Mua sắm</h3>
            <ul className="space-y-4 text-gray-400">
              <li>
                <Link href="/products" className="hover:text-blue-400 transition-colors">
                  Tất cả sản phẩm
                </Link>
              </li>
              <li>
                <Link href="/products?category=electronics" className="hover:text-blue-400 transition-colors">
                  Điện tử
                </Link>
              </li>
              <li>
                <Link href="/products?category=clothing" className="hover:text-blue-400 transition-colors">
                  Thời trang
                </Link>
              </li>
              <li>
                <Link href="/products?category=home" className="hover:text-blue-400 transition-colors">
                  Nhà cửa & Đời sống
                </Link>
              </li>
            </ul>
          </div>

          <div>
            <h3 className="mb-6 text-lg font-semibold">Tài khoản</h3>
            <ul className="space-y-4 text-gray-400">
              <li>
                <Link href="/profile" className="hover:text-blue-400 transition-colors">
                  Tài khoản của tôi
                </Link>
              </li>
              <li>
                <Link href="/orders" className="hover:text-blue-400 transition-colors">
                  Đơn hàng
                </Link>
              </li>
              <li>
                <Link href="/cart" className="hover:text-blue-400 transition-colors">
                  Giỏ hàng
                </Link>
              </li>
              <li>
                <Link href="/wishlist" className="hover:text-blue-400 transition-colors">
                  Danh sách yêu thích
                </Link>
              </li>
            </ul>
          </div>

          <div>
            <h3 className="mb-6 text-lg font-semibold">Liên hệ</h3>
            <ul className="space-y-4 text-gray-400">
              <li className="flex items-start">
                <MapPin className="mr-3 h-5 w-5 text-blue-400 flex-shrink-0" />
                <span>123 Đường Lê Lợi, Quận 1, TP. Hồ Chí Minh</span>
              </li>
              <li className="flex items-center">
                <Phone className="mr-3 h-5 w-5 text-blue-400 flex-shrink-0" />
                <span>+84 123 456 789</span>
              </li>
              <li className="flex items-center">
                <Mail className="mr-3 h-5 w-5 text-blue-400 flex-shrink-0" />
                <span>info@cuahangonline.com</span>
              </li>
            </ul>
          </div>
        </div>

        <div className="mt-12 border-t border-gray-800 pt-8">
          <div className="flex flex-col md:flex-row items-center justify-between gap-4">
            <p className="text-sm text-gray-400">
              &copy; {new Date().getFullYear()} Cửa hàng trực tuyến. Đã đăng ký bản quyền.
            </p>
            <div className="flex items-center space-x-6">
              <Link href="#" className="text-sm text-gray-400 hover:text-blue-400 transition-colors">
                Chính sách bảo mật
              </Link>
              <Link href="#" className="text-sm text-gray-400 hover:text-blue-400 transition-colors">
                Điều khoản dịch vụ
              </Link>
              <Link href="#" className="text-sm text-gray-400 hover:text-blue-400 transition-colors">
                Chính sách đổi trả
              </Link>
            </div>
          </div>
        </div>
      </div>
    </footer>
  )
}
