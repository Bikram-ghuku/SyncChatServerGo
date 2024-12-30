import { Kafka, Producer } from "kafkajs";

const kafka = new Kafka({
    clientId: 'socketInjest',
    brokers: [process.env.KAFKA_BROKERS || "localhost:29092"],
});

async function initProducer() {
    const producer: Producer = kafka.producer({
        allowAutoTopicCreation: true,
    });

    try {
        await producer.connect();
        return producer;
    } catch (error) {
        console.error("Failed to connect Kafka producer:", error);
        throw error;
    }
}


export default initProducer;