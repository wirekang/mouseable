import React from "react";
import { loadBind } from "./gobind";
import { useAsync } from "react-use";
import FunctionCategories from "./components/FunctionCategories";

function App() {
  const state = useAsync(loadBind);

  if (state.loading || !state.value) {
    return <h4>Loading...</h4>;
  }

  return (
    <div>
      <FunctionCategories def={state.value.functionDefinitions} record={state.value.functionNameKeyMap} />
    </div>
  );
}

export default App;
