import { Range } from "monaco-editor";
import { focusMonacoEditor } from "./editor";
import store from "./store";
import { errNoLoadedConfigName } from "./cnst";
import { showMsg } from "./toast";
import { renderAppliedConfig } from "./config";

export function initHotkeys() {
  window.addEventListener("keydown", (e) => {
    if (e.key === "F1" || e.key === "F2" || e.key === "F3") {
      e.preventDefault();

      switch (e.key) {
        case "F1":
          getNextKey();
          break;
        case "F2":
          saveConfig();
          break;
        case "F3":
          applyConfig();
          break;
      }
    }
  });
}

async function getNextKey() {
  focusMonacoEditor();
  const key = await window.__getNextKey();
  const line = store.editor!.getPosition();
  if (!line) {
    return;
  }

  const range = new Range(line.lineNumber, line.column, line.lineNumber, line.column);
  const id = { major: 1, minor: 1 };
  const op = { identifier: id, range: range, text: key };
  store.editor!.executeEdits("my-source", [op]);
}

async function saveConfig() {
  const value = store.editor!.getValue();
  if (!store.loadedConfigName) {
    throw errNoLoadedConfigName;
  }

  await window.__saveConfig(store.loadedConfigName, value);
  showMsg(`${store.loadedConfigName} Saved.`);
}

async function applyConfig() {
  if (!store.loadedConfigName) {
    throw errNoLoadedConfigName;
  }

  await saveConfig();
  await window.__applyConfig(store.loadedConfigName);
  showMsg(`${store.loadedConfigName} Applied.`);
  renderAppliedConfig(store.loadedConfigName);
}
