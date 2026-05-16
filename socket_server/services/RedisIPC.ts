import { createClient, RedisClientType } from 'redis'

export const redisStream = process.env.REDIS_STREAM || 'sync-chat-stream'
const instanceId = process.env.SOCKET_INSTANCE_ID || process.pid
export const redisGroup =
	process.env.REDIS_GROUP || `sync-chat-group-${instanceId}`
export const redisConsumerName =
	process.env.REDIS_CONSUMER || `consumer-${instanceId}`
const url = process.env.REDIS_URL || 'redis://localhost:6379'

export async function initRedisProducer() {
	const producerClient: RedisClientType = createClient({
		url,
	})
	try {
		await producerClient.connect()
		return producerClient
	} catch (err) {
		console.log(err)
		throw Error('Error init Redis Producer')
	}
}

export async function initRedisConsumer() {
	const consumerClient: RedisClientType = createClient({
		url,
	})
	try {
		await consumerClient.connect()
		return consumerClient
	} catch (err) {
		console.log(err)
		throw Error('Error init Redis Consumer')
	}
}

export async function ensureRedisStreamGroup(client: RedisClientType) {
	try {
		await client.xGroupCreate(redisStream, redisGroup, '0', {
			MKSTREAM: true,
		})
	} catch (err: any) {
		if (!String(err?.message || '').includes('BUSYGROUP')) {
			throw err
		}
	}
}
