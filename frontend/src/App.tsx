import { BrowserRouter, Routes, Route, useLocation, Navigate} from "react-router";
import Home from "./pages/Home";
import Blackjack from "./pages/Blackjack";
import Account from "./pages/Account"
import {AppLayout} from "./components/AppLayout"

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
        <Route element={<AppLayout/>}>
        <Route path="/blackjack" element={<Blackjack />} />
        <Route path="/blackjack/:gameId" element={<Blackjack />} />
        <Route path="/account" element={<Account/>}/>
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
