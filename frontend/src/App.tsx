import {useEffect, useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {ConnectedToOrca, TryCreateClient} from "../wailsjs/go/main/App";

function App() {
  const [connected, setConnected] = useState(false);
  const [connectedMessage, setConnectedMessage] = useState('');

  useEffect(() => {
    // Run once on startup
    TryCreateClient().then(res => setConnectedMessage(res));

    // Poll every 5 seconds
    const interval = setInterval(() => {
      ConnectedToOrca().then(res => setConnected(res));
    }, 5000);

    // Cleanup
    return () => clearInterval(interval);
  }, []);

  return (
    <div id="App">
      <img src={logo} id="logo" alt="logo"/>
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
