import { Context } from '@temporalio/activity';

const sleep = (ms: number) => new Promise<void>((resolve) => setTimeout(resolve, ms));

export async function sendGreeting(receipt: string): Promise<string> {
  const logger = Context.current().log;
  logger.info('sendGreeting started', { receipt });

  await sleep(500 + Math.floor(Math.random() * 1000));

  if (Math.random() < 0.2) {
    logger.error('delivery service unavailable (simulated transient failure)');
    throw new Error('delivery service unavailable (simulated transient failure)');
  }

  const deliveryId = `dlv-${Math.random().toString(36).slice(2, 10)}`;
  logger.info('greeting sent', { deliveryId });
  return deliveryId;
}
