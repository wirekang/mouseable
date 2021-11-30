import React, { useState } from "react";
import TextButton from "./TextButton";

interface Props {
  title: string;
  defaultClose?: boolean;
  children: React.ReactNode;
}

export default function GroupBox(props: Props): JSX.Element {
  const [isOpen, setIsOpen] = useState(!props.defaultClose);
  return (
    <div
      style={{
        padding: 5,
      }}
    >
      <TextButton
        onClick={() => {
          setIsOpen((v) => !v);
        }}
        style={{
          marginLeft: 5,
        }}
      >
        {props.title + (isOpen ? " ▲" : " ▼")}
      </TextButton>

      <div
        style={{
          transform: `scaleY(${isOpen ? 1 : 0})`,
          height: isOpen ? "100%" : 0,
          transition: "all 0.2s",
          transformOrigin: "top left",
          border: "1px solid #ccc",
          margin: "2px 0 10px",
        }}
      >
        {props.children}
      </div>
    </div>
  );
}
