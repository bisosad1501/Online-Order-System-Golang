"use client"

import { useEffect, useState } from "react"
import { useParams, useRouter } from "next/navigation"
import Image from "next/image"
import Link from "next/link"
import { orderService } from "@/services/order-service"
import ProductReviewForm from "@/components/product/product-review-form"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Separator } from "@/components/ui/separator"
import { Badge } from "@/components/ui/badge"
import { useToast } from "@/hooks/use-toast"
import { ArrowLeft, CheckCircle } from "lucide-react"

type OrderItem = {
  id: string
  order_id: string
  product_id: string
  quantity: number
  price: number
  product_name?: string
}

export default function OrderReviewPage() {
  const params = useParams()
  const router = useRouter()
  const { toast } = useToast()
  const [orderItems, setOrderItems] = useState<OrderItem[]>([])
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState("")
  const [completedReviews, setCompletedReviews] = useState<string[]>([])

  const orderId = params.id as string

  useEffect(() => {
    async function fetchOrderDetails() {
      try {
        setLoading(true)
        const orderData = await orderService.getOrderById(orderId)

        if (orderData.status !== "DELIVERED") {
          toast({
            title: "Order not delivered",
            description: "You can only review products from delivered orders",
            variant: "destructive",
          })
          router.push(`/orders/${orderId}`)
          return
        }

        setOrderItems(orderData.items)
        if (orderData.items.length > 0) {
          setActiveTab(orderData.items[0].id)
        }
      } catch (error) {
        toast({
          title: "Error",
          description: "Failed to load order details",
          variant: "destructive",
        })
        router.push("/orders")
      } finally {
        setLoading(false)
      }
    }

    fetchOrderDetails()
  }, [orderId, router, toast])

  const handleReviewSuccess = (itemId: string) => {
    setCompletedReviews([...completedReviews, itemId])

    // Find the next item to review
    const currentIndex = orderItems.findIndex(item => item.id === itemId)
    const nextItem = orderItems[currentIndex + 1]

    if (nextItem) {
      setActiveTab(nextItem.id)
    } else {
      // All items reviewed
      toast({
        title: "All reviews submitted",
        description: "Thank you for reviewing all your products!",
      })

      // Redirect after a short delay
      setTimeout(() => {
        router.push(`/orders/${orderId}`)
      }, 2000)
    }
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center">
          <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-t-2 border-blue-600"></div>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-6">
        <Button variant="ghost" asChild className="mb-4">
          <Link href={`/orders/${orderId}`}>
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Order
          </Link>
        </Button>
        <h1 className="text-3xl font-bold">Review Your Purchase</h1>
        <p className="text-muted-foreground">
          Share your experience with the products you purchased in order #{orderId.substring(0, 8)}
        </p>
      </div>

      {orderItems.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12 text-center">
            <CheckCircle className="mb-4 h-16 w-16 text-green-500" />
            <h2 className="mb-2 text-xl font-semibold">No items to review</h2>
            <p className="mb-6 text-muted-foreground">
              There are no items in this order that need to be reviewed.
            </p>
            <Button asChild>
              <Link href={`/orders/${orderId}`}>Return to Order Details</Link>
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
          <div className="lg:col-span-2">
            <Card>
              <CardHeader>
                <CardTitle>Product Reviews</CardTitle>
                <CardDescription>
                  {completedReviews.length} of {orderItems.length} products reviewed
                </CardDescription>
              </CardHeader>
              <CardContent>
                <Tabs value={activeTab} onValueChange={setActiveTab}>
                  <TabsList className="mb-4 w-full grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4">
                    {orderItems.map((item) => (
                      <TabsTrigger key={item.id} value={item.id} className="relative">
                        {item.product_name || `Product ${item.product_id.substring(0, 4)}`}
                        {completedReviews.includes(item.id) && (
                          <CheckCircle className="absolute -right-1 -top-1 h-4 w-4 text-green-500" />
                        )}
                      </TabsTrigger>
                    ))}
                  </TabsList>

                  {orderItems.map((item) => (
                    <TabsContent key={item.id} value={item.id}>
                      {completedReviews.includes(item.id) ? (
                        <div className="flex flex-col items-center justify-center py-8 text-center">
                          <CheckCircle className="mb-4 h-16 w-16 text-green-500" />
                          <h2 className="mb-2 text-xl font-semibold">Review Submitted</h2>
                          <p className="mb-6 text-muted-foreground">
                            Thank you for reviewing this product!
                          </p>
                        </div>
                      ) : (
                        <ProductReviewForm
                          productId={item.product_id}
                          orderId={orderId}
                          onSuccess={() => handleReviewSuccess(item.id)}
                        />
                      )}
                    </TabsContent>
                  ))}
                </Tabs>
              </CardContent>
            </Card>
          </div>

          <div>
            <Card>
              <CardHeader>
                <CardTitle>Order Summary</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {orderItems.map((item) => (
                    <div key={item.id} className="flex items-start gap-4">
                      <div className="relative h-16 w-16 overflow-hidden rounded-md bg-gray-100">
                        <Image
                          src={`/placeholder.svg?height=64&width=64&text=${encodeURIComponent(
                            item.product_id.substring(0, 4)
                          )}`}
                          alt={item.product_name || `Product ${item.product_id.substring(0, 8)}`}
                          fill
                          className="object-cover"
                        />
                        {completedReviews.includes(item.id) && (
                          <div className="absolute inset-0 flex items-center justify-center bg-black/30">
                            <CheckCircle className="h-8 w-8 text-white" />
                          </div>
                        )}
                      </div>
                      <div className="flex-1">
                        <p className="font-medium">
                          {item.product_name || `Product ${item.product_id.substring(0, 8)}`}
                        </p>
                        <p className="text-sm text-muted-foreground">Quantity: {item.quantity}</p>
                        <p className="text-sm text-muted-foreground">${item.price.toFixed(2)}</p>
                      </div>
                      {completedReviews.includes(item.id) ? (
                        <Badge variant="outline" className="bg-green-50 text-green-700">
                          Reviewed
                        </Badge>
                      ) : (
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => setActiveTab(item.id)}
                        >
                          Review
                        </Button>
                      )}
                    </div>
                  ))}
                </div>

                <Separator className="my-4" />

                <div className="mt-4 text-center">
                  <p className="text-sm text-muted-foreground">
                    Your reviews help other customers make better purchase decisions.
                    Thank you for taking the time to share your experience!
                  </p>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      )}
    </div>
  )
}
