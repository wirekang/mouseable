import React from "react";
import Editor, { OnMount } from "@monaco-editor/react";
import { editor } from "monaco-editor";

interface Props {
  value?: string;
  onChange: (v?: string) => void;
  onMount: (e: editor.IStandaloneCodeEditor) => void;
}

export default function MyEditor(props: Props): JSX.Element {
  const onMount: OnMount = (editor, monaco) => {
    props.onMount(editor);
  };

  return (
    <Editor
      height="80vh"
      options={{ wordWrap: "off", formatOnPaste: true, formatOnType: true, minimap: { enabled: false } }}
      language="json"
      value={props.value}
      onChange={props.onChange}
      onMount={onMount}
    />
  );
}
