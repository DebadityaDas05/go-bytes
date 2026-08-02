'use strict';

const fs = require('fs');
const path = require('path');

// 1. Load JS reference implementation and Go Koffi bridge implementation
const jsBytes = require('./index.js');
const goBytes = require('../index.js');

const logFilePath = path.join(__dirname, 'log.txt');
const logStream = fs.createWriteStream(logFilePath, { flags: 'w' });

function log(msg) {
    console.log(msg);
    logStream.write(msg + '\n');
}

// Command-line arguments or default 60s run
const args = process.argv.slice(2);
let targetDurationSec = 60;
for (let i = 0; i < args.length; i++) {
    if (args[i] === '--duration' && args[i + 1]) {
        targetDurationSec = parseFloat(args[i + 1]);
    }
}

log(`================================================================`);
log(`  DIFFERENTIAL FUZZ HARNESS FOR BYTES-GO`);
log(`================================================================`);
log(`Start Time: ${new Date().toISOString()}`);
log(`Target Duration: ${targetDurationSec} seconds`);
log(`JS Reference: fuzz/index.js`);
log(`Go Bridge: bytes-go/index.js (libbytes.dll)`);
log(`================================================================\n`);

// Value & Option Generators
const units = ['b', 'kb', 'mb', 'gb', 'tb', 'pb', 'B', 'KB', 'MB', 'GB', 'TB', 'PB', 'Kb', 'kB', 'mB', 'invalid', ''];
const thousandSeps = ['', ',', '.', ' ', '_', null, undefined];
const unitSeps = ['', ' ', '\t', null, undefined];
const decimalPlacesList = [0, 1, 2, 3, 4, 5, 10, undefined, null];
const fixedDecimalsList = [true, false, undefined, null];

const boundaryNumbers = [
    0, 1, -1, 100, -100,
    1023, 1024, 1025,
    1048575, 1048576, 1048577,
    1073741823, 1073741824, 1073741825,
    Math.pow(1024, 4) - 1, Math.pow(1024, 4), Math.pow(1024, 4) + 1,
    Math.pow(1024, 5) - 1, Math.pow(1024, 5), Math.pow(1024, 5) + 1,
    1.5, -1.5, 0.5, -0.5, 0.0001, 1.0001, 1005.1005,
    NaN, Infinity, -Infinity
];

function getRandomItem(arr) {
    return arr[Math.floor(Math.random() * arr.length)];
}

function getRandomNumber() {
    if (Math.random() < 0.3) {
        return getRandomItem(boundaryNumbers);
    }
    const mag = Math.pow(10, Math.random() * 15);
    const sign = Math.random() < 0.5 ? 1 : -1;
    return sign * mag;
}

function getRandomString() {
    const r = Math.random();
    if (r < 0.2) {
        const num = getRandomNumber();
        const unit = getRandomItem(units);
        const sep = Math.random() < 0.5 ? '' : ' ';
        return `${num}${sep}${unit}`;
    } else if (r < 0.4) {
        const boundary = getRandomItem(boundaryNumbers);
        const unit = getRandomItem(units);
        return `  ${boundary} ${unit}  `;
    } else if (r < 0.6) {
        return getRandomItem(['', '   ', 'abc', '1.2.3kb', '0x11', 'foobar', 'null', 'undefined', '10.5.5']);
    } else if (r < 0.8) {
        return String(getRandomNumber());
    } else {
        return Math.random().toString(36).substring(2);
    }
}

function getRandomOptions() {
    if (Math.random() < 0.2) return undefined;
    if (Math.random() < 0.1) return null;
    const opts = {};
    if (Math.random() < 0.7) opts.decimalPlaces = getRandomItem(decimalPlacesList);
    if (Math.random() < 0.7) opts.fixedDecimals = getRandomItem(fixedDecimalsList);
    if (Math.random() < 0.7) opts.thousandsSeparator = getRandomItem(thousandSeps);
    if (Math.random() < 0.7) opts.unitSeparator = getRandomItem(unitSeps);
    if (Math.random() < 0.7) opts.unit = getRandomItem(units);
    return opts;
}

function getRandomInput() {
    const r = Math.random();
    if (r < 0.4) return getRandomString();
    if (r < 0.8) return getRandomNumber();
    if (r < 0.85) return null;
    if (r < 0.90) return undefined;
    if (r < 0.95) return Math.random() < 0.5;
    return {};
}

function areResultsEqual(a, b) {
    if (a === b) return true;
    if (typeof a === 'number' && typeof b === 'number') {
        if (Number.isNaN(a) && Number.isNaN(b)) return true;
        return a === b;
    }
    return false;
}

// README Specification Seed Examples
const readmeSeedCases = [
    { mode: 'bytes', val: 1024, opts: undefined },
    { mode: 'bytes', val: '1KB', opts: undefined },
    { mode: 'format', val: 1024, opts: undefined },
    { mode: 'format', val: 1000, opts: undefined },
    { mode: 'format', val: 1000, opts: { thousandsSeparator: ' ' } },
    { mode: 'format', val: 1024 * 1.7, opts: { decimalPlaces: 0 } },
    { mode: 'format', val: 1024, opts: { unitSeparator: ' ' } },
    { mode: 'parse', val: '1KB' },
    { mode: 'parse', val: '1024' },
    { mode: 'parse', val: 1024 }
];

