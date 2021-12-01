import React from "react";
import { DataDefinition, DataType } from "../gobind";
import Row from "./Row";
import DataInput from "./DataInput";

interface Props {
  def: DataDefinition;
  value: string;
}

export default function DataRow(props: Props): JSX.Element {
  return (
    <Row
      name={props.def.Name}
      column2={<DataInput def={props.def} value={props.value} />}
      column3={<span>{DataType[props.def.Type]}</span>}
      description={props.def.Description}
    />
  );
}
