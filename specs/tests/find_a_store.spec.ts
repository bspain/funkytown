import { test, expect } from '@playwright/test';

test ('find a store by zip code', async ({ page, isMobile }) => {
    await page.goto('https://www.target.com/store-locator/find-stores');


    const filterButton = page.locator('[data-test="@store-locator/StoreLocatorPage/FilterByLocationBtn"]');
    await expect(filterButton).toBeVisible();

    if (isMobile) {
        await filterButton.tap();
    } else {
        await filterButton.click();
    }
    await page.locator('[data-test="@store-locator/StoreSearchForm"] #zipcode').type('55403', { delay: 25 });

    if (isMobile) {
        await page.locator('[data-test="@store-locator/StoreSearchForm/FindStoreBtn"]').tap();
    } else {
        await page.locator('[data-test="@store-locator/StoreSearchForm/FindStoreBtn"]').click();
    }

    const store = await page.locator('[data-test="@store-locator/StoreCard"] >> nth=0');
    const storeTitle = await store.locator('[data-test="@store-locator/StoreCard/StoreCardTitle"]').textContent();

    expect(storeTitle).toContain("Mpls Nicollet Mall");
})