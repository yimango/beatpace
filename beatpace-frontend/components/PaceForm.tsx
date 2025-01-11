'use client'

import { useState } from 'react'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

interface PaceFormProps {
  accessToken: string
}

const PaceForm: React.FC<PaceFormProps> = ({ accessToken }) => {
  const [paceUnit, setPaceUnit] = useState<'km' | 'mile'>('km')
  const [minutes, setMinutes] = useState('')
  const [seconds, setSeconds] = useState('')
  const [gender, setGender] = useState('male')
  const [height, setHeight] = useState('')
  const [heightUnit, setHeightUnit] = useState<'cm' | 'in'>('cm')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
  
    const paceInSeconds = parseInt(minutes) * 60 + parseInt(seconds)
  
    const data = {
      accessToken,
      paceUnit,
      paceInSeconds,
      gender,
      height: parseFloat(height),
      heightUnit,
    }
  
    try {
      const backendUrl = process.env.NEXT_PUBLIC_BACKEND_URL || 'http://localhost:3001'
      console.log('Using backend URL:', backendUrl)
  
      const response = await fetch(`${backendUrl}/api/generate-playlist`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      })
  
      if (!response.ok) {
        throw new Error('Failed to generate playlist')
      }
  
      const result = await response.json()
      console.log('Playlist generated:', result)
    } catch (error) {
      console.error('Error generating playlist:', error)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6 w-full max-w-md">
      <div>
        <Label htmlFor="paceUnit">Pace Unit</Label>
        <RadioGroup id="paceUnit" value={paceUnit} onValueChange={(value) => setPaceUnit(value as 'km' | 'mile')} className="flex space-x-4">
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="km" id="km" />
            <Label htmlFor="km">min/km</Label>
          </div>
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="mile" id="mile" />
            <Label htmlFor="mile">min/mile</Label>
          </div>
        </RadioGroup>
      </div>
      <div className="flex space-x-4">
        <div className="flex-1">
          <Label htmlFor="minutes">Minutes</Label>
          <Input
            type="number"
            id="minutes"
            value={minutes}
            onChange={(e) => setMinutes(e.target.value)}
            min="0"
            required
          />
        </div>
        <div className="flex-1">
          <Label htmlFor="seconds">Seconds</Label>
          <Input
            type="number"
            id="seconds"
            value={seconds}
            onChange={(e) => setSeconds(e.target.value)}
            min="0"
            max="59"
            required
          />
        </div>
      </div>
      <div>
        <Label>Gender</Label>
        <RadioGroup value={gender} onValueChange={setGender} className="flex mt-1">
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="male" id="male" />
            <Label htmlFor="male">Male</Label>
          </div>
          <div className="flex items-center space-x-2 ml-4">
            <RadioGroupItem value="female" id="female" />
            <Label htmlFor="female">Female</Label>
          </div>
        </RadioGroup>
      </div>
      <div className="flex space-x-4">
        <div className="flex-1">
          <Label htmlFor="height">Height</Label>
          <Input
            type="number"
            id="height"
            value={height}
            onChange={(e) => setHeight(e.target.value)}
            min="0"
            step="0.1"
            required
          />
        </div>
        <div className="flex-1">
          <Label htmlFor="heightUnit">Unit</Label>
          <Select value={heightUnit} onValueChange={(value) => setHeightUnit(value as 'cm' | 'in')}>
            <SelectTrigger id="heightUnit">
              <SelectValue placeholder="Select unit" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="cm">cm</SelectItem>
              <SelectItem value="in">in</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>
      <Button type="submit" className="w-full">Generate Playlist</Button>
    </form>
  )
}

export default PaceForm
