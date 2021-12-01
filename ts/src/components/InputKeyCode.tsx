import React from "react";

interface Props {
  keyCode: number;
  onChange: (c: number) => void;
}

export default function InputKeyCode(props: Props): JSX.Element {
  const onChange = (c: number) => {
    props.onChange(c);
  };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "row",
      }}
    >
      <input
        autoFocus
        value={props.keyCode}
        size={14}
        onKeyDown={(e) => {
          e.preventDefault();
          e.stopPropagation();
          console.log(e.code, e.key);
        }}
      />
      <button
        onClick={() => {
          onChange(0);
        }}
      >
        X
      </button>
    </div>
  );
}
