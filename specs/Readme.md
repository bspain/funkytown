# Funkytown specs
Specs for the funkytown POC using playwright

## Install dependencies
Requires:
- Node v16.15.0

Run
```
npm install
```

## Run Playwright tests
Run all specs with
```
npx playwright test
```

Or individual browsers
```
npx playwright test --project=chromium
npx playwright test --project=firefox
npx playwright test --project=webkit
```