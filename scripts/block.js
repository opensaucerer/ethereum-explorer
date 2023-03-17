import http from 'k6/http';
import { check, sleep } from 'k6';

export default function () {
  const response = http.get(__ENV.API_ENDPOINT + '/block', {
    headers: { Accepts: 'application/json' },
  });
  check(response, { 'status is 200': (r) => r.status === 200 });
  sleep(0.3);
}
