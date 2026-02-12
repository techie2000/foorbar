---
description: 'Frontend and UI development guidelines for the Axiom project'
applyTo: 'frontend/**/*.tsx,frontend/**/*.ts,frontend/**/*.jsx,frontend/**/*.js'
---

# Frontend and UI Development Guidelines

## Date and Time Formatting

### ISO 8601 Date Format (Required)
**Always use ISO 8601 date format (yyyy-mm-dd) for displaying dates to users.**

This international standard is unambiguous and widely recognized globally, avoiding confusion between DD/MM/YYYY (European) and MM/DD/YYYY (American) formats.

#### ✅ CORRECT Examples

```typescript
// Date only (yyyy-mm-dd) with Go zero date handling
const formatDate = (dateString: string | null) => {
  if (!dateString || dateString.startsWith('0001-')) return '-'
  return new Date(dateString).toISOString().split('T')[0]
}
// Output: "2026-02-12" or "-" for invalid dates

// Date and time (yyyy-mm-dd HH:mm:ss) with Go zero date handling
const formatDateTime = (dateString: string | null) => {
  if (!dateString || dateString.startsWith('0001-')) return 'Never'
  return new Date(dateString).toISOString().replace('T', ' ').substring(0, 19)
}
// Output: "2026-02-12 14:30:45" or "Never" for invalid dates

// In React components - inline version
<td>
  {record.last_update_date && !record.last_update_date.startsWith('0001-')
    ? new Date(record.last_update_date).toISOString().split('T')[0]
    : '-'}
</td>
```

**Important**: Go's `time.Time` zero value serializes to `"0001-01-01T00:00:00Z"` instead of `null`. 
Always check for dates starting with `'0001-'` and treat them as "no date" by displaying `-` or `Never`.

#### ❌ INCORRECT Examples

```typescript
// DON'T use locale-specific formatting
new Date(dateString).toLocaleDateString()  // ❌ Outputs: "2/12/2026" (ambiguous!)
new Date(dateString).toLocaleString()     // ❌ Locale-dependent
new Date(dateString).toDateString()       // ❌ Outputs: "Wed Feb 12 2026" (verbose)
```

### Number Formatting
For numbers (not dates), using `.toLocaleString()` is acceptable for thousand separators:

```typescript
// ✅ CORRECT for numbers
<p>{totalRecords.toLocaleString()}</p>  // 3,211,232
<p>{amount.toLocaleString()}</p>        // 1,234,567.89
```

## React and TypeScript Best Practices

### Component Structure
- Use functional components with TypeScript
- Define interfaces for all props and data structures
- Use `'use client'` directive when component needs client-side interactivity

### State Management
- Use `useState` for local component state
- Use `useEffect` for side effects and data fetching
- Clean up effects with return functions when necessary

### API Calls
- Always use environment variables for API base URLs
- Handle loading, error, and success states explicitly
- Use try-catch blocks for all async operations

```typescript
const API_BASE_URL = typeof window !== 'undefined' 
  ? (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:18080')
  : 'http://backend:8080'
```

### Error Handling
- Display user-friendly error messages
- Differentiate between warning notices and critical errors
- Provide actionable guidance when possible

## Styling Guidelines

### Tailwind CSS Usage
- Use Tailwind utility classes consistently
- Follow the glassmorphism design pattern: `bg-white/5 backdrop-blur-sm border-2 border-white/10`
- Use opacity utilities for secondary text: `opacity-70`

### Dark Mode Support
- All components must support dark mode by default
- Use transparent backgrounds with opacity: `bg-white/5`
- Avoid hardcoded light-mode colors like `bg-white`, `bg-gray-50`, `text-gray-900`
- Include `<ThemeToggle />` component in page headers
- **Dropdowns/Select Elements**: Add explicit dark styling to both select and option elements:
  ```tsx
  <select className="bg-white/5 text-white border-white/20">
    <option className="bg-gray-800 text-white">Option 1</option>
    <option className="bg-gray-800 text-white">Option 2</option>
  </select>
  ```

### Responsive Design
- Use responsive grid classes: `grid-cols-1 md:grid-cols-2 lg:grid-cols-3`
- Mobile-first approach with breakpoints
- Ensure tables are scrollable on mobile: `overflow-x-auto`

