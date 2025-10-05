import {initKafkaConsumer, initKafkaProducer} from "./services/KafkaClient.ts";
import { Producer } from "kafkajs";
import express from "express";
import { createServer } from "node:http";
import { SocketServer } from "./services/SocketServer.ts";
import { initRedisConsumer, initRedisProducer, redisTopic } from "./services/RedisIPC.ts";
import { RedisClientType } from "redis";

let kafkaProducer: Producer;
let redisProducerClient: RedisClientType;
let redisConsumerClient: RedisClientType;
const app = express();
const httpServer = createServer(app);

const socketServer = new SocketServer();
socketServer.init(httpServer);
const io = socketServer.getIO();

initRedisProducer().then((redCon : RedisClientType) => {
    redisProducerClient = redCon;
    console.log("Redis Client Connected");
    initKafkaProducer().then((kafkaProd) => {
        kafkaProducer = kafkaProd;
        console.log("Kafka Client connected!");
        socketServer.runListeners(kafkaProd, redCon);
    });
});

initRedisConsumer().then((redisConsumer : RedisClientType) => {
    redisConsumerClient = redisConsumer;
    console.log("Listening for messages from Redis Publisher");
    redisConsumer.subscribe(redisTopic, (message) => {
        socketServer.emitMessage(message);
    })
})

initKafkaConsumer().then(() => {
    console.log("Kafka Consumer is running, but not doing anything yet.");
}).catch((err) => {
    console.error("Error initializing Kafka Consumer:", err);
    process.exit(1);
});

const PORT = process.env.SOCKET_PORT || 4000;
httpServer.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
    console.log("Socket.IO server is ready.");
});

process.on("SIGINT", async () => {
    console.log("Shutting down gracefully...");
    if (kafkaProducer) {
        await kafkaProducer.disconnect();
    }

    if(redisConsumerClient){
        await redisConsumerClient.disconnect();
    }
    
    if(redisProducerClient){
        await redisProducerClient.disconnect();
    }

    io.close(() => {
        process.exit(0);
    });
});