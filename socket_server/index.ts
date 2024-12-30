import initProducer from "./KafkaClient.ts";
import { Producer, CompressionTypes } from "kafkajs";
import express from "express";
import { createServer } from "node:http";
import { SocketServer } from "./SocketServer.ts";

let kafkaProducer: Producer;
const app = express();
const httpServer = createServer(app);


const socketServer = new SocketServer();
socketServer.init(httpServer);
const io = socketServer.getIO();

initProducer().then((kafkaProd) => {
    kafkaProducer = kafkaProd;
    console.log("Kafka Client conencted!");
    socketServer.runListeners(kafkaProducer);
});


const PORT = process.env.PORT || 4000;
httpServer.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
    console.log("Socket.IO server is ready.");
});

process.on("SIGINT", async () => {
    console.log("Shutting down gracefully...");
    if (kafkaProducer) {
        await kafkaProducer.disconnect();
        console.log("Kafka producer disconnected.");
    }
    io.close(() => {
        console.log("Socket.IO server closed.");
        process.exit(0);
    });
});