## Accessibility

### ARIA Labels
- Use semantic HTML elements (`<button>`, `<nav>`, `<main>`)
- Add `aria-label` for icon-only buttons
- Use proper heading hierarchy (`<h1>`, `<h2>`, etc.)

### Keyboard Navigation
- Ensure all interactive elements are keyboard accessible
- Use `:focus` styles for focus indicators
- Disable buttons appropriately with `disabled` attribute

### Form Accessibility
- Use `<label>` elements for all form inputs
- Include placeholder text as guidance
- Show validation errors clearly

## Performance

### Component Optimization
- Use React.memo for expensive re-renders
- Implement pagination for large data sets
- Lazy load components when appropriate

### Asset Optimization
- Use Next.js Image component for images
- Minimize bundle size by importing only what's needed
- Use dynamic imports for heavy components

## Code Organization

### File Structure
- One component per file
- Co-locate related components in subdirectories
- Use descriptive, kebab-case filenames

### Import Organization
```typescript
// 1. React and Next.js imports
import { useState, useEffect } from 'react'
import Link from 'next/link'

// 2. Third-party imports
import ThemeToggle from '../components/ThemeToggle'

// 3. Types and interfaces
interface MyData {
  id: string
  name: string
}

// 4. Component definition
export default function MyComponent() {
  // ...
}
```

### Naming Conventions
- Components: PascalCase (e.g., `ThemeToggle.tsx`)
- Variables and functions: camelCase (e.g., `fetchRecords`, `currentPage`)
- Constants: UPPER_SNAKE_CASE (e.g., `API_BASE_URL`, `PAGE_SIZE`)
- Interfaces: PascalCase with descriptive names (e.g., `LEIRecord`, `ProcessingStatus`)

## Testing Guidelines

### Component Testing
- Test user interactions
- Test loading and error states
- Test accessibility features

### API Integration Testing
- Mock API responses in tests
- Test error handling
- Verify data transformations

## Documentation

### Code Comments
- Follow the self-explanatory code guidelines
- Document complex business logic
- Explain non-obvious TypeScript types

### Component Documentation
- Add JSDoc comments for reusable components
- Document props with descriptions
- Include usage examples for shared components

## Security

### XSS Prevention
- Never use `dangerouslySetInnerHTML` without sanitization
- Validate and sanitize all user inputs
- Use parameterized queries for API calls

### Authentication
- Store tokens securely (httpOnly cookies preferred)
- Include authentication headers in protected API calls
- Handle token expiration gracefully

## Common Patterns

### Pagination
```typescript
const [currentPage, setCurrentPage] = useState(1)
const [pageSize] = useState(50)

const fetchRecords = async () => {
  const offset = (currentPage - 1) * pageSize
  const response = await fetch(
    `${API_BASE_URL}/api/v1/resource?limit=${pageSize}&offset=${offset}`
  )
  // ...
}
```

### Search and Filters
```typescript
const [searchTerm, setSearchTerm] = useState('')
const [filters, setFilters] = useState({})

useEffect(() => {
  fetchRecords()
}, [currentPage, searchTerm, filters])  // Re-fetch when filters change
```

### Loading States
```typescript
if (loading && records.length === 0) {
  return (
    <div className="text-center py-20">
      <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      <p className="mt-4 opacity-70">Loading...</p>
    </div>
  )
}
```

## References

- [Next.js Documentation](https://nextjs.org/docs)
- [React TypeScript Cheatsheet](https://react-typescript-cheatsheet.netlify.app/)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [ISO 8601 Date Format](https://en.wikipedia.org/wiki/ISO_8601)
- [Web Content Accessibility Guidelines (WCAG)](https://www.w3.org/WAI/WCAG21/quickref/)

---

## Summary

- **Dates**: Always ISO 8601 format (yyyy-mm-dd) - NEVER use toLocaleDateString()
- **Styling**: Glassmorphism dark mode by default
- **TypeScript**: Strong typing for all data structures
- **Accessibility**: Semantic HTML and keyboard navigation
- **Performance**: Pagination and lazy loading for large datasets
