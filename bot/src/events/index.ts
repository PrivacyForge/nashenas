import { Markup } from 'telegraf'
import { redisClient } from '../services/redis'
import { bot } from '../services/telegram'

export const initiateRedisEventsListener = async () => {
  const WEB_APP_URL = process.env.WEB_APP_URL!

  const subscriber = redisClient.duplicate()
  await subscriber.connect()

  await subscriber.subscribe('message', (id) => {
    bot.telegram.sendMessage(id, 'یک پیام جدید داری', {
      reply_markup: Markup.inlineKeyboard([
        Markup.button.webApp('مشاهده پیام', `${WEB_APP_URL}/`),
      ]).reply_markup,
    })
  })
}
