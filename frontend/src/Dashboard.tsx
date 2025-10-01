/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect, useState } from "react";
import api from "./api";

interface DashboardProps {
  onLogout: () => void;
}

interface Category {
  id: string;
  name: string;
}

interface Expense {
  id: string;
  description: string;
  amount: number;
  category_id: string;
  category_name?: string;
}

export default function Dashboard({ onLogout }: DashboardProps) {
  const [categories, setCategories] = useState<Category[]>([]);
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [catName, setCatName] = useState("");
  const [exp, setExp] = useState<{ amount: string; category_id: string; description: string }>({
    amount: "",
    category_id: "",
    description: "",
  });

  // Load data from backend
  async function load() {
    try {
      const cats = await api.get<Category[]>("/categories");
      const exps = await api.get("/expenses");
      setCategories(cats.data);
      // Map backend expense response to frontend format
      const expsData: Expense[] = exps.data.expenses.map((e: any) => ({
        id: e.id,
        description: e.description,
        amount: parseFloat(e.amount),
        category_id: e.category_id ?? "",
        category_name: e.category_name ?? "",
      }));
      setExpenses(expsData);
    } catch (err) {
      console.error(err);
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function addCategory() {
    if (!catName.trim()) return;
    await api.post("/categories", { name: catName });
    setCatName("");
    load();
  }

  async function addExpense() {
    // Prevent adding incomplete forms
    if (!exp.description || !exp.amount || !exp.category_id) return;

    await api.post("/expenses", {
      description: exp.description,
      amount: parseFloat(exp.amount),
      category_id: exp.category_id,
    });

    setExp({ description: "", amount: "", category_id: "" });
    load();
  }

  async function logout() {
    await api.post("/auth/logout");
    onLogout();
  }

  // Compute total spent per category
  const categoryTotals: Record<string, number> = {};
  expenses.forEach(e => {
    if (!e.category_id) return;
    categoryTotals[e.category_id] = (categoryTotals[e.category_id] || 0) + e.amount;
  });

  return (
    <div className="container">
      <h2>Dashboard</h2>
      <button onClick={logout} style={{ marginBottom: "1rem" }}>Logout</button>

      {/* Categories */}
      <div>
        <h3>Categories</h3>
        <ul className="list">
          {categories.map(c => (
            <li key={c.id}>
              {c.name} - Total spent: ${categoryTotals[c.id] ?? 0}
            </li>
          ))}
        </ul>
        <div className="form-inline">
          <input
            placeholder="New category"
            value={catName}
            onChange={e => setCatName(e.target.value)}
          />
          <button onClick={addCategory} disabled={!catName.trim()}>Add</button>
        </div>
      </div>

      {/* Expenses */}
      <div>
        <h3>Expenses</h3>
        <ul className="list">
          {expenses.map(e => (
            <li key={e.id}>
              {e.description} - ${e.amount} ({e.category_name})
            </li>
          ))}
        </ul>
        <div className="form-inline">
          <input
            placeholder="Description"
            value={exp.description}
            onChange={e => setExp({ ...exp, description: e.target.value })}
          />
          <input
            type="number"
            placeholder="Amount"
            value={exp.amount}
            onChange={e => setExp({ ...exp, amount: e.target.value })}
          />
          <select
            value={exp.category_id}
            onChange={e => setExp({ ...exp, category_id: e.target.value })}
          >
            <option value="">Select Category</option>
            {categories.map(c => (
              <option key={c.id} value={c.id}>{c.name}</option>
            ))}
          </select>
          <button
            onClick={addExpense}
            disabled={!exp.description || !exp.amount || !exp.category_id}
          >
            Add
          </button>
        </div>
      </div>
    </div>
  );
}
