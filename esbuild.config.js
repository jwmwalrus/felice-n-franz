/* eslint-disable import/no-extraneous-dependencies */
/* eslint-disable global-require */
const esbuild = require('esbuild');
const ncp = require('ncp');

const buildOptions = {
    entryPoints: ['./assets/js/index.js'],
    bundle: true,
    platform: 'neutral',
    outdir: 'public/static',
    sourcemap: true,
    target: 'node12',
    external: Object.keys(require('./package.json').dependencies),
};

const vendorDirs = [
    { source: 'node_modules/bootstrap/dist', dest: 'public/static/vendor/bootstrap' },
    { source: 'node_modules/jquery/dist', dest: 'public/static/vendor/jquery' },
    { source: 'node_modules/@fortawesome/fontawesome-free', dest: 'public/static/vendor/fontawesome-free' },
    { source: 'node_modules/simple-line-icons', dest: 'public/static/vendor/simple-line-icons' },
];

const apply = async () => {
    try {
        const result = await esbuild.build(buildOptions);

        // copy vendorDirs
        for await (const d of vendorDirs) {
            ncp(d.source, d.dest, (err) => {
                if (err) {
                    console.error(err);
                }
            });
        }
    } catch (e) {
        console.error(e);
    }
};

apply();
