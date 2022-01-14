import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";
import { isDev } from "./gobind";
import toast from "react-simple-toasts";
import { focusMonacoEditor } from "./util/dom";

if (isDev) {
  console.log("DEV MODE");
}

window.addEventListener("keydown", (e) => {
  if (e.key === "F1" || e.key === "F2" || e.key === "F3") {
    e.preventDefault();
    toast("Press again. (Changed focus to editor)");
    focusMonacoEditor();
  }
});

window.addEventListener("error", (e) => {
  toast(e.message);
});

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById("root"),
);
