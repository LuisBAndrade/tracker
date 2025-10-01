import axios from "axios";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  withCredentials: true, // allow cookies/sessions
  headers: { "Content-Type": "application/json" },
});

export default api;

