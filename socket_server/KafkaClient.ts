import { Kafka, Producer } from "kafkajs";

const kafka = new Kafka({
    clientId: 'socketInjest',
    brokers: ['localhost:9092'],
});

async function initProducer(){
    const producer:Producer = kafka.producer();
    await producer.connect();

    return producer;
}

export default initProducer;