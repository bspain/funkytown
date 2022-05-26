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
npx playwright test --project=desktop-chrome
npx playwright test --project=desktop-webkit
npx playwright test --project=desktop-firefox
npx playwright test --project=mobile-chrome
npx playwright test --project=mobile-webkit
```