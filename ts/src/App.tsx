import React, { useEffect, useState } from "react";
import MyEditor from "./components/MyEditor";
import TopRow from "./components/TopRow";
import { editor } from "monaco-editor";
import EditorHelp from "./components/EditorHelp";
import TMI from "./components/TMI";
import { useAsyncFn, useInterval } from "react-use";
import { ping } from "./gobind";
import DelayedNotRunning from "./components/DelayedNotRunning";

function App() {
  const [pingState, requestPing] = useAsyncFn(ping);
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
      <MyEditor value={jsonValue} onChange={setJsonValue} onMount={setEditor} schema={"{}"} />
      <TMI />
    </div>
  );
}

export default App;
