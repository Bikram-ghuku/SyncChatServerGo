import { Socket, Server } from "socket.io";
import initProducer from "./KafkaClient.ts";
import { Producer } from "kafkajs";
import express from "express";
import { createServer } from "node:http";
import cors from "cors";


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

const PORT = 4000;
httpServer.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
    console.log("Socket.IO server is ready.");
});