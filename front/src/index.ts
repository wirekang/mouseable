import "toastify-js/src/toastify.css";
import "./style.css";
import { showError } from "./utils/toast";
import "./gobind";
import { initHotkeys } from "./utils/hotkey";
import { initEditor } from "./utils/editor";
import { initButtons, initErrorHandle, initModal, initPing, renderVersion } from "./utils/ect";
import { initConfig } from "./utils/config";

$(() => {
  (async () => {
    try {
      initButtons();
      initPing();
      initErrorHandle();
      initModal();
      initConfig();
      renderVersion();
      await initEditor();
      initHotkeys();
    } catch (e) {
      showError(e);
    }
  })();
});
