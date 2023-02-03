// setup jest
import 'jest';
import {
    Config,
    NeverPullPolicy,
    saveToCache,
    New,
    WithCacheDirectory,
    WithPullPolicy,
    WithApiKey
} from '../../src/registry';
import {ScaleFile} from "@loopholelabs/scalefile";


test('TestPulldownCache', () => {
    const config: Config = {
        cacheDirectory: "testCache",
        pullPolicy: NeverPullPolicy,
        apiKey: "123",
    }

    const scaleFile = new ScaleFile();
    scaleFile.Name = "Test1";
    scaleFile.Signature = "signature1";
    scaleFile.Source = "Hello world";

    const func = "Testfunction"
    saveToCache(func, scaleFile, config)

    // Try reading a scalefile from the cache
    const newScaleFile = New(
        func,
        WithCacheDirectory("testCache"),
        WithPullPolicy(config.pullPolicy),
        WithApiKey(config.apiKey)
    );

    expect(newScaleFile).toEqual(scaleFile);
});
