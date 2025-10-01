/* eslint-disable @typescript-eslint/no-explicit-any */
import { useState } from "react";
import api from "./api";

interface AuthProps {
  onLogin: () => void;
}

export default function Auth({ onLogin }: AuthProps) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  async function login(e: React.FormEvent) {
    e.preventDefault();
    try {
      await api.post("/auth/login", { email, password });
      onLogin();
    } catch (err: any){
      console.error(err.response?.data || err.message);
      alert("Login failed" + (err.response?.data?.error || err.message));
    }
  }

  async function register(e: React.FormEvent) {
    e.preventDefault();
    try {
      await api.post("/auth/register", { email, password });
      alert("Registered! Now log in.");
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (err: any){
      console.error(err.response?.data || err.message);
      alert("Register failed" + (err.response?.data?.error || err.message));
    }
  }

  return (
    <div className="container">
      <h2>Login / Register</h2>
      <form onSubmit={login}>
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={e => setEmail(e.target.value)}
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={e => setPassword(e.target.value)}
        />
        <button type="submit">Login</button>
        <button type="button" onClick={register}>Register</button>
      </form>
    </div>
  );
}
