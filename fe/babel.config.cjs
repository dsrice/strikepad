module.exports = {
    presets: [
        ['@babel/preset-env', {
            targets: {node: 'current'},
            modules: 'commonjs'
        }],
        '@babel/preset-typescript',
        ['@babel/preset-react', {runtime: 'automatic'}]
    ],
    plugins: [
        function () {
            return {
                visitor: {
                    MetaProperty(path) {
                        if (path.node.meta.name === 'import' && path.node.property.name === 'meta') {
                            path.replaceWithSourceString('({env: globalThis.import?.meta?.env || {}})');
                        }
                    }
                }
            };
        }
    ]
};