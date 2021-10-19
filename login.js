import {check, sleep} from 'k6';
import {SharedArray} from "k6/data";
import exec from "k6/execution";
import http from "k6/http";

const data = new SharedArray("singed-up-users", function () {
    return JSON.parse(open('signed-up-users.json'));
})


export let options = {
    scenarios: {
        "first-wave": {
            executor: "shared-iterations",
            vus: data.length,
            iterations: data.length,
            maxDuration: "1h"
        }
    }
}

export default function () {
    var user = data[exec.scenario.iterationInTest];
    var url = `${__ENV.HOSTNAME}/api/v1/user/login/`
    var payload = JSON.stringify({
        username: user.username,
        password: user.password,
    });

    var params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const resp = http.post(url, payload, params);
    sleep(3);

    const checkRes = check(resp, {
        'status is 200': (r) => r.status === 200
    });
}
