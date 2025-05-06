"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import Image from "next/image"
import { motion } from "framer-motion"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { authService } from "@/services/auth-service"
import { Button } from "@/components/ui/button"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from "@/components/ui/card"
import { useToast } from "@/hooks/use-toast"
import { ThemeToggle } from "@/components/theme-toggle"
import { ArrowRight, Mail, User, Phone, Home, CheckCircle } from "lucide-react"

// Schema cho bước 1 (thông tin cơ bản)
const step1Schema = z.object({
  email: z.string().email({ message: "Email không hợp lệ" }),
  first_name: z.string().min(2, { message: "Họ là bắt buộc" }),
  last_name: z.string().min(2, { message: "Tên là bắt buộc" }),
  phone: z.string().optional(),
  address: z.string().optional(),
})

// Schema cho bước 2 (thông tin liên hệ)
const step2Schema = z.object({
  email: z.string().email({ message: "Email không hợp lệ" }),
  first_name: z.string().min(2, { message: "Họ là bắt buộc" }),
  last_name: z.string().min(2, { message: "Tên là bắt buộc" }),
  phone: z.string().min(10, { message: "Số điện thoại hợp lệ là bắt buộc" }),
  address: z.string().min(5, { message: "Địa chỉ là bắt buộc" }),
})

// Schema chung cho toàn bộ form
const registerSchema = step2Schema

type RegisterFormValues = z.infer<typeof registerSchema>

