# Banbot Strategies
This repo contains free strategies for [banbot](https://github.com/banbox/banbot).

## Disclaimer
These strategies are for educational purposes only. Do not risk money which you are afraid to lose. USE THE SOFTWARE AT YOUR OWN RISK. THE AUTHORS AND ALL AFFILIATES ASSUME NO RESPONSIBILITY FOR YOUR TRADING RESULTS.

Always start by testing strategies with a backtesting then run the trading bot in Dry-run. Do not engage money before you understand how it works and what profit/loss you should expect.

We strongly recommend you to read the [document](https://www.banbot.site) and understand the mechanism of this bot.

Some only work in specific market conditions, while others are more "general purpose" strategies. It's noteworthy that depending on the exchange and Pairs used, further optimization can bring better results.

Please keep in mind, results will heavily depend on the symbols, timeframe and timerange used to backtest - so please run your own backtests that mirror your usecase, to evaluate each strategy for yourself.

## Share your own strategies
We welcome contributions of classic algorithmic trading strategies so that more people can quickly get started with banbot.

## FAQ
### What is banbot?
banbot is a free and open source crypto trading bot written in golang. It aims to provide a simple, easy-to-use, high-performance quantitative backtesting experience. 
It contains web ui, backtesting, plotting and money management tools as well as strategy optimization.

### How to test a strategy?
You can use this repository directly as your project, or copy the strategy go file to your strategy project.

Then you can execute `go build -o bot` to compile banbot with your strategy into a single executable file, and run this file directly to access the web ui from `http://localhost:8000`!

### How to install a strategy?
Just run the single executable file, and access `http://localhost:8000/trade`

## Advanced Usage
You can visit [here](/demo/README.md) to view more advanced examples of banbot usage
