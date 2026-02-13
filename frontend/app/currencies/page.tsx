'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import ThemeToggle from '../components/ThemeToggle'

interface Currency {
  id: string
  code: string
  name: string
  symbol: string
  numeric_code: string
}

export default function CurrenciesPage() {
  const [currencies, setCurrencies] = useState<Currency[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [searchTerm, setSearchTerm] = useState('')

  const API_BASE_URL = typeof window !== 'undefined'
    ? (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:18080')
    : 'http://backend:8080'

  useEffect(() => {
    if (typeof window !== 'undefined') {
      fetchCurrencies()
    }
  }, [])

  const fetchCurrencies = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/currencies`, {
        headers: {
          'Accept': 'application/json'
        }
      })

      if (response.ok) {
        const data = await response.json()
        console.log('Currencies API response:', data)
        setCurrencies(data || [])
        if (!data || data.length === 0) {
          setError('No currencies data available yet. The database may need to be populated with reference data.')
        } else {
          setError(null)
        }
      } else {
        setError(`API returned ${response.status}: ${response.statusText}`)
      }
    } catch (err) {
      console.error('Currencies fetch error:', err)
      setError('Unable to connect to backend API. Please ensure the backend service is running at ' + API_BASE_URL)
    } finally {
      setLoading(false)
    }
  }

  const filteredCurrencies = currencies.filter(currency =>
    currency.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    currency.code.toLowerCase().includes(searchTerm.toLowerCase()) ||
    currency.symbol.toLowerCase().includes(searchTerm.toLowerCase())
  )

  if (loading) {
    return (
      <div className="min-h-screen p-8">
        <div className="max-w-7xl mx-auto">
          <div className="text-center py-20">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            <p className="mt-4 opacity-70">Loading currencies...</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen p-8">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8 flex justify-between items-start">
          <div>
            <Link href="/" className="text-blue-400 hover:text-blue-300 mb-4 inline-block">
              ← Back to Home
            </Link>
            <h1 className="text-4xl font-bold mb-2">Currencies</h1>
            <p className="opacity-70">Browse ISO 4217 currency codes and reference data</p>
          </div>
          <ThemeToggle />
        </div>

        {/* Info/Error Alert */}
        {error && (
          <div className={`mb-6 p-4 rounded-lg border ${
            error.includes('No currencies data') 
              ? 'bg-yellow-50 border-yellow-200' 
              : 'bg-red-50 border-red-200'
          }`}>
            <p className={error.includes('No currencies data') ? 'text-yellow-800' : 'text-red-800'}>
              <span className="font-semibold">
                {error.includes('No currencies data') ? '📋 Notice:' : '⚠️ Error:'}
              </span> {error}
            </p>
            {error.includes('No currencies data') && (
              <p className="text-sm text-yellow-700 mt-2">
                💡 Tip: Currencies data is typically loaded during initial system setup. Contact your administrator if this data should be available.
              </p>
            )}
          </div>
        )}

        {/* Search */}
        <div className="mb-6">
          <input
            type="text"
            placeholder="Search by name, code, or symbol..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
          <div className="bg-white dark:bg-white/5 rounded-lg shadow p-6 border-2 border-gray-200 dark:border-white/10">
            <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">Total Currencies</h3>
            <p className="text-3xl font-bold text-gray-900 dark:text-white mt-2">{currencies.length}</p>
          </div>
          <div className="bg-white dark:bg-white/5 rounded-lg shadow p-6 border-2 border-gray-200 dark:border-white/10">
            <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">Filtered Results</h3>
            <p className="text-3xl font-bold text-gray-900 dark:text-white mt-2">{filteredCurrencies.length}</p>
          </div>
          <div className="bg-white dark:bg-white/5 rounded-lg shadow p-6 border-2 border-gray-200 dark:border-white/10">
            <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">Data Standard</h3>
            <p className="text-3xl font-bold text-gray-900 dark:text-white mt-2">ISO 4217</p>
          </div>
        </div>

        {/* Currencies Table */}
        <div className="bg-white dark:bg-white/5 rounded-lg shadow overflow-hidden border-2 border-gray-200 dark:border-white/10">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-white/10">
            <thead className="bg-gray-50 dark:bg-white/5">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Name
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Code
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Symbol
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Numeric
                </th>
              </tr>
            </thead>
            <tbody className="bg-white dark:bg-white/5 divide-y divide-gray-200 dark:divide-white/10">
              {filteredCurrencies.length > 0 ? (
                filteredCurrencies.map((currency) => (
                  <tr key={currency.id} className="hover:bg-gray-50 dark:hover:bg-white/10">
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
                      {currency.name}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      <span className="px-2 py-1 bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 rounded font-mono">
                        {currency.code}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400 text-xl">
                      {currency.symbol}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400 font-mono">
                      {currency.numeric_code}
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={4} className="px-6 py-4 text-center text-sm text-gray-500 dark:text-gray-400">
                    No currencies found matching your search
                  </td>
                </tr>
              )}

            </tbody>
          </table>
        </div>

        {/* Footer Note */}
        <div className="mt-6 text-center text-sm text-gray-500">
          <p>Data source: ISO 4217 Currency Codes • Public reference data</p>
        </div>
      </div>
    </div>
  )
}
