# jade-news
News Aggregator inspired by github.com/jamesroutley/news.routley.io

## Actions
On commit of new code, run test suite and generate release if it passes
On cron run release to update feeds

## Future
state tracking file of links, locically split fetch and render
age links out of tracking file, so we gracefully handle 
articles that get fall off the end of source feed before our expiration

if feed fetch fails, log error, but keep going
if feed fetch returns zero links, log warning but keep going