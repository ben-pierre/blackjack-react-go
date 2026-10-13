import { BrowserRouter, Routes, Route } from "react-router";
import Home from "./pages/Home";
import Blackjack from "./pages/Blackjack";

import "./App.css";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/blackjack" element={<Blackjack />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
