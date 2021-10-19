import http from 'k6/http';
import {check, sleep} from 'k6';

export let options = {
    stages: [
        {duration: '30s', target: 200},
        {duration: '1m', target: 700},
        {duration: '1m', target: 200},
        {duration: '30s', target: 0},
    ],
    thresholds: {
        http_req_duration: ['p(99)<1500'],
    },
};


const BASE_URL = `${__ENV.HOSTNAME}/api/v1/book/`

export default function () {
    let url = BASE_URL
    for (let i = 0; i < 10; i++) {
        let resp = http.get(url)
        if (resp.status === 200) {
            url = resp.json().next
        }
        const checkRes = check(resp, {
            'status is 200': (r) => r.status === 200
        });
    }
    sleep(1)
}
