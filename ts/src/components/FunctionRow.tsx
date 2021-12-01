import React from "react";
import KeyHolder from "./KeyHolder";
import { FunctionDefinition, FunctionKey, When } from "../gobind";
import Row from "./Row";

interface Props {
  def: FunctionDefinition;
  fKey: FunctionKey;
}

export default function FunctionRow(props: Props): JSX.Element {
  return (
    <Row
      name={props.def.Name}
      description={props.def.Description}
      column2={<KeyHolder name={props.def.Name} functionKey={props.fKey} />}
      column3={
        <span
          style={{
            fontSize: 8,
            transform: "scaleX(0.6)",
          }}
        >
          {When[props.def.When]}
        </span>
      }
    />
  );
}
