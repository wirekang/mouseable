import React, { useEffect, useState } from "react";
import MyEditor from "./components/MyEditor";
import TopRow from "./components/TopRow";
import { editor } from "monaco-editor";
import EditorHelp from "./components/EditorHelp";
import { useAsync, useAsyncFn, useInterval } from "react-use";
import DelayedNotRunning from "./components/DelayedNotRunning";

function App() {
  const versionState = useAsync(window.__getVersion__);
  const [pingState, requestPing] = useAsyncFn(window.__ping__);
  const schemaState = useAsync(window.__getSchema__);
  const [jsonValue, setJsonValue] = useState<string>();
  const [editor, setEditor] = useState<editor.IStandaloneCodeEditor | null>(null);

  useInterval(() => {
    requestPing();
  }, 2000);

  return (
    <div style={{ height: "100%" }}>
      {pingState.loading && <DelayedNotRunning delay={1000} />}
      <TopRow version={""} />
      <EditorHelp />
      {schemaState.value && (
        <MyEditor value={jsonValue} onChange={setJsonValue} onMount={setEditor} schema={schemaState.value} />
      )}
    </div>
  );
}

export default App;
