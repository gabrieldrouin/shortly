import http from "k6/http";
import { check } from "k6";

const BASE_URL = __ENV.BASE_URL || "http://nginx";

// Seed short URLs during setup so the redirect scenario has targets to hit.
export function setup() {
  const codes = [];
  for (let i = 0; i < 100; i++) {
    const res = http.post(
      `${BASE_URL}/api/shorten`,
      JSON.stringify({ url: `https://example.com/load-test/page/${i}` }),
      { headers: { "Content-Type": "application/json" } }
    );
    if (res.status === 201) {
      const body = JSON.parse(res.body);
      codes.push(body.short_url.split("/").pop());
    }
  }
  console.log(`Seeded ${codes.length} short URLs`);
  return { codes };
}

export const options = {
  scenarios: {
    redirects: {
      executor: "ramping-arrival-rate",
      startRate: 100,
      timeUnit: "1s",
      preAllocatedVUs: 200,
      maxVUs: 2000,
      stages: [
        { target: 1160, duration: "30s" }, // ramp to average
        { target: 1160, duration: "1m" },  // sustain average
        { target: 3500, duration: "30s" }, // ramp to peak
        { target: 3500, duration: "30s" }, // sustain peak
        { target: 0, duration: "30s" },    // cool down
      ],
      exec: "redirect",
    },
  },
  thresholds: {
    http_req_duration: ["p(95)<200"],
    http_req_failed: ["rate<0.01"],
  },
};

export function redirect(data) {
  const code = data.codes[Math.floor(Math.random() * data.codes.length)];
  const res = http.get(`${BASE_URL}/${code}`, { redirects: 0 });
  check(res, {
    "status is 302": (r) => r.status === 302,
  });
}
