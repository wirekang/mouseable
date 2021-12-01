import { FunctionKey } from "../gobind";
import { fromVKCode } from "win-vk";

export function functionKeyToText(k: FunctionKey): string {
  let s = "";

  if (k.IsWin) {
    s += "<W> + ";
  }

  if (k.IsControl) {
    s += "<C> + ";
  }

  if (k.IsAlt) {
    s += "<A> + ";
  }

  if (k.IsShift) {
    s += "<S> + ";
  }

  if (k.KeyCode !== 0) {
    s += fromVKCode(k.KeyCode);
  } else {
    s = s.substring(0, s.length - 3);
  }

  if (k.IsDouble) {
    s += " X2";
  }
  return s;
}
