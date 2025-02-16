// Простой нагрузочный тест. Регистрируются 2 юзера, один делает другому 2 перевода, покупает 4 вещи и получает инфо одного пользователя

import http from 'k6/http';
import { check, sleep } from 'k6';

const password = 'password';

const apiUrl = 'http://localhost:8080';

function getRandomUsername() {
    return `user${Math.floor(Math.random() * 1000000) + 1}`;
}

export let options = {
    duration: '1m', 
    vus: 30, 
    thresholds: {
        'http_req_failed': ['rate<0.001'], 
    },
};

export default function () {
    const user1 = getRandomUsername();
    const user2 = getRandomUsername();
    let authRes1 = http.post(`${apiUrl}/api/auth`, JSON.stringify({ username: user1, password: password }), { headers: { 'Content-Type': 'application/json' } });
    check(authRes1, { 'auth success': (r) => r.status === 200 });
    let token1 = JSON.parse(authRes1.body).token;
    let authRes2 = http.post(`${apiUrl}/api/auth`, JSON.stringify({ username: user2, password: password }), { headers: { 'Content-Type': 'application/json' } });
    check(authRes2, { 'auth success': (r) => r.status === 200 });
    let token2 = JSON.parse(authRes2.body).token;
    let sendCoinData1 = { toUser: user2, amount: 10 };
    let sendRes1 = http.post(`${apiUrl}/api/sendCoin`, JSON.stringify(sendCoinData1), { headers: { 'Authorization': `Bearer ${token1}`, 'Content-Type': 'application/json' } });
    check(sendRes1, { 'send success': (r) => r.status === 200 });
    let sendCoinData2 = { toUser: user2, amount: 10 };
    let sendRes2 = http.post(`${apiUrl}/api/sendCoin`, JSON.stringify(sendCoinData2), { headers: { 'Authorization': `Bearer ${token1}`, 'Content-Type': 'application/json' } });
    check(sendRes2, { 'send success': (r) => r.status === 200 });
    let infoRes = http.get(`${apiUrl}/api/info`, { headers: { 'Authorization': `Bearer ${token2}` } });
    check(infoRes, { 'info success': (r) => r.status === 200 });
    let info = JSON.parse(infoRes.body);
    const itemsToBuy = ['t-shirt', 'cup', 'book', 'pen'];
    let totalSpent = 0;
    itemsToBuy.forEach(item => {
        let buyRes = http.get(`${apiUrl}/api/buy/${item}`, { headers: { 'Authorization': `Bearer ${token2}` } });
        check(buyRes, { 'buy success': (r) => r.status === 200 });
    });

}

