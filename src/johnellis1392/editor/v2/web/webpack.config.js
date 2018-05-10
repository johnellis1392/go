const path = require('path')
const webpack = require('webpack')

module.exports = {
  mode: 'development',
  context: path.resolve(__dirname, './src'),
  entry: {
    app: './app.js',
  },
  output: {
    filename: '[name].bundle.js',
    path: path.resolve(__dirname, './dist'),
    publicPath: '/assets/',
  },
  devServer: {
    host: '0.0.0.0',
    port: 8080,
    contentBase: path.resolve(__dirname, './src'),
  },

  module: {
    rules: [{
        test: /\.js$/,
        exclude: [/node_modules/],
        use: [{
          loader: 'babel-loader',
          options: {
            presets: ['es2015', 'stage-0', 'react'],
          },
        }]
      },
      {
        test: /\.scss$/,
        exclude: /node_modules/,
        use: [
          'style-loader',
          'css-loader?importLoaders=1&modules&localIdentName=[local]__[name]__[path]__[hash:base64:5]',
          'sass-loader'
        ]
      },
      { test: /\.json$/, loader: 'json-loader' },

      // { test: /\.css$/, loader: 'style-loader!css-loader' },
      // { test: /\.(woff2?|ttf|eot|svg)$/, loader: 'url?limit=10000' },
      // { test: /bootstrap\/dist\/js\/umd\//, loader: 'imports?jQuery=jquery' }, // Load Bootstrap Content
      // { test: /\.png$/, loader: 'file-loader?name=/assets/[name].[ext]' },
    ],
  },

  resolve: {
    modules: [
      path.resolve('./src'),
      path.resolve('./node_modules'),
    ],
  },

  plugins: [
    new webpack.ProvidePlugin({
      jQuery: 'jquery',
      $: 'jquery',
      jquery: 'jquery',
    }),
  ],
}
