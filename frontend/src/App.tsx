import { useEffect, useState } from 'react';
import './App.css';
import { ConnectionStatus, GetDisplayServerType, GetHotKeys } from "../wailsjs/go/main/App";
import { main } from '../wailsjs/go/models';

function App() {
  const [connected, setConnected] = useState(false);
  const [connectedMessage, setConnectedMessage] = useState('');
  const [displayServerType, setDisplayServerType] = useState("unknown");

  const [hotkeys, setHotkeys] = useState<Array<string>>  ([]);

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

      try {
        const displayServerType = await GetDisplayServerType();
        if (!isMounted) return;
        setDisplayServerType(displayServerType);
      } catch (err) {
        console.error("Error fetching display server type:", err);
        if (isMounted) {
          setDisplayServerType("unknown");
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
      <a href='https://github.com/C-Loftus/orca-helper' target="_blank">
            <h1>Orca Helper</h1>
      </a>
      <p className="result">
        {connected ? 'Connected to Orca' : 'Not connected to Orca'}
      </p>
      <p className="result">
        {connectedMessage}
      </p>
      <h2> Keyboard Shortcuts </h2>
      {displayServerType === "x11" ? (
      <table style={{margin: 'auto'}}>
        <thead>
          <tr>
            <th style={{textAlign: 'center'}}>Description</th>
            <th style={{textAlign: 'center'}}>Key</th>
          </tr>
        </thead>
        <tbody>
          {hotkeys.map((hotkey) => {
            const [description, keyName] = hotkey.split(':');
            return (
              <tr key={hotkey}>
                <td style={{textAlign: 'center'}}>{description}</td>
                <td style={{textAlign: 'center'}}>{keyName}</td>
              </tr>
            )
          })}
        </tbody>
      </table>
      ) : (
        <p className="result">
          Your system is running Wayland but Wayland does not support global keyboard shortcuts. Please switch to X11 for keyboard shortcuts to work.
        </p>
      )}
      
    </div>
  );
}

export default App;
