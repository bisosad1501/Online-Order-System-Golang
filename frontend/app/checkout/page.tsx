"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { useCart } from "@/context/cart-context"
import { useAuth } from "@/context/auth-context"
import { orderService } from "@/services/order-service"
import CheckoutSummary from "@/components/checkout/checkout-summary"
import PaymentMethodSelector from "@/components/checkout/payment-method-selector"
import { Button } from "@/components/ui/button"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { useToast } from "@/hooks/use-toast"

const checkoutSchema = z.object({
  shipping_address: z.string().min(5, { message: "Địa chỉ giao hàng là bắt buộc" }),
  payment_method: z.enum(["card"], {
    required_error: "Vui lòng chọn phương thức thanh toán",
  }),
  card_number: z.string().refine(
    (val) => !val || /^[0-9]{13,16}$/.test(val.replace(/\s/g, '')),
    { message: "Số thẻ phải có từ 13-16 chữ số" }
  ).optional(),
  expiry_month: z.string().refine(
    (val) => !val || (/^(0[1-9]|1[0-2])$/.test(val)),
    { message: "Tháng không hợp lệ (01-12)" }
  ).optional(),
  expiry_year: z.string().refine(
    (val) => !val || (/^[0-9]{2}$/.test(val)),
    { message: "Năm không hợp lệ (2 chữ số)" }
  ).optional(),
  cvv: z.string().refine(
    (val) => !val || (/^[0-9]{3}$/.test(val)),
    { message: "CVV không hợp lệ (3 chữ số)" }
  ).optional(),
}).refine((data) => {
  // Nếu phương thức thanh toán là card, các trường thẻ tín dụng là bắt buộc
  if (data.payment_method === "card") {
    return !!data.card_number && !!data.expiry_month && !!data.expiry_year && !!data.cvv;
  }
  return true;
}, {
  message: "Vui lòng điền đầy đủ thông tin thẻ tín dụng",
  path: ["card_number"], // Hiển thị lỗi ở trường card_number
})

type CheckoutFormValues = z.infer<typeof checkoutSchema>

