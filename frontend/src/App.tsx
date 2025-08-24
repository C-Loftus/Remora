import { useEffect, useState } from 'react';
import './App.css';
import { ConnectionStatus, GetHotKeys } from "../wailsjs/go/main/App";
import { main } from '../wailsjs/go/models';

function App() {
  const [connected, setConnected] = useState(false);
  const [connectedMessage, setConnectedMessage] = useState('');

  const [hotkeys, setHotkeys] = useState<Array<main.HotkeyWithMetadata>>  ([]);

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

      try {
        const hotkeys = await GetHotKeys();
        if (!isMounted) return;
        setHotkeys(hotkeys);
      } catch (err) {
        console.error("Error fetching hotkeys:", err);
        if (isMounted) {
          setHotkeys([]);
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
      <a href='https://github.com/C-Loftus/orca-helper'>github</a>
      <h1>Orca Helper</h1>
      <p className="result">
        {connected ? 'Connected to Orca' : 'Not connected to Orca'}
      </p>
      <p className="result">
        {connectedMessage}
      </p>
      <h2> Keyboard Shortcuts </h2>
      <ul>
        {hotkeys.map((hotkey: main.HotkeyWithMetadata) => (
          // JSON SERIALIZED HOTKEY
          <p key={Math.random()}>
            {Object.keys(hotkey)}
          </p>
        ))}
      </ul>
    </div>
  );
}

export default App;
