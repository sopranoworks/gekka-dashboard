import type { Theme } from '../theme'

interface Props {
  theme: Theme
  onToggle: () => void
}

export function ThemeToggle({ theme, onToggle }: Props) {
  return (
    <button
      onClick={onToggle}
      style={{
        border: `1px solid ${theme.border}`,
        borderRadius: 8,
        background: theme.bgPanel,
        color: theme.text,
        padding: '6px 12px',
        cursor: 'pointer',
        fontSize: 14,
      }}
    >
      {theme.name === 'dark' ? '\u2600\uFE0F' : '\uD83C\uDF19'}
    </button>
  )
}
