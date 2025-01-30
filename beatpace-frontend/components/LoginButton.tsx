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
      const clientId = process.env.NEXT_PUBLIC_CLIENT_ID
      const redirectUri = process.env.NEXT_PUBLIC_REDIRECT_URI
      const scope = 'playlist-modify-public playlist-modify-private'
      window.location.href = `https://accounts.spotify.com/authorize?client_id=${clientId}&scope=${encodeURIComponent(scope)}&response_type=code&redirect_uri=${redirectUri}&show_dialog=true`;
    } catch (error) {
      console.error('Error during Spotify authentication:', error);
      toast({
        title: "Authentication Error",
        description: "There was a problem connecting to Spotify. Please try again.",
        variant: "destructive",
      })
    } finally {
      setIsLoading(false)
    }
  };

  if (!isMounted) {
    return null;
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
  );
}

