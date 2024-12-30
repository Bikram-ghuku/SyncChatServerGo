import { Server } from "socket.io";
import { Producer, CompressionTypes } from "kafkajs";
import { Server as HTTPServer } from "node:http";

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

    runListeners(kafkaProducer : Producer){
        this.io.on("connection", (socket) => {
            console.log(socket.id);
            socket.on("message", (data) => {
                kafkaProducer.send({
                    compression: CompressionTypes.GZIP,
                    topic: process.env.KAFKA_TOPIC || "sync-chat-msg",
                    messages: [{value: JSON.stringify(data)}]
                }).then(() => {
                }).catch((err) => console.log("Error: ",err));
            });

            socket.on("disconnect", () => {
                console.log(`Client disconnected: ${socket.id}`);
            });
        });
    }
}
