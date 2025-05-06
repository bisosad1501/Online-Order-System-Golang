"use client"

import { useRouter, useSearchParams } from "next/navigation"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

const categories = [
  { id: "all", name: "Tất cả danh mục" },
  { id: "electronics", name: "Điện tử" },
  { id: "clothing", name: "Thời trang" },
  { id: "home", name: "Nhà cửa & Đời sống" },
  { id: "books", name: "Sách" },
  { id: "toys", name: "Đồ chơi & Trò chơi" },
]

export default function CategoryFilter() {
  const router = useRouter()
  const searchParams = useSearchParams()

  const currentCategory = searchParams.get("category") || "all"

  const handleCategoryChange = (value: string) => {
    const params = new URLSearchParams(searchParams)

    if (value && value !== "all") {
      params.set("category", value)
    } else {
      params.delete("category")
    }

    // Reset to page 1 when changing category
    params.set("page", "1")

    router.push(`/products?${params.toString()}`)
  }

  return (
    <Select value={currentCategory} onValueChange={handleCategoryChange}>
      <SelectTrigger className="w-[180px] bg-white">
        <SelectValue placeholder="Chọn danh mục" />
      </SelectTrigger>
      <SelectContent>
        {categories.map((category) => (
          <SelectItem key={category.id} value={category.id}>
            {category.name}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}
