import { useState } from "react";

function App() {
  const [apiStatus, setApiStatus] = useState("Not checked");

  async function checkApi() {
    const url = "/api/health";
    setApiStatus("Checking...");
    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
      }
      const result = await response.json();
      setApiStatus(result.status);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      setApiStatus(`Error: ${message}`);
    }
  }

  return (
    <>
      <h1>Gym app</h1>
      <button onClick={checkApi}>Check API</button>
      <p>{apiStatus}</p>
    </>
  );
}

export default App;
