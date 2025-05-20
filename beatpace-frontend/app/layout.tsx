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
    <html lang="en" className="h-full">
      <body 
        className={`${inter.className} relative min-h-screen bg-gradient-to-br from-blue-600 via-purple-600 to-pink-500`}
        style={{
          backgroundAttachment: 'fixed'
        }}
      >
        <div className="absolute inset-0 bg-black/30" />
        <div className="relative">
          <div className="container mx-auto min-h-screen px-4 py-8">
            {children}
          </div>
          <Toaster />
        </div>
      </body>
    </html>
  )
}

