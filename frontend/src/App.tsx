import { useState } from 'react'
import { darkTheme, lightTheme, type Theme } from './theme'
import { IconRail } from './components/IconRail'
import { GridPanel } from './components/GridPanel'
import { ThemeToggle } from './components/ThemeToggle'
import { useWebSocket, type ConnectionStatus } from './hooks/useWebSocket'

function StatusDot({ status, theme }: { status: ConnectionStatus; theme: Theme }) {
  const color = status === 'connected' ? theme.success
    : status === 'connecting' ? theme.warning
    : theme.error
  return (
    <span style={{ display: 'inline-flex', alignItems: 'center', gap: 6, fontSize: 12, color: theme.textMuted }}>
      <span style={{ width: 8, height: 8, borderRadius: '50%', background: color }} />
      {status}
    </span>
  )
}

export function App() {
  const [theme, setTheme] = useState<Theme>(darkTheme)
  const [activeView, setActiveView] = useState('cluster')
  const wsUrl = `ws://${window.location.host}/ws`
  const { status } = useWebSocket(wsUrl)

  const toggleTheme = () => {
    setTheme(t => t.name === 'dark' ? lightTheme : darkTheme)
  }

  return (
    <div style={{
      display: 'flex',
      height: '100vh',
      background: theme.bg,
      color: theme.text,
      fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    }}>
      <IconRail theme={theme} active={activeView} onNavigate={setActiveView} />
      <main style={{ flex: 1, padding: 24, overflow: 'auto' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <h1 style={{ fontSize: 20, fontWeight: 600 }}>
              {activeView.charAt(0).toUpperCase() + activeView.slice(1)}
            </h1>
            <StatusDot status={status} theme={theme} />
          </div>
          <ThemeToggle theme={theme} onToggle={toggleTheme} />
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: 16 }}>
          <GridPanel theme={theme} title="Nodes">
            <p style={{ color: theme.textMuted }}>Cluster data will appear here</p>
          </GridPanel>
          <GridPanel theme={theme} title="Events">
            <p style={{ color: theme.textMuted }}>Real-time events will appear here</p>
          </GridPanel>
          <GridPanel theme={theme} title="Shards">
            <p style={{ color: theme.textMuted }}>Shard distribution will appear here</p>
          </GridPanel>
          <GridPanel theme={theme} title="Health">
            <p style={{ color: theme.textMuted }}>Health metrics will appear here</p>
          </GridPanel>
        </div>
      </main>
    </div>
  )
}
