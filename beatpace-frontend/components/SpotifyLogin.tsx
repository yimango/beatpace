import { useState, useEffect } from 'react'
import { Button } from "@/components/ui/button"

export default function SpotifyLogin() {
  const [isLoggedIn, setIsLoggedIn] = useState(false)

  useEffect(() => {
    // Check if the user is logged in (you'll need to implement this logic)
    // For now, we'll just use a placeholder
    const checkLoginStatus = async () => {
      // Placeholder: replace with actual login check
      setIsLoggedIn(false)
    }
    checkLoginStatus()
  }, [])

  const handleLogin = async () => {
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
    }
  }

  return (
    <div className="mb-8">
      {isLoggedIn ? (
        <p className="text-green-500">Logged in to Spotify</p>
      ) : (
        <Button onClick={handleLogin}>Log in to Spotify</Button>
      )}
    </div>
  )
}