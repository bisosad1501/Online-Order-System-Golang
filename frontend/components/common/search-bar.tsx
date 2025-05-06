"use client"

import type React from "react"

import { useState, useTransition } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Search } from "lucide-react"

export default function SearchBar() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [searchQuery, setSearchQuery] = useState(searchParams.get("search") || "")
  const [isPending, startTransition] = useTransition()

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()

    const params = new URLSearchParams(searchParams)
    if (searchQuery) {
      params.set("search", searchQuery)
    } else {
      params.delete("search")
    }

    startTransition(() => {
      router.push(`/products?${params.toString()}`)
    })
  }

  return (
    <form onSubmit={handleSearch} className="relative w-full max-w-lg">
      <Input
        type="text"
        placeholder="Tìm kiếm sản phẩm..."
        value={searchQuery}
        onChange={(e) => setSearchQuery(e.target.value)}
        className="pr-12 bg-white/90 backdrop-blur-sm border-transparent focus:border-blue-500 h-12 rounded-full shadow-lg"
      />
      <Button
        type="submit"
        size="icon"
        className="absolute right-1 top-1 h-10 w-10 rounded-full bg-blue-600 hover:bg-blue-700"
        disabled={isPending}
      >
        <Search className="h-4 w-4 text-white" />
        <span className="sr-only">Tìm kiếm</span>
      </Button>
    </form>
  )
}
