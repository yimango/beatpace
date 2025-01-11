import dynamic from 'next/dynamic'
import { Music, MonitorIcon as Running } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

const LoginButton = dynamic(() => import('@/components/LoginButton'), { ssr: false })
const PaceForm = dynamic(() => import('@/components/PaceForm'), { ssr: false })

export default function Home() {
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
          <LoginButton />
          <PaceForm accessToken="" />
        </CardContent>
      </Card>
    </main>
  )
}