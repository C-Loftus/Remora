import { useEffect, useState } from 'react';
import './App.css';
import { ConnectionStatus } from "../wailsjs/go/main/App";

function App() {
  const [connected, setConnected] = useState(false);
  const [connectedMessage, setConnectedMessage] = useState('');

  useEffect(() => {
    let isMounted = true;

    async function fetchStatus() {
      try {
        const status = await ConnectionStatus();
        if (!isMounted) return;

        setConnected(status.ConnectedToOrca);
        setConnectedMessage(status.ConnectionMessage);
      } catch (err) {
        console.error("Error fetching connection status:", err);
        if (isMounted) {
          setConnected(false);
          setConnectedMessage("Error connecting to Orca");
        }
      }
    }

    fetchStatus();

    const interval = setInterval(fetchStatus, 1000);

    return () => {
      isMounted = false;
      clearInterval(interval);
    };
  }, []);

  return (
    <div id="App">
      <h1>Orca Helper</h1>
      <div className="result">
        {connected ? 'Connected to Orca' : 'Not connected to Orca'}
      </div>
      <div className="result">
        {connectedMessage}
      </div>
    </div>
  );
}

export default App;
