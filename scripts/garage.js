import http from 'k6/http';
import { check, sleep } from 'k6';
import blockTest from './block.js';
import txTest from './tx.js';
import statsTest from './stats.js';

const default_vus = 5000;

const target_vus_env = `${__ENV.TARGET_VUS}`;
const target_vus = parseInt(target_vus_env) || default_vus;

export let options = {
  stages: [
    // Ramp-up from 1 to TARGET_VUS virtual users (VUs) in 5s
    { duration: '5s', target: target_vus },

    // Stay at rest on TARGET_VUS VUs for 10s
    { duration: '10s', target: target_vus },

    // Ramp-down from TARGET_VUS to 0 VUs for 5s
    { duration: '5s', target: 0 },
  ],
};

export default function () {
  blockTest();
  txTest();
  statsTest();
}
