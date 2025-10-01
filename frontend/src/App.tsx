import { useState, useEffect } from "react";
import Auth from "./Auth";
import Dashboard from "./Dashboard";
import api from "./api";

export default function App() {
  const [logged, setLogged] = useState(false);

  async function checkAuth() {
    try {
      await api.get("/auth/me");
      setLogged(true);
    } catch {
      setLogged(false);
    }
  }

  useEffect(() => {
    checkAuth();
  }, []);

  return (
    <>
      {logged ? <Dashboard onLogout={() => setLogged(false)} /> : <Auth onLogin={checkAuth} />}
    </>
  );
}
