import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 200 },  // Ramp up to 200 users
    { duration: '1m', target: 1000 },  // Ramp up to 1000 concurrent users
    { duration: '2m', target: 1000 },  // Stay at 1000 users for 2 minutes
    { duration: '30s', target: 0 },    // Ramp down to 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],  // 95% of requests must complete within 500ms
    http_req_failed: ['rate<0.01'],    // Error rate must be under 1%
  },
};

const BASE_URL = 'http://localhost:8080/api/v1';

export default function () {
  // 1. Browse Products
  let resProducts = http.get(`${BASE_URL}/products?page=1&limit=12`);
  check(resProducts, {
    'products loaded successfully 200': (r) => r.status === 200,
  });

  // 2. Health Liveness
  let resHealth = http.get('http://localhost:8080/health/liveness');
  check(resHealth, {
    'backend liveness ok 200': (r) => r.status === 200,
  });

  sleep(1);
}
