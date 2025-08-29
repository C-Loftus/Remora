import { useEffect, useState } from 'react';
import './App.css';
import { ConnectionStatus, GetDisplayServerType, GetHotKeys, GetModelName, GetPrompt, LastOcrResponse, LastOllamaResponse, OllamaConnectionStatus, SetPrompt } from "../wailsjs/go/main/App";
import Ollama from './Ollama';
import { default as remora } from './assets/images/remora.png'

function App() {
  const [connected, setConnected] = useState(false);
  const [connectedMessage, setConnectedMessage] = useState('');
  const [displayServerType, setDisplayServerType] = useState("unknown");
  const [lastOcrResponse, setLastOcrResponse] = useState<string | null>(null);
  const [hotkeys, setHotkeys] = useState<Array<string>>([]);

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

      try {
        const lastOcrMessage = await LastOcrResponse();
        if (!isMounted) return;
        setLastOcrResponse(lastOcrMessage);
      } catch (err) {
        console.error("Error fetching last ocr message:", err);
        if (isMounted) {
          setLastOcrResponse("");
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
      <img src={remora} tabIndex={-1} alt="The Remora Logo; an Orca whale with a fish swimming below it" width="30%" style={{float: 'left', marginRight: '1em'}} />
      <a href='https://github.com/C-Loftus/remora' target="_blank">
        <h1>Remora</h1>
      </a>
      <p className="result">
        {connected ? 'Connected to Orca ' : 'Not connected to Orca '}
        {connectedMessage}
      </p>
      <h2> Keyboard Shortcuts </h2>
      {displayServerType === "x11" ? (
        <table style={{ margin: 'auto' }}>
          <thead>
            <tr>
              <th style={{ textAlign: 'center' }}>Description</th>
              <th style={{ textAlign: 'center' }}>Key</th>
            </tr>
          </thead>
          <tbody>
            {hotkeys.map((hotkey) => {
              const [description, keyName] = hotkey.split(':');
              return (
                <tr key={hotkey}>
                  <td style={{ textAlign: 'center' }}>{description}</td>
                  <td style={{ textAlign: 'center' }}>{keyName}</td>
                </tr>
              )
            })}
          </tbody>
        </table>
      ) : (
        <p>
          Error: Your system is running Wayland but Wayland does not support global keyboard shortcuts. Please switch to X11 for keyboard shortcuts to work.
        </p>
      )}

      <Ollama />

      <h2 style={{ marginTop: '40px' }}> OCR </h2>
      <p className="result">
        {lastOcrResponse ? (
          <>
            <button onClick={
              () => {
                if (lastOcrResponse) {
                  navigator.clipboard.writeText(lastOcrResponse)
                }
              }
            } className="btn">Copy to Clipboard</button>
            <p>
              {lastOcrResponse}
            </p>
          </>
        ) : "No response yet"}
      </p>
    </div>
  );
}

export default App;
