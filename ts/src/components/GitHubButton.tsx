import React from "react";

interface Props {}

export default function GitHubButton(props: Props): JSX.Element {
  const onClick = () => {
    window.__openLink__("https://github.com/wirekang/mouseable");
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
      <img alt="github" src="github.png" width={50} height={20} />
    </a>
  );
}
