import React, {useState} from "react";
import TextButton from "./TextButton";

interface Props {
  title: string
  defaultClose?: boolean
  children: React.ReactNode
}

export default function CategoryWrapper(props: Props): JSX.Element {
  const [isOpen, setIsOpen] = useState(!props.defaultClose)
  return (
      <div>
        <TextButton onClick={() => {
          setIsOpen((v) => !v)
        }}>
          {props.title + (isOpen ? " ▲" : " ▼")}
        </TextButton>

        <div style={{
          transform: `scaleY(${isOpen ? 1 : 0})`,
          height: isOpen ? "100%" : 0,
          transition: "all 0.3s",
          transformOrigin: "top left",
          border: "1px solid #ccc",
          marginBottom: 10
        }}>
          {props.children}
        </div>
      </div>
  );
}