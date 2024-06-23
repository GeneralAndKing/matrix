const BASE_URL = 'ws://localhost:8080'

const messageSocket = new WebSocket(`${BASE_URL}/message`)
const businessSocket = new WebSocket(`${BASE_URL}/business`)

export { messageSocket, businessSocket }
