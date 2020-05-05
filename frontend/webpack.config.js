const { VueLoaderPlugin } = require('vue-loader')
module.exports = {
    mode: 'development',
    entry: {
        app: './app.js',
    },
    module: {
        rules: [
            {
                test: /\.vue$/,
                use: 'vue-loader'
            },
            {
                test: /\.(css)$/,
                use: ['vue-style-loader', 'css-loader']
            }
        ]
    },
    plugins: [
        new VueLoaderPlugin()
    ]
}