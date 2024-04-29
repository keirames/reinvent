// import {
//     EventStoreDBClient,
//     FORWARDS,
//     NO_STREAM,
//     START,
//     jsonEvent,
// } from '@eventstore/db-client';

import { Kafka } from 'kafkajs';

// const client = new EventStoreDBClient(
//     {
//         endpoint: 'localhost:2113',
//     },
//     { insecure: true },
// );

// async function main() {
//     const streamName = 'es_supported_clients';

//     const event = jsonEvent({
//         type: 'grpc-client',
//         data: {
//             languages: ['typescript', 'javascript'],
//             runtime: 'NodeJS',
//         },
//     });

//     let appendResult;
//     try {
//         appendResult = await client.appendToStream(streamName, [event]);
//     } catch (err) {
//         console.log('appendToStream issues:', err);
//         return;
//     }
//     console.log(`appendResult`, appendResult);

//     const events = client.readStream(streamName, {
//         fromRevision: START,
//         direction: FORWARDS,
//         maxCount: 10,
//     });

//     for await (const event of events) {
//         console.log(event);
//     }
// }

// main().catch((err) => {
//     console.log(`what is this err: ${err}`);
// });

function getBrokersAddr() {
    return ['localhost:9094'];
}

async function initTopic() {
    const kafka = new Kafka({
        brokers: getBrokersAddr(),
    });

    const admin = kafka.admin();
    try {
        await admin.connect();
        await admin.createTopics({
            topics: [{ topic: 'account_paid', numPartitions: 5 }],
        });
    } catch (err) {
        throw err;
    }

    try {
        await admin.disconnect();
    } catch (err) {
        return;
    }
}

async function main() {
    await initTopic();
}

main().catch((err) => {
    console.log('global catch', err);
});
