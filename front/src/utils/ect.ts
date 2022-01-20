import { showError } from "./toast";
import store from "./store";

export function initModal() {
  $("#more-modal").hide();
  $("#examples-modal").hide();
  $("#more-modal-button").on("click", () => {
    $("#more-modal").toggle();
  });
  $("#examples-button").on("click", () => {
    $("#examples-modal").toggle();
  });
}

export function initErrorHandle() {
  const h = (e: any) => {
    showError(e.message ?? e.reason ?? e);
  };

  window.addEventListener("error", h);
  window.addEventListener("unhandledrejection", h);
  window.onerror = h;
}

export function initPing() {
  const nc = $("#not-connected");
  nc.hide();

  setInterval(async () => {
    await window.__ping();
    store.lastPing = Date.now();
  }, 3000);

  setInterval(() => {
    if (Date.now() - store.lastPing > 10000) {
      nc.show();
    }
  }, 1000);
}

export async function renderVersion() {
  const version = await window.__getVersion();
  $("#version").text(version);
}

export function initButtons() {
  $("#github-button").on("click", () => {
    window.__openLink("github.com/wirekang/mouseable");
  });

  $("#terminate-button").on("click", () => {
    if (window.confirm("Terminate mouseable completely? You can close this window without terminate it.")) {
      window.__terminate();
      window.close();
    }
  });
}
