import React from "react";

interface Props {
  isChecked?: boolean;
  onChange?: (b: boolean) => void;
  children: React.ReactText;
  style?: React.CSSProperties;
}

export default function Checkbox(props: Props): JSX.Element {
  return (
    <div
      style={{
        cursor: "pointer",
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        ...props.style,
      }}
      onClick={() => {
        props.onChange && props.onChange(!props.isChecked);
      }}
    >
      <span style={{}}>{props.children}</span>
      <input
        style={{
          pointerEvents: "none",
        }}
        type="checkbox"
        onClick={(e) => {
          e.preventDefault();
        }}
        checked={props.isChecked}
        readOnly
      />
    </div>
  );
}
