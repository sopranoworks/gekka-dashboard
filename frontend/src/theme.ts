export interface Theme {
  name: 'dark' | 'light'
  bg: string
  bgPanel: string
  bgRail: string
  text: string
  textMuted: string
  accent: string
  accentHover: string
  border: string
  success: string
  warning: string
  error: string
}

export const darkTheme: Theme = {
  name: 'dark',
  bg: '#0f0f1a',
  bgPanel: '#1a1a2e',
  bgRail: '#12122b',
  text: '#e0e0e0',
  textMuted: '#888888',
  accent: '#00897B',
  accentHover: '#00a693',
  border: '#2a2a4a',
  success: '#43A047',
  warning: '#FFA726',
  error: '#E53935',
}

export const lightTheme: Theme = {
  name: 'light',
  bg: '#f5f5f5',
  bgPanel: '#ffffff',
  bgRail: '#fafafa',
  text: '#1a1a1a',
  textMuted: '#666666',
  accent: '#00897B',
  accentHover: '#00695C',
  border: '#e0e0e0',
  success: '#2E7D32',
  warning: '#F57C00',
  error: '#C62828',
}
