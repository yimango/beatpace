"use client"  // Add this line at the top of your file

import { useEffect, useState } from 'react'
import dynamic from 'next/dynamic'
import { Music, MonitorIcon as Running } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

const LoginButton = dynamic(() => import('@/components/LoginButton'), { ssr: false })
const PaceForm = dynamic(() => import('@/components/PaceForm'), { ssr: false })

export default function Home() {
  const [authCode, setAuthCode] = useState<string | null>(null)

  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search)
    const code = urlParams.get('code')
    if (code) {
      setAuthCode(code)
    }
  }, [])

  return (
    <main className="flex min-h-screen flex-col items-center justify-center">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="mx-auto mb-4 flex items-center justify-center space-x-2">
            <Running className="h-6 w-6 text-primary" />
            <Music className="h-6 w-6 text-primary" />
          </div>
          <CardTitle className="text-3xl font-bold">BeatPace</CardTitle>
          <CardDescription>Get the perfect playlist for your running pace</CardDescription>
        </CardHeader>
        <CardContent>
          {/* Only render LoginButton if authCode is not set */}
          {!authCode && <LoginButton />}
          {/* Only render PaceForm once authCode is set */}
          {authCode && <PaceForm authCode={authCode} />}
        </CardContent>
      </Card>
    </main>
  )
}
