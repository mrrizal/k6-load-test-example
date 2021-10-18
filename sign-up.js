import {sleep, check} from 'k6';
import {SharedArray} from "k6/data";
import exec from "k6/execution";
import http from "k6/http";

const data = new SharedArray("users", function () {
    return JSON.parse(open('users.json'));
})


const virtualUsers = 250;

export let options = {
	scenarios: {
        "first-wave": {
            executor: "shared-iterations",
            vus: virtualUsers,
            iterations: data.slice(0, virtualUsers).length,
            maxDuration: "1h"
        }
    }
}

export default function () {
    var user = data[exec.scenario.iterationInTest];
    var url = 'http://localhost:30003/api/v1/user/sign-up/'
    var payload = JSON.stringify({
        username: user.username,
        password: user.password,
        first_name: user.first_name,
        last_name: user.last_name,
    });

    var params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const resp = http.post(url, payload, params);
    sleep(3);

    const checkRes = check(resp, {
        'status is 201': (r) => r.status === 201
    });
}
