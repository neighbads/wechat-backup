module.exports = {
    publicPath: './',

    // 配置跨域请求
    devServer: {
        port: 8082,
        host: 'localhost',
        open: false,
        https: false,
        proxy: {
            '/api': {
                target: 'http://localhost:8081/',
                ws: true,
                changeOrigin: true,
                pathRewrite: {
                    '^/api': '/api'
                }
            },
            '/media': {
                target: 'http://localhost:8081/',
                ws: true,
                changeOrigin: true,
                pathRewrite: {
                    '^/media': '/media'
                }
            }
        }

    }
}