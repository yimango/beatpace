'use client'

import { Button } from "@/components/ui/button"
import { Loader2 } from 'lucide-react'
import { useState, useEffect } from 'react'
import { useToast } from "@/hooks/use-toast"

export default function LoginButton() {
  const [isLoading, setIsLoading] = useState(false)
  const [isMounted, setIsMounted] = useState(false)
  const { toast } = useToast()

  useEffect(() => {
    setIsMounted(true)
  }, [])

  const handleLogin = async () => {
    setIsLoading(true)
    try {
      const backendUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3001'
      const response = await fetch(`${backendUrl}/api/spotify-auth`, {
        method: 'GET',
      })

      if (!response.ok) {
        throw new Error('Failed to get Spotify auth URL')
      }

      const data = await response.json()
      window.location.href = data.authUrl
    } catch (error) {
      console.error('Error during Spotify authentication:', error)
      toast({
        title: "Authentication Error",
        description: "There was a problem connecting to Spotify. Please try again.",
        variant: "destructive",
      })
    } finally {
      setIsLoading(false)
    }
  }

  if (!isMounted) {
    return null
  }

  return (
    <Button onClick={handleLogin} className="w-full mb-6" disabled={isLoading}>
      {isLoading ? (
        <>
          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
          Connecting to Spotify
        </>
      ) : (
        'Login with Spotify'
      )}
    </Button>
  )
}

