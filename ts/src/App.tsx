import React, { useEffect, useState } from "react";
import MyEditor from "./components/MyEditor";
import TopRow from "./components/TopRow";
import EditorHelp from "./components/EditorHelp";
import { useAsync, useAsyncFn, useInterval } from "react-use";
import DelayedNotRunning from "./components/DelayedNotRunning";
import HalfRow from "./components/HalfRow";
import ConfigPanel from "./ConfigPanel";
import { editor } from "monaco-editor";
import { useMonaco } from "@monaco-editor/react";
import { applyConfig, getNextKey, loadConfig, saveConfig } from "./util/go";

function App() {
  const versionState = useAsync(window.__getVersion);
  const [pingState, requestPing] = useAsyncFn(window.__ping);
  const schemaState = useAsync(window.__loadSchema);
  const [appliedConfigNameState, requestAppliedConfigName] = useAsyncFn(window.__loadAppliedConfigName);
  const [configNamesState, requestConfigNames] = useAsyncFn(window.__loadConfigNames);
  const [editorValue, setEditorValue] = useState<string>();
  const [editor, setEditor] = useState<editor.IStandaloneCodeEditor>();
  const [loadedConfigName, setLoadedConfigName] = useState<string>();
  const monaco = useMonaco();

  const handleOnLoadConfig = (configName: string) => {
    loadConfig(configName, setLoadedConfigName, setEditorValue);
  };

  const handleOnGetNextKey = () => {
    if (!editor || !monaco) {
      console.log(editor, monaco);
      return;
    }

    getNextKey(editor, monaco);
  };

  const handleOnSave = () => {
    if (!editor) {
      return;
    }

    saveConfig(loadedConfigName, editorValue, editor);
  };

  const handleOnApply = () => {
    applyConfig(loadedConfigName);
  };

  useEffect(() => {
    requestAppliedConfigName();
    requestConfigNames();
  }, []);

  useEffect(() => {
    if (appliedConfigNameState.value && !loadedConfigName) {
      handleOnLoadConfig(appliedConfigNameState.value);
    }
  }, [appliedConfigNameState.loading, handleOnLoadConfig, loadedConfigName]);

  useInterval(() => {
    requestPing();
  }, 2000);

  useEffect(() => {
    if (appliedConfigNameState.value) {
    }
  }, [appliedConfigNameState.value]);

  useEffect(() => {
    if (!editor || !monaco) {
      return;
    }

    editor.addCommand(monaco.KeyCode.F1, handleOnGetNextKey);
    editor.addCommand(monaco.KeyCode.F2, handleOnSave);
    editor.addCommand(monaco.KeyCode.F3, handleOnApply);
  }, [editor, monaco, handleOnGetNextKey, handleOnSave, handleOnApply]);

  return (
    <div style={{ height: "100%" }}>
      {pingState.loading && <DelayedNotRunning delay={1000} />}
      <TopRow version={versionState.value ?? ""} />
      <HalfRow>
        <EditorHelp />
        <ConfigPanel
          configNames={configNamesState.value}
          loadedConfigName={loadedConfigName}
          onLoadConfig={handleOnLoadConfig}
          appliedConfigName={appliedConfigNameState.value}
        />
      </HalfRow>
      {schemaState.value && (
        <MyEditor
          value={editorValue}
          onChange={setEditorValue}
          schema={schemaState.value}
          onMount={setEditor}
        />
      )}
    </div>
  );
}

export default App;
