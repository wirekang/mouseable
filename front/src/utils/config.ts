import store from "./store";

const configSelect = $<HTMLSelectElement>("#config-select");
async function renderConfigNames() {
  const configNames = await window.__loadConfigNames();
  configSelect.empty();
  configNames.forEach((name) => {
    configSelect.append(`<option value="${name}">${name}</option>`);
  });
}

export async function initConfig() {
  await renderConfigNames();
  configSelect.on("change", async (e) => {
    const configName = e.target.value;
    const json = await window.__loadConfig(configName);
    store.loadedConfigName = configName;
    store.editor!.setValue(json);
  });

  const appliedConfigName = await window.__loadAppliedConfigName();
  store.loadedConfigName = appliedConfigName;
  renderAppliedConfig(appliedConfigName);
  loadConfig(appliedConfigName);
}

export function renderAppliedConfig(appliedConfigName: string) {
  $("#applied").text(appliedConfigName);
}

export function loadConfig(loadedConfigName: string) {
  configSelect.val(loadedConfigName);
  configSelect.trigger("change");
}
