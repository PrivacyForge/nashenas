import { Markup } from 'telegraf'
import { START_MESSAGE } from '../messages'
import { bot } from '../services/telegram'

export const handleCommands = () => {
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
    ctx.setChatMenuButton({
      text: 'باز کردن',
      web_app: { url: WEB_APP_URL },
      type: 'web_app',
    })

    const query = ctx.update.message.text.split(' ')
    if (query.length === 2) {
      const usernameWithHash = query[1].split('-')
      if (usernameWithHash.length === 2) {
        return ctx.reply(
          `الان داری به کاربر ${usernameWithHash[0]} پیام می‌فرستی`,
          Markup.inlineKeyboard([
            Markup.button.webApp('ورود به ربات', `${WEB_APP_URL}/@${query[1]}`),
          ]),
        )
      }
    }

    ctx.reply(
      START_MESSAGE,
      Markup.inlineKeyboard([
        Markup.button.webApp('ورود به ربات', WEB_APP_URL),
      ]),
    )
  })
}