export default function RegisterPage() {
  const router = useRouter()
  const { toast } = useToast()
  const [isLoading, setIsLoading] = useState(false)
  const [mounted, setMounted] = useState(false)
  const [step, setStep] = useState(1)

  useEffect(() => {
    setMounted(true)
  }, [])

  // Sử dụng schema dựa trên bước hiện tại
  const form = useForm<RegisterFormValues>({
    resolver: zodResolver(step === 1 ? step1Schema : step2Schema),
    defaultValues: {
      email: "",
      first_name: "",
      last_name: "",
      phone: "",
      address: "",
    },
    mode: "onChange"
  })

  // Xử lý khi chuyển từ bước 1 sang bước 2
  const handleContinue = () => {
    const { email, first_name, last_name } = form.getValues();

    // Kiểm tra các trường ở bước 1
    if (!email || !first_name || !last_name) {
      toast({
        title: "Thiếu thông tin",
        description: "Vui lòng điền đầy đủ thông tin trước khi tiếp tục",
        variant: "destructive",
      });
      return;
    }

    // Kiểm tra email hợp lệ
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      toast({
        title: "Email không hợp lệ",
        description: "Vui lòng nhập đúng định dạng email",
        variant: "destructive",
      });
      return;
    }

    // Nếu thông tin hợp lệ, chuyển sang bước 2
    setStep(2);
  };

  // Xử lý khi submit form
  async function onSubmit(data: RegisterFormValues) {
    // Nếu đang ở bước 1, chỉ chuyển sang bước 2
    if (step === 1) {
      handleContinue();
      return;
    }

    // Nếu đang ở bước 2, tiến hành đăng ký
    setIsLoading(true);
    try {
      const response = await authService.register(data);

      // Lưu ID vào state để hiển thị trong modal
      setRegisteredUserId(response.id);

      // Hiển thị modal thông báo thành công với ID
      setShowSuccessModal(true);

    } catch (error) {
      toast({
        title: "Đăng ký thất bại",
        description: "Vui lòng kiểm tra thông tin và thử lại",
        variant: "destructive",
      });
      setIsLoading(false);
    }
  }

  // State để lưu ID người dùng đã đăng ký và hiển thị modal
  const [registeredUserId, setRegisteredUserId] = useState<string>("");
  const [showSuccessModal, setShowSuccessModal] = useState<boolean>(false);

  // Xử lý khi người dùng đóng modal và chuyển đến trang đăng nhập
  const handleCloseSuccessModal = () => {
    setShowSuccessModal(false);
    setIsLoading(false);
    router.push("/login");
  }

  if (!mounted) {
    return null
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-blue-900 flex flex-col items-center justify-center p-4">
      {/* Success Modal */}
      {showSuccessModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-gray-800 p-8 rounded-lg shadow-xl max-w-md w-full">
            <h2 className="text-2xl font-bold mb-4 text-center">Đăng ký thành công!</h2>
            <p className="mb-6 text-center">
              Tài khoản của bạn đã được tạo thành công. Vui lòng lưu lại ID của bạn để đăng nhập.
            </p>
            <div className="bg-blue-50 dark:bg-blue-900 p-4 rounded-md mb-6">
              <p className="text-center font-mono text-lg font-bold">{registeredUserId}</p>
            </div>
            <div className="flex justify-center">
              <Button
                onClick={handleCloseSuccessModal}
                className="bg-blue-600 hover:bg-blue-700 text-white"
              >
                Đi đến đăng nhập
              </Button>
            </div>
          </div>
        </div>
      )}

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
              src="/images/register-illustration.png"
              alt="Register Illustration"
              fill
              className="object-contain"
              priority
            />
          </div>
          <div className="text-center mt-8">
            <h2 className="text-2xl font-bold text-gray-800 dark:text-white mb-2">Tham gia cùng chúng tôi!</h2>
            <p className="text-gray-600 dark:text-gray-300">
              Tạo tài khoản để trải nghiệm mua sắm tuyệt vời và nhận nhiều ưu đãi đặc biệt.
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
              <CardTitle className="text-2xl font-bold text-center">Đăng ký tài khoản</CardTitle>
              <CardDescription className="text-center">
                Nhập thông tin của bạn để tạo tài khoản mới
                <span className="block mt-2 text-xs text-red-500">
                  Lưu ý: Hãy ghi nhớ ID của bạn sau khi đăng ký thành công để sử dụng cho việc đăng nhập.
                </span>
              </CardDescription>

              <div className="flex justify-center mt-4">
                <div className="flex items-center">
                  <div
                    className={`w-8 h-8 rounded-full flex items-center justify-center ${step >= 1 ? "bg-blue-600 text-white" : "bg-gray-200 dark:bg-gray-700 text-gray-500"}`}
                  >
                    1
                  </div>
                  <div className={`w-16 h-1 ${step >= 2 ? "bg-blue-600" : "bg-gray-200 dark:bg-gray-700"}`}></div>
                  <div
                    className={`w-8 h-8 rounded-full flex items-center justify-center ${step >= 2 ? "bg-blue-600 text-white" : "bg-gray-200 dark:bg-gray-700 text-gray-500"}`}
                  >
                    2
                  </div>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                  {step === 1 ? (
                    <>
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
                      <div className="grid grid-cols-2 gap-4">
                        <FormField
                          control={form.control}
                          name="first_name"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Họ</FormLabel>
                              <FormControl>
                                <div className="relative">
                                  <User className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
                                  <Input
                                    placeholder="Nguyễn"
                                    className="pl-10 bg-gray-50 dark:bg-gray-700 border-gray-200 dark:border-gray-600"
                                    {...field}
                                  />
                                </div>
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                        <FormField
                          control={form.control}
                          name="last_name"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Tên</FormLabel>
                              <FormControl>
                                <div className="relative">
                                  <User className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
                                  <Input
                                    placeholder="Văn A"
                                    className="pl-10 bg-gray-50 dark:bg-gray-700 border-gray-200 dark:border-gray-600"
                                    {...field}
                                  />
                                </div>
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>
                    </>
                  ) : (
                    <>
                      <FormField
                        control={form.control}
                        name="phone"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Số điện thoại</FormLabel>
                            <FormControl>
                              <div className="relative">
                                <Phone className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
                                <Input
                                  placeholder="+84 123 456 789"
                                  className="pl-10 bg-gray-50 dark:bg-gray-700 border-gray-200 dark:border-gray-600"
                                  {...field}
                                />
                              </div>
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={form.control}
                        name="address"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Địa chỉ</FormLabel>
                            <FormControl>
                              <div className="relative">
                                <Home className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
                                <Input
                                  placeholder="123 Đường ABC, Quận XYZ, TP. HCM"
                                  className="pl-10 bg-gray-50 dark:bg-gray-700 border-gray-200 dark:border-gray-600"
                                  {...field}
                                />
                              </div>
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </>
                  )}

                  <div className="flex justify-between mt-6">
                    {step === 2 && (
                      <Button
                        type="button"
                        variant="outline"
                        onClick={() => setStep(1)}
                        className="bg-white dark:bg-gray-700"
                      >
                        Quay lại
                      </Button>
                    )}
                    {step === 1 ? (
                      <Button
                        type="button"
                        onClick={handleContinue}
                        className="w-full bg-blue-600 hover:bg-blue-700 text-white"
                        size="lg"
                      >
                        Tiếp tục
                      </Button>
                    ) : (
                      <Button
                        type="submit"
                        className="ml-auto bg-blue-600 hover:bg-blue-700 text-white"
                        disabled={isLoading}
                        size="lg"
                      >
                        {isLoading ? "Đang đăng ký..." : "Đăng ký"}
                      </Button>
                    )}
                  </div>
                </form>
              </Form>

              <div className="mt-6 space-y-4">
                <div className="flex items-start space-x-2">
                  <CheckCircle className="h-5 w-5 text-green-500 mt-0.5" />
                  <p className="text-sm text-gray-600 dark:text-gray-300">
                    Nhận thông báo về sản phẩm mới và khuyến mãi đặc biệt
                  </p>
                </div>
                <div className="flex items-start space-x-2">
                  <CheckCircle className="h-5 w-5 text-green-500 mt-0.5" />
                  <p className="text-sm text-gray-600 dark:text-gray-300">
                    Theo dõi đơn hàng và lịch sử mua sắm của bạn
                  </p>
                </div>
                <div className="flex items-start space-x-2">
                  <CheckCircle className="h-5 w-5 text-green-500 mt-0.5" />
                  <p className="text-sm text-gray-600 dark:text-gray-300">Tích điểm thưởng và nhận ưu đãi đặc biệt</p>
                </div>
              </div>
            </CardContent>
            <CardFooter className="flex flex-col space-y-4">
              <div className="text-center text-sm">
                Đã có tài khoản?{" "}
                <Link href="/login" className="text-blue-600 hover:underline font-medium">
                  Đăng nhập
                </Link>
              </div>
            </CardFooter>
          </Card>
        </motion.div>
      </div>
    </div>
  )
}
