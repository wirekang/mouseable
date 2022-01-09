import React from "react";
import { openLink } from "../gobind";

interface Props {}

export default function GitHubButton(props: Props): JSX.Element {
  const onClick = () => {
    openLink("https://github.com/wirekang/mouseable");
  };

  return (
    <a
      style={{
        background: "none",
        border: "none",
        margin: 0,
        padding: 3,
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "center",
        cursor: "pointer",
        fontSize: 13,
      }}
      onClick={onClick}
    >
      <span>GitHub </span>
      <img alt="github" src="github.png" width={20} height={20} />
    </a>
  );
}
