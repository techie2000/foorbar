'use client'

import { useEffect, useState } from 'react'

export default function Footer() {
  const [version, setVersion] = useState<string>('loading...')

  useEffect(() => {
    const fetchVersion = async () => {
      try {
        const API_BASE_URL = typeof window !== 'undefined' 
          ? (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:18080')
          : 'http://backend:8080'
        
        const response = await fetch(`${API_BASE_URL}/version`)
        if (response.ok) {
          const data = await response.json()
          setVersion(data.version)
        } else {
          setVersion('unknown')
        }
      } catch (error) {
        console.error('Failed to fetch version:', error)
        setVersion('unknown')
      }
    }

    fetchVersion()
  }, [])

  return (
    <footer className="fixed bottom-0 right-0 px-3 py-1 bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400 text-xs rounded-tl-md border-t border-l border-gray-200 dark:border-gray-700 z-10">
      v{version}
    </footer>
  )
}
