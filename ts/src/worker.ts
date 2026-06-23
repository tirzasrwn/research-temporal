import { NativeConnection, Worker } from '@temporalio/worker';
import * as activities from './activities';

async function run() {
  const connection = await NativeConnection.connect();

  const worker = await Worker.create({
    connection,
    namespace: 'default',
    taskQueue: 'fulfillment-tasks',
    workflowsPath: require.resolve('./workflows'),
    activities,
  });

  await worker.run();
}

run().catch((err) => {
  console.error('worker failed', err);
  process.exit(1);
});