log(`Executing README Specification Seed Suite (${readmeSeedCases.length} cases)...`);
let seedDivergences = 0;
readmeSeedCases.forEach((tc, idx) => {
    let jsRes, goRes;
    if (tc.mode === 'bytes') {
        jsRes = jsBytes(tc.val, tc.opts);
        goRes = goBytes(tc.val, tc.opts);
    } else if (tc.mode === 'format') {
        jsRes = jsBytes.format(tc.val, tc.opts);
        goRes = goBytes.format(tc.val, tc.opts);
    } else {
        jsRes = jsBytes.parse(tc.val);
        goRes = goBytes.parse(tc.val);
    }

    const match = areResultsEqual(jsRes, goRes);
    if (!match) {
        seedDivergences++;
        log(`  [FAIL] Seed #${idx + 1} (${tc.mode}): val=${JSON.stringify(tc.val)} | JS=${jsRes} | Go=${goRes}`);
    } else {
        log(`  [PASS] Seed #${idx + 1} (${tc.mode}): val=${JSON.stringify(tc.val)} => ${jsRes}`);
    }
});

if (seedDivergences === 0) {
    log(`All README Specification Seed Cases PASSED (100% parity).\n`);
} else {
    log(`WARNING: ${seedDivergences} seed cases failed!\n`);
}

// Fuzz Loop Setup
let totalTests = 0;
let parseTests = 0;
let formatTests = 0;
let bytesTests = 0;
let divergences = seedDivergences;
const failures = [];

const startTime = Date.now();
const targetEndTime = startTime + (targetDurationSec * 1000);
let lastReportTime = startTime;

log(`Running randomized fuzz loop for ${targetDurationSec}s...`);

while (Date.now() < targetEndTime) {
    const mode = Math.floor(Math.random() * 3);
    totalTests++;

    if (mode === 0) {
        // Test bytes(val, options)
        bytesTests++;
        const val = getRandomInput();
        const opts = getRandomOptions();
        
        let jsRes, goRes;
        try { jsRes = jsBytes(val, opts); } catch (e) { jsRes = 'ERROR: ' + e.message; }
        try { goRes = goBytes(val, opts); } catch (e) { goRes = 'ERROR: ' + e.message; }

        if (!areResultsEqual(jsRes, goRes)) {
            divergences++;
            failures.push({ mode: 'bytes', val, opts, jsRes, goRes });
        }
    } else if (mode === 1) {
        // Test bytes.parse(val)
        parseTests++;
        const val = getRandomInput();

        let jsRes, goRes;
        try { jsRes = jsBytes.parse(val); } catch (e) { jsRes = 'ERROR: ' + e.message; }
        try { goRes = goBytes.parse(val); } catch (e) { goRes = 'ERROR: ' + e.message; }

        if (!areResultsEqual(jsRes, goRes)) {
            divergences++;
            failures.push({ mode: 'parse', val, jsRes, goRes });
        }
    } else {
        // Test bytes.format(val, options)
        formatTests++;
        const val = getRandomNumber();
        const opts = getRandomOptions();

        let jsRes, goRes;
        try { jsRes = jsBytes.format(val, opts); } catch (e) { jsRes = 'ERROR: ' + e.message; }
        try { goRes = goBytes.format(val, opts); } catch (e) { goRes = 'ERROR: ' + e.message; }

        if (!areResultsEqual(jsRes, goRes)) {
            divergences++;
            failures.push({ mode: 'format', val, opts, jsRes, goRes });
        }
    }

    // Report progress every 10 seconds
    const now = Date.now();
    if (now - lastReportTime >= 10000) {
        const elapsedSec = ((now - startTime) / 1000).toFixed(1);
        const rate = (totalTests / (now - startTime) * 1000).toFixed(0);
        log(`[${elapsedSec}s] Executed ${totalTests.toLocaleString()} iterations (${rate} ops/sec) | Divergences: ${divergences}`);
        lastReportTime = now;
    }
}

const totalDurationSec = ((Date.now() - startTime) / 1000).toFixed(2);
const opsPerSec = (totalTests / totalDurationSec).toFixed(0);

log(`\n================================================================`);
log(`  FUZZING COMPLETE`);
log(`================================================================`);
log(`Total Duration       : ${totalDurationSec} seconds`);
log(`Total Executions     : ${totalTests.toLocaleString()}`);
log(`  - bytes() calls    : ${bytesTests.toLocaleString()}`);
log(`  - parse() calls    : ${parseTests.toLocaleString()}`);
log(`  - format() calls   : ${formatTests.toLocaleString()}`);
log(`Execution Rate       : ${opsPerSec} ops/sec`);
log(`Total Divergences    : ${divergences}`);
log(`================================================================`);

if (divergences === 0) {
    log(`SUCCESS: 0 divergences detected across ${totalTests.toLocaleString()} differential tests.`);
    log(`The Go implementation (bytesutil) has 100% behavioral parity with index.js.`);
} else {
    log(`FAILURE: ${divergences} divergences found!`);
    log(`First 5 Mismatches:`);
    failures.slice(0, 5).forEach((f, idx) => {
        log(`  Match #${idx + 1}: ${JSON.stringify(f)}`);
    });
}

log(`\nEnd Time: ${new Date().toISOString()}`);
logStream.end();
