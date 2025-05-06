"use client"

import { useState, useEffect } from "react"
import Image from "next/image"
import { motion, AnimatePresence } from "framer-motion"
import { Button } from "@/components/ui/button"
import { ShoppingBag, ChevronLeft, ChevronRight } from "lucide-react"
import SearchBar from "@/components/common/search-bar"

const slides = [
  {
    id: 1,
    title: "Khám phá xu hướng mới",
    subtitle: "Cho phong cách sống của bạn",
    description:
      "Sản phẩm chất lượng cao với giá cả cạnh tranh, giao hàng nhanh chóng và dịch vụ khách hàng tuyệt vời.",
    image: "/images/hero-slide-1.jpg",
    color: "from-blue-900/90 to-purple-900/90",
  },
  {
    id: 2,
    title: "Bộ sưu tập mùa hè",
    subtitle: "Giảm giá lên đến 50%",
    description: "Khám phá bộ sưu tập mùa hè mới nhất với nhiều ưu đãi hấp dẫn. Số lượng có hạn!",
    image: "/images/hero-slide-2.jpg",
    color: "from-green-900/90 to-blue-900/90",
  },
  {
    id: 3,
    title: "Công nghệ tiên tiến",
    subtitle: "Nâng tầm trải nghiệm",
    description: "Khám phá các sản phẩm công nghệ mới nhất với nhiều tính năng hiện đại và thiết kế sang trọng.",
    image: "/images/hero-slide-3.jpg",
    color: "from-gray-900/90 to-blue-900/90",
  },
]

export default function HeroCarousel() {
  const [currentSlide, setCurrentSlide] = useState(0)

  const nextSlide = () => {
    setCurrentSlide((prev) => (prev === slides.length - 1 ? 0 : prev + 1))
  }

  const prevSlide = () => {
    setCurrentSlide((prev) => (prev === 0 ? slides.length - 1 : prev - 1))
  }

  useEffect(() => {
    const interval = setInterval(() => {
      nextSlide()
    }, 6000)

    return () => clearInterval(interval)
  }, [])

  return (
    <section className="relative h-[600px] md:h-[700px] overflow-hidden">
      <AnimatePresence mode="wait">
        {slides.map(
          (slide, index) =>
            currentSlide === index && (
              <motion.div
                key={slide.id}
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                transition={{ duration: 1 }}
                className="absolute inset-0"
              >
                <div className={`absolute inset-0 bg-gradient-to-r ${slide.color} z-10`}></div>
                <div className="absolute inset-0">
                  <Image
                    src={slide.image || "/placeholder.svg"}
                    alt={slide.title}
                    fill
                    className="object-cover transform scale-110 animate-slow-zoom"
                    priority
                  />
                </div>

                <div className="container relative z-20 mx-auto px-4 h-full flex flex-col justify-center">
                  <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.8, delay: 0.2 }}
                    className="max-w-2xl"
                  >
                    <h2 className="text-xl md:text-2xl font-medium text-blue-300 mb-2">{slide.subtitle}</h2>
                    <h1 className="text-4xl md:text-6xl font-bold text-white mb-6 leading-tight">
                      {slide.title.split(" ").map((word, i) => (
                        <span key={i} className={i % 2 === 0 ? "" : "text-blue-400"}>
                          {word}{" "}
                        </span>
                      ))}
                    </h1>
                    <p className="text-lg md:text-xl text-white/90 mb-8">{slide.description}</p>
                    <div className="flex flex-col sm:flex-row gap-4">
                      <SearchBar />
                      <Button
                        size="lg"
                        className="bg-white text-blue-600 hover:bg-blue-50 px-8 py-6 rounded-lg shadow-lg hover:shadow-blue-500/30 transition-all duration-300 whitespace-nowrap"
                      >
                        <ShoppingBag className="mr-2 h-5 w-5" /> Mua ngay
                      </Button>
                    </div>
                  </motion.div>
                </div>
              </motion.div>
            ),
        )}
      </AnimatePresence>

      {/* Navigation Buttons */}
      <div className="absolute left-0 right-0 bottom-1/2 z-30 flex justify-between container mx-auto px-4">
        <Button
          variant="ghost"
          size="icon"
          className="rounded-full bg-white/20 backdrop-blur-sm hover:bg-white/40 text-white"
          onClick={prevSlide}
        >
          <ChevronLeft className="h-6 w-6" />
          <span className="sr-only">Previous</span>
        </Button>
        <Button
          variant="ghost"
          size="icon"
          className="rounded-full bg-white/20 backdrop-blur-sm hover:bg-white/40 text-white"
          onClick={nextSlide}
        >
          <ChevronRight className="h-6 w-6" />
          <span className="sr-only">Next</span>
        </Button>
      </div>

      {/* Indicators */}
      <div className="absolute bottom-6 left-0 right-0 z-30 flex justify-center gap-2">
        {slides.map((_, index) => (
          <button
            key={index}
            className={`w-3 h-3 rounded-full transition-all duration-300 ${
              currentSlide === index ? "bg-white w-8" : "bg-white/50"
            }`}
            onClick={() => setCurrentSlide(index)}
          />
        ))}
      </div>
    </section>
  )
}
