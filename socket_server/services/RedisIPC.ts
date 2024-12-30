import {createClient, RedisClientType} from 'redis';

export const redisTopic = process.env.REDIS_TOPIC || "sync-chat-ipc";
const url = process.env.REDIS_URL || 'redis://localhost:6379';

export async function initRedisProducer() {
    const producerClient:RedisClientType = createClient({
        url
    });
    try{
        await producerClient.connect();
        return producerClient;
    }catch(err){
        console.log(err);
        throw Error(err);
    }
}

export async function initRedisConsumer() {
    const consumerClient:RedisClientType = createClient({
        url
    });
    try{
        await consumerClient.connect();
        return consumerClient;
    }catch(err){
        console.log(err);
        throw Error(err);
    }
}