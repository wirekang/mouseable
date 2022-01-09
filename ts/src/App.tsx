import React, { useState } from "react";
import { loadBind } from "./gobind";
import { useAsync } from "react-use";
import GitHubButton from "./components/GitHubButton";
import TerminateButton from "./components/TerminateButton";

function App() {
  const [requesterCounter, setRequesterCounter] = useState(0);
  const goBindState = useAsync(loadBind, [requesterCounter]);

  if (goBindState.loading || !goBindState.value) {
    return <h4>Loading...</h4>;
  }

  const goBind = goBindState.value!;

  return (
    <div>
      <div
        style={{
          display: "flex",
          flexDirection: "row",
          alignItems: "center",
          justifyContent: "space-between",
          margin: 3,
        }}
      >
        <span>Version: {goBind.version}</span>
        <GitHubButton />
        <TerminateButton />
      </div>
    </div>
  );
}

export default App;
