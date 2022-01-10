import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";
import { isDev } from "./gobind";
import toast from "react-simple-toasts";

if (isDev) {
  console.log("DEV MODE");
}

window.addEventListener("keydown", (e) => {
  if (e.key === "F1" || e.key === "F2") {
    e.preventDefault();
    toast("Please focus on the editor. (Press Tab)");
  }
});

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById("root"),
);
