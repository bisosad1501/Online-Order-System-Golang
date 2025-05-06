import { Suspense } from "react"
import ProductList from "@/components/product/product-list"
import SearchBar from "@/components/common/search-bar"
import CategoryFilter from "@/components/product/category-filter"
import ProductPagination from "@/components/product/product-pagination"
import ProductLimitSelector from "@/components/product/product-limit-selector"
import { Skeleton } from "@/components/ui/skeleton"

export default function ProductsPage({
  searchParams,
}: {
  searchParams: { page?: string; limit?: string; category?: string; search?: string }
}) {
  // Server component có thể truy cập trực tiếp vào searchParams
  const page = Number(searchParams.page) || 1
  const limit = Number(searchParams.limit) || 12

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="mb-6 text-3xl font-bold">Products</h1>

      <div className="mb-6 space-y-4">
        <SearchBar />
        <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <CategoryFilter />
          <ProductLimitSelector currentLimit={limit} />
        </div>
      </div>

      <Suspense fallback={<ProductListSkeleton />}>
        <ProductList page={page} limit={limit} category={searchParams.category} search={searchParams.search} />
      </Suspense>

      <div className="mt-8 flex justify-center">
        <ProductPagination currentPage={page} />
      </div>
    </div>
  )
}

function ProductListSkeleton() {
  return (
    <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
      {[...Array(12)].map((_, i) => (
        <div key={i} className="rounded-lg border p-4">
          <Skeleton className="h-40 w-full rounded-md" />
          <Skeleton className="mt-4 h-6 w-3/4" />
          <Skeleton className="mt-2 h-4 w-1/2" />
        </div>
      ))}
    </div>
  )
}
