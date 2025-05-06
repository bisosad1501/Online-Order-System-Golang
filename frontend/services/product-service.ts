import api from "./api"

export const productService = {
  // Get products with pagination, filtering, and search
  async getProducts({
    page = 1,
    limit = 12,
    category,
    search,
  }: {
    page?: number
    limit?: number
    category?: string
    search?: string
  } = {}) {
    try {
      const params = new URLSearchParams()
      params.append("page", page.toString())
      params.append("limit", limit.toString())

      if (category) {
        params.append("category", category)
      }

      if (search) {
        params.append("search", search)
      }

      const response = await api.get(`/products?${params.toString()}`)
      return response.data
    } catch (error) {
      console.error("Error fetching products:", error)
      // Throw error instead of returning mock data
      throw new Error("Không thể lấy dữ liệu sản phẩm từ server. Vui lòng thử lại sau.")
    }
  },

  // Get a single product by ID
  async getProductById(productId: string) {
    try {
      const response = await api.get(`/products/${productId}`)
      return response.data
    } catch (error) {
      console.error("Error fetching product:", error)
      throw new Error("Không thể lấy thông tin sản phẩm từ server. Vui lòng thử lại sau.")
    }
  },

  // Get product recommendations
  async getRecommendations(productId: string) {
    try {
      const response = await api.get(`/recommendations/product/${productId}`)
      return response.data
    } catch (error) {
      console.error("Error fetching recommendations:", error)
      throw new Error("Không thể lấy đề xuất sản phẩm từ server. Vui lòng thử lại sau.")
    }
  },
}

// Function to get a product by ID (for server components)
export async function getProductById(productId: string) {
  try {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:9090"}/products/${productId}`)

    if (!response.ok) {
      throw new Error(`Failed to fetch product: ${response.status}`)
    }

    return await response.json()
  } catch (error) {
    console.error("Error fetching product:", error)
    throw new Error("Không thể lấy thông tin sản phẩm từ server. Vui lòng thử lại sau.")
  }
}

// Thêm các hàm tạo dữ liệu mẫu
function generateMockProducts(count: number) {
  return Array(count)
    .fill(0)
    .map((_, index) => ({
      id: `mock-product-${index}`,
      name: `Sản phẩm mẫu ${index + 1}`,
      description:
        "Mô tả chi tiết về sản phẩm này sẽ được hiển thị ở đây. Đây là một sản phẩm chất lượng cao với nhiều tính năng hữu ích.",
      category_id: "electronics",
      price: 99.99 + index * 10,
      tags: ["electronics", "gadget"],
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }))
}

function generateMockProduct(id: string) {
  return {
    id,
    name: `Sản phẩm chi tiết ${id}`,
    description:
      "Đây là mô tả chi tiết về sản phẩm. Sản phẩm này được thiết kế với chất lượng cao và độ bền tuyệt vời. Phù hợp với nhiều mục đích sử dụng khác nhau.",
    category_id: "electronics",
    price: 129.99,
    tags: ["premium", "electronics", "bestseller"],
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  }
}
