import React from "react";
import MyEditor from "./components/MyEditor";
import TopRow from "./components/TopRow";
import EditorHelp from "./components/EditorHelp";
import { useAsync, useAsyncFn, useInterval } from "react-use";
import DelayedNotRunning from "./components/DelayedNotRunning";

function App() {
  const versionState = useAsync(window.__getVersion__);
  const [pingState, requestPing] = useAsyncFn(window.__ping__);
  const schemaState = useAsync(window.__getSchema__);

  useInterval(() => {
    requestPing();
  }, 2000);

  return (
    <div style={{ height: "100%" }}>
      {pingState.loading && <DelayedNotRunning delay={1000} />}
      <TopRow version={versionState.value ?? ""} />
      <EditorHelp />
      {schemaState.value && <MyEditor schema={schemaState.value} />}
    </div>
  );
}

export default App;
