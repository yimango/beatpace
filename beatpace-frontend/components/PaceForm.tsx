"use client";

import { useState } from 'react';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { useToast } from "@/hooks/use-toast";
import { useRouter } from 'next/navigation';
import { Loader2 } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export default function PaceForm() {
  const [paceUnit, setPaceUnit] = useState<'km' | 'mile'>('km');
  const [minutes, setMinutes] = useState('');
  const [seconds, setSeconds] = useState('');
  const [gender, setGender] = useState<'male' | 'female'>('male');
  const [height, setHeight] = useState('');
  const [heightUnit, setHeightUnit] = useState<'cm' | 'in'>('cm');
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();
  const router = useRouter();
  const [playlistUrl, setPlaylistUrl] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    const paceInSeconds = Number(minutes) * 60 + Number(seconds);
    const heightInCm = heightUnit === 'in' ? Number(height) * 2.54 : Number(height);
    const data = {
      paceInSeconds,
      gender,
      height: Math.round(heightInCm),
    };

    try {
      const token = localStorage.getItem('token');
      if (!token) {
        toast({
          title: "Error",
          description: "Please log in again",
          variant: "destructive",
        });
        router.push('/');
        return;
      }

      // Check if token is expired
      const expires = localStorage.getItem('tokenExpires');
      if (expires && new Date(expires) < new Date()) {
        toast({
          title: "Session Expired",
          description: "Please log in again",
          variant: "destructive",
        });
        localStorage.removeItem('token');
        localStorage.removeItem('tokenExpires');
        localStorage.removeItem('user');
        router.push('/');
        return;
      }

      const backendUrl = process.env.NEXT_PUBLIC_BACKEND_URL ?? 'http://localhost:3001';
      console.log('Making request to:', `${backendUrl}/api/generate-playlist`);
      console.log('Request data:', data);
      console.log('Using token:', token);
      
      const response = await fetch(`${backendUrl}/api/generate-playlist`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error('Error response:', response.status, errorText);
        throw new Error('Failed to generate playlist');
      }

      const result = await response.json();
      console.log('Playlist generated:', result);
      
      // Store the playlist URL in state
      setPlaylistUrl(result.url);
      
      toast({
        title: "Success",
        description: "Your playlist has been generated!",
      });
    } catch (error) {
      console.error('Error generating playlist:', error);
      toast({
        title: "Error",
        description: "Failed to generate playlist. Please try again.",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="space-y-8">
      <form onSubmit={handleSubmit} className="w-full max-w-xl mx-auto space-y-8">
        {/* Pace Unit */}
        <div className="space-y-4">
          <Label htmlFor="paceUnit" className="text-base font-medium text-white">Pace Unit</Label>
          <RadioGroup
            id="paceUnit"
            value={paceUnit}
            onValueChange={(v) => setPaceUnit(v as 'km' | 'mile')}
            className="flex flex-wrap gap-6"
          >
            <div className="flex items-center space-x-3">
              <RadioGroupItem value="km" id="pace-km" className="h-5 w-5 border-2 border-white/80 text-white data-[state=checked]:bg-white data-[state=checked]:text-background" />
              <Label htmlFor="pace-km" className="text-base text-white/90">min/km</Label>
            </div>
            <div className="flex items-center space-x-3">
              <RadioGroupItem value="mile" id="pace-mile" className="h-5 w-5 border-2 border-white/80 text-white data-[state=checked]:bg-white data-[state=checked]:text-background" />
              <Label htmlFor="pace-mile" className="text-base text-white/90">min/mile</Label>
            </div>
          </RadioGroup>
        </div>

        {/* Minutes & Seconds */}
        <div className="grid grid-cols-2 gap-6">
          <div className="space-y-4">
            <Label htmlFor="minutes" className="text-base font-medium text-white">Minutes</Label>
            <Input
              type="number"
              id="minutes"
              value={minutes}
              onChange={(e) => setMinutes(e.target.value)}
              min="0"
              required
              className="h-12 text-lg bg-white/10 border-2 border-white/20 text-white placeholder:text-white/50 focus-visible:ring-2 focus-visible:ring-white/50 focus-visible:border-white/50"
            />
          </div>
          <div className="space-y-4">
            <Label htmlFor="seconds" className="text-base font-medium text-white">Seconds</Label>
            <Input
              type="number"
              id="seconds"
              value={seconds}
              onChange={(e) => setSeconds(e.target.value)}
              min="0"
              max="59"
              required
              className="h-12 text-lg bg-white/10 border-2 border-white/20 text-white placeholder:text-white/50 focus-visible:ring-2 focus-visible:ring-white/50 focus-visible:border-white/50"
            />
          </div>
        </div>

        {/* Gender */}
        <div className="space-y-4">
          <Label className="text-base font-medium text-white">Gender</Label>
          <RadioGroup
            value={gender}
            onValueChange={(v) => setGender(v as 'male' | 'female')}
            className="flex flex-wrap gap-6"
          >
            <div className="flex items-center space-x-3">
              <RadioGroupItem value="male" id="gender-male" className="h-5 w-5 border-2 border-white/80 text-white data-[state=checked]:bg-white data-[state=checked]:text-background" />
              <Label htmlFor="gender-male" className="text-base text-white/90">Male</Label>
            </div>
            <div className="flex items-center space-x-3">
              <RadioGroupItem value="female" id="gender-female" className="h-5 w-5 border-2 border-white/80 text-white data-[state=checked]:bg-white data-[state=checked]:text-background" />
              <Label htmlFor="gender-female" className="text-base text-white/90">Female</Label>
            </div>
          </RadioGroup>
        </div>

        {/* Height */}
        <div className="grid grid-cols-2 gap-6">
          <div className="space-y-4">
            <Label htmlFor="height" className="text-base font-medium text-white">Height</Label>
            <Input
              type="number"
              id="height"
              value={height}
              onChange={(e) => setHeight(e.target.value)}
              min="0"
              step="0.1"
              required
              className="h-12 text-lg bg-white/10 border-2 border-white/20 text-white placeholder:text-white/50 focus-visible:ring-2 focus-visible:ring-white/50 focus-visible:border-white/50"
            />
          </div>
          <div className="space-y-4">
            <Label htmlFor="heightUnit" className="text-base font-medium text-white">Unit</Label>
            <Select
              value={heightUnit}
              onValueChange={(v) => setHeightUnit(v as 'cm' | 'in')}
            >
              <SelectTrigger id="heightUnit" className="h-12 text-lg bg-white/10 border-2 border-white/20 text-white focus:ring-2 focus:ring-white/50 focus:border-white/50">
                <SelectValue placeholder="Select unit" />
              </SelectTrigger>
              <SelectContent className="bg-background/95 border-white/20">
                <SelectItem value="cm" className="text-white focus:bg-white/20 focus:text-white">Centimeters</SelectItem>
                <SelectItem value="in" className="text-white focus:bg-white/20 focus:text-white">Inches</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        <Button 
          type="submit" 
          className="w-full h-12 text-lg font-medium mt-8 bg-gradient-to-r from-blue-500 to-purple-500 hover:from-blue-600 hover:to-purple-600 text-white border-0 shadow-lg" 
          disabled={isLoading}
        >
          {isLoading ? (
            <>
              <Loader2 className="mr-2 h-5 w-5 animate-spin" />
              Generating...
            </>
          ) : (
            'Generate Playlist'
          )}
        </Button>
      </form>

      {/* Display the playlist embed if available */}
      {playlistUrl && (
        <div className="w-full max-w-xl mx-auto mt-8">
          <iframe
            src={`https://open.spotify.com/embed/playlist/${playlistUrl.split('/').pop()}`}
            width="100%"
            height="352"
            frameBorder="0"
            allowFullScreen
            allow="autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture"
            loading="lazy"
            className="rounded-lg shadow-lg"
          />
        </div>
      )}
    </div>
  );
}
