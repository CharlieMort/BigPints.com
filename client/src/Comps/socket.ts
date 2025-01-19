const  debug = true;
export const SOCKURL = debug ? "ws://localhost:80/ws" : "wss://bigpints.com/ws" 
export const API_URL = debug ? "http://localhost:80" : "https://bigpints.com"
export let socket = new WebSocket(SOCKURL)