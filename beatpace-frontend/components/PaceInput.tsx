'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'

export default function PaceInput() {
  const [pace, setPace] = useState('')
  const [unit, setUnit] = useState('min/km')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    // Here we'll make the API call (to be implemented later)
    console.log(`Submitting pace: ${pace} ${unit}`)
    // Placeholder for API call
    // await fetch('/api/generate-playlist', { method: 'POST', body: JSON.stringify({ pace, unit }) })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="flex space-x-2">
        <Input
          type="text"
          value={pace}
          onChange={(e) => setPace(e.target.value)}
          placeholder="Enter your pace"
          className="w-40"
        />
        <Select value={unit} onValueChange={setUnit}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Select unit" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="min/km">min/km</SelectItem>
            <SelectItem value="min/mile">min/mile</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <Button type="submit">Generate Playlist</Button>
    </form>
  )
}

