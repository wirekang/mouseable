import { editor } from "monaco-editor";
import toast from "react-simple-toasts";

export function registerHotkeys(e: editor.IStandaloneCodeEditor, m: any) {
  e.addCommand(m.KeyCode.F1, () => insertKey(e));
  e.addCommand(m.KeyCode.F2, () => save(e));
}

async function insertKey(e: editor.IStandaloneCodeEditor) {
  try {
    toast("Press any key. You can double press, or press the Ctrl, Shift, Alt key together.", { time: 5000 });
    const key = await window.__getNextKey__();
    e.trigger("keyboard", "type", { text: `"${key}"` });
  } catch (err) {
    toast(`${err}`);
  }
}

async function save(e: editor.IStandaloneCodeEditor) {
  try {
    e.trigger("anyString", "editor.action.formatDocument", null);
    await window.__saveConfig__(e.getValue());
  } catch (err) {
    toast(`${err}`);
  }
}
