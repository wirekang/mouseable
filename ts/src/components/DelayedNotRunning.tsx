import React from "react";
import { useTimeout } from "react-use";

interface Props {
  delay: number;
}

export default function DelayedNotRunning(props: Props): JSX.Element {
  const [isReady] = useTimeout(props.delay);
  if (isReady()) {
    return (
      <h2
        style={{
          position: "absolute",
          color: "red",
          width: "100vw",
          height: "100vh",
          backgroundColor: "white",
        }}
      >
        Mouseable is not running.
      </h2>
    );
  }

  return <React.Fragment />;
}
