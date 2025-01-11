import './globals.css'
import { Inter } from 'next/font/google'
import { Toaster } from "@/components/ui/toaster"

const inter = Inter({ subsets: ['latin'] })

export const metadata = {
  title: 'BeatPace',
  description: 'Get the perfect playlist for your running pace',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={`${inter.className} bg-gradient-to-br from-purple-600 to-blue-600 min-h-screen`}>
        <div className="container mx-auto px-4 py-8">
          {children}
        </div>
        <Toaster />
      </body>
    </html>
  )
}

