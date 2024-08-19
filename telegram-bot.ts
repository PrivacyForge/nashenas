import { Markup, Telegraf } from 'telegraf'

const bot = new Telegraf(process.env.BOT_TOKEN!)

const WEB_APP_URL = 'https://nshnas.ir'

const START_MESSAGE = `به ربات ناشناس امن خوش اومدی ⚡️

ما توی این ربات تمام مکالمات رو با ترکیب الگوریتم‌های RSA و AES رمزنگاری می‌کنیم تا کسی جز خودت به داده‌ها دسترسی نداشته باشه.

کد ربات به صورت اوپن‌سورس روی گیت‌هاب هست، حتی می تونی خودت توی امن کردنش مشارکت کنی ^^🌱`

const SEND_MESSAGE = (username: string) =>
  `الان داری به کاربر ${username} پیام می‌فرستی`

bot.start((ctx) => {
  const query = ctx.update.message.text.split(' ')

  if (query.length === 2) {
    return ctx.reply(
      SEND_MESSAGE(query[1]),
      Markup.inlineKeyboard([
        Markup.button.webApp('ورود به ربات', `${WEB_APP_URL}/@${query[1]}`),
      ]),
    )
  }

  ctx.reply(
    START_MESSAGE,
    Markup.inlineKeyboard([Markup.button.webApp('ورود به ربات', WEB_APP_URL)]),
  )
})

bot.launch()
