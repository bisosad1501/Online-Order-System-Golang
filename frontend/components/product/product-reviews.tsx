"use client"

import { useState } from "react"
import Image from "next/image"
import { motion, AnimatePresence } from "framer-motion"
import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"
import { Progress } from "@/components/ui/progress"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { useAuth } from "@/context/auth-context"
import { useToast } from "@/hooks/use-toast"
import { Star, ThumbsUp, MessageSquare, ImageIcon } from "lucide-react"

// Dữ liệu mẫu cho đánh giá
const MOCK_REVIEWS = [
  {
    id: "1",
    user: {
      id: "user1",
      name: "Nguyễn Văn A",
      avatar: "/placeholder.svg?height=40&width=40&text=NVA",
    },
    rating: 5,
    date: "2023-10-15T08:30:00Z",
    content: "Sản phẩm rất tốt, đóng gói cẩn thận, giao hàng nhanh. Tôi rất hài lòng với chất lượng và sẽ mua lại.",
    images: [
      "/placeholder.svg?height=100&width=100&text=Review+Image+1",
      "/placeholder.svg?height=100&width=100&text=Review+Image+2",
    ],
    likes: 12,
    replies: 3,
  },
  {
    id: "2",
    user: {
      id: "user2",
      name: "Trần Thị B",
      avatar: "/placeholder.svg?height=40&width=40&text=TTB",
    },
    rating: 4,
    date: "2023-10-10T14:20:00Z",
    content: "Sản phẩm khá ổn, chất lượng tốt. Tuy nhiên, màu sắc hơi khác so với hình ảnh trên website.",
    images: [],
    likes: 5,
    replies: 1,
  },
  {
    id: "3",
    user: {
      id: "user3",
      name: "Lê Văn C",
      avatar: "/placeholder.svg?height=40&width=40&text=LVC",
    },
    rating: 5,
    date: "2023-10-05T09:45:00Z",
    content: "Tuyệt vời! Sản phẩm đúng như mô tả, giao hàng nhanh, đóng gói cẩn thận. Sẽ ủng hộ shop lần sau.",
    images: ["/placeholder.svg?height=100&width=100&text=Review+Image+3"],
    likes: 8,
    replies: 0,
  },
  {
    id: "4",
    user: {
      id: "user4",
      name: "Phạm Thị D",
      avatar: "/placeholder.svg?height=40&width=40&text=PTD",
    },
    rating: 3,
    date: "2023-09-28T16:15:00Z",
    content: "Sản phẩm tạm được, không quá xuất sắc nhưng cũng không tệ. Giao hàng hơi chậm.",
    images: [],
    likes: 2,
    replies: 1,
  },
  {
    id: "5",
    user: {
      id: "user5",
      name: "Hoàng Văn E",
      avatar: "/placeholder.svg?height=40&width=40&text=HVE",
    },
    rating: 5,
    date: "2023-09-20T11:30:00Z",
    content: "Sản phẩm chất lượng cao, đáng đồng tiền. Nhân viên tư vấn nhiệt tình, giao hàng nhanh chóng.",
    images: [
      "/placeholder.svg?height=100&width=100&text=Review+Image+4",
      "/placeholder.svg?height=100&width=100&text=Review+Image+5",
      "/placeholder.svg?height=100&width=100&text=Review+Image+6",
    ],
    likes: 15,
    replies: 2,
  },
]

// Thống kê đánh giá
const REVIEW_STATS = {
  average: 4.5,
  total: 124,
  distribution: [
    { rating: 5, count: 80, percentage: 65 },
    { rating: 4, count: 30, percentage: 24 },
    { rating: 3, count: 10, percentage: 8 },
    { rating: 2, count: 3, percentage: 2 },
    { rating: 1, count: 1, percentage: 1 },
  ],
}

type ProductReviewsProps = {
  productId: string
}

