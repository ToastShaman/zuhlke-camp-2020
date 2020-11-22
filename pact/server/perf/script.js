import {check, sleep} from 'k6';
import http from 'k6/http';

// noinspection JSUnusedGlobalSymbols
export const options = {
    stages: [
        {duration: '5m', target: 100}, // simulate ramp-up of traffic from 1 to 100 users over 5 minutes.
        {duration: '10m', target: 100}, // stay at 100 users for 10 minutes
        {duration: '5m', target: 0}, // ramp-down to 0 users
    ],
    thresholds: {
        http_req_duration: ['p(99)<1500'], // 99% of requests must complete below 1.5s
    },
};

const BASE_URL = `http://${__ENV.MY_HOSTNAME}`;

function createTodo() {
    return http.post(`${BASE_URL}/todo`, JSON.stringify({
        "category": "shopping",
        "text": "Don't forget the milk",
        "status": "TODO",
        "tags": [
            "groceries",
            "food"
        ]
    }), {
        tags: {name: 'CreateTodo'},
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        }
    });
}

function getTodo(id) {
    return http.get(`${BASE_URL}/todo/${id}`, {
        headers: {'Accept': 'application/json'},
        tags: {name: 'GetTodo'}
    });
}

function deleteTodo(id) {
    return http.del(`${BASE_URL}/todo/${id}`, null, {
        headers: {'Accept': 'application/json'},
        tags: {name: 'DeleteTodo'}
    });
}

// noinspection JSUnusedGlobalSymbols
export default () => {
    let r = createTodo()
    check(r, {'status is 200': (r) => r.status === 200});
    let id = r.json().id

    r = getTodo(id)
    check(r, {'status is 200': (r) => r.status === 200});

    r = deleteTodo(id)
    check(r, {'status is 200': (r) => r.status === 200});

    sleep(1);
}
