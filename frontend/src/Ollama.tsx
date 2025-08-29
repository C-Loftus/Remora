import { useEffect, useState } from "react";
import {
    GetModelName,
    GetPrompt,
    LastOllamaResponse,
    OllamaConnectionStatus,
    OllamaProcessing,
    SetPrompt
} from "../wailsjs/go/main/App";
import { default as loading } from "../src/assets/loading.svg";

export default function Ollama() {
    const [lastOllamaResponse, setLastOllamaResponse] = useState<string | null>(null);
    const [ollamaStatusMessage, setOllamaStatusMessage] = useState<string | null>(null);
    const [currentPrompt, setCurrentPrompt] = useState<string>('');
    const [inputPrompt, setInputPrompt] = useState<string>('');
    const [currentVisionModel, setCurrentVisionModel] = useState<string>('');
    const [processing, setProcessing] = useState(false);

    // Fetch prompt once on startup
    useEffect(() => {
        let isMounted = true;
        async function fetchPrompt() {
            try {
                const prompt = await GetPrompt();
                if (!isMounted) return;
                setCurrentPrompt(prompt);
                setInputPrompt(prompt);
            } catch (err) {
                console.error("Error fetching current prompt:", err);
                if (isMounted) {
                    setCurrentPrompt("");
                    setInputPrompt("");
                }
            }
        }

        fetchPrompt();
        return () => { isMounted = false; };
    }, []);

    // Poll other status fields every 3 seconds
    useEffect(() => {
        let isMounted = true;

        async function fetchStatus() {
            try {
                const status = await OllamaConnectionStatus();
                if (!isMounted) return;
                setOllamaStatusMessage(status);
            } catch (err) {
                console.error("Error fetching Ollama status:", err);
                if (isMounted) setOllamaStatusMessage("");
            }

            try {
                const lastMessage = await LastOllamaResponse();
                if (!isMounted) return;
                setLastOllamaResponse(lastMessage);
            } catch (err) {
                console.error("Error fetching last Ollama message:", err);
                if (isMounted) setLastOllamaResponse("");
            }

            try {
                const model = await GetModelName();
                if (!isMounted) return;
                setCurrentVisionModel(model);
            } catch (err) {
                console.error("Error fetching model name:", err);
                if (isMounted) setCurrentVisionModel("");
            }

            try {
                const processingStatus = await OllamaProcessing();
                if (!isMounted) return;
                setProcessing(processingStatus);
            } catch (err) {
                console.error("Error fetching processing status:", err);
                if (isMounted) setProcessing(false);
            }
        }

        fetchStatus();
        const interval = setInterval(fetchStatus, 3000);

        return () => {
            isMounted = false;
            clearInterval(interval);
        };
    }, []);

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        await SetPrompt(inputPrompt);
        setCurrentPrompt(inputPrompt);
    };

    const handleCopy = () => {
        if (lastOllamaResponse) {
            navigator.clipboard.writeText(lastOllamaResponse)
        }
    };

    return (
        <>
            <h2>Ollama</h2>

            {processing ? <>
                <p>Currently Processing <img src={loading} width="30" height="30" alt="Loading" aria-hidden /></p>
            </> : (
                <p className="result">
                    {ollamaStatusMessage === null ? (
                        "Not connected; please install Ollama and make sure it is running"
                    ) : (
                        <>
                            {ollamaStatusMessage} using model '{currentVisionModel}'
                            <form onSubmit={handleSubmit}>
                                <input
                                    style={{ marginLeft: '10px', width: '300px' }}
                                    type="text"
                                    value={inputPrompt}
                                    onChange={(e) => setInputPrompt(e.target.value)}
                                    placeholder="Edit this to change the prompt for Ollama; this will be used for all subsequent OCR runs"
                                />
                                <button className="btn" type="submit">Set Prompt</button>
                            </form>
                        </>
                    )}
                </p>
            )}

            {currentPrompt && <p>Current prompt: {currentPrompt}</p>}

            {lastOllamaResponse && (
                <>
                    <h3>Last Ollama Response</h3>
                    <button onClick={handleCopy} className="btn">Copy to Clipboard</button>
                    <p>{lastOllamaResponse}</p>
                </>
            )}
        </>
    );
}
