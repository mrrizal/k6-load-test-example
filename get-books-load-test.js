import http from 'k6/http';
import {check, sleep} from 'k6';

export let options = {
    stages: [
        {duration: '30s', target: 100},
        {duration: '1m', target: 300},
        {duration: '1m', target: 100},
        {duration: '30s', target: 0},
    ],
    thresholds: {
        http_req_duration: ['p(99)<1500'],
    },
};


const BASE_URL = 'http://localhost:3000/api/v1/book/'

export default function () {
    let url = BASE_URL
    for (let i = 0; i < 10; i++) {
        let resp = http.get(url)
        url = resp.json().next
        const checkRes = check(resp, {
            'status is 200': (r) => r.status === 200
        });
        sleep(0.1)
    }
    sleep(1)
}