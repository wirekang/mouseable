import React from 'react';
import {loadBind} from "./gobind";
import CategoryWrapper from "./components/CategoryWrapper";

function App() {
  loadBind().then((v) => {
    console.log(v)
  })
  return (<div>
    <CategoryWrapper title="Test">
      <h1>asdfasdf</h1>
      <h4>asdfasdf</h4>
      <h3>asdfasdf</h3>
      <h2>asdfasdf</h2>
    </CategoryWrapper>
    <CategoryWrapper title="Test-2">
      <h1>asdfasdf</h1>
      <h4>asdfasdf</h4>
      <h3>asdfasdf</h3>
      <h2>asdfasdf</h2>
    </CategoryWrapper>
  </div>);
}

export default App;
