import { Context, Markup, Telegraf } from 'telegraf'
import { START_MESSAGE } from '../messages'

export const handleCommands = (bot: Telegraf<Context>) => {
  const WEB_APP_URL = process.env.WEB_APP_URL!

  bot.telegram.setMyCommands([
    {
      command: 'start',
      description: 'شروع',
    },
  ])
  bot.help(async (ctx) => {
    ctx.reply('برای اجرا شدن از دستور /start استفاده کنید')
  })
  bot.start(async (ctx) => {
    const query = ctx.update.message.text.split(' ')
    if (query.length === 2) {
      const username = query[1]
      return ctx.reply(
        `الان داری به کاربر ${username} پیام می‌فرستی`,
        Markup.inlineKeyboard([
          Markup.button.webApp('ورود به ربات', `${WEB_APP_URL}/@${username}`),
        ]),
      )
    }
    ctx.reply(
      START_MESSAGE,
      Markup.inlineKeyboard([
        Markup.button.webApp('ورود به ربات', WEB_APP_URL),
      ]),
    )
  })
}
