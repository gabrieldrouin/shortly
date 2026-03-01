import { FormEvent, useState } from "react";

interface AnalyticsResponse {
  short_code: string;
  original_url: string;
  click_count: number;
}

export default function App() {
  const [input, setInput] = useState("");
  const [data, setData] = useState<AnalyticsResponse | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  function extractCode(value: string): string {
    const trimmed = value.trim();
    try {
      const url = new URL(trimmed);
      return url.pathname.replace(/^\//, "");
    } catch {
      return trimmed;
    }
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setData(null);
    setLoading(true);

    const code = extractCode(input);

    try {
      const res = await fetch(`/api/analytics/${code}`);

      if (!res.ok) {
        const body = await res.json().catch(() => null);
        throw new Error(body?.error || `Request failed (${res.status})`);
      }

      const json: AnalyticsResponse = await res.json();
      setData(json);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Something went wrong");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="card">
      <h1>Analytics</h1>
      <p className="subtitle">Look up click stats for a short link</p>

      <form onSubmit={handleSubmit}>
        <div className="input-group">
          <input
            type="text"
            placeholder="Short code or full URL"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            required
          />
          <button className="btn btn-primary" type="submit" disabled={loading}>
            {loading ? "..." : "Lookup"}
          </button>
        </div>
      </form>

      {error && <p className="error">{error}</p>}

      {data && (
        <div className="stats">
          <div className="stat-row">
            <span className="stat-label">Short Code</span>
            <code>{data.short_code}</code>
          </div>
          <div className="stat-row">
            <span className="stat-label">Original URL</span>
            <a
              className="original-url"
              href={data.original_url}
              target="_blank"
              rel="noopener noreferrer"
            >
              {data.original_url}
            </a>
          </div>
          <div className="stat-row">
            <span className="stat-label">Clicks</span>
            <span className="click-count">{data.click_count}</span>
          </div>
        </div>
      )}

      <nav>
        <a href="/">Shorten a URL</a>
      </nav>
    </div>
  );
}
