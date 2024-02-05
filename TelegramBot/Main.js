const { Telegraf } = require('telegraf');
const axios = require('axios');

const bot = new Telegraf('TOKEN');

// Обработка команды /start
bot.start((ctx) => {
  ctx.reply('Привет! Я твой телеграм-бот для HTTP-запросов.');
});

// Обработка текстовых сообщений
bot.on('text', async (ctx) => {
  const command = ctx.message.text;

  if (command.startsWith('/get') || command.startsWith('/post') || command.startsWith('/put') || command.startsWith('/delete')) {
    const [method, url] = command.split(' ');

    try {
      let response;
      if (method === '/get') {
        response = await axios.get(url);
      } else if (method === '/post') {
        response = await axios.post(url);
      } else if (method === '/put') {
        response = await axios.put(url);
      } else if (method === '/delete') {
        response = await axios.delete(url);
      }

      ctx.reply(JSON.stringify(response.data, null, 2));
    } catch (error) {
      ctx.reply(`Произошла ошибка: ${error.message}`);
    }
  } else {
    ctx.reply('Неправильный формат команды. Используйте /get, /post, /put или /delete, а затем URL.');
  }
});

// Обработка команды /help
bot.help((ctx) => {
  ctx.reply('Используйте команды /get, /post, /put или /delete, а затем укажите URL для выполнения соответствующего HTTP-запроса.');
});

// Запуск бота
bot.launch().then(() => {
  console.log('Бот запущен');
});
