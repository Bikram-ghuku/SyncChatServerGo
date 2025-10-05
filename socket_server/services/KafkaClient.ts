import { Consumer, Kafka, Producer } from 'kafkajs'
import DB from './PostgresDB'

const kafka = new Kafka({
	clientId: 'socketInjest',
	brokers: [process.env.KAFKA_BROKERS || 'localhost:9092'],
})

export const kafkaTopic = process.env.KAFKA_TOPIC || 'sync-chat-msg'

export async function initKafkaProducer() {
	const producer: Producer = kafka.producer({
		allowAutoTopicCreation: true,
	})

	try {
		await producer.connect()
		return producer
	} catch (error) {
		console.error('Failed to connect Kafka producer:', error)
		throw error
	}
}

export async function initKafkaConsumer() {
	const consumer: Consumer = kafka.consumer({
		groupId: 'socketInjest-group',
		allowAutoTopicCreation: true,
	})
	try {
		await consumer.connect()
		await consumer.subscribe({ topic: kafkaTopic, fromBeginning: false })
		consumer.on('consumer.crash', e => {
			console.error('Kafka consumer crashed:', e)
			process.exit(1)
		})
		await consumer.run({
			eachMessage: async ({ topic, partition, message }) => {
				const { jwt, msg, chatId, timeStamp } = JSON.parse(
					message.value?.toString() || '{}'
				)
				const createdAtISO = new Date(timeStamp).toISOString()
				DB.InsertIntoDB(msg, chatId, jwt, createdAtISO).catch(err => {
					console.error('Failed to insert message into DB:', err)
				})
				console.log({
					value: message.value?.toString(),
					topic,
					partition,
				})
			},
		})
	} catch (error) {
		console.error('Failed to connect Kafka consumer:', error)
		throw error
	}
}
