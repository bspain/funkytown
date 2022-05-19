import { test, expect } from '@playwright/test';

test ('Search for album using search button', async ({ page, isMobile }) => {
    await page.goto('https://www.target.com');


    const stofPcr = page.locator('[data-test="@web/SiteTopOfFunnel/PageContentRenderer"]');
    await expect(stofPcr).toBeVisible();

    if (isMobile) {
        await page.locator('[data-test="@web/SearchInputMobile"]').tap();
        await page.locator('[data-test="@web/SearchInputMobile"]').type('and the band marches on without slowing', { delay: 25 });
        await page.locator('[data-test="@web/SearchButtonOverlayMobile"]').click();
    
    } else {
        await page.locator('[data-test="@web/Search/SearchInput"]').type('and the band marches on without slowing', { delay: 25 });
        await page.locator('[data-test="@web/Search/SearchButton"]').click();
    }
    const results = page.locator('[data-test="resultsHeading"]');

    await expect(results).toContainText('results for');
    await expect(results).toContainText('band marches on without slowing');
})