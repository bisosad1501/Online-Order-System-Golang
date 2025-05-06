"use client"

import { CreditCard, Wallet } from "lucide-react"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { Label } from "@/components/ui/label"

export default function PaymentMethodSelector({
  value,
  onChange,
}: {
  value: string
  onChange: (value: string) => void
}) {
  // Hiện tại chỉ hỗ trợ thanh toán bằng thẻ tín dụng
  return (
    <RadioGroup value={value} onValueChange={onChange} className="grid grid-cols-1 gap-4">
      <div>
        <RadioGroupItem value="card" id="card" className="peer sr-only" defaultChecked />
        <Label
          htmlFor="card"
          className="flex cursor-pointer flex-col items-center justify-between rounded-md border-2 border-muted bg-popover p-4 hover:bg-accent hover:text-accent-foreground peer-data-[state=checked]:border-primary [&:has([data-state=checked])]:border-primary"
        >
          <CreditCard className="mb-3 h-6 w-6" />
          <div className="text-center">
            <p className="font-medium">Thẻ tín dụng</p>
            <p className="text-sm text-muted-foreground">Thanh toán bằng Visa, Mastercard, v.v.</p>
          </div>
        </Label>
      </div>
      {/* Có thể thêm các phương thức thanh toán khác ở đây trong tương lai */}
    </RadioGroup>
  )
}
