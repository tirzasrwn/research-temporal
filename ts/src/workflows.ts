import { proxyActivities, ApplicationFailure } from '@temporalio/workflow';
import type * as activities from './activities';

const { sendGreeting } = proxyActivities<typeof activities>({
  startToCloseTimeout: '10s',
  retry: {
    initialInterval: '1s',
    backoffCoefficient: 2,
    maximumAttempts: 5,
  },
});

export async function fulfillGreeting(receipt: string): Promise<string> {
  if (!receipt || !receipt.startsWith('Hello ') || !receipt.includes(' | ')) {
    throw ApplicationFailure.nonRetryable(
      `invalid greeting receipt: ${JSON.stringify(receipt)} (expected "Hello <name>! <notif> | <receipt>")`,
      'InvalidGreetingFormat',
    );
  }

  const deliveryId = await sendGreeting(receipt);
  return `${receipt} -> ${deliveryId}`;
}
