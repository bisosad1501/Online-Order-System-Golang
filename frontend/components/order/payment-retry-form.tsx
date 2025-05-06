"use client"

import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { orderService } from "@/services/order-service"
import { Button } from "@/components/ui/button"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { useToast } from "@/hooks/use-toast"

const paymentSchema = z.object({
  payment_method: z.enum(["credit_card", "paypal"]),
  card_number: z
    .string()
    .optional()
    .refine((val) => !val || val.length >= 16, { message: "Card number must be at least 16 digits" }),
  expiry_month: z
    .string()
    .optional()
    .refine((val) => !val || (Number.parseInt(val) >= 1 && Number.parseInt(val) <= 12), {
      message: "Month must be between 1 and 12",
    }),
  expiry_year: z
    .string()
    .optional()
    .refine((val) => !val || Number.parseInt(val) >= new Date().getFullYear() % 100, {
      message: "Year must be current or future",
    }),
  cvv: z
    .string()
    .optional()
    .refine((val) => !val || (val.length >= 3 && val.length <= 4), { message: "CVV must be 3 or 4 digits" }),
})

type PaymentFormValues = z.infer<typeof paymentSchema>

export default function PaymentRetryForm({
  orderId,
  onSuccess,
  onCancel,
}: {
  orderId: string
  onSuccess: () => void
  onCancel: () => void
}) {
  const { toast } = useToast()
  const [isLoading, setIsLoading] = useState(false)

  const form = useForm<PaymentFormValues>({
    resolver: zodResolver(paymentSchema),
    defaultValues: {
      payment_method: "credit_card",
      card_number: "",
      expiry_month: "",
      expiry_year: "",
      cvv: "",
    },
  })

  const paymentMethod = form.watch("payment_method")

  async function onSubmit(data: PaymentFormValues) {
    setIsLoading(true)
    try {
      await orderService.retryPayment(orderId, data)
      toast({
        title: "Payment successful",
        description: "Your payment has been processed successfully",
      })
      onSuccess()
    } catch (error) {
      toast({
        title: "Payment failed",
        description: "There was an error processing your payment. Please try again.",
        variant: "destructive",
      })
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="rounded-lg border p-4">
      <h3 className="mb-4 text-lg font-semibold">Retry Payment</h3>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <FormField
            control={form.control}
            name="payment_method"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Payment Method</FormLabel>
                <Select onValueChange={field.onChange} defaultValue={field.value}>
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Select payment method" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value="credit_card">Credit Card</SelectItem>
                    <SelectItem value="paypal">PayPal</SelectItem>
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />

          {paymentMethod === "credit_card" && (
            <>
              <FormField
                control={form.control}
                name="card_number"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Card Number</FormLabel>
                    <FormControl>
                      <Input placeholder="1234 5678 9012 3456" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className="grid grid-cols-3 gap-4">
                <FormField
                  control={form.control}
                  name="expiry_month"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Expiry Month</FormLabel>
                      <FormControl>
                        <Input placeholder="MM" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="expiry_year"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Expiry Year</FormLabel>
                      <FormControl>
                        <Input placeholder="YY" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="cvv"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>CVV</FormLabel>
                      <FormControl>
                        <Input placeholder="123" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </>
          )}

          <div className="flex justify-end gap-2">
            <Button type="button" variant="outline" onClick={onCancel} disabled={isLoading}>
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading ? "Processing..." : "Process Payment"}
            </Button>
          </div>
        </form>
      </Form>
    </div>
  )
}
