/* eslint-disable import/no-extraneous-dependencies */
/* eslint-disable global-require */
const esbuild = require('esbuild');
const glob = require('glob');
const path = require('path');
const mv = require('mv');
const ncp = require('ncp');

const buildOptions = {
    entryPoints: ['./assets/js/index.js'],
    bundle: true,
    platform: 'neutral',
    // outfile: 'public/js/index.js',
    outdir: 'public/js',
    // loader: ['css'],
    sourcemap: true,
    target: 'node12',
    assetNames: '[dir]/[name]',
    external: Object.keys(require('./package.json').dependencies),
};

const vendorDirs = [
    { source: 'node_modules/bootstrap/dist', dest: 'public/vendor/bootstrap' },
    { source: 'node_modules/jquery/dist', dest: 'public/vendor/jquery' },
    { source: 'node_modules/@fortawesome/fontawesome-free', dest: 'public/vendor/fontawesome-free' },
    { source: 'node_modules/simple-line-icons', dest: 'public/vendor/simple-line-icons' },
];

const apply = async () => {
    try {
        const result = await esbuild.build(buildOptions);

        // move js/*.css to css/
        await glob(
            'public/js/*.css',
            {},
            (err, files) => {
                if (err) {
                    console.error(err);
                    return;
                }
                files.forEach(
                    (f) => mv(
                        f,
                        'public/css/' + path.basename(f),
                        console.error,
                    ),
                );
            },
        );

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
