import { Socket, Server } from "socket.io";
import initProducer from "./KafkaClient.ts";
import { Producer } from "kafkajs";
import express from "express";
import { createServer } from "node:http";

let kafkaProducer: Producer;
const app = express();
const httpServer = createServer(app);

initProducer().then((kafkaProd) => {
    kafkaProducer = kafkaProd;
    console.log("Kafka Client conencted!");
});

const io = new Server(httpServer, {
    cors: {
        origin: "*",
        methods: ["GET", "POST"],
    },
});

io.on("connection", (socket) => {
    console.log(socket.id);
    socket.on("message", (data) => {
        kafkaProducer.send({
            topic: process.env.KAFKA_TOPIC || "sync-chat-msg",
            messages: [{value: JSON.stringify(data)}]
        }).then(() => {
            console.log("Sent to kafka successfully");
        }).catch((err) => console.log("Error: ",err));
    })
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