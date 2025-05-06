import Link from "next/link"
import Image from "next/image"
import { productService } from "@/services/product-service"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardFooter } from "@/components/ui/card"
import AddToCartButton from "@/components/cart/add-to-cart-button"

export default async function ProductRecommendations({
  productId,
}: {
  productId: string
}) {
  const products = await productService.getRecommendations(productId)

  if (!products.length) {
    return null
  }

  return (
    <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
      {products.map((product) => (
        <Card key={product.id} className="overflow-hidden">
          <div className="relative aspect-square">
            <Image
              src={product.image_url || `/default-product-blue.svg`}
              alt={product.name}
              fill
              className="object-cover transition-transform hover:scale-105"
            />
          </div>
          <CardContent className="p-4">
            <div className="mb-2 flex flex-wrap gap-1">
              {product.tags.slice(0, 2).map((tag) => (
                <Badge key={tag} variant="secondary" className="text-xs">
                  {tag}
                </Badge>
              ))}
            </div>
            <h3 className="line-clamp-1 text-lg font-semibold">{product.name}</h3>
            <p className="mt-2 text-lg font-bold text-blue-600">${product.price.toFixed(2)}</p>
          </CardContent>
          <CardFooter className="flex gap-2 p-4 pt-0">
            <Button asChild variant="outline" className="flex-1">
              <Link href={`/products/${product.id}`}>View</Link>
            </Button>
            <AddToCartButton product={product} variant="icon" />
          </CardFooter>
        </Card>
      ))}
    </div>
  )
}
