import _ from 'underscore'

const context = require.context(__dirname, true, /\/.+\/index\.js$/)
const data = _.chain(context.keys())
  .map((key) => {
    const fileData = context(key)
    const branchName = key.match(/\.\/(.+)\/index.js/i)[1]
    return [
      branchName,
      fileData.default
    ]
  })
  .object()
  .value()

module.exports = data
