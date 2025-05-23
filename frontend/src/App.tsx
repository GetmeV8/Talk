import React from 'react'
import { SocketProvider } from './contexts/SocketContext'
import { Chat } from './components/Chat'

function App() {
  return (
    <SocketProvider>
      <Chat />
    </SocketProvider>
  )
}

export default App 