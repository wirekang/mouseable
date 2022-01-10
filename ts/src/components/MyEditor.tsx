import React from "react";
import Editor, { OnMount } from "@monaco-editor/react";
import { editor, Uri } from "monaco-editor";

interface Props {
  value?: string;
  onChange: (v?: string) => void;
  onMount: (e: editor.IStandaloneCodeEditor) => void;
  schema?: string;
}

export default function MyEditor(props: Props): JSX.Element {
  const onMount: OnMount = (editor, monaco) => {
    const modelUri = monaco.Uri.parse("a://b/foo.json");
    const model = monaco.editor.createModel("", "json", modelUri);
    monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
      validate: true,
      schemaValidation: "error",
      schemas: [
        {
          uri: "a://b/foo.json",
          schema: JSON.parse(props.schema ?? "{}"),
          fileMatch: [modelUri.toString()],
        },
      ],
    });
    editor.setModel(model);
    props.onMount(editor);
  };

  return (
    <Editor
      height="80vh"
      options={{
        wordWrap: "off",
        formatOnPaste: true,
        formatOnType: true,
        minimap: { enabled: false },
        stablePeek: true,
        suggest: {
          insertMode: "replace",
          preview: true,
        },
      }}
      language="json"
      value={props.value}
      onChange={props.onChange}
      onMount={onMount}
      overrideServices={{
        storageService: {
          get() {},
          getNumber() {},
          getBoolean(key: any) {
            return key === "expandSuggestionDocs";
          },
          remove() {},
          store() {},
          onDidChangeStorage() {},
          onWillSaveState() {},
        },
      }}
    />
  );
}
