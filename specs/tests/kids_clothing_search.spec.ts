import { test, expect } from '@playwright/test';

test ('Search for kids clothing using search button', async ({ page }) => {
    await page.goto('https://www.target.com');


    const stofPcr = page.locator('[data-test="@web/SiteTopOfFunnel/PageContentRenderer"]');
    await expect(stofPcr).toBeVisible();

    await page.locator('[data-test="@web/SearchInputMobile"]').tap();
    await page.locator('[data-test="@web/SearchInputMobile"]').type('cat and jack jumper', { delay: 25 });
    await page.locator('[data-test="@web/SearchButtonOverlayMobile"]').click();
    
    const results = page.locator('[data-test="resultsHeading"]');

    await expect(results).toContainText('results for');
    await expect(results).toContainText('jack jumper');
})