export default function CheckoutPage() {
  const { cart, clearCart } = useCart()
  const { user } = useAuth()
  const router = useRouter()
  const { toast } = useToast()
  const [isLoading, setIsLoading] = useState(false)

  const form = useForm<CheckoutFormValues>({
    resolver: zodResolver(checkoutSchema),
    defaultValues: {
      shipping_address: user?.address || "",
      payment_method: "card",
      card_number: "",
      expiry_month: "",
      expiry_year: "",
      cvv: "",
    },
  })

  const paymentMethod = form.watch("payment_method")

  async function onSubmit(data: CheckoutFormValues) {
    if (!cart || !user) return

    setIsLoading(true)
    try {
      console.log("Checkout form data:", data);
      console.log("User:", user);
      console.log("Cart:", cart);

      const orderData = {
        customer_id: user.id,
        shipping_address: data.shipping_address,
        items: cart.items.map((item) => ({
          product_id: item.product_id,
          quantity: item.quantity,
          price: item.price,
        })),
      }

      console.log("Submitting order data:", JSON.stringify(orderData, null, 2));
      const order = await orderService.createOrder(orderData)

      // Cấu trúc dữ liệu phù hợp với backend
      const cardNumber = data.card_number?.replace(/\s/g, '') || "";
      // Đảm bảo số thẻ không quá dài
      const formattedCardNumber = cardNumber.length > 16 ? cardNumber.substring(0, 16) : cardNumber;

      const paymentData = {
        order_id: order.id,
        amount: cart.items.reduce((total, item) => total + item.price * item.quantity, 0),
        payment_method: "card", // Luôn sử dụng "card" để phù hợp với backend
        card_number: formattedCardNumber,
        expiry_month: data.expiry_month || "",
        expiry_year: data.expiry_year || "",
        cvv: data.cvv || "",
      }

      console.log("Order created:", order);
      console.log("Submitting payment data:", JSON.stringify(paymentData, null, 2));
      await orderService.processPayment(paymentData)

      // Clear cart and redirect to order confirmation
      clearCart()
      toast({
        title: "Đặt hàng thành công",
        description: "Đơn hàng của bạn đã được đặt và đang được xử lý",
      })
      router.push(`/orders/${order.id}`)
    } catch (error: any) {
      console.error("Checkout error:", error);
      console.error("Error type:", typeof error);
      console.error("Error message:", error.message);
      console.error("Error stack:", error.stack);

      if (error.response) {
        console.error("Error response status:", error.response.status);
        console.error("Error response headers:", error.response.headers);
        console.error("Error response data:", error.response.data);

        const errorMessage = error.response.data?.error || "";
        console.error("Error message from response:", errorMessage);

        if (error.response.status === 500 && (errorMessage.includes("verify customer") || error.message.includes("verify customer"))) {
          toast({
            title: "Lỗi xác minh người dùng",
            description: "Không thể xác minh thông tin người dùng. Vui lòng đăng nhập lại hoặc liên hệ hỗ trợ.",
            variant: "destructive",
          });
        } else if (errorMessage.includes("inventory") || errorMessage.includes("out of stock") || errorMessage.includes("not available")) {
          toast({
            title: "Sản phẩm không có sẵn",
            description: "Một số sản phẩm trong đơn hàng của bạn hiện không có sẵn trong kho. Vui lòng kiểm tra lại giỏ hàng của bạn.",
            variant: "destructive",
          });
        } else if (errorMessage.includes("payment") || errorMessage.includes("card") || errorMessage.includes("value too long")) {
          // Xử lý lỗi thanh toán cụ thể
          let paymentErrorDescription = "Vui lòng kiểm tra thông tin thanh toán của bạn và thử lại.";

          if (errorMessage.includes("value too long")) {
            paymentErrorDescription = "Thông tin thẻ tín dụng không hợp lệ. Vui lòng kiểm tra lại số thẻ, ngày hết hạn và mã CVV.";

            // Hiển thị thông báo cụ thể hơn dựa trên lỗi
            if (errorMessage.includes("character varying(4)")) {
              paymentErrorDescription = "Số thẻ tín dụng quá dài. Vui lòng nhập tối đa 16 chữ số.";
            } else if (errorMessage.includes("character varying(2)")) {
              paymentErrorDescription = "Tháng hoặc năm hết hạn không hợp lệ. Vui lòng nhập đúng định dạng (MM/YY).";
            } else if (errorMessage.includes("character varying(3)")) {
              paymentErrorDescription = "Mã CVV không hợp lệ. Vui lòng nhập đúng 3 chữ số.";
            }
          } else if (errorMessage.includes("card number")) {
            paymentErrorDescription = "Số thẻ tín dụng không hợp lệ. Vui lòng kiểm tra lại.";
          } else if (errorMessage.includes("expiry")) {
            paymentErrorDescription = "Ngày hết hạn thẻ không hợp lệ. Vui lòng kiểm tra lại.";
          } else if (errorMessage.includes("cvv")) {
            paymentErrorDescription = "Mã CVV không hợp lệ. Vui lòng kiểm tra lại.";
          }

          toast({
            title: "Thanh toán không thành công",
            description: paymentErrorDescription,
            variant: "destructive",
          });
        } else {
          toast({
            title: "Đặt hàng không thành công",
            description: "Đã xảy ra lỗi khi xử lý đơn hàng của bạn. Vui lòng thử lại sau.",
            variant: "destructive",
          });
        }
      } else {
        console.error("Error during checkout (no response):", error);
        toast({
          title: "Đặt hàng không thành công",
          description: "Đã xảy ra lỗi khi xử lý đơn hàng của bạn. Vui lòng thử lại sau.",
          variant: "destructive",
        });
      }
    } finally {
      setIsLoading(false)
    }
  }

  if (!cart || cart.items.length === 0) {
    router.push("/cart")
    return null
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="mb-6 text-3xl font-bold">Thanh toán</h1>

      {/* Hiển thị các bước trong quy trình thanh toán */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          <div className="flex flex-col items-center">
            <div className="flex h-8 w-8 items-center justify-center rounded-full bg-blue-600 text-white">1</div>
            <span className="mt-2 text-sm font-medium text-blue-600">Giỏ hàng</span>
          </div>
          <div className="h-0.5 w-16 bg-blue-600"></div>
          <div className="flex flex-col items-center">
            <div className="flex h-8 w-8 items-center justify-center rounded-full bg-blue-600 text-white">2</div>
            <span className="mt-2 text-sm font-medium text-blue-600">Thông tin</span>
          </div>
          <div className="h-0.5 w-16 bg-blue-600"></div>
          <div className="flex flex-col items-center">
            <div className="flex h-8 w-8 items-center justify-center rounded-full bg-blue-600 text-white">3</div>
            <span className="mt-2 text-sm font-medium text-blue-600">Thanh toán</span>
          </div>
          <div className="h-0.5 w-16 bg-gray-200"></div>
          <div className="flex flex-col items-center">
            <div className="flex h-8 w-8 items-center justify-center rounded-full bg-gray-200 text-gray-600">4</div>
            <span className="mt-2 text-sm font-medium text-gray-600">Xác nhận</span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
        <div className="lg:col-span-2">
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
              <div className="rounded-lg border p-6">
                <h2 className="mb-4 text-xl font-semibold">Thông tin giao hàng</h2>

                <FormField
                  control={form.control}
                  name="shipping_address"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Địa chỉ giao hàng</FormLabel>
                      <FormControl>
                        <Textarea placeholder="Nhập địa chỉ giao hàng đầy đủ của bạn" className="min-h-[100px]" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <div className="rounded-lg border p-6">
                <h2 className="mb-4 text-xl font-semibold">Phương thức thanh toán</h2>

                <FormField
                  control={form.control}
                  name="payment_method"
                  render={({ field }) => (
                    <FormItem>
                      <FormControl>
                        <PaymentMethodSelector value={field.value} onChange={field.onChange} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                {paymentMethod === "card" && (
                  <div className="mt-6 space-y-6 border-t pt-6">
                    <h3 className="text-lg font-medium">Thông tin thẻ tín dụng</h3>
                    <p className="text-sm text-muted-foreground mb-4">
                      Thông tin thẻ của bạn được bảo mật và mã hóa
                    </p>

                    <FormField
                      control={form.control}
                      name="card_number"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>
                            <span className="flex items-center gap-1">
                              Số thẻ
                              <span className="text-red-500">*</span>
                            </span>
                          </FormLabel>
                          <FormControl>
                            <Input
                              placeholder="4242 4242 4242 4242"
                              {...field}
                              className="font-mono"
                              maxLength={16}
                              required
                              onChange={(e) => {
                                // Chỉ cho phép nhập số
                                const value = e.target.value.replace(/\D/g, '');
                                // Giới hạn độ dài tối đa 16 chữ số
                                const truncated = value.substring(0, 16);
                                field.onChange(truncated);
                              }}
                            />
                          </FormControl>
                          <FormMessage />
                          <p className="text-xs text-muted-foreground mt-1">Nhập số thẻ không có khoảng trắng</p>
                        </FormItem>
                      )}
                    />

                    <div className="grid grid-cols-3 gap-4">
                      <FormField
                        control={form.control}
                        name="expiry_month"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>
                              <span className="flex items-center gap-1">
                                Tháng hết hạn
                                <span className="text-red-500">*</span>
                              </span>
                            </FormLabel>
                            <FormControl>
                              <Input
                                placeholder="MM"
                                {...field}
                                maxLength={2}
                                className="font-mono"
                                required
                              />
                            </FormControl>
                            <FormMessage />
                            <p className="text-xs text-muted-foreground mt-1">01-12</p>
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={form.control}
                        name="expiry_year"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>
                              <span className="flex items-center gap-1">
                                Năm hết hạn
                                <span className="text-red-500">*</span>
                              </span>
                            </FormLabel>
                            <FormControl>
                              <Input
                                placeholder="YY"
                                {...field}
                                maxLength={2}
                                className="font-mono"
                                required
                              />
                            </FormControl>
                            <FormMessage />
                            <p className="text-xs text-muted-foreground mt-1">Chỉ nhập 2 số cuối (VD: 25 cho 2025)</p>
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={form.control}
                        name="cvv"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>
                              <span className="flex items-center gap-1">
                                CVV
                                <span className="text-red-500">*</span>
                              </span>
                            </FormLabel>
                            <FormControl>
                              <Input
                                placeholder="123"
                                {...field}
                                maxLength={3}
                                className="font-mono"
                                type="password"
                                required
                              />
                            </FormControl>
                            <FormMessage />
                            <p className="text-xs text-muted-foreground mt-1">3 chữ số ở mặt sau thẻ</p>
                          </FormItem>
                        )}
                      />
                    </div>
                  </div>
                )}
              </div>



              <div className="space-y-4">
                <div className="rounded-lg border border-gray-200 bg-gray-50 p-4">
                  <div className="flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-green-600">
                      <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10" />
                    </svg>
                    <p className="text-xs text-gray-600">
                      Thông tin thanh toán của bạn được bảo mật và mã hóa
                    </p>
                  </div>
                </div>

                <Button type="submit" className="w-full h-12 text-base" disabled={isLoading}>
                  {isLoading ? (
                    <span className="flex items-center gap-2">
                      <svg className="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      Đang xử lý...
                    </span>
                  ) : (
                    <span className="flex items-center gap-2">
                      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <rect width="20" height="14" x="2" y="5" rx="2" />
                        <line x1="2" x2="22" y1="10" y2="10" />
                      </svg>
                      Hoàn tất đặt hàng
                    </span>
                  )}
                </Button>
              </div>
            </form>
          </Form>
        </div>

        <div>
          <CheckoutSummary items={cart.items} />
        </div>
      </div>
    </div>
  )
}
