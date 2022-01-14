import { editor } from "monaco-editor";

type Store = {
  loadedConfigName?: string;
  lastPing: number;
  editor?: editor.IStandaloneCodeEditor;
};

const store: Store = {
  lastPing: Date.now(),
};

export default store;
