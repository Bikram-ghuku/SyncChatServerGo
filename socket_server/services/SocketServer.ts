import { WebSocketServer, WebSocket } from 'ws'
import { Producer, CompressionTypes } from 'kafkajs'
import { Server as HTTPServer } from 'node:http'
import { redisStream } from './RedisIPC.ts'
import { kafkaTopic } from './KafkaClient.ts'
import { RedisClientType } from 'redis'

export class SocketServer {
	private wss: WebSocketServer | null
	constructor() {
		this.wss = null
	}

	init(httpServer: HTTPServer) {
		this.wss = new WebSocketServer({ server: httpServer })
	}

	getWSS() {
		if (!this.wss) {
			throw new Error('WebSocket server not initialized')
		}
		return this.wss
	}

	runListeners(kafkaProducer: Producer, redisProducer: RedisClientType) {
		const wss = this.getWSS()
		wss.on('connection', socket => {
			socket.on('message', raw => {
				let parsed: any
				try {
					parsed = JSON.parse(raw.toString())
				} catch {
					return
				}
				const payload = parsed?.type ? parsed.data : parsed

				kafkaProducer
					.send({
						compression: CompressionTypes.GZIP,
						topic: kafkaTopic,
						messages: [{ value: JSON.stringify(payload) }],
					})
					.catch(err => console.log('Error: ', err))

				redisProducer
					.xAdd(redisStream, '*', { payload: JSON.stringify(payload) })
					.catch(err => console.log(err))
			})

			socket.on('close', () => {
				console.log('Client disconnected')
			})
		})
	}

	emitMessage(mesage: string) {
		const wss = this.getWSS()
		let payload: any
		try {
			payload = JSON.parse(mesage)
		} catch {
			return
		}
		const data = JSON.stringify({ type: 'message', data: payload })
		wss.clients.forEach(client => {
			if (client.readyState === WebSocket.OPEN) {
				client.send(data)
			}
		})
	}
}
