import React, { useEffect, useRef, useState } from "react";
import { useAsyncFn } from "react-use";
import { getKeyCode, getKeyText } from "../gobind";

interface Props {
  keyCode: number;
  onChange: (c: number) => void;
}

export default function InputKeyCode(props: Props): JSX.Element {
  const [keyTextState, requestKeyText] = useAsyncFn(getKeyText.bind(null, props.keyCode));
  const [text, setText] = useState("");
  const [keyCodeState, requestKeyCode] = useAsyncFn(getKeyCode);
  const [isSent, setIsSent] = useState(false);
  const ref = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!isSent) {
      return;
    }

    if (keyCodeState.loading || !keyCodeState.value) {
      setText("...");
      return;
    }

    props.onChange(keyCodeState.value);
    ref.current?.blur();
    requestKeyText();
  }, [keyCodeState.loading, props.onChange, requestKeyText, isSent, keyCodeState.value]);

  useEffect(() => {
    if (keyTextState.value !== undefined) {
      setText(keyTextState.value);
    }
  }, [keyTextState.loading, setText, keyCodeState.value]);

  const request = () => {
    setIsSent(true);
    requestKeyCode();
  };

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
      <input ref={ref} autoFocus readOnly value={text} size={14} onFocus={request} />
      <button
        onClick={() => {
          onChange(0);
          setText("-");
        }}
      >
        X
      </button>
    </div>
  );
}
