import type { ReactNode } from 'react'
import type { Theme } from '../theme'

interface Props {
  theme: Theme
  title: string
  children: ReactNode
}

export function GridPanel({ theme, title, children }: Props) {
  return (
    <div style={{
      background: theme.bgPanel,
      border: `1px solid ${theme.border}`,
      borderRadius: 12,
      padding: 20,
    }}>
      <h3 style={{ fontSize: 14, fontWeight: 600, marginBottom: 12, color: theme.text }}>{title}</h3>
      {children}
    </div>
  )
}
