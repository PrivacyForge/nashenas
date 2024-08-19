import 'dotenv/config'
import { Telegraf } from 'telegraf'
import { handleCommands } from './bot/commands'
;(async function main() {
  try {
    const BOT_TOKEN = process.env.BOT_TOKEN
    const WEB_APP_URL = process.env.WEB_APP_URL

    if (!BOT_TOKEN)
      return console.log(new Error('BOT_TOKEN Not found!, initiating the bot'))

    if (!WEB_APP_URL)
      return console.log(
        new Error('WEB_APP_URL Not found!, initiating the bot'),
      )

    const bot = new Telegraf(BOT_TOKEN)
    console.log('Bot initiated')
    handleCommands(bot)

    bot.launch().then(() => console.log('bot launched'))
    process.once('SIGINT', () => bot.stop('SIGINT'))
    process.once('SIGTERM', () => bot.stop('SIGTERM'))
    return bot
  } catch (error) {
    console.log(error)
  }
})()
