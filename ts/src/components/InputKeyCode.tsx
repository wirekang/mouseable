import React, { KeyboardEventHandler, useEffect, useRef, useState } from "react";
import { useAsyncFn } from "react-use";
import { getKeyCode } from "../gobind";
import { fromVKCode } from "win-vk";

interface Props {
  keyCode: number;
  onChange: (c: number) => void;
}

export default function InputKeyCode(props: Props): JSX.Element {
  const [text, setText] = useState(fromVKCode(props.keyCode) ?? `${props.keyCode}`);
  const [state, doRequest] = useAsyncFn(getKeyCode);
  const [isSent, setIsSent] = useState(false);
  const ref = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!isSent) {
      return;
    }

    if (state.loading || !state.value) {
      setText("...");
      return;
    }

    props.onChange(state.value);
    ref.current?.blur();
    const key = fromVKCode(state.value);
    if (key) {
      setText(key);
    }
  }, [state.loading]);

  const request = () => {
    setIsSent(true);
    doRequest();
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
