import {SharedArray} from "k6/data";
import exec from "k6/execution";
import {check, sleep} from 'k6';
import http from "k6/http";


const data = new SharedArray("post-book", function () {
    return JSON.parse(open('tokens.json'));
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
    var token = data[exec.scenario.iterationInTest];
    var url = `${__ENV.HOSTNAME}/api/v1/book/`

    for (let i = 0; i < 50; i++) {
        let title = "book" + i
        var payload = JSON.stringify({
            title: title
        });

        var params = {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token
            },
        };

        const resp = http.post(url, payload, params);
        sleep(0.25)
        const checkRes = check(resp, {
            'status is 201': (r) => r.status === 201
        });
    }
    sleep(3);


}
