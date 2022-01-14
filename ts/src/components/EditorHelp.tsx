import React from "react";

interface Props {}

export default function EditorHelp(props: Props): JSX.Element {
  return (
    <div
      style={{
        fontSize: 12,
      }}
    >
      <ul>
        <li>
          Press <b>Ctrl+I</b> to show suggestions.
        </li>
        <li>
          Press <b>F1</b> to insert key.
        </li>
        <li>
          Press <b>F2</b> to save.
        </li>
        <li>
          Press <b>F3</b> to apply.
        </li>
      </ul>
    </div>
  );
}
