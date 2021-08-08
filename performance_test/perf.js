import http from 'k6/http';
import { check, group, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '1s', target: 30 }, // stay at 100 users for 10 minutes
    { duration: '20s', target: 30 }, // stay at 100 users for 10 minutes
    { duration: '1s', target: 0 }, // stay at 100 users for 10 minutes
  ]
};

//const BASE_URL = 'http://192.168.0.201:30390';
const BASE_URL = 'http://192.168.0.161:8081';
//const BASE_URL = 'http://192.168.0.226';
//const BASE_URL = 'http://shorturl.bw';

function uuidv4() {
  return 'http://xxxxx4xxxyx.xxxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

function getRandomInt(max) {
  return Math.floor(Math.random() * max);
}

export let history_list = []

export default () => {
	let url = uuidv4()
	let createShortURL = http.post(`${BASE_URL}/add`, JSON.stringify({
		"url": url
	}), { 
		headers: { 
			'Content-Type': 'application/json' 
		} 
	});

	let shortName = createShortURL.json('short')
	history_list.push(shortName)

	// get just created url
	let getLongURL = http.get(`${BASE_URL}/${shortName}`)
	let nowStatus = getLongURL.json('status')
	let long = getLongURL.json('url')
	
	// get history random url
	let getHistoryLongURL = http.get(`${BASE_URL}/${history_list[getRandomInt(history_list.length)]}`)
	let nowHistoryStatus = getHistoryLongURL.json('status')
	let longHistory = getHistoryLongURL.json('url')
	
	sleep(0.5);
};
