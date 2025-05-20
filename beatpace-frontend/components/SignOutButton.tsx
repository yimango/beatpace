'use client'

import { Button } from "@/components/ui/button"
import { useToast } from "@/hooks/use-toast"
import { useRouter } from "next/navigation"
import { LogOut } from "lucide-react"

export default function SignOutButton() {
  const { toast } = useToast()
  const router = useRouter()

  const handleSignOut = async () => {
    try {
      const token = localStorage.getItem('token')
      if (!token) {
        router.push('/')
        return
      }

      const backendUrl = process.env.NEXT_PUBLIC_BACKEND_URL || 'http://localhost:3001'
      console.log('Making sign-out request to:', `${backendUrl}/api/signout`)
      const response = await fetch(`${backendUrl}/api/signout`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      })
      console.log('Sign-out response status:', response.status)

      if (!response.ok) {
        throw new Error('Failed to sign out')
      }

      // Clear local storage
      localStorage.removeItem('token')
      localStorage.removeItem('tokenExpires')
      localStorage.removeItem('user')

      toast({
        title: "Success",
        description: "Successfully signed out",
      })

      // Reload the page to show signed out state
      window.location.href = '/'
    } catch (error) {
      console.error('Error signing out:', error)
      toast({
        title: "Error",
        description: "Failed to sign out. Please try again.",
        variant: "destructive",
      })
    }
  }

  return (
    <Button 
      onClick={handleSignOut}
      variant="ghost"
      className="text-white hover:text-white/80 hover:bg-white/10"
    >
      <LogOut className="h-5 w-5 mr-2" />
      Sign Out
    </Button>
  )
} 