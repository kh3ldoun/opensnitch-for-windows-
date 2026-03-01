import React from "react";
import ReactDOM from "react-dom/client";

function App() {
  const [response, setResponse] = React.useState("");

  const ping = async () => {
    // Requires @tauri-apps/api/core
    // const { invoke } = await import('@tauri-apps/api/core');
    // const res = await invoke("ping_backend");
    // setResponse(res as string);
  };

  return (
    <div style={{ fontFamily: "sans-serif", padding: "20px" }}>
      <h1>WinSnitch Interactive Firewall</h1>
      <button onClick={ping}>Ping Backend</button>
      <p>{response}</p>
    </div>
  );
}

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
