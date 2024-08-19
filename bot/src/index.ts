import 'dotenv/config'
import { handleCommands } from './commands/commands'
import { initiateRedisEventsListener } from './events'
import { bot } from './services/telegram'

async function main() {
  try {
    const WEB_APP_URL = process.env.WEB_APP_URL

    if (!WEB_APP_URL)
      return console.log(
        new Error('WEB_APP_URL Not found!, initiating the bot'),
      )

    handleCommands()
    initiateRedisEventsListener()

    bot.launch()
    process.once('SIGINT', async () => {
      bot.stop('SIGINT')
    })
    process.once('SIGTERM', async () => {
      bot.stop('SIGTERM')
    })
    return bot
  } catch (error) {
    console.log(error)
  }
}
main()
