
# stock_report

A tool that sends a daily report to a user's Slack account with lists of the stocks in the SNP500 that are undervalued and overvalued.

A stock is determined to be undervalued if its market price is five or more standard deviations below its 200 day average price.
Similarly, a stock is determined to be overvalued if its market price is five or more standard deviations above its 200 day average price.

The report is sent through a webhook to the user's slack channel of choice at 9am Monday through Friday.


## Installation

Simply download the repository and run the main file. The command line will prompt you for a url to send the webhook to.


## Usage

To use stock_report, you will need a unique webhook url. This can be created by following the instructions here: https://slack.com/help/articles/115005265063-Incoming-webhooks-for-Slack.

If you would like to edit the frequency or time of the report, you can edit this line:

Markup : 'code(c.AddFunc("0 9 * * 1-5", func() { RunAnalysis(webhookURL) }))'

This link can help you with crontab customization: https://crontab.guru/

## Contributions

I made use of the following Go libraries and packages in this project
* https://github.com/piquette/finance-go
* https://github.com/robfig/cron
* https://github.com/slack-go/slack

## Acknowledgements

I would like to thank my TA Hanbang Wang for all his teaching and guidance throughout this project.

## Authors

- [@maxed8](https://github.com/maxed8)

## License
MIT © 2022 Max Edelstein
