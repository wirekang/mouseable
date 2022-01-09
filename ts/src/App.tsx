import React, { useState } from "react";
import { loadBind } from "./gobind";
import { useAsync } from "react-use";
import MyEditor from "./components/MyEditor";
import TopRow from "./components/TopRow";
import { editor } from "monaco-editor";
import EditorHelp from "./components/EditorHelp";

function App() {
  const [requesterCounter, setRequesterCounter] = useState(0);
  const goBindState = useAsync(loadBind, [requesterCounter]);
  const [value, setValue] = useState<string>();
  const [editor, setEditor] = useState<editor.IStandaloneCodeEditor | null>(null);

  if (goBindState.loading || !goBindState.value) {
    return <h4>Loading...</h4>;
  }

  const goBind = goBindState.value!;

  return (
    <div style={{ height: "100%" }}>
      <TopRow version={goBind.version} />
      <EditorHelp />
      <MyEditor value={value} onChange={setValue} onMount={setEditor} />
    </div>
  );
}

export default App;
