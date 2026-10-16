import { BrowserRouter, Routes, Route, useLocation } from "react-router";
import Home from "./pages/Home";
import Blackjack from "./pages/Blackjack";

import "./App.css";

function RouteDebug() {
  const location = useLocation();
  console.log("Current route:", location.pathname);
  return null;
}

function App() {
  return (
    <BrowserRouter>
      <RouteDebug />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/blackjack" element={<Blackjack />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
