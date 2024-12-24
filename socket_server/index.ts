import { Socket } from "socket.io";
import initProducer from "./KafkaClient.ts";
import { Producer } from "kafkajs";

let kafkaProducer: Producer;

initProducer().then((kafkaProd) => {
    kafkaProducer = kafkaProd;
    console.log("Kafka Client conencted!");
})