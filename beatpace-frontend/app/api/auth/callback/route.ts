import { NextResponse } from 'next/server'

export async function POST(request: Request) {
  console.log('Backend URL:', process.env.BACKEND_URL) // 🔍 Debug log

  const body = await request.json()
  const token = request.headers.get('Authorization')

  try {
    const backendUrl = process.env.BACKEND_URL || 'http://localhost:3001'
    console.log('Using Backend URL:', backendUrl) // 🔍 Another debug log

    const backendResponse = await fetch(`${backendUrl}/api/generate-playlist`, {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': token || '',
      },
      body: JSON.stringify(body),
    })

    if (!backendResponse.ok) {
      const errorText = await backendResponse.text()
      throw new Error(`Backend request failed: ${errorText}`)
    }

    const data = await backendResponse.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Error calling Go backend:', error)
    return NextResponse.json({ error: 'Failed to generate playlist' }, { status: 500 })
  }
}
