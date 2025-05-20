// app/page.tsx
"use client";

import { useEffect, useState } from "react";
import dynamic from "next/dynamic";
import { Music, MonitorIcon as Running, Loader2 } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useSearchParams } from 'next/navigation';
import { useToast } from "@/hooks/use-toast";
import LoginButton from '@/components/LoginButton';
import SignOutButton from '@/components/SignOutButton';

const PaceForm = dynamic(() => import("@/components/PaceForm"), {
  ssr: false,
});

interface User {
  id: string;
  spotify_user_id: string;
  email?: string;
}

interface MeResponse {
  user: User;
  spotify_token?: {
    access_token: string;
    expires_at: string;
  };
}

export default function Page() {
  const [loading, setLoading] = useState(true);
  const [me, setMe] = useState<MeResponse | null>(null);
  const searchParams = useSearchParams();
  const { toast } = useToast();

  useEffect(() => {
    const checkAuth = async () => {
      try {
        // Check if we have a token
        const token = localStorage.getItem('token');
        if (!token) {
          setMe(null);
          setLoading(false);
          return;
        }

        // Check if token is expired
        const expires = localStorage.getItem('tokenExpires');
        if (expires && new Date(expires) < new Date()) {
          // Token is expired, clear storage
          localStorage.removeItem('token');
          localStorage.removeItem('tokenExpires');
          localStorage.removeItem('user');
          setMe(null);
          setLoading(false);
          return;
        }

        // Try to get cached user data first
        const cachedUser = localStorage.getItem('user');
        if (cachedUser) {
          setMe(JSON.parse(cachedUser));
        }

        // Fetch fresh user data
        const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/api/me`, {
          method: "GET",
          headers: {
            'Authorization': `Bearer ${token}`,
          },
          credentials: "include",
        });

        if (!response.ok) {
          throw new Error("unauthenticated");
        }

        const data = await response.json();
        setMe(data);
        localStorage.setItem('user', JSON.stringify(data));
      } catch (error) {
        console.error('Auth check failed:', error);
        // Clear storage on error
        localStorage.removeItem('token');
        localStorage.removeItem('tokenExpires');
        localStorage.removeItem('user');
        setMe(null);
      } finally {
        setLoading(false);
      }
    };

    checkAuth();
  }, []);

  useEffect(() => {
    // Handle token if present in URL
    const token = searchParams.get('token');
    const expires = searchParams.get('expires');
    
    if (token && expires) {
      // Store the token
      localStorage.setItem('token', token);
      localStorage.setItem('tokenExpires', expires);

      // Show success message
      toast({
        title: "Success",
        description: "Successfully connected to Spotify",
      });

      // Clean up URL
      window.history.replaceState({}, '', '/');
    }
  }, [searchParams, toast]);

  // 1) Still checking?
  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center p-4">
        <div className="text-center">
          <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4 text-white" />
          <p className="text-white/90 text-lg">Loading...</p>
        </div>
      </div>
    );
  }

  // 2) Not logged in → show Spotify login
  if (!me) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center p-4 sm:p-8">
        <Card className="w-[90%] max-w-md mx-auto shadow-xl border-0 bg-white/10 backdrop-blur supports-[backdrop-filter]:bg-white/5">
          <CardHeader className="text-center space-y-4">
            <div className="flex items-center justify-center space-x-3">
              <Running className="h-8 w-8 text-white" />
              <Music className="h-8 w-8 text-white" />
            </div>
            <CardTitle className="text-3xl sm:text-4xl font-bold tracking-tight text-white">BeatPace</CardTitle>
            <CardDescription className="text-base sm:text-lg text-white/80">
              Get the perfect playlist for your running pace
            </CardDescription>
          </CardHeader>
          <CardContent className="pb-8">
            <LoginButton />
          </CardContent>
        </Card>
      </div>
    );
  }

  // 3) Logged in → show the pace form
  return (
    <div className="container mx-auto min-h-screen px-4 py-6 sm:px-8 sm:py-12">
      <div className="flex justify-end mb-4">
        <SignOutButton />
      </div>
      <Card className="w-[95%] max-w-2xl mx-auto shadow-xl border-0 bg-white/10 backdrop-blur supports-[backdrop-filter]:bg-white/5">
        <CardHeader className="text-center sm:text-left space-y-2">
          <div className="flex items-center justify-center sm:justify-start space-x-3 mb-2">
            <Running className="h-6 w-6 text-white" />
            <Music className="h-6 w-6 text-white" />
          </div>
          <CardTitle className="text-2xl sm:text-3xl font-bold text-white">Generate Your Running Playlist</CardTitle>
          <CardDescription className="text-base sm:text-lg text-white/80">
            Enter your details to get a personalized playlist that matches your running pace
          </CardDescription>
        </CardHeader>
        <CardContent>
          <PaceForm />
        </CardContent>
      </Card>
    </div>
  );
}
