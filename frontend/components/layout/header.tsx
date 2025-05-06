"use client"

import { useState, useEffect } from "react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { useAuth } from "@/context/auth-context"
import { useCart } from "@/context/cart-context"
import { useNotifications } from "@/context/notification-context"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet"
import { Badge } from "@/components/ui/badge"
import { ThemeToggle } from "@/components/theme-toggle"
import { cn } from "@/lib/utils"
import { Bell, Heart, LogOut, Menu, Package, Search, ShoppingCart, User } from "lucide-react"

export default function Header() {
  const pathname = usePathname()
  const { user, isAuthenticated, logout } = useAuth()
  const { cart } = useCart()
  const { notifications } = useNotifications()
  const [isScrolled, setIsScrolled] = useState(false)
  const [showSearchBar, setShowSearchBar] = useState(false)

  const unreadNotifications = notifications.filter((notification) => notification.status === "UNREAD").length

  const cartItemCount = cart && cart.items ? cart.items.reduce((total, item) => total + item.quantity, 0) : 0

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 10)
    }

    window.addEventListener("scroll", handleScroll)
    return () => window.removeEventListener("scroll", handleScroll)
  }, [])

  const navigation = [
    { name: "Trang chủ", href: "/" },
    { name: "Sản phẩm", href: "/products" },
    { name: "Khuyến mãi", href: "/promotions" },
    { name: "Về chúng tôi", href: "/about" },
    { name: "Liên hệ", href: "/contact" },
  ]

  return (
    <header
      className={cn(
        "sticky top-0 z-50 w-full backdrop-blur transition-all duration-300",
        isScrolled ? "bg-background/90 shadow-md" : "bg-transparent",
      )}
    >
      <div className="container flex h-20 items-center px-4">
        <Sheet>
          <SheetTrigger asChild>
            <Button variant="ghost" size="icon" className="md:hidden">
              <Menu className="h-6 w-6" />
              <span className="sr-only">Mở menu</span>
            </Button>
          </SheetTrigger>
          <SheetContent side="left" className="w-[300px]">
            <div className="flex flex-col gap-6 py-6">
              <Link href="/" className="flex items-center space-x-2">
                <Package className="h-6 w-6 text-blue-600" />
                <span className="font-bold text-xl">Cửa hàng trực tuyến</span>
              </Link>

              <div className="flex flex-col gap-4">
                {navigation.map((item) => (
                  <Link
                    key={item.name}
                    href={item.href}
                    className={cn(
                      "text-lg font-medium transition-colors hover:text-blue-600",
                      pathname === item.href ? "text-blue-600" : "text-foreground",
                    )}
                  >
                    {item.name}
                  </Link>
                ))}
              </div>

              <div className="border-t pt-4">
                {isAuthenticated && (
                  <>
                    <div className="mb-4 flex items-center">
                      <div className="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center mr-3">
                        <User className="h-5 w-5 text-blue-600" />
                      </div>
                      <div>
                        <p className="font-medium">
                          {user?.first_name} {user?.last_name}
                        </p>
                        <p className="text-sm text-muted-foreground">{user?.email}</p>
                      </div>
                    </div>
                    <div className="flex flex-col gap-2">
                      <Link
                        href="/profile"
                        className="flex items-center text-foreground hover:text-blue-600 transition-colors"
                      >
                        <User className="mr-2 h-4 w-4" />
                        Hồ sơ
                      </Link>
                      <Link
                        href="/orders"
                        className="flex items-center text-foreground hover:text-blue-600 transition-colors"
                      >
                        <Package className="mr-2 h-4 w-4" />
                        Đơn hàng
                      </Link>
                      <Link
                        href="/notifications"
                        className="flex items-center text-foreground hover:text-blue-600 transition-colors"
                      >
                        <Bell className="mr-2 h-4 w-4" />
                        Thông báo
                        {unreadNotifications > 0 && (
                          <Badge variant="destructive" className="ml-2 h-5 w-5 rounded-full p-0 text-xs">
                            {unreadNotifications}
                          </Badge>
                        )}
                      </Link>
                      <button
                        onClick={logout}
                        className="flex items-center text-red-500 hover:text-red-600 transition-colors"
                      >
                        <LogOut className="mr-2 h-4 w-4" />
                        Đăng xuất
                      </button>
                    </div>
                  </>
                )}
                {!isAuthenticated && (
                  <div className="flex flex-col gap-2">
                    <Link
                      href="/login"
                      className="w-full rounded-lg bg-white border border-blue-600 text-blue-600 py-2 px-4 text-center font-medium hover:bg-blue-50 transition-colors"
                    >
                      Đăng nhập
                    </Link>
                    <Link
                      href="/register"
                      className="w-full rounded-lg bg-blue-600 text-white py-2 px-4 text-center font-medium hover:bg-blue-700 transition-colors"
                    >
                      Đăng ký
                    </Link>
                  </div>
                )}
              </div>
            </div>
          </SheetContent>
        </Sheet>

        <Link href="/" className="mr-6 flex items-center space-x-2">
          <Package className="h-7 w-7 text-blue-600" />
          <span className="hidden font-bold text-xl sm:inline-block">Cửa hàng trực tuyến</span>
        </Link>

        <nav className="hidden md:flex md:flex-1 md:items-center md:justify-between">
          <div className="flex items-center space-x-6">
            {navigation.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                  "relative font-medium transition-colors hover:text-blue-600",
                  pathname === item.href
                    ? "text-blue-600 after:absolute after:bottom-[-20px] after:left-0 after:h-1 after:w-full after:bg-blue-600"
                    : "text-foreground",
                )}
              >
                {item.name}
              </Link>
            ))}
          </div>
        </nav>

        <div className="flex flex-1 items-center justify-end space-x-4">
          <ThemeToggle />

          <Button
            variant="ghost"
            size="icon"
            onClick={() => setShowSearchBar(!showSearchBar)}
            className="rounded-full hover:bg-blue-100 hover:text-blue-600 transition-colors"
          >
            <Search className="h-5 w-5" />
            <span className="sr-only">Tìm kiếm</span>
          </Button>

          <Link href="/wishlist" className="hidden sm:block">
            <Button
              variant="ghost"
              size="icon"
              className="rounded-full hover:bg-blue-100 hover:text-blue-600 transition-colors"
            >
              <Heart className="h-5 w-5" />
              <span className="sr-only">Danh sách yêu thích</span>
            </Button>
          </Link>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className="rounded-full hover:bg-blue-100 hover:text-blue-600 transition-colors relative"
              >
                <ShoppingCart className="h-5 w-5" />
                {cartItemCount > 0 && (
                  <Badge variant="destructive" className="absolute -right-1 -top-1 h-5 w-5 rounded-full p-0 text-xs">
                    {cartItemCount}
                  </Badge>
                )}
                <span className="sr-only">Giỏ hàng</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-80">
              <DropdownMenuLabel className="font-normal">
                <div className="flex justify-between items-center">
                  <p className="text-sm font-medium leading-none">Giỏ hàng của bạn</p>
                  <p className="text-xs text-muted-foreground">{cartItemCount} sản phẩm</p>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />

              {!cart || !cart.items || cart.items.length === 0 ? (
                <div className="py-6 text-center">
                  <ShoppingCart className="h-10 w-10 text-muted-foreground mx-auto mb-2" />
                  <p className="text-sm text-muted-foreground">Giỏ hàng của bạn đang trống</p>
                </div>
              ) : (
                <>
                  <div className="max-h-[300px] overflow-auto py-1">
                    {cart.items.slice(0, 3).map((item) => (
                      <div key={item.id} className="flex items-center gap-3 px-2 py-2 hover:bg-accent">
                        <div className="h-12 w-12 rounded bg-secondary flex items-center justify-center text-xs overflow-hidden">
                          {item.product_id.substring(0, 4)}
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="text-sm font-medium truncate">Sản phẩm {item.product_id.substring(0, 8)}</p>
                          <p className="text-xs text-muted-foreground">{item.quantity} x ${item.price.toFixed(2)}</p>
                        </div>
                        <p className="text-sm font-medium">${(item.quantity * item.price).toFixed(2)}</p>
                      </div>
                    ))}

                    {cart.items.length > 3 && (
                      <div className="px-2 py-1 text-center text-xs text-muted-foreground">
                        + {cart.items.length - 3} sản phẩm khác
                      </div>
                    )}
                  </div>

                  <DropdownMenuSeparator />
                  <div className="p-2">
                    <div className="flex justify-between py-1 text-sm">
                      <span>Tổng cộng:</span>
                      <span className="font-medium">
                        ${cart.items?.reduce((total, item) => total + item.price * item.quantity, 0).toFixed(2)}
                      </span>
                    </div>
                    <div className="flex gap-2 mt-2">
                      <Button asChild variant="outline" className="flex-1 h-9">
                        <Link href="/cart">Xem giỏ hàng</Link>
                      </Button>
                      <Button asChild className="flex-1 h-9 bg-blue-600 hover:bg-blue-700">
                        <Link href="/checkout">Thanh toán</Link>
                      </Button>
                    </div>
                  </div>
                </>
              )}
            </DropdownMenuContent>
          </DropdownMenu>

          {isAuthenticated && (
            <>
              <Link href="/notifications" className="relative">
                <Button
                  variant="ghost"
                  size="icon"
                  className="rounded-full hover:bg-blue-100 hover:text-blue-600 transition-colors"
                >
                  <Bell className="h-5 w-5" />
                  {unreadNotifications > 0 && (
                    <Badge variant="destructive" className="absolute -right-1 -top-1 h-5 w-5 rounded-full p-0 text-xs">
                      {unreadNotifications}
                    </Badge>
                  )}
                  <span className="sr-only">Thông báo</span>
                </Button>
              </Link>

              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="rounded-full hover:bg-blue-100 hover:text-blue-600 transition-colors"
                  >
                    <User className="h-5 w-5" />
                    <span className="sr-only">Menu người dùng</span>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-56">
                  <DropdownMenuLabel className="font-normal">
                    <div className="flex flex-col space-y-1">
                      <p className="text-sm font-medium leading-none">
                        {user?.first_name} {user?.last_name}
                      </p>
                      <p className="text-xs leading-none text-muted-foreground">{user?.email}</p>
                    </div>
                  </DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem asChild>
                    <Link href="/profile" className="cursor-pointer">
                      <User className="mr-2 h-4 w-4" /> Hồ sơ
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuItem asChild>
                    <Link href="/orders" className="cursor-pointer">
                      <Package className="mr-2 h-4 w-4" /> Đơn hàng
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={logout} className="text-red-500 focus:text-red-500 cursor-pointer">
                    <LogOut className="mr-2 h-4 w-4" />
                    Đăng xuất
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </>
          )}

          {!isAuthenticated && (
            <div className="hidden md:flex items-center space-x-2">
              <Link href="/login">
                <Button
                  variant="outline"
                  size="sm"
                  className="rounded-full border-blue-600 text-blue-600 hover:bg-blue-50 hover:text-blue-700"
                >
                  Đăng nhập
                </Button>
              </Link>
              <Link href="/register">
                <Button size="sm" className="rounded-full bg-blue-600 hover:bg-blue-700">
                  Đăng ký
                </Button>
              </Link>
            </div>
          )}
        </div>
      </div>

      {/* Search bar overlay */}
      {showSearchBar && (
        <div className="absolute inset-x-0 top-full bg-background/95 backdrop-blur-sm shadow-lg p-4 animate-fade-in">
          <div className="container mx-auto">
            <div className="relative">
              <input
                type="text"
                placeholder="Tìm kiếm sản phẩm..."
                className="w-full px-4 py-3 pr-10 rounded-lg border border-input focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                autoFocus
              />
              <Button
                variant="ghost"
                size="icon"
                className="absolute right-2 top-1/2 transform -translate-y-1/2"
                onClick={() => setShowSearchBar(false)}
              >
                <Search className="h-5 w-5" />
              </Button>
            </div>
          </div>
        </div>
      )}
    </header>
  )
}
