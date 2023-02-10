// setup jest
import 'jest';
import {
    Config,
    NeverPullPolicy,
    saveToCache,
    New,
    WithCacheDirectory,
    WithPullPolicy,
    WithApiKey, WithApiBaseUrl, WithOrganization, computeSha256
} from './index';
import {ScaleFunc} from "@loopholelabs/scalefile";

const testingApiBaseUrl = "https://api.dev.scale.sh/v1"

/**
 * @jest-environment node
 */
test('TestPulldownCache', async () => {
    const config: Config = {
        cacheDirectory: "testCache",
        pullPolicy: NeverPullPolicy,
        apiKey: "123",
        apiBaseUrl: testingApiBaseUrl,
        organization: "testOrg",
    }

    const scaleFunc = new ScaleFunc(
        "v1alpha",
        "TestFunction",
        "1",
        "signature1",
        "go",
        Buffer.from("Hello world"),
    );

    const func = "TestFunction"
    const tag = "1"
    const hash = await computeSha256(scaleFunc.Function);
    saveToCache(func, tag, hash, config, scaleFunc)

    // Try reading a scalefile from the cache
    const newScaleFunc = await New(
        func,
        tag,
        WithCacheDirectory("testCache"),
        WithPullPolicy(config.pullPolicy),
        WithApiKey(config.apiKey)
    );

    expect(newScaleFunc.Version).toBe(scaleFunc.Version);
    expect(newScaleFunc.Name).toBe(scaleFunc.Name);
    expect(newScaleFunc.Tag).toBe(scaleFunc.Tag);
    expect(newScaleFunc.Signature).toBe(scaleFunc.Signature);
    expect(newScaleFunc.Language).toBe(scaleFunc.Language);
    expect(newScaleFunc.Function).toEqual(scaleFunc.Function);
    expect(newScaleFunc.Size).toBe(61);
});

test('TestRegistryDownload', async () => {
   // get api key from environment
    const apiKey = process.env.SCALE_API_KEY;
    if (apiKey == undefined) {
        console.log("SCALE_API_KEY not set, skipping test")
        return
    }

    const scaleFunc = await New(
        "TestRegistryDownload",
        "1",
        WithApiKey(apiKey),
        WithApiBaseUrl(testingApiBaseUrl),
        WithOrganization("alex"),
    );

    expect(scaleFunc.Version).toBe("v1alpha");
    expect(scaleFunc.Name).toBe("TestRegistryDownload");
    expect(scaleFunc.Tag).toBe("1");
    expect(scaleFunc.Signature).toBe("signature1");
    expect(scaleFunc.Language).toBe("go");
});
