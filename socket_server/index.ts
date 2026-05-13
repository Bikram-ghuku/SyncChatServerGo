import { initKafkaConsumer, initKafkaProducer } from './services/KafkaClient.ts'
import { Producer } from 'kafkajs'
import express from 'express'
import { createServer } from 'node:http'
import { SocketServer } from './services/SocketServer.ts'
import {
	ensureRedisStreamGroup,
	initRedisConsumer,
	initRedisProducer,
	redisConsumerName,
	redisGroup,
	redisStream,
} from './services/RedisIPC.ts'
import { RedisClientType } from 'redis'

let kafkaProducer: Producer
let redisProducerClient: RedisClientType
let redisConsumerClient: RedisClientType
const app = express()
const httpServer = createServer(app)

const socketServer = new SocketServer()
socketServer.init(httpServer)
const wss = socketServer.getWSS()

initRedisProducer().then((redCon: RedisClientType) => {
	redisProducerClient = redCon
	console.log('Redis Client Connected')
	initKafkaProducer().then(kafkaProd => {
		kafkaProducer = kafkaProd
		console.log('Kafka Client connected!')
		socketServer.runListeners(kafkaProd, redCon)
	})
})

initRedisConsumer().then(async (redisConsumer: RedisClientType) => {
	redisConsumerClient = redisConsumer
	await ensureRedisStreamGroup(redisConsumer)
	console.log('Listening for messages from Redis Stream')
	const readLoop = async () => {
		while (true) {
			const response = await redisConsumer.xReadGroup(
				redisGroup,
				redisConsumerName,
				[{ key: redisStream, id: '>' }],
				{ COUNT: 100, BLOCK: 5000 }
			)
			if (!response) continue
			for (const stream of response) {
				for (const entry of stream.messages) {
					const payload = entry.message?.payload
					if (payload) {
						socketServer.emitMessage(payload)
						await redisConsumer.xAck(redisStream, redisGroup, entry.id)
					}
				}
			}
		}
	}
	readLoop().catch(err => {
		console.error('Redis stream read error:', err)
		process.exit(1)
	})
})

initKafkaConsumer()
	.then(() => {
		console.log('Kafka Consumer is running, but not doing anything yet.')
	})
	.catch(err => {
		console.error('Error initializing Kafka Consumer:', err)
		process.exit(1)
	})

const PORT = process.env.SOCKET_PORT || 4000
httpServer.listen(PORT, () => {
	console.log(`Server is running on http://localhost:${PORT}`)
	console.log('WebSocket server is ready.')
})

process.on('SIGINT', async () => {
	console.log('Shutting down gracefully...')
	if (kafkaProducer) {
		await kafkaProducer.disconnect()
	}

	if (redisConsumerClient) {
		await redisConsumerClient.disconnect()
	}

	if (redisProducerClient) {
		await redisProducerClient.disconnect()
	}

	wss.close(() => {
		process.exit(0)
	})
})
