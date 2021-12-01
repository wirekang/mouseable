import React from "react";
import GroupBox from "./GroupBox";
import { openGitHub } from "../gobind";

interface Props {
  version: string;
}

export default function InfoGroupBox(props: Props): JSX.Element {
  return (
    <GroupBox title="Mouseable">
      <div
        style={{
          display: "flex",
          flexDirection: "row",
          alignItems: "center",
          justifyContent: "space-between",
          margin: 3,
        }}
      >
        <span>{props.version}</span>

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
          onClick={openGitHub}
        >
          <span>GitHub</span>
          <img alt="github" src="github.png" width={20} height={20} />
        </a>
      </div>
    </GroupBox>
  );
}
