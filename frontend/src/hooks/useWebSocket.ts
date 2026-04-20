import { useEffect, useRef, useState } from 'react'

export type ConnectionStatus = 'connecting' | 'connected' | 'disconnected'

export interface WebSocketMessage {
  type: string
  data?: unknown
}

export function useWebSocket(url: string) {
  const [status, setStatus] = useState<ConnectionStatus>('connecting')
  const [lastMessage, setLastMessage] = useState<WebSocketMessage | null>(null)
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    const ws = new WebSocket(url)
    wsRef.current = ws

    ws.onopen = () => setStatus('connected')
    ws.onclose = () => {
      setStatus('disconnected')
      setTimeout(() => setStatus('connecting'), 3000)
    }
    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data) as WebSocketMessage
        setLastMessage(msg)
      } catch {
        // ignore non-JSON
      }
    }

    return () => { ws.close() }
  }, [url])

  return { status, lastMessage }
}
