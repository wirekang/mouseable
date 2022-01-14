import * as monaco from "monaco-editor";
import { KeyCode } from "monaco-editor";
import store from "./store";

export async function initEditor() {
  const schema = await window.__loadSchema();
  // @ts-ignore
  self.MonacoEnvironment = {
    getWorkerUrl: function (moduleId: any, label: string) {
      if (label === "json") {
        return "./json.worker.bundle.js";
      }
      if (label === "css" || label === "scss" || label === "less") {
        return "./css.worker.bundle.js";
      }
      if (label === "html" || label === "handlebars" || label === "razor") {
        return "./html.worker.bundle.js";
      }
      if (label === "typescript" || label === "javascript") {
        return "./ts.worker.bundle.js";
      }
      return "./editor.worker.bundle.js";
    },
  };

  const modelUri = monaco.Uri.parse("a://b/foo.json");
  const model = monaco.editor.createModel("", "json", modelUri);
  monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
    validate: true,
    schemaValidation: "error",
    schemas: [
      {
        uri: "a://b/foo.json",
        schema: JSON.parse(schema),
        fileMatch: [modelUri.toString()],
      },
    ],
  });

  const editor = monaco.editor.create(
    document.getElementById("my-monaco")!,
    {
      fontFamily: "D2Coding",
      language: "json",
      formatOnPaste: true,
      formatOnType: true,
      minimap: { enabled: false },
      stablePeek: true,
      suggest: {
        preview: true,
      },
    },
    {
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
    },
  );

  editor.setModel(model);
  editor.addCommand(KeyCode.F1, () => {
    downKey("F1");
  });
  editor.addCommand(KeyCode.F2, () => {
    downKey("F2");
  });
  editor.addCommand(KeyCode.F3, () => {
    downKey("F3");
  });
  store.editor = editor;
}

export function focusMonacoEditor() {
  $("textarea.monaco-mouse-cursor-text").trigger("focus");
}

function downKey(key: string) {
  window.dispatchEvent(new KeyboardEvent("keydown", { key }));
}
