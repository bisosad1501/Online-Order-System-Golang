"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import Image from "next/image"
import { motion } from "framer-motion"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { useAuth } from "@/context/auth-context"
import { Button } from "@/components/ui/button"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from "@/components/ui/card"
import { useToast } from "@/hooks/use-toast"
import { ThemeToggle } from "@/components/theme-toggle"
import { Facebook, Twitter, Github, Mail, Lock, User, ArrowRight } from "lucide-react"
import { authService } from "@/services/auth-service"

const loginSchema = z.object({
  email: z.string().email({ message: "Email không hợp lệ" }),
})

type LoginFormValues = z.infer<typeof loginSchema>

export default function LoginPage() {
  const { login } = useAuth()
  const router = useRouter()
  const { toast } = useToast()
  const [isLoading, setIsLoading] = useState(false)
  const [mounted, setMounted] = useState(false)
  const [redirectPath, setRedirectPath] = useState("/")

  useEffect(() => {
    setMounted(true)

    // Kiểm tra xem có tham số redirect trong URL không
    const searchParams = new URLSearchParams(window.location.search)
    const redirect = searchParams.get("redirect")
    if (redirect) {
      setRedirectPath(redirect)
    }
  }, [])

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
    },
    mode: "onChange"
  })

  async function onSubmit(data: LoginFormValues) {
    setIsLoading(true)
    try {
      await login(data)
      toast({
        title: "Đăng nhập thành công",
        description: "Chào mừng bạn quay trở lại!",
      })

      // Chuyển hướng người dùng đến trang trước đó hoặc trang chủ
      router.push(redirectPath)
    } catch (error) {
      toast({
        title: "Đăng nhập thất bại",
        description: "Vui lòng kiểm tra email của bạn",
        variant: "destructive",
      })
    } finally {
      setIsLoading(false)
    }
  }

  if (!mounted) {
    return null
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-blue-900 flex flex-col items-center justify-center p-4">
      <div className="absolute top-4 right-4">
        <ThemeToggle />
      </div>

      <Link
        href="/"
        className="absolute top-4 left-4 flex items-center space-x-2 text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 transition-colors"
      >
        <ArrowRight className="h-4 w-4 rotate-180" />
        <span>Quay lại trang chủ</span>
      </Link>

      <div className="w-full max-w-6xl grid grid-cols-1 md:grid-cols-2 gap-8 items-center">
        <motion.div
          initial={{ opacity: 0, x: -50 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.5 }}
          className="hidden md:block"
        >
          <div className="relative h-[500px] w-full">
            <Image
              src="/images/login-illustration.png"
              alt="Login Illustration"
              fill
              className="object-contain"
              priority
            />
          </div>
          <div className="text-center mt-8">
            <h2 className="text-2xl font-bold text-gray-800 dark:text-white mb-2">Chào mừng trở lại!</h2>
            <p className="text-gray-600 dark:text-gray-300">
              Đăng nhập để tiếp tục mua sắm và theo dõi đơn hàng của bạn.
            </p>
          </div>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.2 }}
        >
          <Card className="border-none shadow-xl bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm">
            <CardHeader className="space-y-1">
              <div className="flex justify-center mb-4">
                <div className="w-12 h-12 rounded-full bg-blue-100 dark:bg-blue-900 flex items-center justify-center">
                  <User className="h-6 w-6 text-blue-600 dark:text-blue-400" />
                </div>
              </div>
              <CardTitle className="text-2xl font-bold text-center">Đăng nhập</CardTitle>
              <CardDescription className="text-center">
                Nhập email của bạn để đăng nhập.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                  <FormField
                    control={form.control}
                    name="email"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Email</FormLabel>
                        <FormControl>
                          <div className="relative">
                            <Mail className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
                            <Input
                              placeholder="email@example.com"
                              className="pl-10 bg-gray-50 dark:bg-gray-700 border-gray-200 dark:border-gray-600"
                              {...field}
                            />
                          </div>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />


                  <Button
                    type="submit"
                    className="w-full bg-blue-600 hover:bg-blue-700 text-white"
                    disabled={isLoading}
                    size="lg"
                  >
                    {isLoading ? "Đang đăng nhập..." : "Đăng nhập"}
                  </Button>
                </form>
              </Form>

              <div className="relative my-6">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-gray-300 dark:border-gray-600"></div>
                </div>
                <div className="relative flex justify-center text-sm">
                  <span className="px-2 bg-white dark:bg-gray-800 text-gray-500 dark:text-gray-400">
                    Hoặc đăng nhập với
                  </span>
                </div>
              </div>

              <div className="grid grid-cols-3 gap-3">
                <Button variant="outline" className="bg-white dark:bg-gray-700">
                  <Facebook className="h-5 w-5 text-blue-600" />
                </Button>
                <Button variant="outline" className="bg-white dark:bg-gray-700">
                  <Twitter className="h-5 w-5 text-blue-400" />
                </Button>
                <Button variant="outline" className="bg-white dark:bg-gray-700">
                  <Github className="h-5 w-5" />
                </Button>
              </div>
            </CardContent>
            <CardFooter className="flex flex-col space-y-4">
              <div className="text-center text-sm">
                Chưa có tài khoản?{" "}
                <Link href="/register" className="text-blue-600 hover:underline font-medium">
                  Đăng ký ngay
                </Link>
              </div>
            </CardFooter>
          </Card>
        </motion.div>
      </div>
    </div>
  )
}
