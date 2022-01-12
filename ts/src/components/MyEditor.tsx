import React from "react";
import Editor, { OnMount } from "@monaco-editor/react";
import { registerHotkeys } from "../util/editor";

interface Props {
  schema: string;
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
          schema: JSON.parse(props.schema),
          fileMatch: [modelUri.toString()],
        },
      ],
    });
    editor.setModel(model);
    registerHotkeys(editor, monaco);
  };

  return (
    <Editor
      height="80vh"
      options={{
        formatOnPaste: true,
        formatOnType: true,
        minimap: { enabled: false },
        stablePeek: true,
        suggest: {
          preview: true,
        },
      }}
      language="json"
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
