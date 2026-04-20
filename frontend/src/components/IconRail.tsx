import type { Theme } from '../theme'

interface Props {
  theme: Theme
  active: string
  onNavigate: (view: string) => void
}

const navItems = [
  { id: 'cluster', icon: '\u2B21', label: 'Cluster' },
  { id: 'nodes', icon: '\u25C9', label: 'Nodes' },
  { id: 'shards', icon: '\u25EB', label: 'Shards' },
  { id: 'alerts', icon: '\u26A1', label: 'Alerts' },
  { id: 'settings', icon: '\u2699', label: 'Settings' },
]

export function IconRail({ theme, active, onNavigate }: Props) {
  return (
    <nav style={{
      width: 56,
      background: theme.bgRail,
      borderRight: `1px solid ${theme.border}`,
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      paddingTop: 16,
      gap: 8,
    }}>
      <div style={{
        fontSize: 20,
        color: theme.accent,
        marginBottom: 24,
        fontWeight: 'bold',
      }}>G</div>
      {navItems.map(item => (
        <button
          key={item.id}
          onClick={() => onNavigate(item.id)}
          title={item.label}
          style={{
            width: 40,
            height: 40,
            border: 'none',
            borderRadius: 8,
            background: active === item.id ? theme.accent : 'transparent',
            color: active === item.id ? '#fff' : theme.textMuted,
            fontSize: 18,
            cursor: 'pointer',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            transition: 'background 0.15s',
          }}
        >
          {item.icon}
        </button>
      ))}
    </nav>
  )
}
