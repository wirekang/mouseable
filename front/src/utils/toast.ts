import * as Toastify from "toastify-js";

export function showMsg(msg: any) {
  show(msg, {});
}

export function showError(msg: any) {
  show(msg, { background: "#880000" });
}

function show(msg: any, style: any) {
  const text = `${msg}`;
  Toastify({
    text,
    gravity: "top",
    position: "center",
    duration: Math.max(Math.min(text.length * 200, 5000), 1500),
    style,
  }).showToast();
}
