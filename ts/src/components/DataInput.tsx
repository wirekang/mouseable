import React, { useState } from "react";
import { changeData, DataDefinition, DataType } from "../gobind";

interface Props {
  def: DataDefinition;
  value: string;
}

export default function DataInput(props: Props): JSX.Element {
  const [buffer, setBuffer] = useState(props.value);
  const [isError, setIsError] = useState(false);

  const onApply = () => {
    if (props.def.Type === DataType.Bool) {
      changeData(props.def.Name, buffer);
      return;
    }

    let n = 0;
    if (props.def.Type === DataType.Int) {
      n = parseInt(buffer, 10);
    } else {
      n = parseFloat(buffer);
    }

    if (isNaN(n)) {
      setIsError(true);
    } else {
      setIsError(false);
      setBuffer(`${n}`);
      changeData(props.def.Name, `${n}`);
    }
  };

  const text = (
    <input
      style={{
        fontSize: 11,
      }}
      type="text"
      size={3}
      value={buffer}
      onChange={(e) => {
        setBuffer(e.target.value);
      }}
    />
  );
  const checkbox = (
    <input
      type="checkbox"
      checked={buffer === "true"}
      readOnly
      onClick={() => {
        setBuffer(buffer === "true" ? "false" : "true");
      }}
    />
  );

  const input = props.def.Type === DataType.Bool ? checkbox : text;

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "row",
        backgroundColor: isError ? "#f88" : "#fff",
        justifyContent: "space-between",
        alignItems: "center",
        padding: 1,
        margin: "0 10px",
      }}
    >
      {input}
      <input
        type="button"
        style={{ cursor: "pointer", fontSize: 11, border: "none" }}
        value="Apply"
        onClick={onApply}
      />
    </div>
  );
}
