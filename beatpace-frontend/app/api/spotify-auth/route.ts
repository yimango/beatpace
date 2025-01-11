import { NextResponse } from 'next/server'

export async function POST() {
  try {
    const clientId = process.env.SPOTIFY_CLIENT_ID
    const redirectUri = process.env.SPOTIFY_REDIRECT_URI

    if (!clientId || !redirectUri) {
      throw new Error('Missing Spotify credentials')
    }

    const scope = 'user-read-private user-read-email playlist-modify-public playlist-modify-private'

    const authUrl = `https://accounts.spotify.com/authorize?client_id=${clientId}&response_type=code&redirect_uri=${encodeURIComponent(redirectUri)}&scope=${encodeURIComponent(scope)}`

    return NextResponse.json({ authUrl })
  } catch (error) {
    console.error('Error in Spotify authentication:', error)
    return NextResponse.json({ error: 'Failed to initialize Spotify authentication' }, { status: 500 })
  }
}

