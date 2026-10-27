import { Outlet } from "react-router";
import AuthBar from "./AuthBar"

export function AppLayout() {
  return (
    <>
      <AuthBar />
      <Outlet />
    </>
  );
}