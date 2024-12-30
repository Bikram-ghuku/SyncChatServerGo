import { Server } from "socket.io";
import { Producer, CompressionTypes } from "kafkajs";
import { Server as HTTPServer } from "node:http";
import { redisTopic } from "./RedisIPC.ts";
import { kafkaTopic } from "./KafkaClient.ts";
import { RedisClientType } from "redis";

export class SocketServer {
    private io : Server;
    constructor(){
        this.io = new Server({
            cors: {
                origin: "*",
                methods: ["GET", "POST"],
            },
        });
    }

    init(httpServer: HTTPServer){
        this.io.attach(httpServer);
    }

    getIO() {
        return this.io;
    }

    runListeners(kafkaProducer : Producer, redisProducer:RedisClientType){
        this.io.on("connection", (socket) => {
            console.log(socket.id);
            socket.on("message", (data) => {
                kafkaProducer.send({
                    compression: CompressionTypes.GZIP,
                    topic: kafkaTopic,
                    messages: [{value: JSON.stringify(data)}]
                }).catch((err) => console.log("Error: ",err));

                redisProducer.publish(redisTopic, JSON.stringify(data)).catch((err) => console.log(err));
            });

            socket.on("disconnect", () => {
                console.log(`Client disconnected: ${socket.id}`);
            });
        });
    }

    emitMessage(mesage:string){
        this.io.emit("message", JSON.parse(mesage));
    }
}
