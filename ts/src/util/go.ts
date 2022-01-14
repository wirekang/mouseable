import { editor } from "monaco-editor";
import toast from "react-simple-toasts";
import { Monaco } from "@monaco-editor/react";

export async function getNextKey(editor: editor.IStandaloneCodeEditor, monaco: Monaco): Promise<void> {
  try {
    toast("Press any keys include double press.");
    const key = await window.__getNextKey();
    if (!key) {
      return;
    }

    const line = editor.getPosition();
    if (!line) {
      return;
    }

    const range = new monaco.Range(line.lineNumber, line.column, line.lineNumber, line.column);
    const id = { major: 1, minor: 1 };
    const op = { identifier: id, range: range, text: `"${key}"` };
    editor.executeEdits("my-source", [op]);
  } catch (e) {
    toast(`${e}`);
  }
}

export async function saveConfig(
  loadedConfigName: string | undefined,
  editorValue: string | undefined,
  editor: editor.IStandaloneCodeEditor,
) {
  try {
    editor.trigger("anyString", "editor.action.formatDocument", null);
    if (!loadedConfigName) {
      toast("Nothing was loaded.");
      return;
    }
    await window.__saveConfig(loadedConfigName, editorValue ?? "");
    toast("Saved");
  } catch (e) {
    toast(`${e}`);
  }
}

export async function loadConfig(
  name: string,
  setLoadedConfigName: (v: string) => void,
  setEditorValue: (v: string) => void,
) {
  try {
    const json = await window.__loadConfig(name);
    setLoadedConfigName(name);
    setEditorValue(json);
  } catch (e) {
    toast(`${e}`);
  }
}

export async function applyConfig(name: string | undefined) {
  try {
    if (!name) {
      toast("Nothing was loaded.");
      return;
    }

    await window.__applyConfig(name);
    toast(`${name} Applied`);
  } catch (e) {
    toast(`${e}`);
  }
}
