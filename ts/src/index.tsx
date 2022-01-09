import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";

window.addEventListener("keydown", (e) => {
  if (e.key === "F1") {
    e.preventDefault();
  }
});

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById("root"),
);
