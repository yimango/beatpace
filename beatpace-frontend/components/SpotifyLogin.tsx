import { useState, useEffect } from 'react'
import { Button } from "@/components/ui/button"
import { loadEnvFile } from 'process'

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
    // Redirect to Spotify login
    loadEnvFile()
    const clientId = process.env.NEXT_PUBLIC_CLIENT_ID
    const redirectUri = process.env.NEXT_PUBLIC_REDIRECT_URI
    const scope = 'user-read-private user-read-email playlist-modify-private'
    const spotifyAuthUrl = `https://accounts.spotify.com/authorize?client_id=${clientId}&scope=${encodeURIComponent(scope)}&response_type=code&redirect_uri=${encodeURIComponent(redirectUri ? redirectUri : "")}&show_dialog=true`
    window.location.href = spotifyAuthUrl
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