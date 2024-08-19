import 'dotenv/config'
import { Telegraf } from 'telegraf'
import { handleCommands } from './commands/commands'
import { redisClient } from './services/redis'
import { initiateRedisEventsListener } from './events'
import { bot } from './services/telegram'

async function main() {
  try {
    const WEB_APP_URL = process.env.WEB_APP_URL

    if (!WEB_APP_URL)
      return console.log(
        new Error('WEB_APP_URL Not found!, initiating the bot'),
      )

    await redisClient.connect()

    handleCommands()
    initiateRedisEventsListener()

    bot.launch().then(() => console.log('bot launched'))
    process.once('SIGINT', () => bot.stop('SIGINT'))
    process.once('SIGTERM', () => bot.stop('SIGTERM'))
    return bot
  } catch (error) {
    console.log(error)
  }
}
main()
