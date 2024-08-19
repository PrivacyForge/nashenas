import { Context, Telegraf } from 'telegraf'

export const handleCommands = (bot: Telegraf<Context>) => {
  const WEB_APP_URL = process.env.WEB_APP_URL!

  bot.telegram.setMyCommands([
    {
      command: 'start',
      description: 'Start',
    },
  ])
  bot.help(async (ctx) => {
    ctx.reply('Run the /start command to Start')
  })
  bot.start(async (ctx) => {
    ctx.setChatMenuButton({
      text: 'Launch',
      web_app: { url: WEB_APP_URL },
      type: 'web_app',
    })
    ctx.sendMessage('Welcome to nashenas, click to open.', {
      reply_markup: {
        inline_keyboard: [
          [
            {
              text: 'Start Playing ⛏️',
              web_app: {
                url: WEB_APP_URL,
              },
            },
          ],
        ],
      },
    })
  })
}
