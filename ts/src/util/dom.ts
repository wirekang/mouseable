export function focusMonacoEditor() {
  (document.querySelector("textarea.monaco-mouse-cursor-text") as any).focus();
}
