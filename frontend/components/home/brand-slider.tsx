"use client"

import { useRef } from "react"
import Image from "next/image"
import { motion } from "framer-motion"

const brands = [
  { id: 1, name: "Apple", logo: "/images/placeholder.svg" },
  { id: 2, name: "Samsung", logo: "/images/placeholder.svg" },
  { id: 3, name: "Sony", logo: "/images/placeholder.svg" },
  { id: 4, name: "Microsoft", logo: "/images/placeholder.svg" },
  { id: 5, name: "Google", logo: "/images/placeholder.svg" },
  { id: 6, name: "LG", logo: "/images/placeholder.svg" },
  { id: 7, name: "Dell", logo: "/images/placeholder.svg" },
  { id: 8, name: "HP", logo: "/images/placeholder.svg" },
  { id: 9, name: "Asus", logo: "/images/placeholder.svg" },
  { id: 10, name: "Lenovo", logo: "/images/placeholder.svg" },
]

export default function BrandSlider() {
  const containerRef = useRef<HTMLDivElement>(null)

  return (
    <div className="overflow-hidden" ref={containerRef}>
      <motion.div
        className="flex gap-8 items-center"
        animate={{
          x: ["0%", "-50%"],
        }}
        transition={{
          x: {
            repeat: Number.POSITIVE_INFINITY,
            repeatType: "loop",
            duration: 20,
            ease: "linear",
          },
        }}
      >
        {/* Duplicate the brands to create a seamless loop */}
        {[...brands, ...brands].map((brand, index) => (
          <div
            key={`${brand.id}-${index}`}
            className="flex-shrink-0 w-32 h-20 flex items-center justify-center bg-white dark:bg-gray-800 rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow duration-300"
          >
            <Image
              src={brand.logo || "/placeholder.svg"}
              alt={brand.name}
              width={80}
              height={40}
              className="object-contain"
            />
          </div>
        ))}
      </motion.div>
    </div>
  )
}
