"use client"

import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { Button } from "@/components/ui/button"
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { Textarea } from "@/components/ui/textarea"
import { useToast } from "@/hooks/use-toast"
import { Star } from "lucide-react"

const reviewSchema = z.object({
  rating: z.number().min(1, "Please select a rating").max(5),
  comment: z.string().min(10, "Comment must be at least 10 characters"),
})

type ReviewFormValues = z.infer<typeof reviewSchema>

type ProductReviewFormProps = {
  productId: string
  orderId: string
  onSuccess: () => void
}

export default function ProductReviewForm({ productId, orderId, onSuccess }: ProductReviewFormProps) {
  const { toast } = useToast()
  const [isLoading, setIsLoading] = useState(false)
  const [hoveredRating, setHoveredRating] = useState(0)

  const form = useForm<ReviewFormValues>({
    resolver: zodResolver(reviewSchema),
    defaultValues: {
      rating: 0,
      comment: "",
    },
  })

  const rating = form.watch("rating")

  async function onSubmit(data: ReviewFormValues) {
    setIsLoading(true)
    try {
      // In a real app, this would call an API to submit the review
      console.log("Submitting review:", { productId, orderId, ...data })
      
      // Simulate API call
      await new Promise((resolve) => setTimeout(resolve, 1000))
      
      toast({
        title: "Review submitted",
        description: "Thank you for your feedback!",
      })
      
      onSuccess()
    } catch (error) {
      toast({
        title: "Failed to submit review",
        description: "There was an error submitting your review. Please try again.",
        variant: "destructive",
      })
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <FormField
          control={form.control}
          name="rating"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Rating</FormLabel>
              <FormControl>
                <div className="flex items-center">
                  {[1, 2, 3, 4, 5].map((star) => (
                    <button
                      key={star}
                      type="button"
                      className="p-1"
                      onMouseEnter={() => setHoveredRating(star)}
                      onMouseLeave={() => setHoveredRating(0)}
                      onClick={() => field.onChange(star)}
                    >
                      <Star
                        className={`h-8 w-8 ${
                          star <= (hoveredRating || field.value)
                            ? "fill-yellow-400 text-yellow-400"
                            : "text-gray-300 dark:text-gray-600"
                        } transition-colors`}
                      />
                    </button>
                  ))}
                </div>
              </FormControl>
              <FormDescription>How would you rate this product?</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="comment"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Your Review</FormLabel>
              <FormControl>
                <Textarea
                  placeholder="Share your experience with this product..."
                  className="min-h-[120px] resize-none"
                  {...field}
                />
              </FormControl>
              <FormDescription>
                Your review will help other customers make better purchase decisions.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type="submit" disabled={isLoading}>
          {isLoading ? "Submitting..." : "Submit Review"}
        </Button>
      </form>
    </Form>
  )
}
