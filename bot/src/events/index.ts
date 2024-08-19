import { Markup } from 'telegraf'
import { redisClient } from '../services/redis'
import { bot } from '../services/telegram'
import { NEW_MESSAGE } from '../messages'

export const initiateRedisEventsListener = async () => {
  const WEB_APP_URL = process.env.WEB_APP_URL!

  const subscriber = redisClient
  await subscriber.connect()

  await subscriber.subscribe('message', (id) => {
    console.log(id);
    bot.telegram.sendMessage(id, NEW_MESSAGE, {
      reply_markup: Markup.inlineKeyboard([
        Markup.button.webApp('مشاهده پیام📬', `${WEB_APP_URL}/`),
      ]).reply_markup,
    })
  })

  process.on('exit', async () => {
    await subscriber.disconnect()
  })

  process.on('SIGINT', async () => {
    await subscriber.disconnect()
  })
}
