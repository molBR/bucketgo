import http from 'k6/http';
import { sleep } from 'k6';


export const options = {
    vus: 10,
    duration: '10m',
  };
export default function () {
  const response = http.get('http://localhost:8080/');
  console.log("Status Code: ", response.status)
  sleep(1);
}