export default function ProductReviews({ productId }: ProductReviewsProps) {
  const { isAuthenticated, user } = useAuth()
  const { toast } = useToast()
  const [reviews, setReviews] = useState(MOCK_REVIEWS)
  const [newReview, setNewReview] = useState("")
  const [rating, setRating] = useState(5)
  const [filter, setFilter] = useState("all")
  const [showReviewForm, setShowReviewForm] = useState(false)
  const [hoveredStar, setHoveredStar] = useState(0)

  const handleSubmitReview = () => {
    if (!isAuthenticated) {
      toast({
        title: "Bạn cần đăng nhập",
        description: "Vui lòng đăng nhập để đánh giá sản phẩm",
        variant: "destructive",
      })
      return
    }

    if (!newReview.trim()) {
      toast({
        title: "Nội dung trống",
        description: "Vui lòng nhập nội dung đánh giá",
        variant: "destructive",
      })
      return
    }

    if (rating === 0) {
      toast({
        title: "Chưa chọn số sao",
        description: "Vui lòng chọn số sao đánh giá",
        variant: "destructive",
      })
      return
    }

    // Thêm đánh giá mới
    const newReviewObject = {
      id: `new-${Date.now()}`,
      user: {
        id: user?.id || "guest",
        name: user ? `${user.first_name} ${user.last_name}` : "Khách",
        avatar: `/placeholder.svg?height=40&width=40&text=${user?.first_name?.charAt(0) || "G"}`,
      },
      rating,
      date: new Date().toISOString(),
      content: newReview,
      images: [],
      likes: 0,
      replies: 0,
    }

    setReviews([newReviewObject, ...reviews])
    setNewReview("")
    setRating(5)
    setShowReviewForm(false)

    toast({
      title: "Đánh giá thành công",
      description: "Cảm ơn bạn đã đánh giá sản phẩm",
    })
  }

  const handleLikeReview = (reviewId: string) => {
    if (!isAuthenticated) {
      toast({
        title: "Bạn cần đăng nhập",
        description: "Vui lòng đăng nhập để thích đánh giá",
        variant: "destructive",
      })
      return
    }

    setReviews(reviews.map((review) => (review.id === reviewId ? { ...review, likes: review.likes + 1 } : review)))
  }

  const filteredReviews =
    filter === "all" ? reviews : reviews.filter((review) => review.rating === Number.parseInt(filter))

  return (
    <div>
      {/* Review Statistics */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-8">
        <div className="md:col-span-1 flex flex-col items-center justify-center">
          <div className="text-5xl font-bold text-blue-600">{REVIEW_STATS.average}</div>
          <div className="flex mt-2">
            {[...Array(5)].map((_, i) => (
              <Star
                key={i}
                className={`h-5 w-5 ${
                  i < Math.floor(REVIEW_STATS.average) ? "text-yellow-400 fill-yellow-400" : "text-gray-300"
                }`}
              />
            ))}
          </div>
          <div className="text-sm text-gray-500 mt-2">{REVIEW_STATS.total} đánh giá</div>
        </div>

        <div className="md:col-span-2">
          {REVIEW_STATS.distribution.map((item) => (
            <div key={item.rating} className="flex items-center mb-2">
              <div className="w-12 text-sm">{item.rating} sao</div>
              <Progress value={item.percentage} className="h-2 flex-1 mx-2" />
              <div className="w-12 text-sm text-right">{item.count}</div>
            </div>
          ))}
        </div>
      </div>

      {/* Review Actions */}
      <div className="flex flex-wrap items-center justify-between gap-4 mb-6">
        <div className="flex flex-wrap gap-2">
          <Button variant={filter === "all" ? "default" : "outline"} size="sm" onClick={() => setFilter("all")}>
            Tất cả
          </Button>
          {[5, 4, 3, 2, 1].map((star) => (
            <Button
              key={star}
              variant={filter === star.toString() ? "default" : "outline"}
              size="sm"
              onClick={() => setFilter(star.toString())}
            >
              {star} sao
            </Button>
          ))}
        </div>

        <Button onClick={() => setShowReviewForm(!showReviewForm)}>
          {showReviewForm ? "Hủy đánh giá" : "Viết đánh giá"}
        </Button>
      </div>

      {/* Review Form */}
      <AnimatePresence>
        {showReviewForm && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: "auto" }}
            exit={{ opacity: 0, height: 0 }}
            className="overflow-hidden"
          >
            <div className="bg-gray-50 rounded-xl p-6 mb-8">
              <h3 className="font-semibold mb-4">Đánh giá sản phẩm</h3>

              <div className="flex items-center mb-4">
                <span className="mr-2">Đánh giá:</span>
                <div className="flex">
                  {[1, 2, 3, 4, 5].map((star) => (
                    <Star
                      key={star}
                      className={`h-6 w-6 cursor-pointer ${
                        star <= (hoveredStar || rating) ? "text-yellow-400 fill-yellow-400" : "text-gray-300"
                      }`}
                      onMouseEnter={() => setHoveredStar(star)}
                      onMouseLeave={() => setHoveredStar(0)}
                      onClick={() => setRating(star)}
                    />
                  ))}
                </div>
              </div>

              <Textarea
                placeholder="Chia sẻ trải nghiệm của bạn về sản phẩm này..."
                value={newReview}
                onChange={(e) => setNewReview(e.target.value)}
                className="min-h-[120px] mb-4"
              />

              <div className="flex items-center justify-between">
                <Button variant="outline" size="sm" className="flex items-center gap-1">
                  <ImageIcon className="h-4 w-4" />
                  Thêm ảnh
                </Button>
                <Button onClick={handleSubmitReview}>Gửi đánh giá</Button>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Reviews List */}
      <div className="space-y-6">
        {filteredReviews.length > 0 ? (
          filteredReviews.map((review) => (
            <div key={review.id} className="border-b pb-6">
              <div className="flex items-start gap-4">
                <Avatar>
                  <AvatarImage src={review.user.avatar || "/placeholder.svg"} alt={review.user.name} />
                  <AvatarFallback>{review.user.name.charAt(0)}</AvatarFallback>
                </Avatar>

                <div className="flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <h4 className="font-medium">{review.user.name}</h4>
                    <div className="flex">
                      {[...Array(5)].map((_, i) => (
                        <Star
                          key={i}
                          className={`h-4 w-4 ${
                            i < review.rating ? "text-yellow-400 fill-yellow-400" : "text-gray-300"
                          }`}
                        />
                      ))}
                    </div>
                  </div>

                  <p className="text-sm text-gray-500 mb-2">
                    {new Date(review.date).toLocaleDateString("vi-VN", {
                      year: "numeric",
                      month: "long",
                      day: "numeric",
                    })}
                  </p>

                  <p className="text-gray-700 mb-4">{review.content}</p>

                  {review.images.length > 0 && (
                    <div className="flex flex-wrap gap-2 mb-4">
                      {review.images.map((image, index) => (
                        <div key={index} className="relative h-20 w-20 overflow-hidden rounded-md">
                          <Image
                            src={image || "/placeholder.svg"}
                            alt={`Review image ${index + 1}`}
                            fill
                            className="object-cover"
                          />
                        </div>
                      ))}
                    </div>
                  )}

                  <div className="flex items-center gap-4">
                    <Button
                      variant="ghost"
                      size="sm"
                      className="flex items-center gap-1 text-gray-500 hover:text-blue-600"
                      onClick={() => handleLikeReview(review.id)}
                    >
                      <ThumbsUp className="h-4 w-4" />
                      <span>{review.likes}</span>
                    </Button>

                    <Button
                      variant="ghost"
                      size="sm"
                      className="flex items-center gap-1 text-gray-500 hover:text-blue-600"
                    >
                      <MessageSquare className="h-4 w-4" />
                      <span>{review.replies}</span>
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          ))
        ) : (
          <div className="text-center py-8">
            <p className="text-gray-500">Không có đánh giá nào phù hợp với bộ lọc.</p>
          </div>
        )}
      </div>
    </div>
  )
}
