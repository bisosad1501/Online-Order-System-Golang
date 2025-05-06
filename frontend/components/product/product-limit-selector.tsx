"use client"

import { useRouter, usePathname, useSearchParams } from "next/navigation"

export default function ProductLimitSelector({ currentLimit }: { currentLimit: number }) {
  const router = useRouter()
  const pathname = usePathname()
  const searchParams = useSearchParams()

  const handleLimitChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const newLimit = e.target.value
    
    // Tạo một đối tượng URLSearchParams mới từ searchParams hiện tại
    const params = new URLSearchParams(searchParams.toString())
    
    // Cập nhật limit và đặt lại page về 1
    params.set("limit", newLimit)
    params.set("page", "1")
    
    // Chuyển hướng đến URL mới
    router.push(`${pathname}?${params.toString()}`)
  }

  return (
    <div className="flex items-center gap-2">
      <span className="text-sm text-muted-foreground">Items per page:</span>
      <select
        className="rounded-md border border-input bg-background px-3 py-1 text-sm"
        defaultValue={currentLimit}
        onChange={handleLimitChange}
      >
        <option value="12">12</option>
        <option value="24">24</option>
        <option value="36">36</option>
        <option value="48">48</option>
      </select>
    </div>
  )
}
